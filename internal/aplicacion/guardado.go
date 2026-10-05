package aplicacion

import (
	"fmt"

	"github.com/ETurriza/juego_futbol/internal/liga"
	"github.com/ETurriza/juego_futbol/internal/modelo"
)

// ResultadoGuardado es un partido jugado, con equipos por índice.
type ResultadoGuardado struct {
	Local          int
	Visitante      int
	GolesLocal     int
	GolesVisitante int
	// Detalle son las alineaciones y los sucesos. Es vacío en los partidos
	// guardados antes de que existieran las estadísticas.
	Detalle modelo.DetallePartido
}

// Guardado es una foto de una carrera hecha solo con tipos simples, lista para
// persistir. Guarda el estado real (equipos, jugadores y resultados) y no solo
// la semilla, porque los fichajes harán que las plantillas dejen de poder
// derivarse de ella. El calendario no se guarda: es determinista.
type Guardado struct {
	Semilla int64
	// Usuario es el índice, en Equipos, del equipo del usuario.
	Usuario int
	// Numero es el número de la temporada en curso, desde 1.
	Numero int
	// ProximoID es la ID del próximo jugador que se cree; siempre mayor que la
	// de cualquier jugador existente.
	ProximoID int
	// Historial tiene un resumen por cada temporada terminada.
	Historial []ResumenTemporada
	// Archivo tiene las estadísticas de cada jugador de la liga en cada
	// temporada terminada (solo de quienes jugaron).
	Archivo []EstadisticaTemporada
	Equipos []modelo.Equipo
	// Resultados[i] son los partidos de la jornada i, en el orden del
	// calendario; su longitud es la cantidad de jornadas jugadas.
	Resultados [][]ResultadoGuardado
}

// Exportar devuelve una copia independiente del estado de la carrera.
func (c *Carrera) Exportar() Guardado {
	g := Guardado{
		Semilla:   c.Semilla,
		Usuario:   c.Usuario,
		Numero:    c.Numero,
		ProximoID: c.ProximoID,
		Historial: append([]ResumenTemporada(nil), c.Historial...),
		Archivo:   append([]EstadisticaTemporada(nil), c.Archivo...),
		Equipos:   clonarEquipos(c.Temporada.Equipos),
	}
	for _, jornada := range c.Temporada.Resultados {
		rs := make([]ResultadoGuardado, len(jornada))
		for i, r := range jornada {
			rs[i] = ResultadoGuardado{
				Local: r.Local, Visitante: r.Visitante,
				GolesLocal: r.GolesLocal, GolesVisitante: r.GolesVisitante,
				Detalle: clonarDetalle(r.Detalle),
			}
		}
		g.Resultados = append(g.Resultados, rs)
	}
	return g
}

// Importar reconstruye una carrera a partir de un Guardado. Valida que los
// datos sean coherentes: equipos y jugadores válidos, y resultados que
// correspondan al calendario.
func Importar(g Guardado) (*Carrera, error) {
	if g.Usuario < 0 || g.Usuario >= len(g.Equipos) {
		return nil, fmt.Errorf("guardado invalido: usuario %d fuera de rango para %d equipos",
			g.Usuario, len(g.Equipos))
	}
	if g.Numero < 1 {
		return nil, fmt.Errorf("guardado invalido: numero de temporada %d, debe ser al menos 1", g.Numero)
	}
	if len(g.Historial) != g.Numero-1 {
		return nil, fmt.Errorf("guardado invalido: temporada %d con %d resumenes en el historial, se esperaban %d",
			g.Numero, len(g.Historial), g.Numero-1)
	}
	for i, h := range g.Historial {
		if h.Numero != i+1 {
			return nil, fmt.Errorf("guardado invalido: el resumen %d del historial es de la temporada %d",
				i+1, h.Numero)
		}
	}
	for i, a := range g.Archivo {
		if a.Temporada < 1 || a.Temporada >= g.Numero {
			return nil, fmt.Errorf("guardado invalido: la estadistica %d es de la temporada %d, que no esta terminada",
				i+1, a.Temporada)
		}
		if a.Jugador <= 0 || a.Partidos < 1 {
			return nil, fmt.Errorf("guardado invalido: la estadistica %d del archivo no tiene jugador o partidos", i+1)
		}
	}
	ids := map[int]string{}
	equipoDe := map[int]int{} // ID de jugador -> índice de su equipo
	for i, e := range g.Equipos {
		if err := e.Validar(); err != nil {
			return nil, fmt.Errorf("guardado invalido: equipo %d: %w", i, err)
		}
		for _, j := range e.Plantilla {
			if otro, ok := ids[j.ID]; ok {
				return nil, fmt.Errorf("guardado invalido: ID de jugador %d repetida en %q y %q",
					j.ID, otro, e.Nombre)
			}
			ids[j.ID] = e.Nombre
			equipoDe[j.ID] = i
			if j.ID >= g.ProximoID {
				return nil, fmt.Errorf("guardado invalido: el jugador %d de %q no es menor que ProximoID (%d)",
					j.ID, e.Nombre, g.ProximoID)
			}
		}
	}

	temporada, err := liga.Nueva(clonarEquipos(g.Equipos))
	if err != nil {
		return nil, fmt.Errorf("guardado invalido: %w", err)
	}
	if len(g.Resultados) > len(temporada.Calendario) {
		return nil, fmt.Errorf("guardado invalido: %d jornadas jugadas, el calendario tiene %d",
			len(g.Resultados), len(temporada.Calendario))
	}
	for i, jornada := range g.Resultados {
		esperados := temporada.Calendario[i]
		if len(jornada) != len(esperados) {
			return nil, fmt.Errorf("guardado invalido: jornada %d con %d partidos, se esperaban %d",
				i+1, len(jornada), len(esperados))
		}
		rs := make([]liga.Resultado, len(jornada))
		for k, r := range jornada {
			p := liga.Partido{Local: r.Local, Visitante: r.Visitante}
			if p != esperados[k] {
				return nil, fmt.Errorf("guardado invalido: jornada %d partido %d no coincide con el calendario",
					i+1, k+1)
			}
			if r.GolesLocal < 0 || r.GolesVisitante < 0 {
				return nil, fmt.Errorf("guardado invalido: jornada %d partido %d con goles negativos",
					i+1, k+1)
			}
			if err := validarDetalle(r, equipoDe); err != nil {
				return nil, fmt.Errorf("guardado invalido: jornada %d partido %d: %w", i+1, k+1, err)
			}
			rs[k] = liga.Resultado{
				Partido: p, GolesLocal: r.GolesLocal, GolesVisitante: r.GolesVisitante,
				Detalle: clonarDetalle(r.Detalle),
			}
		}
		temporada.Resultados = append(temporada.Resultados, rs)
	}
	return &Carrera{
		Semilla:   g.Semilla,
		Temporada: temporada,
		Usuario:   g.Usuario,
		Numero:    g.Numero,
		Historial: append([]ResumenTemporada(nil), g.Historial...),
		Archivo:   append([]EstadisticaTemporada(nil), g.Archivo...),
		ProximoID: g.ProximoID,
	}, nil
}

