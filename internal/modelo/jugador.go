package modelo

import (
	"errors"
	"fmt"
)

// Límites de edad de un jugador.
const (
	EdadMin = 16
	EdadMax = 45
)

// Jugador es un futbolista con sus atributos.
type Jugador struct {
	ID        int
	Nombre    string
	Edad      int
	Posicion  Posicion
	Atributos Atributos
}

// Valoracion devuelve la media del jugador en su posición.
func (j Jugador) Valoracion() int {
	return j.Atributos.Media(j.Posicion)
}

// Validar comprueba que el jugador sea consistente.
func (j Jugador) Validar() error {
	if j.Nombre == "" {
		return errors.New("jugador sin nombre")
	}
	if j.Edad < EdadMin || j.Edad > EdadMax {
		return fmt.Errorf("jugador %q: edad fuera de rango [%d,%d]: %d",
			j.Nombre, EdadMin, EdadMax, j.Edad)
	}
	if !j.Posicion.Valida() {
		return fmt.Errorf("jugador %q: posicion invalida: %v", j.Nombre, j.Posicion)
	}
	if err := j.Atributos.Validar(); err != nil {
		return fmt.Errorf("jugador %q: %w", j.Nombre, err)
	}
	return nil
}
