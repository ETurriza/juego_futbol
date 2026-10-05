package persistencia

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"time"

	_ "modernc.org/sqlite" // driver "sqlite", Go puro

	"github.com/ETurriza/juego_futbol/internal/aplicacion"
	"github.com/ETurriza/juego_futbol/internal/modelo"
)

// Repositorio guarda las partidas en una base SQLite. Es seguro para uso
// concurrente: usa una sola conexión, por lo que las operaciones se serializan.
type Repositorio struct {
	db *sql.DB
	// ahora da la hora de cada guardado; las pruebas pueden reemplazarla.
	ahora func() time.Time
}

var _ aplicacion.RepositorioPartidas = (*Repositorio)(nil)

// Abrir abre (y crea si hace falta, junto con su carpeta) la base de datos en
// ruta y la lleva al esquema más reciente.
func Abrir(ctx context.Context, ruta string) (*Repositorio, error) {
	abs, err := filepath.Abs(ruta)
	if err != nil {
		return nil, err
	}
	if err := os.MkdirAll(filepath.Dir(abs), 0o755); err != nil {
		return nil, fmt.Errorf("crear la carpeta de la base de datos: %w", err)
	}

	q := url.Values{}
	q.Add("_pragma", "foreign_keys(1)")
	q.Add("_pragma", "busy_timeout(5000)")
	q.Add("_pragma", "journal_mode(WAL)")
	dsn := (&url.URL{Scheme: "file", Path: abs, RawQuery: q.Encode()}).String()

	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(1)
	if err := migrar(ctx, db); err != nil {
		db.Close()
		return nil, err
	}
	return &Repositorio{db: db, ahora: time.Now}, nil
}

// Cerrar libera la base de datos.
func (r *Repositorio) Cerrar() error { return r.db.Close() }

// Guardar crea o reemplaza por completo la partida de la ranura.
func (r *Repositorio) Guardar(ctx context.Context, ranura string, g aplicacion.Guardado) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if err := aplicacion.ValidarRanura(ranura); err != nil {
		return err
	}
	resumen, err := g.Resumen(ranura)
	if err != nil {
		return err
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback() //nolint:errcheck // tras Commit es un no-op

	// Borrar primero arrastra, por ON DELETE CASCADE, equipos, jugadores y
	// resultados de una versión anterior.
	if _, err := tx.ExecContext(ctx, "DELETE FROM partidas WHERE ranura = ?", ranura); err != nil {
		return err
	}
	if _, err := tx.ExecContext(ctx,
		`INSERT INTO partidas (ranura, semilla, usuario, equipo_usuario, jornadas_jugadas, total_jornadas,
		                       actualizada, temporada, proximo_id)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		ranura, g.Semilla, g.Usuario, resumen.Equipo, resumen.Jornada, resumen.TotalJornadas,
		r.ahora().UnixMilli(), g.Numero, g.ProximoID); err != nil {
		return err
	}
	for _, h := range g.Historial {
		if _, err := tx.ExecContext(ctx,
			`INSERT INTO historial (ranura, numero, campeon, puesto_usuario, puntos_usuario)
			 VALUES (?, ?, ?, ?, ?)`,
			ranura, h.Numero, h.Campeon, h.PuestoUsuario, h.PuntosUsuario); err != nil {
			return err
		}
	}

	insEquipo, err := tx.PrepareContext(ctx, "INSERT INTO equipos (ranura, indice, nombre) VALUES (?, ?, ?)")
	if err != nil {
		return err
	}
	defer insEquipo.Close()
	insJugador, err := tx.PrepareContext(ctx,
		`INSERT INTO jugadores (ranura, equipo, orden, id, nombre, edad, posicion,
		                        ritmo, tiro, pase, regate, defensa, fisico, reflejos)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`)
	if err != nil {
		return err
	}
	defer insJugador.Close()
	for i, e := range g.Equipos {
		if _, err := insEquipo.ExecContext(ctx, ranura, i, e.Nombre); err != nil {
			return err
		}
		for orden, j := range e.Plantilla {
			a := j.Atributos
			if _, err := insJugador.ExecContext(ctx, ranura, i, orden, j.ID, j.Nombre, j.Edad, int(j.Posicion),
				a.Ritmo, a.Tiro, a.Pase, a.Regate, a.Defensa, a.Fisico, a.Reflejos); err != nil {
				return err
			}
		}
	}

	insResultado, err := tx.PrepareContext(ctx,
		`INSERT INTO resultados (ranura, jornada, orden, local, visitante, goles_local, goles_visitante)
		 VALUES (?, ?, ?, ?, ?, ?, ?)`)
	if err != nil {
		return err
	}
	defer insResultado.Close()
	for jornada, partidos := range g.Resultados {
		for orden, p := range partidos {
			if _, err := insResultado.ExecContext(ctx, ranura, jornada, orden,
				p.Local, p.Visitante, p.GolesLocal, p.GolesVisitante); err != nil {
				return err
			}
			if err := guardarDetalle(ctx, tx, ranura, jornada, orden, p.Detalle); err != nil {
				return err
			}
		}
	}
	for orden, a := range g.Archivo {
		e := a.Estadisticas
		if _, err := tx.ExecContext(ctx,
			`INSERT INTO estadisticas_temporada (ranura, orden, temporada, jugador, nombre, equipo, posicion, edad,
			        partidos, titularidades, minutos, goles, asistencias, amarillas, rojas, imbatidas, encajados,
			        suma_valoracion)
			 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			ranura, orden, a.Temporada, a.Jugador, a.Nombre, a.Equipo, int(a.Posicion), a.Edad,
			e.Partidos, e.Titularidades, e.Minutos, e.Goles, e.Asistencias, e.Amarillas, e.Rojas,
			e.PorteriasImbatidas, e.GolesEncajados, e.SumaValoracion); err != nil {
			return err
		}
	}
	return tx.Commit()
}

