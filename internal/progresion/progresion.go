package progresion

import (
	"math"
	"math/rand"

	"github.com/ETurriza/juego_futbol/internal/modelo"
)

// Constantes de calibración. Ajustar aquí cambia cómo envejece la liga.
const (
	// DesfasePortero son los años que se corre la curva de un portero: a los 33
	// rinde como un jugador de campo de 30.
	DesfasePortero = 3

	// ruido anual de cada atributo y forma común del año (ambos se acumulan de
	// un año a otro, así que unos jugadores aguantan mejor que otros).
	ruidoAtributo = 1.5
	ruidoForma    = 1.5

	// Rendimientos decrecientes: cuanto más cerca del tope está un atributo, menos
	// crece. Con valor <= techoAtributo-amplitudTecho crece al ritmo normal; al
	// llegar a techoAtributo, solo factorMinimo de lo normal. Solo afecta al
	// crecimiento: la caída por edad no se frena.
	techoAtributo = 92.0
	amplitudTecho = 25.0
	factorMinimo  = 0.10

	// Retiro por nivel: desde edadRetiroPorNivel (edad efectiva) un jugador cuya
	// valoración está por debajo de nivelMinimoTitular suma pesoNivelEnRetiro de
	// probabilidad de retiro por cada punto que le falta. Un veterano que ya no
	// da el nivel se retira antes que uno que sigue rindiendo.
	edadRetiroPorNivel = 33
	nivelMinimoTitular = 62
	pesoNivelEnRetiro  = 0.08

	// "Último prime": entre estas edades, cada año hay una probabilidad de un
	// gran año que suma un bono a los atributos técnicos.
	granAnoDesde = 31
	granAnoHasta = 34
	probGranAno  = 0.20
	bonoGranAno  = 3.0
)

// tramo es el cambio medio anual de atributos técnicos y físicos desde una
// edad (efectiva) en adelante.
type tramo struct {
	desde   int
	tecnico float64
	fisico  float64
}

// tramos va de menor a mayor edad; se aplica el último con desde <= edad.
var tramos = []tramo{
	{0, +4.0, +4.0},
	{22, +2.5, +2.0},
	{25, +1.0, 0.0},
	{28, +0.5, -1.0}, // meseta
	{32, 0.0, -2.5},
	{34, -3.0, -5.0},
	{35, -4.5, -6.0},
	{36, -6.0, -7.0},
}

// retiros es la probabilidad de retirarse a cada edad (efectiva); por debajo de
// la primera entrada es 0 y desde la última es 1.
var retiros = []struct {
	edad int
	prob float64
}{
	{34, 0.05},
	{35, 0.15},
	{36, 0.30},
	{37, 0.50},
	{38, 0.75},
	{39, 1.00},
}

// edadEfectiva corre la curva de los porteros unos años.
func edadEfectiva(j modelo.Jugador) int {
	if j.Posicion == modelo.Portero {
		return j.Edad - DesfasePortero
	}
	return j.Edad
}

func buscarTramo(edad int) tramo {
	t := tramos[0]
	for _, c := range tramos {
		if edad >= c.desde {
			t = c
		}
	}
	return t
}

// CambioMedio devuelve el cambio medio anual esperado de los atributos
// técnicos y físicos de un jugador de esa edad y posición, sin contar el ruido
// ni el bono del último prime.
func CambioMedio(edad int, p modelo.Posicion) (tecnico, fisico float64) {
	t := buscarTramo(edadEfectiva(modelo.Jugador{Edad: edad, Posicion: p}))
	return t.tecnico, t.fisico
}

// ProbabilidadRetiro es la probabilidad de que un jugador de esa edad y
// posición se retire al final de la temporada.
func ProbabilidadRetiro(edad int, p modelo.Posicion) float64 {
	e := edadEfectiva(modelo.Jugador{Edad: edad, Posicion: p})
	prob := 0.0
	for _, r := range retiros {
		if e >= r.edad {
			prob = r.prob
		}
	}
	return prob
}

// ProbabilidadRetiroDe es la probabilidad de que este jugador se retire al final
// de la temporada: la que corresponde a su edad y posición, más un adelanto si es
// un veterano que ya no da el nivel (ver edadRetiroPorNivel).
func ProbabilidadRetiroDe(j modelo.Jugador) float64 {
	prob := ProbabilidadRetiro(j.Edad, j.Posicion)
	if edadEfectiva(j) >= edadRetiroPorNivel {
		if falta := nivelMinimoTitular - j.Valoracion(); falta > 0 {
			prob += float64(falta) * pesoNivelEnRetiro
		}
	}
	return math.Min(prob, 1)
}

// SeRetira decide si el jugador se retira al final de la temporada, según su
// edad y su nivel.
func SeRetira(r *rand.Rand, j modelo.Jugador) bool {
	// Siempre se consume un número aleatorio, para que la secuencia no dependa
	// de la edad ni del nivel.
	x := r.Float64()
	return x < ProbabilidadRetiroDe(j)
}

// Envejecer devuelve al jugador un año mayor, con sus atributos evolucionados.
// No modifica el original.
func Envejecer(r *rand.Rand, j modelo.Jugador) modelo.Jugador {
	t := buscarTramo(edadEfectiva(j))
	tecnico, fisico := t.tecnico, t.fisico

	if e := edadEfectiva(j); e >= granAnoDesde && e <= granAnoHasta && r.Float64() < probGranAno {
		tecnico += bonoGranAno
	}
	forma := r.NormFloat64() * ruidoForma

	mover := func(valor int, cambio float64) int {
		if cambio > 0 {
			cambio *= factorCrecimiento(valor)
		}
		v := float64(valor) + cambio + forma + r.NormFloat64()*ruidoAtributo
		return min(max(int(math.Round(v)), modelo.AtributoMin), modelo.AtributoMax)
	}

	a := j.Atributos
	j.Edad++
	j.Atributos = modelo.Atributos{
		Ritmo:   mover(a.Ritmo, fisico),
		Fisico:  mover(a.Fisico, fisico),
		Tiro:    mover(a.Tiro, tecnico),
		Pase:    mover(a.Pase, tecnico),
		Regate:  mover(a.Regate, tecnico),
		Defensa: mover(a.Defensa, tecnico),
		// Los reflejos solo evolucionan en los porteros; en el resto son
		// irrelevantes y se dejan como están.
		Reflejos: a.Reflejos,
	}
	if j.Posicion == modelo.Portero {
		j.Atributos.Reflejos = mover(a.Reflejos, tecnico)
	}
	return j
}

// factorCrecimiento es la fracción del crecimiento normal que conserva un
// atributo con el valor dado: 1 lejos del tope y factorMinimo cerca de él.
func factorCrecimiento(valor int) float64 {
	return math.Min(math.Max((techoAtributo-float64(valor))/amplitudTecho, factorMinimo), 1)
}
