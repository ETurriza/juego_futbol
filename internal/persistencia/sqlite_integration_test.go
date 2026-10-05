//go:build integration

package persistencia

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/ETurriza/juego_futbol/internal/aplicacion"
	"github.com/ETurriza/juego_futbol/internal/aplicacion/contrato"
	"github.com/ETurriza/juego_futbol/internal/modelo"
)

func abrirTemporal(t *testing.T) *Repositorio {
	t.Helper()
	r, err := Abrir(context.Background(), filepath.Join(t.TempDir(), "juego.db"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { r.Cerrar() })
	return r
}

func guardadoDePrueba(t *testing.T, semilla int64, jornadas int) aplicacion.Guardado {
	t.Helper()
	c, err := aplicacion.NuevaCarrera(semilla, 10)
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < jornadas; i++ {
		if _, err := c.AvanzarJornada(); err != nil {
			t.Fatal(err)
		}
	}
	return c.Exportar()
}

func TestRepositorioSQLiteCumpleElContrato(t *testing.T) {
	contrato.Ejecutar(t, func(t *testing.T) aplicacion.RepositorioPartidas {
		return abrirTemporal(t)
	})
}

func TestLosDatosSobrevivenAReabrir(t *testing.T) {
	ctx := context.Background()
	ruta := filepath.Join(t.TempDir(), "juego.db")
	original := guardadoDePrueba(t, 11, 9)

	r, err := Abrir(ctx, ruta)
	if err != nil {
		t.Fatal(err)
	}
	if err := r.Guardar(ctx, "principal", original); err != nil {
		t.Fatal(err)
	}
	if err := r.Cerrar(); err != nil {
		t.Fatal(err)
	}

	r, err = Abrir(ctx, ruta)
	if err != nil {
		t.Fatal(err)
	}
	defer r.Cerrar()
	cargado, err := r.Cargar(ctx, "principal")
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(original, cargado) {
		t.Error("tras reabrir la base, lo cargado no es identico a lo guardado")
	}
}

func TestCarreraGuardadaSeReanudaIgualQueLaContinua(t *testing.T) {
	ctx := context.Background()
	r := abrirTemporal(t)
	continua, _ := aplicacion.NuevaCarrera(21, 10)
	for i := 0; i < 5; i++ {
		continua.AvanzarJornada()
	}
	if err := aplicacion.GuardarCarrera(ctx, r, "p", continua); err != nil {
		t.Fatal(err)
	}
	reanudada, err := aplicacion.CargarCarrera(ctx, r, "p")
	if err != nil {
		t.Fatal(err)
	}
	for !continua.Terminada() {
		a, _ := continua.AvanzarJornada()
		b, err := reanudada.AvanzarJornada()
		if err != nil || !reflect.DeepEqual(a, b) {
			t.Fatalf("la jornada %d difiere tras cargar desde SQLite: %v", continua.Jornada(), err)
		}
	}
}

func TestAbrirCreaLaCarpetaYAceptaRutasConCaracteresEspeciales(t *testing.T) {
	ruta := filepath.Join(t.TempDir(), "carpeta nueva", "mi base #1 ?ñ", "juego.db")
	r, err := Abrir(context.Background(), ruta)
	if err != nil {
		t.Fatal(err)
	}
	defer r.Cerrar()
	if _, err := os.Stat(ruta); err != nil {
		t.Errorf("la base deberia existir en %q: %v", ruta, err)
	}
	if err := r.Guardar(context.Background(), "p", guardadoDePrueba(t, 1, 1)); err != nil {
		t.Error(err)
	}
}

func TestAbrirFallaSiLaRutaNoSePuedeCrear(t *testing.T) {
	archivo := filepath.Join(t.TempDir(), "es-un-archivo")
	if err := os.WriteFile(archivo, []byte("x"), 0o600); err != nil {
		t.Fatal(err)
	}
	// Una carpeta que cuelga de un archivo regular no se puede crear.
	if _, err := Abrir(context.Background(), filepath.Join(archivo, "juego.db")); err == nil {
		t.Error("deberia fallar")
	}
}

func versionEsquema(t *testing.T, r *Repositorio) int {
	t.Helper()
	var v int
	if err := r.db.QueryRow("PRAGMA user_version").Scan(&v); err != nil {
		t.Fatal(err)
	}
	return v
}

func TestMigraciones(t *testing.T) {
	ctx := context.Background()
	ruta := filepath.Join(t.TempDir(), "juego.db")

	r, err := Abrir(ctx, ruta)
	if err != nil {
		t.Fatal(err)
	}
	if v := versionEsquema(t, r); v != len(migraciones) {
		t.Errorf("version = %d, se esperaba %d", v, len(migraciones))
	}
	g := guardadoDePrueba(t, 1, 2)
	if err := r.Guardar(ctx, "p", g); err != nil {
		t.Fatal(err)
	}

	// Migrar de nuevo no cambia nada ni borra datos.
	if err := migrar(ctx, r.db); err != nil {
		t.Fatalf("migrar de nuevo: %v", err)
	}
	if _, err := r.Cargar(ctx, "p"); err != nil {
		t.Errorf("los datos deberian seguir ahi: %v", err)
	}

	// Una base de una versión posterior se rechaza al abrirla.
	if _, err := r.db.Exec(fmt.Sprintf("PRAGMA user_version = %d", len(migraciones)+1)); err != nil {
		t.Fatal(err)
	}
	r.Cerrar()
	if _, err := Abrir(ctx, ruta); err == nil || !strings.Contains(err.Error(), "mas nueva") {
		t.Errorf("err = %v, se esperaba rechazo de una version mas nueva", err)
	}
}

func TestBorrarEliminaTodasLasTablas(t *testing.T) {
	ctx := context.Background()
	r := abrirTemporal(t)
	if err := r.Guardar(ctx, "a", guardadoDePrueba(t, 1, 3)); err != nil {
		t.Fatal(err)
	}
	if err := r.Guardar(ctx, "b", guardadoDePrueba(t, 2, 3)); err != nil {
		t.Fatal(err)
	}
	if err := r.Borrar(ctx, "a"); err != nil {
		t.Fatal(err)
	}
	for _, tabla := range []string{"partidas", "equipos", "jugadores", "resultados"} {
		var huerfanas int
		if err := r.db.QueryRow("SELECT COUNT(*) FROM " + tabla + " WHERE ranura = 'a'").Scan(&huerfanas); err != nil {
			t.Fatal(err)
		}
		if huerfanas != 0 {
			t.Errorf("%s: quedaron %d filas de la ranura borrada", tabla, huerfanas)
		}
		var restantes int
		r.db.QueryRow("SELECT COUNT(*) FROM " + tabla + " WHERE ranura = 'b'").Scan(&restantes)
		if restantes == 0 {
			t.Errorf("%s: se perdieron las filas de la otra ranura", tabla)
		}
	}
}

func TestGuardarFallidoNoDanaLoAnterior(t *testing.T) {
	ctx := context.Background()
	r := abrirTemporal(t)
	original := guardadoDePrueba(t, 1, 3)
	if err := r.Guardar(ctx, "p", original); err != nil {
		t.Fatal(err)
	}

	// Un disparador que falla al insertar resultados hace que el guardado
	// falle a mitad de camino, despues de haber borrado y reinsertado datos.
	if _, err := r.db.Exec(`CREATE TRIGGER falla BEFORE INSERT ON resultados
		BEGIN SELECT RAISE(ABORT, 'falla simulada'); END`); err != nil {
		t.Fatal(err)
	}
	if err := r.Guardar(ctx, "p", guardadoDePrueba(t, 2, 6)); err == nil {
		t.Fatal("el guardado deberia fallar")
	}
	if _, err := r.db.Exec("DROP TRIGGER falla"); err != nil {
		t.Fatal(err)
	}

	cargado, err := r.Cargar(ctx, "p")
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(original, cargado) {
		t.Error("un guardado fallido dejo la partida anterior alterada")
	}
}

func TestCargarDetectaDatosCorruptos(t *testing.T) {
	casos := map[string]string{
		"falta un equipo":           "DELETE FROM equipos WHERE ranura = 'p' AND indice = 3",
		"falta un jugador":          "DELETE FROM jugadores WHERE ranura = 'p' AND equipo = 2 AND orden = 5",
		"jugador de otro equipo":    "UPDATE jugadores SET equipo = 77 WHERE ranura = 'p' AND equipo = 1 AND orden = 0",
		"falta un partido":          "DELETE FROM resultados WHERE ranura = 'p' AND jornada = 1 AND orden = 2",
		"resultado de mas jornadas": "UPDATE resultados SET jornada = 40 WHERE ranura = 'p' AND jornada = 0 AND orden = 0",
	}
	for nombre, sentencia := range casos {
		t.Run(nombre, func(t *testing.T) {
			ctx := context.Background()
			r := abrirTemporal(t)
			if err := r.Guardar(ctx, "p", guardadoDePrueba(t, 3, 4)); err != nil {
				t.Fatal(err)
			}
			// Se desactivan las claves foráneas solo para fabricar el daño.
			if _, err := r.db.Exec("PRAGMA foreign_keys = OFF"); err != nil {
				t.Fatal(err)
			}
			if _, err := r.db.Exec(sentencia); err != nil {
				t.Fatal(err)
			}
			_, err := r.Cargar(ctx, "p")
			if err == nil || !strings.Contains(err.Error(), "datos corruptos") {
				t.Errorf("err = %v, se esperaba 'datos corruptos'", err)
			}
		})
	}
}

func TestCargarGuardadoAlteradoLoRechazaAplicacion(t *testing.T) {
	// La base puede estar bien formada pero con datos incoherentes para el
	// juego; Importar debe rechazarlo.
	ctx := context.Background()
	r := abrirTemporal(t)
	if err := r.Guardar(ctx, "p", guardadoDePrueba(t, 3, 4)); err != nil {
		t.Fatal(err)
	}
	if _, err := r.db.Exec("UPDATE jugadores SET edad = 0 WHERE ranura = 'p' AND equipo = 0 AND orden = 0"); err != nil {
		t.Fatal(err)
	}
	if _, err := aplicacion.CargarCarrera(ctx, r, "p"); err == nil {
		t.Error("una edad invalida deberia rechazarse al cargar la carrera")
	}
}

func TestUsoConcurrente(t *testing.T) {
	ctx := context.Background()
	r := abrirTemporal(t)
	const ranuras = 8
	g := guardadoDePrueba(t, 5, 4)

	var wg sync.WaitGroup
	errs := make(chan error, ranuras*3)
	for i := 0; i < ranuras; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			ranura := fmt.Sprintf("r%d", i)
			for k := 0; k < 3; k++ {
				if err := r.Guardar(ctx, ranura, g); err != nil {
					errs <- fmt.Errorf("%s guardar: %w", ranura, err)
					return
				}
				cargado, err := r.Cargar(ctx, ranura)
				if err != nil || !reflect.DeepEqual(g, cargado) {
					errs <- fmt.Errorf("%s cargar: %v", ranura, err)
					return
				}
				if _, err := r.Listar(ctx); err != nil {
					errs <- fmt.Errorf("%s listar: %w", ranura, err)
					return
				}
			}
		}(i)
	}
	wg.Wait()
	close(errs)
	for err := range errs {
		t.Error(err)
	}
	lista, err := r.Listar(ctx)
	if err != nil || len(lista) != ranuras {
		t.Errorf("Listar = %d partidas, %v; se esperaban %d", len(lista), err, ranuras)
	}
}