// guardarDetalle escribe las alineaciones (omitiendo los puestos sin jugador) y
// los sucesos de un partido. Un detalle vacío no escribe nada.
func guardarDetalle(ctx context.Context, tx *sql.Tx, ranura string, jornada, orden int, d modelo.DetallePartido) error {
	for local, titulares := range map[int][modelo.TitularesPorEquipo]int{1: d.TitularesLocal, 0: d.TitularesVisitante} {
		for puesto, id := range titulares {
			if id == 0 {
				continue
			}
			if _, err := tx.ExecContext(ctx,
				`INSERT INTO alineaciones (ranura, jornada, orden, local, puesto, jugador) VALUES (?, ?, ?, ?, ?, ?)`,
				ranura, jornada, orden, local, puesto, id); err != nil {
				return err
			}
		}
	}
	for i, e := range d.Eventos {
		local := 0
		if e.Local {
			local = 1
		}
		if _, err := tx.ExecContext(ctx,
			`INSERT INTO eventos (ranura, jornada, orden, indice, minuto, tipo, local, jugador, otro)
			 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			ranura, jornada, orden, i, e.Minuto, int(e.Tipo), local, e.Jugador, e.Otro); err != nil {
			return err
		}
	}
	return nil
}

// Cargar devuelve la partida de la ranura. Devuelve aplicacion.ErrPartidaNoExiste
// (envuelto) si no hay ninguna, y un error descriptivo si los datos guardados
// son incoherentes.
func (r *Repositorio) Cargar(ctx context.Context, ranura string) (aplicacion.Guardado, error) {
	if err := ctx.Err(); err != nil {
		return aplicacion.Guardado{}, err
	}
	// Una transacción de solo lectura da una vista consistente de las tablas.
	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{ReadOnly: true})
	if err != nil {
		return aplicacion.Guardado{}, err
	}
	defer tx.Rollback() //nolint:errcheck

	var g aplicacion.Guardado
	var jornadas int
	err = tx.QueryRowContext(ctx,
		"SELECT semilla, usuario, jornadas_jugadas, temporada, proximo_id FROM partidas WHERE ranura = ?", ranura).
		Scan(&g.Semilla, &g.Usuario, &jornadas, &g.Numero, &g.ProximoID)
	if errors.Is(err, sql.ErrNoRows) {
		return aplicacion.Guardado{}, fmt.Errorf("%w: %q", aplicacion.ErrPartidaNoExiste, ranura)
	}
	if err != nil {
		return aplicacion.Guardado{}, err
	}

	if g.Historial, err = cargarHistorial(ctx, tx, ranura, g.Numero-1); err != nil {
		return aplicacion.Guardado{}, err
	}
	if g.Equipos, err = cargarEquipos(ctx, tx, ranura); err != nil {
		return aplicacion.Guardado{}, err
	}
	if g.Resultados, err = cargarResultados(ctx, tx, ranura, jornadas); err != nil {
		return aplicacion.Guardado{}, err
	}
	if err = cargarDetalles(ctx, tx, ranura, g.Resultados); err != nil {
		return aplicacion.Guardado{}, err
	}
	if g.Archivo, err = cargarArchivo(ctx, tx, ranura); err != nil {
		return aplicacion.Guardado{}, err
	}
	return g, nil
}

// cargarDetalles completa los resultados con sus alineaciones y sucesos.
func cargarDetalles(ctx context.Context, tx *sql.Tx, ranura string, resultados [][]aplicacion.ResultadoGuardado) error {
	partido := func(jornada, orden int) (*aplicacion.ResultadoGuardado, error) {
		if jornada < 0 || jornada >= len(resultados) || orden < 0 || orden >= len(resultados[jornada]) {
			return nil, fmt.Errorf("datos corruptos en %q: detalle de un partido inexistente (jornada %d, partido %d)",
				ranura, jornada+1, orden+1)
		}
		return &resultados[jornada][orden], nil
	}

	filas, err := tx.QueryContext(ctx,
		`SELECT jornada, orden, local, puesto, jugador FROM alineaciones WHERE ranura = ?
		 ORDER BY jornada, orden, local, puesto`, ranura)
	if err != nil {
		return err
	}
	defer filas.Close()
	for filas.Next() {
		var jornada, orden, local, puesto, jugador int
		if err := filas.Scan(&jornada, &orden, &local, &puesto, &jugador); err != nil {
			return err
		}
		p, err := partido(jornada, orden)
		if err != nil {
			return err
		}
		if puesto < 0 || puesto >= modelo.TitularesPorEquipo {
			return fmt.Errorf("datos corruptos en %q: puesto %d fuera de la alineacion", ranura, puesto)
		}
		if local == 1 {
			p.Detalle.TitularesLocal[puesto] = jugador
		} else {
			p.Detalle.TitularesVisitante[puesto] = jugador
		}
	}
	if err := filas.Err(); err != nil {
		return err
	}

	evs, err := tx.QueryContext(ctx,
		`SELECT jornada, orden, indice, minuto, tipo, local, jugador, otro FROM eventos WHERE ranura = ?
		 ORDER BY jornada, orden, indice`, ranura)
	if err != nil {
		return err
	}
	defer evs.Close()
	for evs.Next() {
		var jornada, orden, indice, tipo, local int
		var e modelo.Evento
		if err := evs.Scan(&jornada, &orden, &indice, &e.Minuto, &tipo, &local, &e.Jugador, &e.Otro); err != nil {
			return err
		}
		p, err := partido(jornada, orden)
		if err != nil {
			return err
		}
		if indice != len(p.Detalle.Eventos) {
			return fmt.Errorf("datos corruptos en %q: falta un suceso del partido %d de la jornada %d",
				ranura, orden+1, jornada+1)
		}
		e.Tipo, e.Local = modelo.TipoEvento(tipo), local == 1
		p.Detalle.Eventos = append(p.Detalle.Eventos, e)
	}
	return evs.Err()
}

func cargarArchivo(ctx context.Context, tx *sql.Tx, ranura string) ([]aplicacion.EstadisticaTemporada, error) {
	filas, err := tx.QueryContext(ctx,
		`SELECT orden, temporada, jugador, nombre, equipo, posicion, edad, partidos, titularidades, minutos,
		        goles, asistencias, amarillas, rojas, imbatidas, encajados, suma_valoracion
		 FROM estadisticas_temporada WHERE ranura = ? ORDER BY orden`, ranura)
	if err != nil {
		return nil, err
	}
	defer filas.Close()
	var archivo []aplicacion.EstadisticaTemporada
	for filas.Next() {
		var orden, posicion int
		var a aplicacion.EstadisticaTemporada
		e := &a.Estadisticas
		if err := filas.Scan(&orden, &a.Temporada, &a.Jugador, &a.Nombre, &a.Equipo, &posicion, &a.Edad,
			&e.Partidos, &e.Titularidades, &e.Minutos, &e.Goles, &e.Asistencias, &e.Amarillas, &e.Rojas,
			&e.PorteriasImbatidas, &e.GolesEncajados, &e.SumaValoracion); err != nil {
			return nil, err
		}
		if orden != len(archivo) {
			return nil, fmt.Errorf("datos corruptos en %q: falta una estadistica archivada", ranura)
		}
		a.Posicion = modelo.Posicion(posicion)
		archivo = append(archivo, a)
	}
	return archivo, filas.Err()
}

func cargarHistorial(ctx context.Context, tx *sql.Tx, ranura string, esperadas int) ([]aplicacion.ResumenTemporada, error) {
	filas, err := tx.QueryContext(ctx,
		`SELECT numero, campeon, puesto_usuario, puntos_usuario
		 FROM historial WHERE ranura = ? ORDER BY numero`, ranura)
	if err != nil {
		return nil, err
	}
	defer filas.Close()
	var historial []aplicacion.ResumenTemporada
	for filas.Next() {
		var h aplicacion.ResumenTemporada
		if err := filas.Scan(&h.Numero, &h.Campeon, &h.PuestoUsuario, &h.PuntosUsuario); err != nil {
			return nil, err
		}
		if h.Numero != len(historial)+1 {
			return nil, fmt.Errorf("datos corruptos en %q: falta el resumen de la temporada %d", ranura, len(historial)+1)
		}
		historial = append(historial, h)
	}
	if err := filas.Err(); err != nil {
		return nil, err
	}
	if len(historial) != esperadas {
		return nil, fmt.Errorf("datos corruptos en %q: %d resumenes en el historial, se esperaban %d",
			ranura, len(historial), esperadas)
	}
	return historial, nil
}

func cargarEquipos(ctx context.Context, tx *sql.Tx, ranura string) ([]modelo.Equipo, error) {
	filas, err := tx.QueryContext(ctx, "SELECT indice, nombre FROM equipos WHERE ranura = ? ORDER BY indice", ranura)
	if err != nil {
		return nil, err
	}
	defer filas.Close()
	var equipos []modelo.Equipo
	for filas.Next() {
		var indice int
		var e modelo.Equipo
		if err := filas.Scan(&indice, &e.Nombre); err != nil {
			return nil, err
		}
		if indice != len(equipos) {
			return nil, fmt.Errorf("datos corruptos en %q: falta el equipo %d", ranura, len(equipos))
		}
		equipos = append(equipos, e)
	}
	if err := filas.Err(); err != nil {
		return nil, err
	}

	js, err := tx.QueryContext(ctx,
		`SELECT equipo, orden, id, nombre, edad, posicion, ritmo, tiro, pase, regate, defensa, fisico, reflejos
		 FROM jugadores WHERE ranura = ? ORDER BY equipo, orden`, ranura)
	if err != nil {
		return nil, err
	}
	defer js.Close()
	siguiente := map[int]int{} // próximo orden esperado por equipo
	for js.Next() {
		var equipo, orden, posicion int
		var j modelo.Jugador
		a := &j.Atributos
		if err := js.Scan(&equipo, &orden, &j.ID, &j.Nombre, &j.Edad, &posicion,
			&a.Ritmo, &a.Tiro, &a.Pase, &a.Regate, &a.Defensa, &a.Fisico, &a.Reflejos); err != nil {
			return nil, err
		}
		if equipo < 0 || equipo >= len(equipos) {
			return nil, fmt.Errorf("datos corruptos en %q: jugador de un equipo inexistente (%d)", ranura, equipo)
		}
		if orden != siguiente[equipo] {
			return nil, fmt.Errorf("datos corruptos en %q: falta un jugador del equipo %d", ranura, equipo)
		}
		siguiente[equipo]++
		j.Posicion = modelo.Posicion(posicion)
		equipos[equipo].Plantilla = append(equipos[equipo].Plantilla, j)
	}
	return equipos, js.Err()
}

func cargarResultados(ctx context.Context, tx *sql.Tx, ranura string, jornadas int) ([][]aplicacion.ResultadoGuardado, error) {
	filas, err := tx.QueryContext(ctx,
		`SELECT jornada, orden, local, visitante, goles_local, goles_visitante
		 FROM resultados WHERE ranura = ? ORDER BY jornada, orden`, ranura)
	if err != nil {
		return nil, err
	}
	defer filas.Close()
	var resultados [][]aplicacion.ResultadoGuardado
	if jornadas > 0 {
		resultados = make([][]aplicacion.ResultadoGuardado, jornadas)
	}
	for filas.Next() {
		var jornada, orden int
		var p aplicacion.ResultadoGuardado
		if err := filas.Scan(&jornada, &orden, &p.Local, &p.Visitante, &p.GolesLocal, &p.GolesVisitante); err != nil {
			return nil, err
		}
		if jornada < 0 || jornada >= jornadas {
			return nil, fmt.Errorf("datos corruptos en %q: resultado de la jornada %d, solo hay %d jugadas",
				ranura, jornada+1, jornadas)
		}
		if orden != len(resultados[jornada]) {
			return nil, fmt.Errorf("datos corruptos en %q: falta un partido de la jornada %d", ranura, jornada+1)
		}
		resultados[jornada] = append(resultados[jornada], p)
	}
	return resultados, filas.Err()
}

// Listar devuelve las partidas guardadas, la más reciente primero.
func (r *Repositorio) Listar(ctx context.Context) ([]aplicacion.ResumenPartida, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	filas, err := r.db.QueryContext(ctx,
		`SELECT ranura, temporada, equipo_usuario, jornadas_jugadas, total_jornadas, actualizada
		 FROM partidas ORDER BY actualizada DESC, ranura ASC`)
	if err != nil {
		return nil, err
	}
	defer filas.Close()
	var out []aplicacion.ResumenPartida
	for filas.Next() {
		var s aplicacion.ResumenPartida
		var ms int64
		if err := filas.Scan(&s.Ranura, &s.Temporada, &s.Equipo, &s.Jornada, &s.TotalJornadas, &ms); err != nil {
			return nil, err
		}
		s.Actualizada = time.UnixMilli(ms)
		out = append(out, s)
	}
	return out, filas.Err()
}

// Borrar elimina la partida de la ranura. Devuelve aplicacion.ErrPartidaNoExiste
// (envuelto) si no existe.
func (r *Repositorio) Borrar(ctx context.Context, ranura string) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	res, err := r.db.ExecContext(ctx, "DELETE FROM partidas WHERE ranura = ?", ranura)
	if err != nil {
		return err
	}
	if n, err := res.RowsAffected(); err != nil {
		return err
	} else if n == 0 {
		return fmt.Errorf("%w: %q", aplicacion.ErrPartidaNoExiste, ranura)
	}
	return nil
}
