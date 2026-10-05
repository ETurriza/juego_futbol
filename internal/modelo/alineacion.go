package modelo

import "fmt"

// Alineacion es el once que sale a jugar y el orden del banquillo.
type Alineacion struct {
	Formacion Formacion
	// Titulares son las IDs de los once jugadores, en el orden de
	// Formacion.Puestos(): portero, defensas, medios y delanteros.
	Titulares [TitularesPorEquipo]int
	// Banquillo son las IDs de los suplentes en orden de preferencia: los cambios
	// automáticos usan primero a los primeros.
	Banquillo []int
}

// Validar comprueba que la alineación sea utilizable con la plantilla del
// equipo:
//   - la formación es válida y hay once titulares distintos de la plantilla;
//   - el portero es un portero (y solo en una emergencia, cuando ninguno de los
//     porteros de la plantilla está disponible, puede serlo otro jugador);
//   - un portero no juega de jugador de campo;
//   - nadie de noDisponibles (lesionados o sancionados) está alineado;
//   - el banquillo no repite jugadores ni incluye a titulares ni a gente ajena.
func (a Alineacion) Validar(e Equipo, noDisponibles map[int]bool) error {
	if !a.Formacion.Valida() {
		return fmt.Errorf("formacion invalida: %d", int(a.Formacion))
	}
	plantilla := map[int]Jugador{}
	for _, j := range e.Plantilla {
		plantilla[j.ID] = j
	}
	hayPortero := false
	for _, j := range e.Plantilla {
		if j.Posicion == Portero && !noDisponibles[j.ID] {
			hayPortero = true
		}
	}

	puestos := a.Formacion.Puestos()
	vistos := map[int]bool{}
	for i, id := range a.Titulares {
		j, ok := plantilla[id]
		switch {
		case !ok:
			return fmt.Errorf("el titular del puesto %d (ID %d) no es de la plantilla", i+1, id)
		case vistos[id]:
			return fmt.Errorf("%s esta alineado dos veces", j.Nombre)
		case noDisponibles[id]:
			return fmt.Errorf("%s no esta disponible", j.Nombre)
		case puestos[i] == Portero && j.Posicion != Portero && hayPortero:
			return fmt.Errorf("%s no es portero y hay porteros disponibles", j.Nombre)
		case puestos[i] != Portero && j.Posicion == Portero:
			return fmt.Errorf("el portero %s no puede jugar de %v", j.Nombre, puestos[i])
		}
		vistos[id] = true
	}

	banca := map[int]bool{}
	for _, id := range a.Banquillo {
		j, ok := plantilla[id]
		switch {
		case !ok:
			return fmt.Errorf("el suplente con ID %d no es de la plantilla", id)
		case vistos[id]:
			return fmt.Errorf("%s es titular y esta en el banquillo", j.Nombre)
		case banca[id]:
			return fmt.Errorf("%s esta dos veces en el banquillo", j.Nombre)
		}
		banca[id] = true
	}
	return nil
}

// Puesto devuelve la posición del puesto i del once.
func (a Alineacion) Puesto(i int) Posicion {
	if i < 0 || i >= TitularesPorEquipo {
		return Portero
	}
	return a.Formacion.Puestos()[i]
}

// Vacia indica si no hay ningún titular (la alineación sin elegir).
func (a Alineacion) Vacia() bool {
	return a.Titulares == [TitularesPorEquipo]int{}
}
