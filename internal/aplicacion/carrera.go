package aplicacion

import (
	"fmt"
	"math/rand"
	"sort"

	"github.com/ETurriza/juego_futbol/internal/generador"
	"github.com/ETurriza/juego_futbol/internal/liga"
	"github.com/ETurriza/juego_futbol/internal/modelo"
)

// Límites de equipos de una liga. El máximo lo marca la cantidad de nombres de
// club distintos que puede inventar el generador.
const (
	MinEquipos = 2
	MaxEquipos = 40
)

// ErrTemporadaTerminada se devuelve al avanzar una carrera cuya temporada ya
// terminó.
var ErrTemporadaTerminada = liga.ErrTemporadaTerminada

// Carrera es una partida en modo carrera: una liga en curso y el equipo que
// dirige el usuario.
type Carrera struct {
	Semilla   int64
	Temporada *liga.Temporada
	// Usuario es el índice, en Temporada.Equipos, del equipo del usuario.
	Usuario int
	// Numero es el número de la temporada en curso, desde 1.
	Numero int
	// Historial tiene un resumen por cada temporada ya terminada, en orden.
	Historial []ResumenTemporada
	// Archivo tiene las estadísticas de los jugadores de la liga en cada
	// temporada terminada (solo de quienes jugaron).
	Archivo []EstadisticaTemporada
	// ProximoID es la ID que recibirá el próximo jugador creado. Nunca se
	// reutiliza una ID, ni la de un jugador retirado.
	ProximoID int
}

// ResultadoPartido es un partido jugado, con los nombres de los equipos.
type ResultadoPartido struct {
	Local          string
	Visitante      string
	GolesLocal     int
	GolesVisitante int
	// EsDelUsuario indica si juega el equipo del usuario.
	EsDelUsuario bool
}

// FilaTabla es una línea de la tabla de posiciones con su puesto.
type FilaTabla struct {
	Posicion int // desde 1
	liga.Fila
	EsDelUsuario bool
}

// NuevaCarrera crea una carrera con numEquipos equipos inventados y asigna uno
// al usuario. La misma semilla y cantidad dan siempre la misma carrera.
func NuevaCarrera(semilla int64, numEquipos int) (*Carrera, error) {
	if numEquipos < MinEquipos || numEquipos > MaxEquipos {
		return nil, fmt.Errorf("numero de equipos fuera de rango [%d,%d]: %d",
			MinEquipos, MaxEquipos, numEquipos)
	}
	r := rand.New(rand.NewSource(semilla))

	equipos := make([]modelo.Equipo, 0, numEquipos)
	nombres := make(map[string]bool, numEquipos)
	for i := 0; i < numEquipos; i++ {
		nombre := nombreNuevo(r, nombres)
		nombres[nombre] = true
		// IDs de jugador únicas en toda la liga.
		equipos = append(equipos, generador.Equipo(r, nombre, i*generador.TamanoPlantilla+1))
	}

	temporada, err := liga.Nueva(equipos)
	if err != nil {
		return nil, err
	}
	return &Carrera{
		Semilla:   semilla,
		Temporada: temporada,
		Usuario:   r.Intn(numEquipos),
		Numero:    1,
		ProximoID: numEquipos*generador.TamanoPlantilla + 1,
	}, nil
}

// nombreNuevo inventa un nombre de club que no esté en usados.
func nombreNuevo(r *rand.Rand, usados map[string]bool) string {
	for {
		if n := generador.NombreEquipo(r); !usados[n] {
			return n
		}
	}
}

// semillaJornada deriva la semilla de una jornada a partir de la semilla de la
// carrera (mezcla tipo splitmix64), de modo que jornadas distintas no
// comparten secuencia aleatoria.
func semillaJornada(semilla int64, jornada int) int64 {
	return mezclar(uint64(semilla) + (uint64(jornada)+1)*0x9E3779B97F4A7C15)
}

// mezclar es el paso final de splitmix64: dispersa los bits de z.
func mezclar(z uint64) int64 {
	z = (z ^ (z >> 30)) * 0xBF58476D1CE4E5B9
	z = (z ^ (z >> 27)) * 0x94D049BB133111EB
	return int64(z ^ (z >> 31))
}

