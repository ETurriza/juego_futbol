package simulacion

import (
	"math"
	"math/rand"

	"github.com/ETurriza/juego_futbol/internal/modelo"
)

// Constantes de calibración del modelo. Ajustar aquí cambia el ritmo de goles
// y cuánto pesa la diferencia de calidad entre los equipos.
const (
	// golesBase es la esperanza de goles de un equipo contra un rival de su
	// misma calidad y sin ventaja de local.
	golesBase = 1.30
	// ventajaLocal multiplica la esperanza de goles del equipo local.
	ventajaLocal = 1.10
	// exponenteFuerza controla cuánto se amplifica la razón ataque/defensa.
	exponenteFuerza = 2.0
	// fuerzaNeutra es la fuerza que se asume cuando falta una línea completa.
	fuerzaNeutra = 50.0

	// límites de seguridad.
	golesMax      = 12
	esperanzaMin  = 0.05
	esperanzaMax  = 6.0
	atributoMinFz = 1.0
)

// Formación fija 4-3-3 para elegir los titulares de cada línea.
const (
	titularesDefensas = 4
	titularesMedios   = 3
	titularesDelant   = 3
)

// Resultado es el marcador de un partido y su detalle: alineaciones y sucesos.
type Resultado struct {
	GolesLocal     int
	GolesVisitante int
	Detalle        modelo.DetallePartido
}

// Simular juega un partido entre local y visitante y devuelve el marcador.
// Es una función pura: no modifica los equipos y toda la aleatoriedad sale de r.
func Simular(r *rand.Rand, local, visitante modelo.Equipo) Resultado {
	ataqueL, defensaL := fuerzas(local)
	ataqueV, defensaV := fuerzas(visitante)

	esperanzaL := golesBase * ventajaLocal * math.Pow(ataqueL/defensaV, exponenteFuerza)
	esperanzaV := golesBase * math.Pow(ataqueV/defensaL, exponenteFuerza)

	golesL := poisson(r, esperanzaL)
	golesV := poisson(r, esperanzaV)
	return Resultado{
		GolesLocal:     golesL,
		GolesVisitante: golesV,
		Detalle:        generarDetalle(r, local, visitante, golesL, golesV),
	}
}

// fuerzas devuelve la fuerza de ataque y de defensa de un equipo (en la escala
// de los atributos) a partir de sus titulares.
func fuerzas(e modelo.Equipo) (ataque, defensa float64) {
	porteros := mejores(e, modelo.Portero, 1)
	defensas := mejores(e, modelo.Defensa, titularesDefensas)
	medios := mejores(e, modelo.Mediocampista, titularesMedios)
	delanteros := mejores(e, modelo.Delantero, titularesDelant)

	var a, d media
	for _, j := range delanteros {
		a.sumar(2, (float64(j.Atributos.Tiro)+float64(j.Atributos.Regate)+
			float64(j.Atributos.Ritmo)+float64(j.Atributos.Pase))/4)
	}
	for _, j := range medios {
		a.sumar(1, (float64(j.Atributos.Pase)+float64(j.Atributos.Regate)+
			float64(j.Atributos.Tiro))/3)
		d.sumar(0.5, float64(j.Atributos.Defensa))
	}
	for _, j := range defensas {
		d.sumar(1, (float64(j.Atributos.Defensa)+float64(j.Atributos.Fisico))/2)
	}
	for _, j := range porteros {
		d.sumar(3, float64(j.Atributos.Reflejos))
	}
	return a.valor(), d.valor()
}

// media acumula un promedio ponderado.
type media struct{ suma, peso float64 }

func (m *media) sumar(peso, valor float64) {
	m.suma += peso * valor
	m.peso += peso
}

func (m media) valor() float64 {
	if m.peso == 0 {
		return fuerzaNeutra
	}
	return math.Max(m.suma/m.peso, atributoMinFz)
}

// mejores devuelve hasta n jugadores de la posición dada, los de mayor
// valoración primero. No modifica la plantilla.
func mejores(e modelo.Equipo, p modelo.Posicion, n int) []modelo.Jugador {
	var candidatos []modelo.Jugador
	for _, j := range e.Plantilla {
		if j.Posicion == p {
			candidatos = append(candidatos, j)
		}
	}
	// selección simple por inserción: las plantillas son pequeñas y el orden
	// debe ser estable y determinista (mayor valoración, luego menor ID).
	for i := 1; i < len(candidatos); i++ {
		for k := i; k > 0 && mejorQue(candidatos[k], candidatos[k-1]); k-- {
			candidatos[k], candidatos[k-1] = candidatos[k-1], candidatos[k]
		}
	}
	if len(candidatos) > n {
		candidatos = candidatos[:n]
	}
	return candidatos
}

func mejorQue(a, b modelo.Jugador) bool {
	if va, vb := a.Valoracion(), b.Valoracion(); va != vb {
		return va > vb
	}
	return a.ID < b.ID
}

// poisson muestrea una variable de Poisson con la esperanza dada (algoritmo de
// Knuth), acotada a golesMax.
func poisson(r *rand.Rand, esperanza float64) int {
	esperanza = math.Min(math.Max(esperanza, esperanzaMin), esperanzaMax)
	limite := math.Exp(-esperanza)
	k, p := 0, 1.0
	for {
		p *= r.Float64()
		if p <= limite || k >= golesMax {
			return k
		}
		k++
	}
}