// Resumen devuelve los datos de listado de la partida guardada en la ranura
// dada, sin la fecha de actualización, que pone el repositorio.
func (g Guardado) Resumen(ranura string) (ResumenPartida, error) {
	if g.Usuario < 0 || g.Usuario >= len(g.Equipos) {
		return ResumenPartida{}, fmt.Errorf("guardado invalido: usuario %d fuera de rango", g.Usuario)
	}
	calendario, err := liga.Calendario(len(g.Equipos))
	if err != nil {
		return ResumenPartida{}, fmt.Errorf("guardado invalido: %w", err)
	}
	return ResumenPartida{
		Ranura:        ranura,
		Temporada:     g.Numero,
		Equipo:        g.Equipos[g.Usuario].Nombre,
		Jornada:       len(g.Resultados),
		TotalJornadas: len(calendario),
	}, nil
}

func (g Guardado) clonar() Guardado {
	c := Guardado{
		Semilla:   g.Semilla,
		Usuario:   g.Usuario,
		Numero:    g.Numero,
		ProximoID: g.ProximoID,
		Historial: append([]ResumenTemporada(nil), g.Historial...),
		Archivo:   append([]EstadisticaTemporada(nil), g.Archivo...),
		Equipos:   clonarEquipos(g.Equipos),
	}
	for _, jornada := range g.Resultados {
		copia := append([]ResultadoGuardado(nil), jornada...)
		for i := range copia {
			copia[i].Detalle = clonarDetalle(copia[i].Detalle)
		}
		c.Resultados = append(c.Resultados, copia)
	}
	return c
}

// clonarDetalle copia los sucesos para no compartir memoria; sin sucesos queda
// nil, igual que en un detalle recién creado.
func clonarDetalle(d modelo.DetallePartido) modelo.DetallePartido {
	d.Eventos = append([]modelo.Evento(nil), d.Eventos...)
	return d
}

// validarDetalle comprueba que el detalle de un partido sea coherente con su
// marcador y con los equipos que lo jugaron. Un detalle vacío (partido anterior
// a las estadísticas) es válido.
func validarDetalle(r ResultadoGuardado, equipoDe map[int]int) error {
	if r.Detalle.Vacio() {
		return nil
	}
	if gl, gv := r.Detalle.Goles(); gl != r.GolesLocal || gv != r.GolesVisitante {
		return fmt.Errorf("los sucesos suman %d-%d y el marcador es %d-%d", gl, gv, r.GolesLocal, r.GolesVisitante)
	}
	partes, err := r.Detalle.Participaciones()
	if err != nil {
		return err
	}
	for _, p := range partes {
		esperado := r.Visitante
		if p.Local {
			esperado = r.Local
		}
		if equipo, ok := equipoDe[p.Jugador]; !ok || equipo != esperado {
			return fmt.Errorf("el jugador %d no es del equipo que lo alinea", p.Jugador)
		}
	}
	return nil
}

func clonarEquipos(equipos []modelo.Equipo) []modelo.Equipo {
	out := make([]modelo.Equipo, len(equipos))
	for i, e := range equipos {
		out[i] = modelo.Equipo{Nombre: e.Nombre, Plantilla: append([]modelo.Jugador(nil), e.Plantilla...)}
	}
	return out
}
