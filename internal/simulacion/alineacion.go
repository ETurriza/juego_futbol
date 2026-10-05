package simulacion

import (
	"math"
	"sort"

	"github.com/ETurriza/juego_futbol/internal/modelo"
)

// perfil es cuánto cambia una formación el ataque y la defensa de un equipo
// respecto del 4-3-3, la referencia. Las formaciones más ofensivas marcan más y
// reciben más, y las defensivas al revés; los perfiles son de efecto neto
// parecido (ninguna formación es mejor que las demás por sí sola): lo que decide
// cuál conviene es la plantilla, es decir, en qué líneas es más fuerte.
var perfil = map[modelo.Formacion]struct{ ataque, defensa float64 }{
	modelo.F433: {1.000, 1.000},
	modelo.F442: {0.994, 1.014},
	modelo.F352: {1.035, 0.995},
	modelo.F451: {0.990, 1.050},
	modelo.F532: {0.951, 1.051},
	modelo.F343: {1.056, 0.956},
}

// Criterios son las preferencias con que se elige una alineación automática.
type Criterios struct {
	// Formacion fija la formación; si es nil se elige la que mejor rinde con la
	// plantilla.
	Formacion *modelo.Formacion
	// NoDisponibles son los jugadores que no pueden jugar (lesionados o
	// sancionados), por su ID.
	NoDisponibles map[int]bool
}

// titular es un jugador alineado en un puesto.
type titular struct {
	j      *modelo.Jugador
	puesto modelo.Posicion
}

// plantel son los jugadores disponibles de un equipo, por posición natural y
// ordenados de mayor a menor valoración; se calcula una sola vez por equipo.
type plantel struct {
	porPosicion [4][]*modelo.Jugador // indexado por modelo.Posicion
	todos       []*modelo.Jugador    // todos, de mayor a menor valoración
}

func nuevoPlantel(e modelo.Equipo, noDisp map[int]bool) plantel {
	type conValor struct {
		j *modelo.Jugador
		v int
	}
	cv := make([]conValor, 0, len(e.Plantilla))
	for i := range e.Plantilla {
		j := &e.Plantilla[i]
		if !noDisp[j.ID] {
			cv = append(cv, conValor{j, j.Valoracion()})
		}
	}
	sort.SliceStable(cv, func(a, b int) bool {
		if cv[a].v != cv[b].v {
			return cv[a].v > cv[b].v
		}
		return cv[a].j.ID < cv[b].j.ID
	})
	var p plantel
	p.todos = make([]*modelo.Jugador, len(cv))
	for i, x := range cv {
		p.todos[i] = x.j
		p.porPosicion[x.j.Posicion] = append(p.porPosicion[x.j.Posicion], x.j)
	}
	return p
}

// AlineacionAutomatica elige el once y el banquillo de un equipo: para cada
// formación toma, en cada línea, a los mejores jugadores disponibles de esa
// posición (y, si faltan, a los mejores de las líneas contiguas) y se queda con la
// formación que mejor rinde. El resultado es determinista.
func AlineacionAutomatica(e modelo.Equipo, c Criterios) modelo.Alineacion {
	p := nuevoPlantel(e, c.NoDisponibles)
	if c.Formacion != nil {
		return p.alineacion(*c.Formacion)
	}
	// Se compara cada formación contra un rival tan fuerte como el propio equipo
	// (en su 4-3-3, la formación de referencia): así ninguna formación gana por
	// defecto y lo que decide es en qué líneas es fuerte la plantilla.
	a0, d0 := fuerzasDe(p.seleccion(modelo.F433), modelo.F433)
	rival := (a0 + d0) / 2
	mejor, mejorPuntaje := modelo.F433, math.Inf(-1)
	for _, f := range modelo.Formaciones {
		if pt := puntaje(p.seleccion(f), f, rival); pt > mejorPuntaje {
			mejor, mejorPuntaje = f, pt
		}
	}
	return p.alineacion(mejor)
}

// alineacionParaFormacion arma el mejor once para la formación dada.
func alineacionParaFormacion(e modelo.Equipo, f modelo.Formacion, noDisp map[int]bool) modelo.Alineacion {
	return nuevoPlantel(e, noDisp).alineacion(f)
}

// puntaje estima lo bien que rinde una alineación: la diferencia de goles
// esperada contra un rival cuyo ataque y defensa valen rival.
func puntaje(tit [modelo.TitularesPorEquipo]*titular, f modelo.Formacion, rival float64) float64 {
	a, d := fuerzasDe(tit, f)
	return math.Pow(a/rival, exponenteFuerza) - math.Pow(rival/d, exponenteFuerza)
}

