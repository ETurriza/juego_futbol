package simulacion

import (
	"math/rand"
	"reflect"
	"testing"

	"github.com/ETurriza/juego_futbol/internal/modelo"
)

func nuevoRand(semilla int64) *rand.Rand {
	return rand.New(rand.NewSource(semilla))
}

// equipoUniforme arma una plantilla de 22 con todos los atributos iguales a v.
func equipoUniforme(nombre string, v int) modelo.Equipo {
	e := modelo.Equipo{Nombre: nombre}
	id := 1
	for _, c := range []struct {
		p modelo.Posicion
		n int
	}{
		{modelo.Portero, 2}, {modelo.Defensa, 7},
		{modelo.Mediocampista, 7}, {modelo.Delantero, 6},
	} {
		for i := 0; i < c.n; i++ {
			e.Plantilla = append(e.Plantilla, modelo.Jugador{
				ID: id, Nombre: "J", Edad: 25, Posicion: c.p,
				Atributos: modelo.Atributos{
					Ritmo: v, Tiro: v, Pase: v, Regate: v, Defensa: v, Fisico: v, Reflejos: v,
				},
			})
			id++
		}
	}
	return e
}

func TestSimularReproducible(t *testing.T) {
	a, b := equipoUniforme("A", 60), equipoUniforme("B", 55)
	for semilla := int64(1); semilla <= 20; semilla++ {
		r1 := Simular(nuevoRand(semilla), a, b)
		r2 := Simular(nuevoRand(semilla), a, b)
		if !reflect.DeepEqual(r1, r2) {
			t.Fatalf("semilla %d: %v != %v", semilla, r1, r2)
		}
	}
}

func TestSimularNoModificaEquipos(t *testing.T) {
	a, b := equipoUniforme("A", 60), equipoUniforme("B", 55)
	copiaA, copiaB := equipoUniforme("A", 60), equipoUniforme("B", 55)
	r := nuevoRand(1)
	for i := 0; i < 50; i++ {
		Simular(r, a, b)
	}
	if !reflect.DeepEqual(a, copiaA) || !reflect.DeepEqual(b, copiaB) {
		t.Error("Simular modifico los equipos")
	}
}

type estadistica struct {
	victoriasLocal, empates, victoriasVisitante int
	goles                                       int
}

func jugar(n int, local, visitante modelo.Equipo) estadistica {
	r := nuevoRand(12345)
	var s estadistica
	for i := 0; i < n; i++ {
		res := Simular(r, local, visitante)
		s.goles += res.GolesLocal + res.GolesVisitante
		switch {
		case res.GolesLocal > res.GolesVisitante:
			s.victoriasLocal++
		case res.GolesLocal < res.GolesVisitante:
			s.victoriasVisitante++
		default:
			s.empates++
		}
	}
	return s
}

func TestEquipoMejorGanaMas(t *testing.T) {
	const n = 2000
	fuerte, debil := equipoUniforme("Fuerte", 75), equipoUniforme("Debil", 50)

	// El mejor gana como local y como visitante.
	s := jugar(n, fuerte, debil)
	if s.victoriasLocal*100 < n*60 {
		t.Errorf("fuerte de local gano %d/%d, se esperaba mas del 60%%", s.victoriasLocal, n)
	}
	s = jugar(n, debil, fuerte)
	if s.victoriasVisitante*100 < n*60 {
		t.Errorf("fuerte de visitante gano %d/%d, se esperaba mas del 60%%", s.victoriasVisitante, n)
	}
}

func TestEquiposIgualesConVentajaLocal(t *testing.T) {
	const n = 5000
	a, b := equipoUniforme("A", 60), equipoUniforme("B", 60)
	s := jugar(n, a, b)
	if s.victoriasLocal <= s.victoriasVisitante {
		t.Errorf("el local deberia ganar mas: local %d, visitante %d",
			s.victoriasLocal, s.victoriasVisitante)
	}
	// Pero no de forma abrumadora: la ventaja es pequeña.
	if s.victoriasLocal*100 > n*55 {
		t.Errorf("ventaja de local excesiva: %d/%d", s.victoriasLocal, n)
	}
	if s.empates == 0 {
		t.Error("deberia haber empates entre equipos iguales")
	}
}

func TestPromedioDeGolesRealista(t *testing.T) {
	const n = 5000
	for _, par := range [][2]int{{60, 60}, {70, 55}, {50, 50}} {
		s := jugar(n, equipoUniforme("A", par[0]), equipoUniforme("B", par[1]))
		prom := float64(s.goles) / n
		if prom < 2.0 || prom > 3.6 {
			t.Errorf("%v: promedio de %.2f goles por partido fuera de [2.0, 3.6]", par, prom)
		}
	}
}

func TestGolesAcotados(t *testing.T) {
	// El extremo mas desparejo posible no debe pasar del tope ni ser negativo.
	fuerte, debil := equipoUniforme("F", modelo.AtributoMax), equipoUniforme("D", modelo.AtributoMin)
	r := nuevoRand(99)
	for i := 0; i < 2000; i++ {
		res := Simular(r, fuerte, debil)
		if res.GolesLocal < 0 || res.GolesVisitante < 0 ||
			res.GolesLocal > golesMax || res.GolesVisitante > golesMax {
			t.Fatalf("resultado fuera de rango: %+v", res)
		}
	}
}

func TestEquiposVaciosNoFallan(t *testing.T) {
	r := nuevoRand(5)
	vacio := modelo.Equipo{Nombre: "Vacio"}
	for i := 0; i < 200; i++ {
		res := Simular(r, vacio, vacio)
		if res.GolesLocal < 0 || res.GolesVisitante < 0 {
			t.Fatalf("goles negativos: %+v", res)
		}
	}
	// Un equipo vacío contra uno normal tampoco debe fallar.
	Simular(r, vacio, equipoUniforme("A", 60))
	Simular(r, equipoUniforme("A", 60), vacio)
}

func TestMejoresEsDeterministaYNoMuta(t *testing.T) {
	e := equipoUniforme("A", 60)
	e.Plantilla[2].Atributos.Defensa = 90 // un defensa claramente mejor
	antes := append([]modelo.Jugador(nil), e.Plantilla...)

	m1 := mejores(e, modelo.Defensa, titularesDefensas)
	m2 := mejores(e, modelo.Defensa, titularesDefensas)
	if !reflect.DeepEqual(m1, m2) {
		t.Error("mejores no es determinista")
	}
	if len(m1) != titularesDefensas {
		t.Errorf("se esperaban %d titulares, hay %d", titularesDefensas, len(m1))
	}
	if m1[0].ID != e.Plantilla[2].ID {
		t.Errorf("el mejor defensa deberia ser ID %d, es %d", e.Plantilla[2].ID, m1[0].ID)
	}
	if !reflect.DeepEqual(antes, e.Plantilla) {
		t.Error("mejores modifico la plantilla")
	}
}
