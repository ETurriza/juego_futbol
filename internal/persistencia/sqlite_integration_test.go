//go:build integration

package persistencia

import (
	"context"
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
