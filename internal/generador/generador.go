package generador

import (
	"math"
	"math/rand"

	"github.com/ETurriza/juego_futbol/internal/modelo"
)

// composicion es la distribución de una plantilla de 22 jugadores.
var composicion = []struct {
	posicion modelo.Posicion
	cantidad int
}{
	{modelo.Portero, 2},
	{modelo.Defensa, 7},
	{modelo.Mediocampista, 7},
	{modelo.Delantero, 6},
}

// TamanoPlantilla es el número de jugadores de un equipo generado.
const TamanoPlantilla = 22

// perfil desplaza cada atributo respecto a la calidad base del jugador según
// su posición.
var perfil = map[modelo.Posicion]modelo.Atributos{
	modelo.Portero:       {Ritmo: -15, Tiro: -35, Pase: -10, Regate: -30, Defensa: -10, Fisico: 0, Reflejos: 25},
	modelo.Defensa:       {Ritmo: 0, Tiro: -20, Pase: -5, Regate: -15, Defensa: 20, Fisico: 10, Reflejos: -50},
	modelo.Mediocampista: {Ritmo: 0, Tiro: 0, Pase: 15, Regate: 10, Defensa: -5, Fisico: 0, Reflejos: -50},
	modelo.Delantero:     {Ritmo: 10, Tiro: 20, Pase: 0, Regate: 10, Defensa: -25, Fisico: 0, Reflejos: -50},
}

func limitar(v int) int {
	return min(max(v, modelo.AtributoMin), modelo.AtributoMax)
}

func elegir(r *rand.Rand, lista []string) string {
	return lista[r.Intn(len(lista))]
}

// NombreJugador devuelve un nombre completo inventado.
func NombreJugador(r *rand.Rand) string {
	return elegir(r, nombres) + " " + elegir(r, apellidos)
}

// NombreEquipo devuelve un nombre de club inventado, como "Deportivo Valmora".
func NombreEquipo(r *rand.Rand) string {
	return elegir(r, prefijos) + " " + elegir(r, ciudades)
}

// Atributos genera atributos para la posición dada: una calidad base común al
// jugador, el perfil de la posición y una pequeña variación por atributo.
func Atributos(r *rand.Rand, p modelo.Posicion) modelo.Atributos {
	calidad := int(math.Round(r.NormFloat64()*8 + 62))
	ajuste := perfil[p]
	attr := func(desplazamiento int) int {
		return limitar(calidad + desplazamiento + r.Intn(13) - 6)
	}
	return modelo.Atributos{
		Ritmo:    attr(ajuste.Ritmo),
		Tiro:     attr(ajuste.Tiro),
		Pase:     attr(ajuste.Pase),
		Regate:   attr(ajuste.Regate),
		Defensa:  attr(ajuste.Defensa),
		Fisico:   attr(ajuste.Fisico),
		Reflejos: attr(ajuste.Reflejos),
	}
}

// Jugador genera un jugador con la ID y la posición dadas.
func Jugador(r *rand.Rand, id int, p modelo.Posicion) modelo.Jugador {
	return modelo.Jugador{
		ID:        id,
		Nombre:    NombreJugador(r),
		Edad:      17 + r.Intn(20), // 17 a 36
		Posicion:  p,
		Atributos: Atributos(r, p),
	}
}

// Equipo genera una plantilla de TamanoPlantilla jugadores con IDs
// consecutivas desde idInicial. Los nombres de jugador no se repiten dentro del
// equipo.
func Equipo(r *rand.Rand, nombre string, idInicial int) modelo.Equipo {
	e := modelo.Equipo{Nombre: nombre, Plantilla: make([]modelo.Jugador, 0, TamanoPlantilla)}
	usados := make(map[string]bool, TamanoPlantilla)
	id := idInicial
	for _, c := range composicion {
		for i := 0; i < c.cantidad; i++ {
			j := Jugador(r, id, c.posicion)
			// Con 1600 combinaciones y 22 jugadores, reintentar basta; el tope
			// evita un bucle infinito si se achican las listas.
			for intento := 0; usados[j.Nombre] && intento < 100; intento++ {
				j.Nombre = NombreJugador(r)
			}
			usados[j.Nombre] = true
			e.Plantilla = append(e.Plantilla, j)
			id++
		}
	}
	return e
}
