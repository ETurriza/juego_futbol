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

// Talento es un rasgo oculto que multiplica el crecimiento de un jugador: 100 es
// un crecimiento normal, 130 uno un 30 % más rápido y 80 uno un 20 % más lento.
// Decide hasta dónde puede llegar un jugador y es lo que separa a los cracks.
const (
	TalentoNeutro = 100
	TalentoMin    = 60
	TalentoMax    = 200
)

// Jugador es un futbolista con sus atributos.
type Jugador struct {
	ID        int
	Nombre    string
	Edad      int
	Posicion  Posicion
	Atributos Atributos
	// Talento es el rasgo oculto de crecimiento (ver TalentoNeutro). Cero, el
	// valor de un jugador sin talento definido, se trata como neutro.
	Talento int
}

// TalentoEfectivo devuelve el talento del jugador; un jugador sin talento
// definido (0) tiene el neutro.
func (j Jugador) TalentoEfectivo() int {
	if j.Talento == 0 {
		return TalentoNeutro
	}
	return j.Talento
}

// Proyeccion es la pista visible del talento de un jugador: lo que un ojeador
// diría de sus posibilidades de crecer.
type Proyeccion int

const (
	ProyeccionLimitada Proyeccion = iota
	ProyeccionNormal
	ProyeccionAlta
	ProyeccionExcepcional
)

// Umbrales de talento de cada proyección.
const (
	talentoProyeccionNormal      = 84
	talentoProyeccionAlta        = 118
	talentoProyeccionExcepcional = 138
)

func (p Proyeccion) String() string {
	switch p {
	case ProyeccionLimitada:
		return "limitada"
	case ProyeccionNormal:
		return "normal"
	case ProyeccionAlta:
		return "alta"
	case ProyeccionExcepcional:
		return "excepcional"
	default:
		return fmt.Sprintf("Proyeccion(%d)", int(p))
	}
}

// Proyeccion es la pista del talento del jugador.
func (j Jugador) Proyeccion() Proyeccion {
	switch t := j.TalentoEfectivo(); {
	case t >= talentoProyeccionExcepcional:
		return ProyeccionExcepcional
	case t >= talentoProyeccionAlta:
		return ProyeccionAlta
	case t >= talentoProyeccionNormal:
		return ProyeccionNormal
	default:
		return ProyeccionLimitada
	}
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
	if j.Talento != 0 && (j.Talento < TalentoMin || j.Talento > TalentoMax) {
		return fmt.Errorf("jugador %q: talento fuera de rango [%d,%d]: %d",
			j.Nombre, TalentoMin, TalentoMax, j.Talento)
	}
	if err := j.Atributos.Validar(); err != nil {
		return fmt.Errorf("jugador %q: %w", j.Nombre, err)
	}
	return nil
}