func TestListarOrdenaPorFechaConRelojControlado(t *testing.T) {
	ctx := context.Background()
	r := abrirTemporal(t)
	reloj := time.Date(2026, 3, 1, 10, 0, 0, 0, time.UTC)
	r.ahora = func() time.Time { return reloj }
	g := guardadoDePrueba(t, 1, 1)
	for _, ranura := range []string{"b", "a", "c"} {
		if err := r.Guardar(ctx, ranura, g); err != nil {
			t.Fatal(err)
		}
		reloj = reloj.Add(time.Minute)
	}
	lista, err := r.Listar(ctx)
	if err != nil {
		t.Fatal(err)
	}
	var orden []string
	for _, s := range lista {
		orden = append(orden, s.Ranura)
	}
	if want := []string{"c", "a", "b"}; !reflect.DeepEqual(orden, want) {
		t.Errorf("orden = %v, se esperaba %v", orden, want)
	}
	if want := time.Date(2026, 3, 1, 10, 2, 0, 0, time.UTC); !lista[0].Actualizada.Equal(want) {
		t.Errorf("Actualizada = %v, se esperaba %v", lista[0].Actualizada, want)
	}
}

func TestErrorSiLaBaseSeCerro(t *testing.T) {
	r := abrirTemporal(t)
	r.Cerrar()
	if err := r.Guardar(context.Background(), "p", guardadoDePrueba(t, 1, 1)); err == nil {
		t.Error("guardar sobre una base cerrada deberia fallar")
	}
	if _, err := r.Cargar(context.Background(), "p"); err == nil || errors.Is(err, aplicacion.ErrPartidaNoExiste) {
		t.Errorf("cargar sobre una base cerrada deberia dar un error distinto de 'no existe': %v", err)
	}
}

