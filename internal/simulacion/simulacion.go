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
	// Cracks: umbralCrack es la valoración a partir de la cual un jugador suma al
	// equipo más que su parte proporcional; cada posición suma cierta cantidad de
	// puntos de fuerza por punto de exceso (más en las líneas donde un jugador pesa
	// menos, para que un crack de cualquier posición se note parecido); topeCracks
	// es el máximo (de saturación) que pueden sumar los cracks de una misma línea.
	umbralCrack     = 85.0
	crackDelantero  = 1.2
	crackMediocampo = 1.6
	crackDefensa    = 3.0
	crackPortero    = 1.9
	topeCracks      = 16.0

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

// Marcador simula solo el resultado de un partido, sin alineaciones ni sucesos.
// Da los mismos goles que Simular con la misma aleatoriedad, pero mucho más
// rápido: sirve para simular muchos partidos cuando no hace falta el detalle.
func Marcador(r *rand.Rand, local, visitante modelo.Equipo) (golesLocal, golesVisitante int) {
	ataqueL, defensaL := fuerzas(local)
	ataqueV, defensaV := fuerzas(visitante)

	esperanzaL := golesBase * ventajaLocal * math.Pow(ataqueL/defensaV, exponenteFuerza)
	esperanzaV := golesBase * math.Pow(ataqueV/defensaL, exponenteFuerza)

	golesLocal = poisson(r, esperanzaL)
	golesVisitante = poisson(r, esperanzaV)
	return golesLocal, golesVisitante
}

// Simular juega un partido entre local y visitante y devuelve el marcador y su
// detalle. Es una función pura: no modifica los equipos y toda la aleatoriedad
// sale de r.
func Simular(r *rand.Rand, local, visitante modelo.Equipo) Resultado {
	golesL, golesV := Marcador(r, local, visitante)
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
		v := (float64(j.Atributos.Tiro) + float64(j.Atributos.Regate) +
			float64(j.Atributos.Ritmo) + float64(j.Atributos.Pase)) / 4
		a.sumar(2, v)
		a.sumarCrack(j, crackDelantero)
	}
	for _, j := range medios {
		v := (float64(j.Atributos.Pase) + float64(j.Atributos.Regate) + float64(j.Atributos.Tiro)) / 3
		a.sumar(1, v)
		a.sumarCrack(j, crackMediocampo)
		d.sumar(0.5, float64(j.Atributos.Defensa))
	}
	for _, j := range defensas {
		v := (float64(j.Atributos.Defensa) + float64(j.Atributos.Fisico)) / 2
		d.sumar(1, v)
		d.sumarCrack(j, crackDefensa)
	}
	for _, j := range porteros {
		v := float64(j.Atributos.Reflejos)
		d.sumar(3, v)
		d.sumarCrack(j, crackPortero)
	}
	return a.valor(), d.valor()
}

// media acumula la fuerza de una línea: la media ponderada de sus jugadores más
// el aporte de sus cracks. Un crack (un jugador cuyo valor de rol supera
// umbralCrack) suma puntosPorCrack por cada punto de exceso; el aporte total de
// la línea se satura en topeCracks, de modo que un crack se nota pero muchos
// cracks no deciden solos el partido.
type media struct{ suma, peso, cracks float64 }

func (m *media) sumar(peso, valor float64) {
	m.suma += peso * valor
	m.peso += peso
}

// sumarCrack acumula el exceso de la valoración del jugador sobre el umbral de
// crack, con el peso de su posición. Se usa la valoración (la que ve el usuario) y
// no el valor del rol: así un 90 cuenta igual sea portero o delantero.
func (m *media) sumarCrack(j modelo.Jugador, puntosPorExceso float64) {
	m.cracks += puntosPorExceso * math.Max(float64(j.Valoracion())-umbralCrack, 0)
}

func (m media) valor() float64 {
	if m.peso == 0 {
		return fuerzaNeutra
	}
	aporte := topeCracks * (1 - math.Exp(-m.cracks/topeCracks))
	return math.Max(m.suma/m.peso+aporte, atributoMinFz)
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
