package aplicacion

import (
	"errors"
	"reflect"
	"testing"

	"github.com/ETurriza/juego_futbol/internal/modelo"
)

func nuevaCarrera(t *testing.T, semilla int64, n int) *Carrera {
	t.Helper()
	c, err := NuevaCarrera(semilla, n)
	if err != nil {
		t.Fatalf("NuevaCarrera(%d, %d): %v", semilla, n, err)
	}
	return c
}

func TestNuevaCarreraValidaNumeroDeEquipos(t *testing.T) {
	for _, n := range []int{-1, 0, 1, MaxEquipos + 1} {
		if _, err := NuevaCarrera(1, n); err == nil {
			t.Errorf("n=%d deberia dar error", n)
		}
	}
	for _, n := range []int{MinEquipos, 10, MaxEquipos} {
		if _, err := NuevaCarrera(1, n); err != nil {
			t.Errorf("n=%d rechazado: %v", n, err)
		}
	}
}

func TestNuevaCarreraReproducible(t *testing.T) {
	a, b := nuevaCarrera(t, 42, 10), nuevaCarrera(t, 42, 10)
	if !reflect.DeepEqual(a, b) {
		t.Fatal("la misma semilla deberia dar la misma carrera")
	}
	if c := nuevaCarrera(t, 43, 10); reflect.DeepEqual(a, c) {
		t.Fatal("semillas distintas deberian dar carreras distintas")
	}
}

func TestNuevaCarreraEquiposValidosYUnicos(t *testing.T) {
	for _, n := range []int{2, 10, MaxEquipos} {
		c := nuevaCarrera(t, 7, n)
		if len(c.Temporada.Equipos) != n {
			t.Fatalf("n=%d: %d equipos", n, len(c.Temporada.Equipos))
		}
		nombres := map[string]bool{}
		ids := map[int]bool{}
		for _, e := range c.Temporada.Equipos {
			if err := e.Validar(); err != nil {
				t.Errorf("n=%d: equipo invalido: %v", n, err)
			}
			if nombres[e.Nombre] {
				t.Errorf("n=%d: nombre de equipo repetido: %s", n, e.Nombre)
			}
			nombres[e.Nombre] = true
			for _, j := range e.Plantilla {
				if ids[j.ID] {
					t.Fatalf("n=%d: ID de jugador repetida en la liga: %d", n, j.ID)
				}
				ids[j.ID] = true
			}
		}
		if c.Usuario < 0 || c.Usuario >= n {
			t.Errorf("n=%d: usuario fuera de rango: %d", n, c.Usuario)
		}
		if c.NombreEquipo() != c.Temporada.Equipos[c.Usuario].Nombre {
			t.Error("NombreEquipo no corresponde al equipo del usuario")
		}
	}
}

func TestSemillaJornada(t *testing.T) {
	if semillaJornada(1, 0) != semillaJornada(1, 0) {
		t.Error("semillaJornada deberia ser determinista")
	}
	vistas := map[int64]bool{}
	for j := 0; j < 200; j++ {
		s := semillaJornada(1, j)
		if vistas[s] {
			t.Fatalf("la jornada %d repite semilla", j)
		}
		vistas[s] = true
	}
	if semillaJornada(1, 3) == semillaJornada(2, 3) {
		t.Error("semillas de carrera distintas deberian dar semillas de jornada distintas")
	}
}

func TestAvanzarJornadaHastaElFinal(t *testing.T) {
	c := nuevaCarrera(t, 5, 10)
	if c.TotalJornadas() != 18 {
		t.Fatalf("TotalJornadas = %d, se esperaban 18", c.TotalJornadas())
	}
	if c.UltimaJornada() != nil {
		t.Error("UltimaJornada deberia ser nil antes de jugar")
	}
	if _, ok := c.Campeon(); ok {
		t.Error("no deberia haber campeon antes de terminar")
	}

	for i := 1; i <= c.TotalJornadas(); i++ {
		resultados, err := c.AvanzarJornada()
		if err != nil {
			t.Fatalf("jornada %d: %v", i, err)
		}
		if c.Jornada() != i {
			t.Errorf("Jornada = %d, se esperaba %d", c.Jornada(), i)
		}
		if len(resultados) != 5 {
			t.Errorf("jornada %d: %d partidos, se esperaban 5", i, len(resultados))
		}
		if !reflect.DeepEqual(resultados, c.UltimaJornada()) {
			t.Errorf("jornada %d: UltimaJornada no coincide con lo devuelto", i)
		}
		delUsuario := 0
		for _, r := range resultados {
			if r.Local == "" || r.Visitante == "" || r.GolesLocal < 0 || r.GolesVisitante < 0 {
				t.Errorf("jornada %d: resultado invalido: %+v", i, r)
			}
			if r.EsDelUsuario {
				delUsuario++
			}
		}
		if delUsuario != 1 {
			t.Errorf("jornada %d: %d partidos del usuario, se esperaba 1", i, delUsuario)
		}
		if i < c.TotalJornadas() && c.Terminada() {
			t.Errorf("jornada %d: no deberia estar terminada", i)
		}
	}

	if !c.Terminada() {
		t.Fatal("la temporada deberia estar terminada")
	}
	if _, err := c.AvanzarJornada(); !errors.Is(err, ErrTemporadaTerminada) {
		t.Errorf("err = %v, se esperaba ErrTemporadaTerminada", err)
	}
	campeon, ok := c.Campeon()
	if !ok || campeon != c.Tabla()[0].Equipo {
		t.Errorf("Campeon = %q, %v; la tabla lidera %q", campeon, ok, c.Tabla()[0].Equipo)
	}
}