// guardarComoV1 escribe un Guardado con el esquema de la migración 1, tal como
// lo hacia la version anterior del programa.
func guardarComoV1(t *testing.T, db *sql.DB, ranura string, g aplicacion.Guardado) {
	t.Helper()
	exec := func(consulta string, args ...any) {
		t.Helper()
		if _, err := db.Exec(consulta, args...); err != nil {
			t.Fatal(err)
		}
	}
	resumen, err := g.Resumen(ranura)
	if err != nil {
		t.Fatal(err)
	}
	exec(`INSERT INTO partidas (ranura, semilla, usuario, equipo_usuario, jornadas_jugadas, total_jornadas, actualizada)
	      VALUES (?, ?, ?, ?, ?, ?, ?)`,
		ranura, g.Semilla, g.Usuario, resumen.Equipo, resumen.Jornada, resumen.TotalJornadas, 1700000000000)
	for i, e := range g.Equipos {
		exec("INSERT INTO equipos (ranura, indice, nombre) VALUES (?, ?, ?)", ranura, i, e.Nombre)
		for orden, j := range e.Plantilla {
			a := j.Atributos
			exec(`INSERT INTO jugadores (ranura, equipo, orden, id, nombre, edad, posicion,
			                             ritmo, tiro, pase, regate, defensa, fisico, reflejos)
			      VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
				ranura, i, orden, j.ID, j.Nombre, j.Edad, int(j.Posicion),
				a.Ritmo, a.Tiro, a.Pase, a.Regate, a.Defensa, a.Fisico, a.Reflejos)
		}
	}
	for jornada, partidos := range g.Resultados {
		for orden, p := range partidos {
			exec(`INSERT INTO resultados (ranura, jornada, orden, local, visitante, goles_local, goles_visitante)
			      VALUES (?, ?, ?, ?, ?, ?, ?)`,
				ranura, jornada, orden, p.Local, p.Visitante, p.GolesLocal, p.GolesVisitante)
		}
	}
}

func TestMigracion2ConservaLasPartidasDeLaVersion1(t *testing.T) {
	ctx := context.Background()
	ruta := filepath.Join(t.TempDir(), "viejo.db")

	// Una base con el esquema de la version 1 y una partida dentro.
	db, err := sql.Open("sqlite", "file:"+ruta+"?_pragma=foreign_keys(1)")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := db.Exec(migraciones[0]); err != nil {
		t.Fatal(err)
	}
	if _, err := db.Exec("PRAGMA user_version = 1"); err != nil {
		t.Fatal(err)
	}
	original := guardadoDePrueba(t, 14, 6)
	// El esquema v1 no guardaba el detalle de los partidos ni el talento.
	sinDetalle(&original)
	sinTalento(&original)
	guardarComoV1(t, db, "vieja", original)
	if err := db.Close(); err != nil {
		t.Fatal(err)
	}

	// Al abrirla se migra a la version actual sin perder nada.
	r, err := Abrir(ctx, ruta)
	if err != nil {
		t.Fatal(err)
	}
	defer r.Cerrar()
	if v := versionEsquema(t, r); v != len(migraciones) {
		t.Errorf("version = %d, se esperaba %d", v, len(migraciones))
	}
	cargado, err := r.Cargar(ctx, "vieja")
	if err != nil {
		t.Fatal(err)
	}
	// Queda en la temporada 1, sin historial, y con ProximoID recalculado a
	// partir de los jugadores existentes.
	if cargado.Numero != 1 || len(cargado.Historial) != 0 {
		t.Errorf("Numero = %d, historial = %d; se esperaba 1 y vacio", cargado.Numero, len(cargado.Historial))
	}
	if !reflect.DeepEqual(original, cargado) {
		t.Error("la partida migrada no es identica a la original")
	}
	c, err := aplicacion.CargarCarrera(ctx, r, "vieja")
	if err != nil {
		t.Fatalf("la partida migrada no se puede cargar como carrera: %v", err)
	}
	// Y la carrera migrada puede seguir: termina la temporada y pasa a la 2.
	for !c.Terminada() {
		if _, err := c.AvanzarJornada(); err != nil {
			t.Fatal(err)
		}
	}
	if _, err := c.SiguienteTemporada(); err != nil {
		t.Fatal(err)
	}
	if err := aplicacion.GuardarCarrera(ctx, r, "vieja", c); err != nil {
		t.Fatal(err)
	}
	segunda, err := aplicacion.CargarCarrera(ctx, r, "vieja")
	if err != nil || segunda.Numero != 2 || len(segunda.Historial) != 1 {
		t.Errorf("tras migrar y avanzar: %v, temporada %d", err, segunda.Numero)
	}
}

func TestCargarDetectaUnHistorialIncoherente(t *testing.T) {
	casos := map[string]string{
		"falta un resumen":      "DELETE FROM historial WHERE ranura = 'p' AND numero = 2",
		"sobra un resumen":      "INSERT INTO historial VALUES ('p', 4, 'X', 1, 1)",
		"resumen mal numerado":  "UPDATE historial SET numero = 9 WHERE ranura = 'p' AND numero = 3",
		"temporada sin resumen": "UPDATE partidas SET temporada = 7 WHERE ranura = 'p'",
	}
	for nombre, sentencia := range casos {
		t.Run(nombre, func(t *testing.T) {
			ctx := context.Background()
			r := abrirTemporal(t)
			c, _ := aplicacion.NuevaCarrera(5, 10)
			for i := 0; i < 3; i++ {
				for !c.Terminada() {
					c.AvanzarJornada()
				}
				if _, err := c.SiguienteTemporada(); err != nil {
					t.Fatal(err)
				}
			}
			if err := aplicacion.GuardarCarrera(ctx, r, "p", c); err != nil {
				t.Fatal(err)
			}
			if _, err := r.db.Exec(sentencia); err != nil {
				t.Fatal(err)
			}
			_, err := r.Cargar(ctx, "p")
			if err == nil || !strings.Contains(err.Error(), "datos corruptos") {
				t.Errorf("err = %v, se esperaba 'datos corruptos'", err)
			}
		})
	}
}

func TestBorrarTambienEliminaElHistorial(t *testing.T) {
	ctx := context.Background()
	r := abrirTemporal(t)
	c, _ := aplicacion.NuevaCarrera(5, 10)
	for !c.Terminada() {
		c.AvanzarJornada()
	}
	if _, err := c.SiguienteTemporada(); err != nil {
		t.Fatal(err)
	}
	if err := aplicacion.GuardarCarrera(ctx, r, "p", c); err != nil {
		t.Fatal(err)
	}
	var antes int
	r.db.QueryRow("SELECT COUNT(*) FROM historial WHERE ranura = 'p'").Scan(&antes)
	if antes != 1 {
		t.Fatalf("deberia haber 1 resumen guardado, hay %d", antes)
	}
	if err := r.Borrar(ctx, "p"); err != nil {
		t.Fatal(err)
	}
	var despues int
	r.db.QueryRow("SELECT COUNT(*) FROM historial WHERE ranura = 'p'").Scan(&despues)
	if despues != 0 {
		t.Errorf("quedaron %d filas de historial huerfanas", despues)
	}
}

// sinDetalle quita el detalle de los partidos, como estaban guardadas las
// partidas de antes de las estadisticas.
func sinDetalle(g *aplicacion.Guardado) {
	for _, jornada := range g.Resultados {
		for k := range jornada {
			jornada[k].Detalle = modelo.DetallePartido{}
		}
	}
}

// sinTalento deja a todos los jugadores con el talento neutro, que es el que
// reciben al migrar las partidas guardadas antes de que existiera.
func sinTalento(g *aplicacion.Guardado) {
	for i := range g.Equipos {
		for k := range g.Equipos[i].Plantilla {
			g.Equipos[i].Plantilla[k].Talento = modelo.TalentoNeutro
		}
	}
}

func TestMigracion3ConservaLasPartidasDeLaVersion2(t *testing.T) {
	ctx := context.Background()
	ruta := filepath.Join(t.TempDir(), "v2.db")

	// Una partida de la version 2: con historial de temporadas, sin detalle.
	c, _ := aplicacion.NuevaCarrera(18, 10)
	for temporada := 0; temporada < 2; temporada++ {
		for !c.Terminada() {
			c.AvanzarJornada()
		}
		if _, err := c.SiguienteTemporada(); err != nil {
			t.Fatal(err)
		}
	}
	for i := 0; i < 5; i++ {
		c.AvanzarJornada()
	}
	original := c.Exportar()
	sinDetalle(&original)
	sinTalento(&original)
	original.Archivo = nil // la v2 tampoco archivaba estadisticas

	db, err := sql.Open("sqlite", "file:"+ruta+"?_pragma=foreign_keys(1)")
	if err != nil {
		t.Fatal(err)
	}
	for _, m := range migraciones[:2] {
		if _, err := db.Exec(m); err != nil {
			t.Fatal(err)
		}
	}
	if _, err := db.Exec("PRAGMA user_version = 2"); err != nil {
		t.Fatal(err)
	}
	guardarComoV1(t, db, "v2", original)
	if _, err := db.Exec("UPDATE partidas SET temporada = ?, proximo_id = ? WHERE ranura = 'v2'",
		original.Numero, original.ProximoID); err != nil {
		t.Fatal(err)
	}
	for _, h := range original.Historial {
		if _, err := db.Exec("INSERT INTO historial VALUES ('v2', ?, ?, ?, ?)",
			h.Numero, h.Campeon, h.PuestoUsuario, h.PuntosUsuario); err != nil {
			t.Fatal(err)
		}
	}
	db.Close()

	r, err := Abrir(ctx, ruta)
	if err != nil {
		t.Fatal(err)
	}
	defer r.Cerrar()
	if v := versionEsquema(t, r); v != len(migraciones) {
		t.Errorf("version = %d, se esperaba %d", v, len(migraciones))
	}
	cargado, err := r.Cargar(ctx, "v2")
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(original, cargado) {
		t.Error("la partida migrada desde la v2 no es identica")
	}

	// Se puede seguir jugando: los partidos nuevos traen detalle y estadisticas, y
	// al terminar la temporada se archivan solo los partidos con detalle.
	carrera, err := aplicacion.CargarCarrera(ctx, r, "v2")
	if err != nil {
		t.Fatal(err)
	}
	for !carrera.Terminada() {
		carrera.AvanzarJornada()
	}
	filas, _ := carrera.EstadisticasJugadores()
	jugaron := 0
	for _, f := range filas {
		if f.Partidos > 0 {
			jugaron++
		}
	}
	if jugaron == 0 {
		t.Error("los partidos nuevos deberian tener estadisticas")
	}
	if _, err := carrera.SiguienteTemporada(); err != nil {
		t.Fatal(err)
	}
	if err := aplicacion.GuardarCarrera(ctx, r, "v2", carrera); err != nil {
		t.Fatal(err)
	}
	vuelta, err := aplicacion.CargarCarrera(ctx, r, "v2")
	if err != nil || len(vuelta.Archivo) == 0 || vuelta.Numero != 4 {
		t.Errorf("tras migrar y avanzar: %v, temporada %d, archivo %d", err, vuelta.Numero, len(vuelta.Archivo))
	}
}

func TestBorrarEliminaDetalleYArchivo(t *testing.T) {
	ctx := context.Background()
	r := abrirTemporal(t)
	c, _ := aplicacion.NuevaCarrera(5, 10)
	for temporada := 0; temporada < 2; temporada++ {
		for !c.Terminada() {
			c.AvanzarJornada()
		}
		c.SiguienteTemporada()
	}
	c.AvanzarJornada()
	otra, _ := aplicacion.NuevaCarrera(6, 10)
	otra.AvanzarJornada()
	for ranura, carrera := range map[string]*aplicacion.Carrera{"p": c, "otra": otra} {
		if err := aplicacion.GuardarCarrera(ctx, r, ranura, carrera); err != nil {
			t.Fatal(err)
		}
	}
	cuenta := func(tabla, ranura string) int {
		var n int
		if err := r.db.QueryRow("SELECT COUNT(*) FROM "+tabla+" WHERE ranura = ?", ranura).Scan(&n); err != nil {
			t.Fatal(err)
		}
		return n
	}
	for _, tabla := range []string{"alineaciones", "eventos", "estadisticas_temporada"} {
		if cuenta(tabla, "p") == 0 || cuenta(tabla, "otra") == 0 && tabla != "estadisticas_temporada" {
			t.Fatalf("%s: faltan filas tras guardar", tabla)
		}
	}
	if err := r.Borrar(ctx, "p"); err != nil {
		t.Fatal(err)
	}
	for _, tabla := range []string{"alineaciones", "eventos", "estadisticas_temporada"} {
		if n := cuenta(tabla, "p"); n != 0 {
			t.Errorf("%s: quedaron %d filas de la ranura borrada", tabla, n)
		}
	}
	if cuenta("alineaciones", "otra") == 0 || cuenta("eventos", "otra") == 0 {
		t.Error("se perdieron las filas de la otra ranura")
	}
}

func TestCargarDetectaUnDetalleOUnArchivoCorruptos(t *testing.T) {
	casos := map[string]string{
		"falta un suceso":                 "DELETE FROM eventos WHERE ranura = 'p' AND jornada = 0 AND orden = 0 AND indice = 1",
		"suceso de un partido ajeno":      "UPDATE eventos SET orden = 77 WHERE ranura = 'p' AND jornada = 0 AND orden = 0 AND indice = 0",
		"alineacion de partido ajeno":     "UPDATE alineaciones SET jornada = 40 WHERE ranura = 'p' AND jornada = 0 AND orden = 0 AND local = 1 AND puesto = 0",
		"puesto fuera de la alineacion":   "UPDATE alineaciones SET puesto = 30 WHERE ranura = 'p' AND jornada = 0 AND orden = 0 AND local = 1 AND puesto = 0",
		"falta una estadistica archivada": "DELETE FROM estadisticas_temporada WHERE ranura = 'p' AND orden = 3",
	}
	for nombre, sentencia := range casos {
		t.Run(nombre, func(t *testing.T) {
			ctx := context.Background()
			r := abrirTemporal(t)
			c, _ := aplicacion.NuevaCarrera(5, 10)
			for !c.Terminada() {
				c.AvanzarJornada()
			}
			c.SiguienteTemporada()
			c.AvanzarJornada()
			c.AvanzarJornada()
			if err := aplicacion.GuardarCarrera(ctx, r, "p", c); err != nil {
				t.Fatal(err)
			}
			if _, err := r.db.Exec("PRAGMA foreign_keys = OFF"); err != nil {
				t.Fatal(err)
			}
			if _, err := r.db.Exec(sentencia); err != nil {
				t.Fatal(err)
			}
			_, err := r.Cargar(ctx, "p")
			if err == nil || !strings.Contains(err.Error(), "datos corruptos") {
				t.Errorf("err = %v, se esperaba 'datos corruptos'", err)
			}
		})
	}
}

func TestMigracion4ConservaLasPartidasDeLaVersion3(t *testing.T) {
	ctx := context.Background()
	ruta := filepath.Join(t.TempDir(), "v3.db")

	// Una partida de la versión 3 (con estadísticas, pero sin talento).
	c, _ := aplicacion.NuevaCarrera(19, 10)
	for i := 0; i < 6; i++ {
		c.AvanzarJornada()
	}
	original := c.Exportar()
	sinDetalle(&original) // el ayudante de la prueba escribe el esquema v1
	sinTalento(&original)

	db, err := sql.Open("sqlite", "file:"+ruta+"?_pragma=foreign_keys(1)")
	if err != nil {
		t.Fatal(err)
	}
	for _, m := range migraciones[:3] {
		if _, err := db.Exec(m); err != nil {
			t.Fatal(err)
		}
	}
	if _, err := db.Exec("PRAGMA user_version = 3"); err != nil {
		t.Fatal(err)
	}
	guardarComoV1(t, db, "v3", original)
	if _, err := db.Exec("UPDATE partidas SET temporada = ?, proximo_id = ? WHERE ranura = 'v3'",
		original.Numero, original.ProximoID); err != nil {
		t.Fatal(err)
	}
	db.Close()

	r, err := Abrir(ctx, ruta)
	if err != nil {
		t.Fatal(err)
	}
	defer r.Cerrar()
	if v := versionEsquema(t, r); v != len(migraciones) {
		t.Errorf("version = %d, se esperaba %d", v, len(migraciones))
	}
	cargado, err := r.Cargar(ctx, "v3")
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(original, cargado) {
		t.Error("la partida migrada desde la v3 no es identica (todos con talento neutro)")
	}
	for _, e := range cargado.Equipos {
		for _, j := range e.Plantilla {
			if j.Talento != modelo.TalentoNeutro {
				t.Fatalf("%s: talento %d, se esperaba el neutro", j.Nombre, j.Talento)
			}
		}
	}

	// Al seguir jugando, los juveniles nuevos traen su talento y se conserva al guardar.
	carrera, err := aplicacion.CargarCarrera(ctx, r, "v3")
	if err != nil {
		t.Fatal(err)
	}
	for !carrera.Terminada() {
		carrera.AvanzarJornada()
	}
	if _, err := carrera.SiguienteTemporada(); err != nil {
		t.Fatal(err)
	}
	if err := aplicacion.GuardarCarrera(ctx, r, "v3", carrera); err != nil {
		t.Fatal(err)
	}
	vuelta, err := aplicacion.CargarCarrera(ctx, r, "v3")
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(carrera.Exportar().Equipos, vuelta.Exportar().Equipos) {
		t.Error("tras migrar y avanzar, las plantillas no se recuperan identicas")
	}
	distintos := 0
	for _, e := range vuelta.Temporada.Equipos {
		for _, j := range e.Plantilla {
			if j.Talento != modelo.TalentoNeutro {
				distintos++
			}
		}
	}
	if distintos == 0 {
		t.Error("los juveniles nuevos deberian traer un talento distinto del neutro")
	}
}

func TestElTalentoViajaPorLaBaseDeDatos(t *testing.T) {
	ctx := context.Background()
	r := abrirTemporal(t)
	c, _ := aplicacion.NuevaCarrera(23, 10)
	if err := aplicacion.GuardarCarrera(ctx, r, "p", c); err != nil {
		t.Fatal(err)
	}
	var distintos, total int
	if err := r.db.QueryRow("SELECT COUNT(*), SUM(talento <> 100) FROM jugadores WHERE ranura = 'p'").Scan(&total, &distintos); err != nil {
		t.Fatal(err)
	}
	if total != 220 || distintos < 100 {
		t.Errorf("de %d jugadores guardados, %d con talento distinto del neutro; se esperaba la gran mayoria", total, distintos)
	}
	cargada, err := aplicacion.CargarCarrera(ctx, r, "p")
	if err != nil {
		t.Fatal(err)
	}
	for i, e := range c.Temporada.Equipos {
		for k, j := range e.Plantilla {
			if got := cargada.Temporada.Equipos[i].Plantilla[k].Talento; got != j.Talento {
				t.Fatalf("%s: talento %d tras cargar, era %d", j.Nombre, got, j.Talento)
			}
		}
	}
}

func TestLaAlineacionYLasFormacionesViajanPorLaBaseDeDatos(t *testing.T) {
	ctx := context.Background()
	r := abrirTemporal(t)
	c, _ := aplicacion.NuevaCarrera(31, 10)
	f := modelo.F352
	if err := c.ElegirAlineacion(c.AlineacionAutomatica(&f)); err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 5; i++ {
		c.AvanzarJornada()
	}
	if err := aplicacion.GuardarCarrera(ctx, r, "a", c); err != nil {
		t.Fatal(err)
	}
	cargada, err := aplicacion.CargarCarrera(ctx, r, "a")
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(c.Exportar(), cargada.Exportar()) {
		t.Error("la carrera cargada no es identica a la guardada")
	}
	al, manual := cargada.AlineacionVigente()
	if !manual || al.Formacion != modelo.F352 || !reflect.DeepEqual(al, c.Alineacion) {
		t.Errorf("la alineacion elegida no se recupero: %v (manual=%v)", al, manual)
	}

	// Las formaciones distintas del 4-3-3 de cada partido se guardan y se recuperan.
	var distintas int
	if err := r.db.QueryRow(`SELECT COUNT(*) FROM resultados WHERE ranura = 'a'
	                         AND (formacion_local <> 0 OR formacion_visitante <> 0)`).Scan(&distintas); err != nil {
		t.Fatal(err)
	}
	if distintas == 0 {
		t.Error("deberia haber partidos guardados con una formacion distinta del 4-3-3")
	}

	// Volver a la automatica borra la alineacion elegida de la base de datos.
	cargada.UsarAlineacionAutomatica()
	if err := aplicacion.GuardarCarrera(ctx, r, "a", cargada); err != nil {
		t.Fatal(err)
	}
	var filas int
	if err := r.db.QueryRow(`SELECT (SELECT COUNT(*) FROM alineacion_titulares WHERE ranura = 'a')
	                              + (SELECT COUNT(*) FROM alineacion_banquillo WHERE ranura = 'a')`).Scan(&filas); err != nil {
		t.Fatal(err)
	}
	if filas != 0 {
		t.Errorf("quedaron %d filas de una alineacion descartada", filas)
	}
	otra, err := aplicacion.CargarCarrera(ctx, r, "a")
	if err != nil {
		t.Fatal(err)
	}
	if _, manual := otra.AlineacionVigente(); manual {
		t.Error("tras descartarla no deberia haber alineacion elegida")
	}

	// Borrar la partida elimina tambien las tablas de la alineacion.
	if err := aplicacion.GuardarCarrera(ctx, r, "a", c); err != nil {
		t.Fatal(err)
	}
	if err := r.Borrar(ctx, "a"); err != nil {
		t.Fatal(err)
	}
	if err := r.db.QueryRow(`SELECT (SELECT COUNT(*) FROM alineacion_titulares)
	                              + (SELECT COUNT(*) FROM alineacion_banquillo)`).Scan(&filas); err != nil {
		t.Fatal(err)
	}
	if filas != 0 {
		t.Errorf("quedaron %d filas de alineacion huerfanas tras borrar", filas)
	}
}
