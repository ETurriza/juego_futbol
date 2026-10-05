package generador

import (
	"math/rand"
	"reflect"
	"testing"

	"github.com/ETurriza/juego_futbol/internal/modelo"
)

func nuevoRand(semilla int64) *rand.Rand {
	return rand.New(rand.NewSource(semilla))
}

func TestEquipoReproducible(t *testing.T) {
	a := Equipo(nuevoRand(42), "Deportivo Prueba", 1)
	b := Equipo(nuevoRand(42), "Deportivo Prueba", 1)
	if !reflect.DeepEqual(a, b) {
		t.Fatal("la misma semilla deberia dar el mismo equipo")
	}
	c := Equipo(nuevoRand(43), "Deportivo Prueba", 1)
	if reflect.DeepEqual(a, c) {
		t.Fatal("semillas distintas deberian dar equipos distintos")
	}
}

func TestEquipoComposicionYValidez(t *testing.T) {
	for semilla := int64(1); semilla <= 50; semilla++ {
		e := Equipo(nuevoRand(semilla), NombreEquipo(nuevoRand(semilla)), 100)
		if err := e.Validar(); err != nil {
			t.Fatalf("semilla %d: equipo invalido: %v", semilla, err)
		}
		if len(e.Plantilla) != TamanoPlantilla {
			t.Fatalf("semilla %d: plantilla de %d", semilla, len(e.Plantilla))
		}
		for _, c := range composicion {
			if got := e.Contar(c.posicion); got != c.cantidad {
				t.Errorf("semilla %d: %d de %v, se esperaban %d", semilla, got, c.posicion, c.cantidad)
			}
		}
	}
}

func TestEquipoIDsConsecutivasYNombresUnicos(t *testing.T) {
	e := Equipo(nuevoRand(7), "Club Prueba", 500)
	nombres := map[string]bool{}
	for i, j := range e.Plantilla {
		if j.ID != 500+i {
			t.Errorf("jugador %d tiene ID %d, se esperaba %d", i, j.ID, 500+i)
		}
		if nombres[j.Nombre] {
			t.Errorf("nombre repetido en el equipo: %s", j.Nombre)
		}
		nombres[j.Nombre] = true
	}
}

func TestAtributosEnRango(t *testing.T) {
	r := nuevoRand(1)
	for i := 0; i < 2000; i++ {
		for _, p := range modelo.Posiciones {
			if err := Atributos(r, p).Validar(); err != nil {
				t.Fatalf("%v: %v", p, err)
			}
		}
	}
}

func TestPerfilesPorPosicion(t *testing.T) {
	r := nuevoRand(2)
	const n = 1000
	var reflejosPortero, reflejosDefensa, tiroDelantero, tiroDefensa int
	for i := 0; i < n; i++ {
		reflejosPortero += Atributos(r, modelo.Portero).Reflejos
		reflejosDefensa += Atributos(r, modelo.Defensa).Reflejos
		tiroDelantero += Atributos(r, modelo.Delantero).Tiro
		tiroDefensa += Atributos(r, modelo.Defensa).Tiro
	}
	if reflejosPortero <= reflejosDefensa {
		t.Error("los porteros deberian tener mas reflejos que los defensas")
	}
	if tiroDelantero <= tiroDefensa {
		t.Error("los delanteros deberian tener mas tiro que los defensas")
	}
}

func TestJugadorValido(t *testing.T) {
	r := nuevoRand(3)
	for i := 0; i < 500; i++ {
		for _, p := range modelo.Posiciones {
			if err := Jugador(r, i, p).Validar(); err != nil {
				t.Fatalf("%v: %v", p, err)
			}
		}
	}
}

func TestListasSinDuplicados(t *testing.T) {
	for nombre, lista := range map[string][]string{
		"nombres": nombres, "apellidos": apellidos, "ciudades": ciudades, "prefijos": prefijos,
	} {
		vistos := map[string]bool{}
		for _, s := range lista {
			if s == "" {
				t.Errorf("%s: entrada vacia", nombre)
			}
			if vistos[s] {
				t.Errorf("%s: duplicado %q", nombre, s)
			}
			vistos[s] = true
		}
	}
}

func TestNombreEquipoReproducible(t *testing.T) {
	if NombreEquipo(nuevoRand(9)) != NombreEquipo(nuevoRand(9)) {
		t.Error("NombreEquipo deberia ser reproducible con la misma semilla")
	}
	if NombreJugador(nuevoRand(9)) != NombreJugador(nuevoRand(9)) {
		t.Error("NombreJugador deberia ser reproducible con la misma semilla")
	}
}
