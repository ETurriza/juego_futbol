package generador

import (
	"math"
	"math/rand"

	"github.com/ETurriza/juego_futbol/internal/modelo"
	"github.com/ETurriza/juego_futbol/internal/progresion"
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
	return atributosDeCalidad(r, p, int(math.Round(r.NormFloat64()*8+62)))
}

// atributosDeCalidad genera los atributos de la posición a partir de una
// calidad base ya elegida.
func atributosDeCalidad(r *rand.Rand, p modelo.Posicion, calidad int) modelo.Atributos {
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

// Jugador genera un jugador de una edad entre 17 y 36 años (hasta 39 los porteros). Se crea como un
// juvenil de 16 y se le hace envejecer con la propia progresión, de modo que la
// calidad de cada edad es la que dejaría la carrera de un jugador de la liga: la
// liga inicial se parece a la que habrá tras muchas temporadas.
func Jugador(r *rand.Rand, id int, p modelo.Posicion) modelo.Jugador {
	// Los porteros juegan DesfasePortero años más, así que su rango de edades se
	// alarga: la liga inicial tiene la misma mezcla de edades que la que resulta
	// tras muchas temporadas.
	rango := 20 // 17 a 36
	if p == modelo.Portero {
		rango += progresion.DesfasePortero // 17 a 39
	}
	edad := 17 + r.Intn(rango)

	// Solo existen los jugadores que han llegado a esa edad sin retirarse: si en
	// algún año el jugador se habría retirado (por edad o por nivel), se parte de
	// uno nuevo. Así la liga inicial no tiene veteranos que en la carrera ya
	// habrían dejado el fútbol.
	var j modelo.Jugador
	for intento := 0; intento < maxIntentosSuperviviente; intento++ {
		j = nuevoJuvenil(r, id, p, EdadJuvenilMin)
		vivo := true
		for j.Edad < edad {
			if progresion.SeRetira(r, j) {
				vivo = false
				break
			}
			j = progresion.Envejecer(r, j)
		}
		if vivo {
			break
		}
	}
	return j
}

// maxIntentosSuperviviente acota los reintentos al crear un jugador veterano;
// con las probabilidades de retiro, bastan unas pocas docenas.
const maxIntentosSuperviviente = 2000

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

// Calidad base de los juveniles: bastante por debajo de la de un jugador
// hecho, porque la progresión por edad los hace crecer varios años seguidos.
const (
	juvenilCalidadMedia = 40.0
	juvenilCalidadDesv  = 5.0
	// Edad de un juvenil: de EdadJuvenilMin a EdadJuvenilMax.
	EdadJuvenilMin = 16
	EdadJuvenilMax = 19
)

// Juvenil genera un jugador de cantera, de EdadJuvenilMin a EdadJuvenilMax
// años, con la ID y la posición dadas. Como cualquier jugador, nace con 16 años y
// crece con la progresión: uno de 19 ya ha pasado tres años de crecimiento.
func Juvenil(r *rand.Rand, id int, p modelo.Posicion) modelo.Jugador {
	edad := EdadJuvenilMin + r.Intn(EdadJuvenilMax-EdadJuvenilMin+1)
	j := nuevoJuvenil(r, id, p, EdadJuvenilMin)
	for j.Edad < edad {
		j = progresion.Envejecer(r, j)
	}
	return j
}

// nuevoJuvenil crea un jugador de la edad dada con la calidad de un juvenil.
func nuevoJuvenil(r *rand.Rand, id int, p modelo.Posicion, edad int) modelo.Jugador {
	calidad := int(math.Round(r.NormFloat64()*juvenilCalidadDesv + juvenilCalidadMedia))
	return modelo.Jugador{
		ID:        id,
		Nombre:    NombreJugador(r),
		Edad:      edad,
		Posicion:  p,
		Atributos: atributosDeCalidad(r, p, calidad),
		Talento:   talentoAleatorio(r),
	}
}

// MinimoPorPosicion es cuántos jugadores de la posición tiene una plantilla
// completa; la renovación de plantillas repone hasta esa cantidad.
func MinimoPorPosicion(p modelo.Posicion) int {
	for _, c := range composicion {
		if c.posicion == p {
			return c.cantidad
		}
	}
	return 0
}

// talentoDesv es la dispersión (logarítmica) del talento: con ella, de cada cien
// jugadores unos tres o cuatro llegan a ser cracks.
const talentoDesv = 0.18

// talentoAleatorio sortea el talento de un jugador con una distribución
// log-normal centrada en TalentoNeutro: la mayoría crece de forma normal y unos
// pocos mucho más (y unos pocos menos).
func talentoAleatorio(r *rand.Rand) int {
	t := int(math.Round(modelo.TalentoNeutro * math.Exp(r.NormFloat64()*talentoDesv)))
	return min(max(t, modelo.TalentoMin), modelo.TalentoMax)
}
