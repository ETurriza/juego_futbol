package modelo

import (
	"errors"
	"fmt"
	"sort"
)

// Reglas de composición de una plantilla.
const (
	MinPorteros  = 2
	MinPlantilla = 16

	// titulares: un portero y diez jugadores de campo.
	titularesCampo = 10
)

// Equipo es un club con su plantilla.
type Equipo struct {
	Nombre    string
	Plantilla []Jugador
}

// Contar devuelve cuántos jugadores de la plantilla juegan en la posición dada.
func (e Equipo) Contar(p Posicion) int {
	n := 0
	for _, j := range e.Plantilla {
		if j.Posicion == p {
			n++
		}
	}
	return n
}

// Valoracion es el promedio (0-99) del mejor portero y los diez mejores
// jugadores de campo. Una plantilla sin portero o con menos de diez jugadores
// de campo promedia solo a los que hay; sin jugadores devuelve 0.
func (e Equipo) Valoracion() int {
	var porteros, campo []int
	for _, j := range e.Plantilla {
		if j.Posicion == Portero {
			porteros = append(porteros, j.Valoracion())
		} else {
			campo = append(campo, j.Valoracion())
		}
	}
	sort.Sort(sort.Reverse(sort.IntSlice(porteros)))
	sort.Sort(sort.Reverse(sort.IntSlice(campo)))

	var titulares []int
	if len(porteros) > 0 {
		titulares = append(titulares, porteros[0])
	}
	if len(campo) > titularesCampo {
		campo = campo[:titularesCampo]
	}
	titulares = append(titulares, campo...)
	if len(titulares) == 0 {
		return 0
	}
	suma := 0
	for _, v := range titulares {
		suma += v
	}
	return (suma + len(titulares)/2) / len(titulares)
}

// Validar comprueba el nombre, los jugadores, las IDs únicas y la composición
// mínima de la plantilla.
func (e Equipo) Validar() error {
	if e.Nombre == "" {
		return errors.New("equipo sin nombre")
	}
	if len(e.Plantilla) < MinPlantilla {
		return fmt.Errorf("equipo %q: plantilla de %d jugadores, minimo %d",
			e.Nombre, len(e.Plantilla), MinPlantilla)
	}
	if n := e.Contar(Portero); n < MinPorteros {
		return fmt.Errorf("equipo %q: %d porteros, minimo %d", e.Nombre, n, MinPorteros)
	}
	vistos := make(map[int]bool, len(e.Plantilla))
	for _, j := range e.Plantilla {
		if err := j.Validar(); err != nil {
			return fmt.Errorf("equipo %q: %w", e.Nombre, err)
		}
		if vistos[j.ID] {
			return fmt.Errorf("equipo %q: ID de jugador repetida: %d", e.Nombre, j.ID)
		}
		vistos[j.ID] = true
	}
	return nil
}
