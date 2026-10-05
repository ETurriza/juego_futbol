package liga

import (
	"fmt"

	"github.com/ETurriza/juego_futbol/internal/modelo"
)

// Valoración de un partido, en décimas (60 es un 6,0).
const (
	valoracionBase      = 60
	valoracionMin       = 10
	valoracionMax       = 100
	bonoGol             = 10
	bonoAsistencia      = 6
	castigoAmarilla     = 3
	castigoRoja         = 15
	bonoImbatida        = 5
	castigoGolPortero   = 3
	castigoGolDefensa   = 2
	bonoVictoria        = 3
	castigoDerrota      = 3
	minutosParaImbatida = 60
)

// EstadisticaEquipo son los números de un equipo en la temporada.
type EstadisticaEquipo struct {
	Equipo    string
	PJ        int
	G, E, P   int
	GF, GC    int
	Imbatidas int // partidos sin recibir gol
	SinMarcar int // partidos sin marcar gol
	Amarillas int
	Rojas     int
}

// Estadisticas calcula las estadísticas de la temporada de cada jugador que
// participó en algún partido, a partir del detalle de los partidos jugados. Los
// partidos sin detalle no cuentan. Devuelve un error si el detalle de algún
// partido es incoherente.
func (t *Temporada) Estadisticas() (map[int]modelo.Estadisticas, error) {
	posicion := map[int]modelo.Posicion{}
	for _, e := range t.Equipos {
		for _, j := range e.Plantilla {
			posicion[j.ID] = j.Posicion
		}
	}
	total := map[int]modelo.Estadisticas{}
	for n, jornada := range t.Resultados {
		for k, res := range jornada {
			if res.Detalle.Vacio() {
				continue
			}
			parcial, err := estadisticasDePartido(res.Detalle, res.GolesLocal, res.GolesVisitante, posicion)
			if err != nil {
				return nil, fmt.Errorf("jornada %d partido %d: %w", n+1, k+1, err)
			}
			for id, e := range parcial {
				total[id] = total[id].Sumar(e)
			}
		}
	}
	return total, nil
}

// estadisticasDePartido calcula lo que hizo cada participante de un partido.
func estadisticasDePartido(d modelo.DetallePartido, golesLocal, golesVisitante int,
	posicion map[int]modelo.Posicion) (map[int]modelo.Estadisticas, error) {

	partes, err := d.Participaciones()
	if err != nil {
		return nil, err
	}
	if gl, gv := d.Goles(); gl != golesLocal || gv != golesVisitante {
		return nil, fmt.Errorf("los sucesos suman %d-%d y el marcador es %d-%d", gl, gv, golesLocal, golesVisitante)
	}

	out := make(map[int]modelo.Estadisticas, len(partes))
	for _, p := range partes {
		if p.Minutos() <= 0 {
			continue
		}
		propios, rivales := golesLocal, golesVisitante
		if !p.Local {
			propios, rivales = golesVisitante, golesLocal
		}
		pos := posicion[p.Jugador]

		e := modelo.Estadisticas{Partidos: 1, Minutos: p.Minutos()}
		if p.Titular() {
			e.Titularidades = 1
		}
		var delta int
		for _, ev := range d.Eventos {
			switch ev.Tipo {
			case modelo.Gol:
				if ev.Local == p.Local && ev.Jugador == p.Jugador {
					e.Goles++
				}
				if ev.Local == p.Local && ev.Otro == p.Jugador {
					e.Asistencias++
				}
				// Gol del rival mientras este jugador estaba en el campo.
				if ev.Local != p.Local && ev.Minuto > p.Desde && ev.Minuto <= p.Hasta {
					e.GolesEncajados++
				}
			case modelo.Amarilla:
				if ev.Jugador == p.Jugador {
					e.Amarillas++
				}
			case modelo.Roja:
				if ev.Jugador == p.Jugador {
					e.Rojas++
				}
			}
		}
		defensivo := pos == modelo.Portero || pos == modelo.Defensa
		imbatida := rivales == 0 && defensivo && p.Minutos() >= minutosParaImbatida
		if imbatida {
			e.PorteriasImbatidas = 1
		}

		delta += bonoGol*e.Goles + bonoAsistencia*e.Asistencias
		delta -= castigoAmarilla*e.Amarillas + castigoRoja*e.Rojas
		if imbatida {
			delta += bonoImbatida
		}
		switch pos {
		case modelo.Portero:
			delta -= castigoGolPortero * e.GolesEncajados
		case modelo.Defensa:
			delta -= castigoGolDefensa * e.GolesEncajados
		}
		switch {
		case propios > rivales:
			delta += bonoVictoria
		case propios < rivales:
			delta -= castigoDerrota
		}
		// Quien juega poco no puede alejarse mucho del 6,0.
		if p.Minutos() < minutosParaImbatida {
			delta = delta * p.Minutos() / minutosParaImbatida
		}
		e.SumaValoracion = min(max(valoracionBase+delta, valoracionMin), valoracionMax)
		out[p.Jugador] = e
	}
	return out, nil
}

// EstadisticasEquipos devuelve los números de cada equipo en la temporada, en el
// mismo orden que Equipos. Los resultados y los goles cuentan siempre; las
// tarjetas, solo en los partidos con detalle.
func (t *Temporada) EstadisticasEquipos() []EstadisticaEquipo {
	out := make([]EstadisticaEquipo, len(t.Equipos))
	for i, e := range t.Equipos {
		out[i].Equipo = e.Nombre
	}
	sumar := func(i, aFavor, enContra int) {
		e := &out[i]
		e.PJ++
		e.GF += aFavor
		e.GC += enContra
		switch {
		case aFavor > enContra:
			e.G++
		case aFavor == enContra:
			e.E++
		default:
			e.P++
		}
		if enContra == 0 {
			e.Imbatidas++
		}
		if aFavor == 0 {
			e.SinMarcar++
		}
	}
	for _, jornada := range t.Resultados {
		for _, res := range jornada {
			sumar(res.Local, res.GolesLocal, res.GolesVisitante)
			sumar(res.Visitante, res.GolesVisitante, res.GolesLocal)
			for _, ev := range res.Detalle.Eventos {
				equipo := res.Visitante
				if ev.Local {
					equipo = res.Local
				}
				switch ev.Tipo {
				case modelo.Amarilla:
					out[equipo].Amarillas++
				case modelo.Roja:
					out[equipo].Rojas++
				}
			}
		}
	}
	return out
}