// semillaTemporada es la semilla base de las jornadas de la temporada en curso.
// La primera temporada usa la semilla de la carrera tal cual, de modo que una
// semilla sigue dando la misma liga de siempre; las siguientes derivan la suya.
func (c *Carrera) semillaTemporada() int64 {
	if c.Numero <= 1 {
		return c.Semilla
	}
	return mezclar(uint64(c.Semilla) ^ (uint64(c.Numero) * 0xD6E8FEB86659FD93))
}

// semillaEvolucion es la semilla con la que evolucionan las plantillas al
// terminar la temporada numero.
func semillaEvolucion(semilla int64, numero int) int64 {
	return mezclar(uint64(semilla) + uint64(numero)*0xA0761D6478BD642F)
}

// AvanzarJornada juega la próxima jornada y devuelve sus resultados. Devuelve
// ErrTemporadaTerminada si ya no quedan jornadas.
func (c *Carrera) AvanzarJornada() ([]ResultadoPartido, error) {
	r := rand.New(rand.NewSource(semillaJornada(c.semillaTemporada(), c.Temporada.JornadaActual())))
	resultados, err := c.Temporada.JugarJornada(r)
	if err != nil {
		return nil, err
	}
	return c.resultados(resultados), nil
}

func (c *Carrera) resultados(rs []liga.Resultado) []ResultadoPartido {
	out := make([]ResultadoPartido, len(rs))
	for i, res := range rs {
		out[i] = ResultadoPartido{
			Local:          c.Temporada.Equipos[res.Local].Nombre,
			Visitante:      c.Temporada.Equipos[res.Visitante].Nombre,
			GolesLocal:     res.GolesLocal,
			GolesVisitante: res.GolesVisitante,
			EsDelUsuario:   res.Local == c.Usuario || res.Visitante == c.Usuario,
		}
	}
	return out
}

// UltimaJornada devuelve los resultados de la última jornada jugada, o nil si
// todavía no se jugó ninguna.
func (c *Carrera) UltimaJornada() []ResultadoPartido {
	if len(c.Temporada.Resultados) == 0 {
		return nil
	}
	return c.resultados(c.Temporada.Resultados[len(c.Temporada.Resultados)-1])
}

// Jornada es el número de jornadas jugadas.
func (c *Carrera) Jornada() int { return c.Temporada.JornadaActual() }

// TotalJornadas es el número de jornadas de la temporada.
func (c *Carrera) TotalJornadas() int { return len(c.Temporada.Calendario) }

// Terminada indica si la temporada ya terminó.
func (c *Carrera) Terminada() bool { return c.Temporada.Terminada() }

// NombreEquipo es el nombre del equipo del usuario.
func (c *Carrera) NombreEquipo() string { return c.equipo().Nombre }

// ValoracionEquipo es la valoración (0-99) del equipo del usuario.
func (c *Carrera) ValoracionEquipo() int { return c.equipo().Valoracion() }

func (c *Carrera) equipo() modelo.Equipo { return c.Temporada.Equipos[c.Usuario] }

// Tabla devuelve la tabla de posiciones actual.
func (c *Carrera) Tabla() []FilaTabla {
	nombre := c.NombreEquipo()
	filas := c.Temporada.Tabla()
	out := make([]FilaTabla, len(filas))
	for i, f := range filas {
		out[i] = FilaTabla{Posicion: i + 1, Fila: f, EsDelUsuario: f.Equipo == nombre}
	}
	return out
}

// Campeon devuelve el nombre del campeón; solo existe cuando la temporada
// terminó.
func (c *Carrera) Campeon() (string, bool) {
	if !c.Terminada() {
		return "", false
	}
	return c.Temporada.Tabla()[0].Equipo, true
}

// Plantilla devuelve una copia de la plantilla del usuario, ordenada por
// posición (portero a delantero), luego por valoración descendente y por ID.
func (c *Carrera) Plantilla() []modelo.Jugador {
	plantilla := append([]modelo.Jugador(nil), c.equipo().Plantilla...)
	sort.SliceStable(plantilla, func(a, b int) bool {
		ja, jb := plantilla[a], plantilla[b]
		switch {
		case ja.Posicion != jb.Posicion:
			return ja.Posicion < jb.Posicion
		case ja.Valoracion() != jb.Valoracion():
			return ja.Valoracion() > jb.Valoracion()
		default:
			return ja.ID < jb.ID
		}
	})
	return plantilla
}
