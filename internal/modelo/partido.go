package modelo

import "fmt"

// Constantes de un partido.
const (
	TitularesPorEquipo = 11
	MinutosPartido     = 90
)

// TipoEvento es la clase de suceso de un partido.
type TipoEvento int

const (
	Gol TipoEvento = iota
	Amarilla
	Roja
	Sustitucion
)

// Valido indica si t es un tipo de suceso conocido.
func (t TipoEvento) Valido() bool { return t >= Gol && t <= Sustitucion }

func (t TipoEvento) String() string {
	switch t {
	case Gol:
		return "Gol"
	case Amarilla:
		return "Amarilla"
	case Roja:
		return "Roja"
	case Sustitucion:
		return "Sustitucion"
	default:
		return fmt.Sprintf("TipoEvento(%d)", int(t))
	}
}

// Evento es un suceso de un partido. Las IDs de jugador son siempre positivas;
// 0 significa "ninguno".
type Evento struct {
	Minuto int // de 1 a MinutosPartido
	Tipo   TipoEvento
	Local  bool // true si el suceso es del equipo local
	// Jugador es el goleador, el jugador amonestado o expulsado, o el que sale
	// en una sustitución.
	Jugador int
	// Otro es el asistente de un gol (0 si no hubo) o el jugador que entra en una
	// sustitución.
	Otro int
}

// DetallePartido son las alineaciones titulares y los sucesos de un partido,
// ordenados por minuto. De él se derivan los minutos de cada jugador y todas
// las estadísticas.
type DetallePartido struct {
	// Titulares de cada equipo; 0 en las posiciones sin jugador (equipos con
	// menos de once jugadores).
	TitularesLocal     [TitularesPorEquipo]int
	TitularesVisitante [TitularesPorEquipo]int
	Eventos            []Evento
}

// Vacio indica si el detalle no tiene información: ni titulares ni sucesos.
// Es el caso de los partidos anteriores a las estadísticas.
func (d DetallePartido) Vacio() bool {
	return len(d.Eventos) == 0 &&
		d.TitularesLocal == [TitularesPorEquipo]int{} &&
		d.TitularesVisitante == [TitularesPorEquipo]int{}
}

// Goles cuenta los goles de cada equipo según los sucesos.
func (d DetallePartido) Goles() (local, visitante int) {
	for _, e := range d.Eventos {
		if e.Tipo != Gol {
			continue
		}
		if e.Local {
			local++
		} else {
			visitante++
		}
	}
	return local, visitante
}

// Participacion es el tramo de un partido que jugó un jugador: estuvo en el
// campo durante los minutos (Desde, Hasta].
type Participacion struct {
	Jugador int
	Local   bool
	Desde   int // 0 para un titular; el minuto de entrada para un suplente
	Hasta   int // MinutosPartido, o el minuto de salida o expulsión
}

// Minutos son los minutos jugados.
func (p Participacion) Minutos() int { return p.Hasta - p.Desde }

// Titular indica si el jugador empezó el partido.
func (p Participacion) Titular() bool { return p.Desde == 0 }

// Participaciones reconstruye, a partir de las alineaciones y los sucesos, qué
// jugadores jugaron y durante cuántos minutos: los titulares primero (locales y
// luego visitantes) y los suplentes en orden de entrada. Devuelve un error si el
// detalle es incoherente: sucesos fuera de orden o de minuto, jugadores que no
// estaban en el campo, repetidos o del equipo equivocado.
func (d DetallePartido) Participaciones() ([]Participacion, error) {
	var out []Participacion
	indice := map[int]int{}
	agregar := func(id int, local bool, desde int) error {
		if id <= 0 {
			return fmt.Errorf("ID de jugador invalida: %d", id)
		}
		if _, ok := indice[id]; ok {
			return fmt.Errorf("el jugador %d aparece dos veces en el partido", id)
		}
		indice[id] = len(out)
		out = append(out, Participacion{Jugador: id, Local: local, Desde: desde, Hasta: MinutosPartido})
		return nil
	}
	for _, grupo := range []struct {
		ids   [TitularesPorEquipo]int
		local bool
	}{{d.TitularesLocal, true}, {d.TitularesVisitante, false}} {
		for _, id := range grupo.ids {
			if id == 0 {
				continue
			}
			if err := agregar(id, grupo.local, 0); err != nil {
				return nil, err
			}
		}
	}

	// enCampo comprueba que el jugador esté jugando en el minuto del suceso.
	enCampo := func(id int, local bool, minuto int) (int, error) {
		i, ok := indice[id]
		if !ok {
			return 0, fmt.Errorf("el jugador %d no participa en el partido", id)
		}
		p := out[i]
		if p.Local != local {
			return 0, fmt.Errorf("el jugador %d es del otro equipo", id)
		}
		if minuto <= p.Desde || minuto > p.Hasta {
			return 0, fmt.Errorf("el jugador %d no estaba en el campo en el minuto %d", id, minuto)
		}
		return i, nil
	}

	anterior := 0
	for n, e := range d.Eventos {
		switch {
		case !e.Tipo.Valido():
			return nil, fmt.Errorf("suceso %d: tipo invalido %v", n+1, e.Tipo)
		case e.Minuto < 1 || e.Minuto > MinutosPartido:
			return nil, fmt.Errorf("suceso %d: minuto %d fuera de [1,%d]", n+1, e.Minuto, MinutosPartido)
		case e.Minuto < anterior:
			return nil, fmt.Errorf("suceso %d: fuera de orden (minuto %d tras el %d)", n+1, e.Minuto, anterior)
		}
		anterior = e.Minuto

		i, err := enCampo(e.Jugador, e.Local, e.Minuto)
		if err != nil {
			return nil, fmt.Errorf("suceso %d (%v): %w", n+1, e.Tipo, err)
		}
		switch e.Tipo {
		case Gol:
			if e.Otro != 0 {
				if e.Otro == e.Jugador {
					return nil, fmt.Errorf("suceso %d: un jugador no se asiste a si mismo", n+1)
				}
				if _, err := enCampo(e.Otro, e.Local, e.Minuto); err != nil {
					return nil, fmt.Errorf("suceso %d (asistente): %w", n+1, err)
				}
			}
		case Roja:
			out[i].Hasta = e.Minuto
		case Sustitucion:
			out[i].Hasta = e.Minuto
			if err := agregar(e.Otro, e.Local, e.Minuto); err != nil {
				return nil, fmt.Errorf("suceso %d (entra): %w", n+1, err)
			}
		}
	}
	return out, nil
}
