package modelo

import "fmt"

// Formacion es la disposición táctica de un equipo: cuántos defensas, medios y
// delanteros alinea además del portero. El valor cero es el 4-3-3, la formación
// de los partidos anteriores a que se pudiera elegir.
type Formacion int

const (
	F433 Formacion = iota
	F442
	F352
	F451
	F532
	F343
)

// Formaciones lista todas las formaciones en orden estable.
var Formaciones = []Formacion{F433, F442, F352, F451, F532, F343}

// lineasDe da los defensas, medios y delanteros de cada formación.
var lineasDe = map[Formacion][3]int{
	F433: {4, 3, 3},
	F442: {4, 4, 2},
	F352: {3, 5, 2},
	F451: {4, 5, 1},
	F532: {5, 3, 2},
	F343: {3, 4, 3},
}

// Valida indica si f es una formación conocida.
func (f Formacion) Valida() bool {
	_, ok := lineasDe[f]
	return ok
}

// Lineas devuelve cuántos defensas, medios y delanteros alinea la formación.
func (f Formacion) Lineas() (defensas, medios, delanteros int) {
	l := lineasDe[f]
	return l[0], l[1], l[2]
}

func (f Formacion) String() string {
	if !f.Valida() {
		return fmt.Sprintf("Formacion(%d)", int(f))
	}
	d, m, del := f.Lineas()
	return fmt.Sprintf("%d-%d-%d", d, m, del)
}

// Puestos devuelve la posición de cada puesto del once, en el orden en que se
// guardan los titulares: el portero, los defensas, los medios y los delanteros.
func (f Formacion) Puestos() [TitularesPorEquipo]Posicion {
	var out [TitularesPorEquipo]Posicion
	d, m, del := f.Lineas()
	i := 0
	out[i] = Portero
	i++
	for _, linea := range []struct {
		p Posicion
		n int
	}{{Defensa, d}, {Mediocampista, m}, {Delantero, del}} {
		for k := 0; k < linea.n && i < len(out); k++ {
			out[i] = linea.p
			i++
		}
	}
	return out
}