// seleccion elige el once de la formación: en cada línea, a los mejores de su
// posición natural y, si faltan, a quien mejor rinda fuera de posición. Los
// puestos que no se pueden cubrir quedan en nil.
func (p plantel) seleccion(f modelo.Formacion) [modelo.TitularesPorEquipo]*titular {
	var tit [modelo.TitularesPorEquipo]*titular
	puestos := f.Puestos()
	var sig [4]int // siguiente candidato natural de cada posición
	var usados [modelo.TitularesPorEquipo]int
	nUsados := 0
	estaUsado := func(id int) bool {
		for _, u := range usados[:nUsados] {
			if u == id {
				return true
			}
		}
		return false
	}
	for i, puesto := range puestos {
		// Saltar a los que ya se usaron como emergencia en otra línea.
		for sig[puesto] < len(p.porPosicion[puesto]) && estaUsado(p.porPosicion[puesto][sig[puesto]].ID) {
			sig[puesto]++
		}
		var j *modelo.Jugador
		if sig[puesto] < len(p.porPosicion[puesto]) {
			j = p.porPosicion[puesto][sig[puesto]]
			sig[puesto]++
		} else {
			var ok bool
			if j, ok = p.deEmergencia(puesto, estaUsado); !ok {
				continue
			}
		}
		usados[nUsados] = j.ID
		nUsados++
		tit[i] = &titular{j: j, puesto: puesto}
	}
	return tit
}

// deEmergencia elige, cuando no queda nadie de la posición natural, al jugador sin
// usar que mejor rinde en el puesto: un portero solo para la portería, y nunca un
// portero en el campo.
func (p plantel) deEmergencia(puesto modelo.Posicion, estaUsado func(int) bool) (*modelo.Jugador, bool) {
	var mejor *modelo.Jugador
	mejorV, hay := -1, false
	for _, j := range p.todos {
		if estaUsado(j.ID) || j.Posicion == puesto || (j.Posicion == modelo.Portero && puesto != modelo.Portero) {
			continue
		}
		if v := j.ValoracionEn(puesto); v > mejorV {
			mejor, mejorV, hay = j, v, true
		}
	}
	return mejor, hay
}

// alineacion arma la alineación de la formación: el once y el banquillo (el resto
// de los disponibles, de mayor a menor valoración).
func (p plantel) alineacion(f modelo.Formacion) modelo.Alineacion {
	al := modelo.Alineacion{Formacion: f}
	usados := map[int]bool{}
	for i, t := range p.seleccion(f) {
		if t != nil {
			al.Titulares[i] = t.j.ID
			usados[t.j.ID] = true
		}
	}
	for _, j := range p.todos {
		if !usados[j.ID] {
			al.Banquillo = append(al.Banquillo, j.ID)
		}
	}
	return al
}

// Fuerzas devuelve la fuerza de ataque y de defensa de un equipo (en la escala de
// los atributos) con la alineación dada: cada titular aporta según los atributos
// del puesto en que juega, por su familiaridad con él, y la formación ajusta el
// resultado con su perfil.
func Fuerzas(e modelo.Equipo, al modelo.Alineacion) (ataque, defensa float64) {
	puestos := al.Formacion.Puestos()
	var tit [modelo.TitularesPorEquipo]*titular
	for i, id := range al.Titulares {
		for k := range e.Plantilla {
			if e.Plantilla[k].ID == id {
				tit[i] = &titular{j: &e.Plantilla[k], puesto: puestos[i]}
				break
			}
		}
	}
	return fuerzasDe(tit, al.Formacion)
}

// fuerzasAutomaticas es la fuerza de un equipo con su mejor once de la formación
// dada, sin construir la alineación.
func fuerzasAutomaticas(e modelo.Equipo, f modelo.Formacion) (ataque, defensa float64) {
	return fuerzasDe(nuevoPlantel(e, nil).seleccion(f), f)
}

// fuerzasDe calcula la fuerza de un once ya elegido.
func fuerzasDe(tit [modelo.TitularesPorEquipo]*titular, f modelo.Formacion) (ataque, defensa float64) {
	var a, d media
	for _, t := range tit {
		if t == nil {
			continue
		}
		j, puesto := *t.j, t.puesto
		fam := modelo.Familiaridad(j.Posicion, puesto)
		at := j.Atributos
		switch puesto {
		case modelo.Delantero:
			v := fam * (float64(at.Tiro) + float64(at.Regate) + float64(at.Ritmo) + float64(at.Pase)) / 4
			a.sumar(2, v)
			a.sumarCrack(j, puesto, crackDelantero)
		case modelo.Mediocampista:
			v := fam * (float64(at.Pase) + float64(at.Regate) + float64(at.Tiro)) / 3
			a.sumar(1, v)
			a.sumarCrack(j, puesto, crackMediocampo)
			d.sumar(0.5, fam*float64(at.Defensa))
		case modelo.Defensa:
			v := fam * (float64(at.Defensa) + float64(at.Fisico)) / 2
			d.sumar(1, v)
			d.sumarCrack(j, puesto, crackDefensa)
		case modelo.Portero:
			v := fam * float64(at.Reflejos)
			d.sumar(3, v)
			d.sumarCrack(j, puesto, crackPortero)
		}
	}
	p := perfil[f]
	return a.valor() * p.ataque, d.valor() * p.defensa
}
