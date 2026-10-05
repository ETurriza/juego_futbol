package simulacion

import (
	"math"
	"math/rand"
	"sort"

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

// Resultado es el marcador de un partido y su detalle: alineaciones y sucesos.
type Resultado struct {
	GolesLocal     int
	GolesVisitante int
	Detalle        modelo.DetallePartido
}

// efectoRoja es cuánto cambia la esperanza de goles por cada partido completo
// que un equipo juega con diez: el que juega con diez marca esa fracción menos y
// su rival marca esa fracción más.
const efectoRoja = 0.18

// Marcador simula solo el resultado de un partido, sin alineaciones ni sucesos y,
// por tanto, sin el efecto de las tarjetas: cada equipo juega un 4-3-3 con sus
// mejores jugadores de cada posición. Es mucho más rápido que Simular y sirve para
// simular muchos partidos cuando no hace falta el detalle.
func Marcador(r *rand.Rand, local, visitante modelo.Equipo) (golesLocal, golesVisitante int) {
	ataqueL, defensaL := fuerzasAutomaticas(local, modelo.F433)
	ataqueV, defensaV := fuerzasAutomaticas(visitante, modelo.F433)

	golesLocal = poisson(r, esperanzaLocal(ataqueL, defensaV, 1))
	golesVisitante = poisson(r, esperanzaVisitante(ataqueV, defensaL, 1))
	return golesLocal, golesVisitante
}

func esperanzaLocal(ataque, defensaRival, factor float64) float64 {
	return golesBase * ventajaLocal * math.Pow(ataque/defensaRival, exponenteFuerza) * factor
}

func esperanzaVisitante(ataque, defensaRival, factor float64) float64 {
	return golesBase * math.Pow(ataque/defensaRival, exponenteFuerza) * factor
}

// Simular juega un partido entre local y visitante y devuelve el marcador y su
// detalle. Cada equipo elige solo la formación y el once que mejor le van. Es una
// función pura: no modifica los equipos y toda la aleatoriedad sale de r.
func Simular(r *rand.Rand, local, visitante modelo.Equipo) Resultado {
	return SimularConAlineaciones(r,
		local, AlineacionAutomatica(local, Criterios{}),
		visitante, AlineacionAutomatica(visitante, Criterios{}))
}

// SimularConAlineaciones juega un partido con las alineaciones dadas (que deben
// ser válidas para sus plantillas). Primero se deciden los cambios y las
// tarjetas; las rojas dejan a un equipo con diez y cambian la esperanza de goles
// del resto del partido, y después se reparten los goles entre quienes estaban en
// el campo.
func SimularConAlineaciones(r *rand.Rand, local modelo.Equipo, alL modelo.Alineacion,
	visitante modelo.Equipo, alV modelo.Alineacion) Resultado {

	tl := lineaDeTiempo(r, true, local, alL)
	tv := lineaDeTiempo(r, false, visitante, alV)

	ataqueL, defensaL := Fuerzas(local, alL)
	ataqueV, defensaV := Fuerzas(visitante, alV)
	conDiezL, conDiezV := tl.fraccionConDiez(), tv.fraccionConDiez()
	factorL := math.Max(1-efectoRoja*conDiezL+efectoRoja*conDiezV, 0.3)
	factorV := math.Max(1-efectoRoja*conDiezV+efectoRoja*conDiezL, 0.3)

	golesL := poisson(r, esperanzaLocal(ataqueL, defensaV, factorL))
	golesV := poisson(r, esperanzaVisitante(ataqueV, defensaL, factorV))

	var eventos []modelo.Evento
	eventos = append(eventos, tl.eventos...)
	eventos = append(eventos, tl.goles(r, golesL)...)
	eventos = append(eventos, tv.eventos...)
	eventos = append(eventos, tv.goles(r, golesV)...)
	// Orden por minuto; a igual minuto se conserva el orden de generación.
	sort.SliceStable(eventos, func(a, b int) bool { return eventos[a].Minuto < eventos[b].Minuto })

	return Resultado{
		GolesLocal:     golesL,
		GolesVisitante: golesV,
		Detalle: modelo.DetallePartido{
			FormacionLocal:     alL.Formacion,
			FormacionVisitante: alV.Formacion,
			TitularesLocal:     alL.Titulares,
			TitularesVisitante: alV.Titulares,
			Eventos:            eventos,
		},
	}
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
// crack, con el peso de su puesto. Se usa la valoración en ese puesto (la que ve el
// usuario) y no el valor del rol: así un 90 cuenta igual sea portero o delantero.
func (m *media) sumarCrack(j modelo.Jugador, puesto modelo.Posicion, puntosPorExceso float64) {
	m.cracks += puntosPorExceso * math.Max(float64(j.ValoracionEn(puesto))-umbralCrack, 0)
}

func (m media) valor() float64 {
	if m.peso == 0 {
		return fuerzaNeutra
	}
	aporte := topeCracks * (1 - math.Exp(-m.cracks/topeCracks))
	return math.Max(m.suma/m.peso+aporte, atributoMinFz)
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