func TestCarreraReproducibleJornadaAJornada(t *testing.T) {
	a, b := nuevaCarrera(t, 9, 10), nuevaCarrera(t, 9, 10)
	for !a.Terminada() {
		ra, _ := a.AvanzarJornada()
		rb, _ := b.AvanzarJornada()
		if !reflect.DeepEqual(ra, rb) {
			t.Fatalf("jornada %d difiere entre carreras con la misma semilla", a.Jornada())
		}
	}
	if !reflect.DeepEqual(a.Tabla(), b.Tabla()) {
		t.Error("las tablas finales deberian coincidir")
	}
}

// La jornada k depende solo de la semilla y de k: una carrera que se reanuda
// con las jornadas previas ya jugadas produce lo mismo que una continua.
func TestJornadaDependeSoloDeSemillaYNumero(t *testing.T) {
	continua := nuevaCarrera(t, 11, 10)
	for i := 0; i < 6; i++ {
		continua.AvanzarJornada()
	}
	esperada, _ := continua.AvanzarJornada() // jornada 7

	reanudada := nuevaCarrera(t, 11, 10)
	reanudada.Temporada.Resultados = append(reanudada.Temporada.Resultados,
		continua.Temporada.Resultados[:6]...)
	obtenida, err := reanudada.AvanzarJornada()
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(esperada, obtenida) {
		t.Error("la jornada 7 deberia ser igual al reanudar desde la jornada 6")
	}
}

func TestEquiposImparesDescansanSinFallar(t *testing.T) {
	c := nuevaCarrera(t, 3, 5)
	for !c.Terminada() {
		resultados, err := c.AvanzarJornada()
		if err != nil {
			t.Fatal(err)
		}
		delUsuario := 0
		for _, r := range resultados {
			if r.EsDelUsuario {
				delUsuario++
			}
		}
		if delUsuario > 1 {
			t.Errorf("jornada %d: %d partidos del usuario", c.Jornada(), delUsuario)
		}
	}
	if c.TotalJornadas() != 10 {
		t.Errorf("TotalJornadas = %d, se esperaban 10", c.TotalJornadas())
	}
}

func TestTablaMarcaAlUsuarioYNumeraPosiciones(t *testing.T) {
	c := nuevaCarrera(t, 2, 10)
	for i := 0; i < 5; i++ {
		c.AvanzarJornada()
	}
	tabla := c.Tabla()
	if len(tabla) != 10 {
		t.Fatalf("la tabla tiene %d filas", len(tabla))
	}
	usuarios := 0
	for i, f := range tabla {
		if f.Posicion != i+1 {
			t.Errorf("fila %d tiene posicion %d", i, f.Posicion)
		}
		if f.EsDelUsuario {
			usuarios++
			if f.Equipo != c.NombreEquipo() {
				t.Errorf("fila marcada como del usuario es %q, no %q", f.Equipo, c.NombreEquipo())
			}
		}
		if f.PJ != 5 {
			t.Errorf("%s: PJ = %d, se esperaban 5", f.Equipo, f.PJ)
		}
	}
	if usuarios != 1 {
		t.Errorf("%d filas del usuario, se esperaba 1", usuarios)
	}
}

func TestPlantillaOrdenadaYCopia(t *testing.T) {
	c := nuevaCarrera(t, 4, 10)
	plantilla := c.Plantilla()
	if len(plantilla) != len(c.Temporada.Equipos[c.Usuario].Plantilla) {
		t.Fatalf("plantilla de %d jugadores", len(plantilla))
	}
	for i := 1; i < len(plantilla); i++ {
		a, b := plantilla[i-1], plantilla[i]
		if a.Posicion > b.Posicion {
			t.Fatalf("posiciones fuera de orden: %v antes de %v", a.Posicion, b.Posicion)
		}
		if a.Posicion == b.Posicion && a.Valoracion() < b.Valoracion() {
			t.Errorf("%v: valoraciones fuera de orden (%d antes de %d)",
				a.Posicion, a.Valoracion(), b.Valoracion())
		}
	}
	if plantilla[0].Posicion != modelo.Portero || plantilla[len(plantilla)-1].Posicion != modelo.Delantero {
		t.Error("deberia empezar por porteros y terminar en delanteros")
	}

	// Modificar la copia no afecta a la carrera.
	original := c.Temporada.Equipos[c.Usuario].Plantilla[0]
	plantilla[0].Nombre = "Alterado"
	plantilla[len(plantilla)-1].Nombre = "Alterado"
	if c.Temporada.Equipos[c.Usuario].Plantilla[0] != original {
		t.Error("Plantilla devolvio una vista, no una copia")
	}
}

func TestValoracionEquipo(t *testing.T) {
	c := nuevaCarrera(t, 6, 10)
	if v := c.ValoracionEquipo(); v < 1 || v > modelo.AtributoMax {
		t.Errorf("valoracion fuera de rango: %d", v)
	}
}
