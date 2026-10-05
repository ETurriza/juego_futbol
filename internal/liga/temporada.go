package liga

import (
	"errors"
	"fmt"
	"math/rand"
	"sort"

	"github.com/ETurriza/juego_futbol/internal/modelo"
	"github.com/ETurriza/juego_futbol/internal/simulacion"
)

// Puntos por resultado.
const (
	PuntosVictoria = 3
	PuntosEmpate   = 1
)

// ErrTemporadaTerminada se devuelve al intentar jugar una jornada cuando ya se
// jugaron todas.
var ErrTemporadaTerminada = errors.New("la temporada ya termino")

// Resultado es un partido ya jugado.
type Resultado struct {
	Partido
	GolesLocal     int
	GolesVisitante int
}

// Fila es una línea de la tabla de posiciones.
type Fila struct {
	Equipo string
	PJ     int // partidos jugados
	G      int // ganados
	E      int // empatados
	P      int // perdidos
	GF     int // goles a favor
	GC     int // goles en contra
	DG     int // diferencia de goles
	Pts    int
}

// Temporada es una liga en curso. La tabla se calcula a partir de los
// resultados, que son la única fuente de verdad.
type Temporada struct {
	Equipos    []modelo.Equipo
	Calendario []Jornada
	// Resultados[i] son los partidos jugados de la jornada i; su longitud es la
	// cantidad de jornadas jugadas.
	Resultados [][]Resultado
}

// Nueva crea una temporada para los equipos dados. Los nombres deben ser
// únicos y no vacíos, y se necesitan al menos dos equipos.
func Nueva(equipos []modelo.Equipo) (*Temporada, error) {
	vistos := make(map[string]bool, len(equipos))
	for i, e := range equipos {
		if e.Nombre == "" {
			return nil, fmt.Errorf("el equipo %d no tiene nombre", i)
		}
		if vistos[e.Nombre] {
			return nil, fmt.Errorf("nombre de equipo repetido: %q", e.Nombre)
		}
		vistos[e.Nombre] = true
	}
	calendario, err := Calendario(len(equipos))
	if err != nil {
		return nil, err
	}
	return &Temporada{
		Equipos:    append([]modelo.Equipo(nil), equipos...),
		Calendario: calendario,
	}, nil
}

// JornadaActual es el número de jornadas ya jugadas (también el índice de la
// próxima jornada).
func (t *Temporada) JornadaActual() int { return len(t.Resultados) }

// Terminada indica si ya se jugaron todas las jornadas.
func (t *Temporada) Terminada() bool { return t.JornadaActual() >= len(t.Calendario) }

// JugarJornada simula los partidos de la próxima jornada y devuelve sus
// resultados. Devuelve ErrTemporadaTerminada si no quedan jornadas.
func (t *Temporada) JugarJornada(r *rand.Rand) ([]Resultado, error) {
	if t.Terminada() {
		return nil, ErrTemporadaTerminada
	}
	jornada := t.Calendario[t.JornadaActual()]
	resultados := make([]Resultado, 0, len(jornada))
	for _, p := range jornada {
		marcador := simulacion.Simular(r, t.Equipos[p.Local], t.Equipos[p.Visitante])
		resultados = append(resultados, Resultado{
			Partido:        p,
			GolesLocal:     marcador.GolesLocal,
			GolesVisitante: marcador.GolesVisitante,
		})
	}
	t.Resultados = append(t.Resultados, resultados)
	return resultados, nil
}

// JugarTemporada juega todas las jornadas que falten.
func (t *Temporada) JugarTemporada(r *rand.Rand) error {
	for !t.Terminada() {
		if _, err := t.JugarJornada(r); err != nil {
			return err
		}
	}
	return nil
}

// Tabla devuelve la tabla de posiciones con los resultados jugados hasta ahora,
// ordenada por puntos, diferencia de goles, goles a favor y nombre.
func (t *Temporada) Tabla() []Fila {
	filas := make([]Fila, len(t.Equipos))
	for i, e := range t.Equipos {
		filas[i].Equipo = e.Nombre
	}
	for _, jornada := range t.Resultados {
		for _, res := range jornada {
			registrar(&filas[res.Local], res.GolesLocal, res.GolesVisitante)
			registrar(&filas[res.Visitante], res.GolesVisitante, res.GolesLocal)
		}
	}
	sort.SliceStable(filas, func(a, b int) bool {
		fa, fb := filas[a], filas[b]
		switch {
		case fa.Pts != fb.Pts:
			return fa.Pts > fb.Pts
		case fa.DG != fb.DG:
			return fa.DG > fb.DG
		case fa.GF != fb.GF:
			return fa.GF > fb.GF
		default:
			return fa.Equipo < fb.Equipo
		}
	})
	return filas
}

// registrar suma un partido a la fila de un equipo desde su punto de vista.
func registrar(f *Fila, aFavor, enContra int) {
	f.PJ++
	f.GF += aFavor
	f.GC += enContra
	f.DG = f.GF - f.GC
	switch {
	case aFavor > enContra:
		f.G++
		f.Pts += PuntosVictoria
	case aFavor == enContra:
		f.E++
		f.Pts += PuntosEmpate
	default:
		f.P++
	}
}
