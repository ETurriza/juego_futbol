package modelo

import "fmt"

// Posicion es el puesto en el que juega un jugador.
type Posicion int

const (
	Portero Posicion = iota
	Defensa
	Mediocampista
	Delantero
)

// Posiciones lista todas las posiciones válidas en orden estable.
var Posiciones = []Posicion{Portero, Defensa, Mediocampista, Delantero}

// Valida indica si p es una posición conocida.
func (p Posicion) Valida() bool {
	return p >= Portero && p <= Delantero
}

func (p Posicion) String() string {
	switch p {
	case Portero:
		return "Portero"
	case Defensa:
		return "Defensa"
	case Mediocampista:
		return "Mediocampista"
	case Delantero:
		return "Delantero"
	default:
		return fmt.Sprintf("Posicion(%d)", int(p))
	}
}
