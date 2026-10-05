package aplicacion

import (
	"fmt"

	"github.com/ETurriza/juego_futbol/internal/modelo"
	"github.com/ETurriza/juego_futbol/internal/simulacion"
)

// equipoUsuario es el equipo que dirige el usuario.
func (c *Carrera) equipoUsuario() modelo.Equipo { return c.Temporada.Equipos[c.Usuario] }

// alineacionesElegidas da las alineaciones que no son automáticas, por índice de
// equipo: la del usuario si la eligió y sigue siendo válida.
func (c *Carrera) alineacionesElegidas() map[int]modelo.Alineacion {
	if al, manual := c.AlineacionVigente(); manual {
		return map[int]modelo.Alineacion{c.Usuario: al}
	}
	return nil
}

// AlineacionVigente devuelve la alineación con que jugará el equipo del usuario y
// si es la que él eligió (true) o la automática (false). La elegida solo vale
// mientras sea válida para la plantilla actual.
func (c *Carrera) AlineacionVigente() (modelo.Alineacion, bool) {
	if c.AlineacionManual && c.Alineacion.Validar(c.equipoUsuario(), nil) == nil {
		return c.Alineacion, true
	}
	return c.AlineacionAutomatica(nil), false
}

// AlineacionElegidaInvalida devuelve el motivo por el que la alineación elegida
// por el usuario ya no vale, o nil si no la eligió o sigue siendo válida. Sirve
// para avisarle de que se está usando la automática.
func (c *Carrera) AlineacionElegidaInvalida() error {
	if !c.AlineacionManual {
		return nil
	}
	return c.Alineacion.Validar(c.equipoUsuario(), nil)
}

// ElegirAlineacion fija la alineación del equipo del usuario. Devuelve un error,
// sin cambiar nada, si no es válida para la plantilla.
func (c *Carrera) ElegirAlineacion(al modelo.Alineacion) error {
	if err := al.Validar(c.equipoUsuario(), nil); err != nil {
		return fmt.Errorf("alineacion invalida: %w", err)
	}
	c.Alineacion = modelo.Alineacion{
		Formacion: al.Formacion,
		Titulares: al.Titulares,
		Banquillo: append([]int(nil), al.Banquillo...),
	}
	c.AlineacionManual = true
	return nil
}

// UsarAlineacionAutomatica descarta la alineación elegida: el equipo vuelve a
// jugar con la automática.
func (c *Carrera) UsarAlineacionAutomatica() {
	c.Alineacion = modelo.Alineacion{}
	c.AlineacionManual = false
}

// AlineacionAutomatica propone la mejor alineación para el equipo del usuario: con
// la formación dada, o con la que mejor rinda si es nil.
func (c *Carrera) AlineacionAutomatica(formacion *modelo.Formacion) modelo.Alineacion {
	return simulacion.AlineacionAutomatica(c.equipoUsuario(), simulacion.Criterios{Formacion: formacion})
}
