package aplicacion

import (
	"fmt"
	"sort"

	"github.com/ETurriza/juego_futbol/internal/liga"
	"github.com/ETurriza/juego_futbol/internal/modelo"
)

// EstadisticaJugador es el jugador, su club y sus números de la temporada.
type EstadisticaJugador struct {
	Jugador modelo.Jugador
	Equipo  string
	modelo.Estadisticas
}

// EstadisticaTemporada es lo que se archiva de un jugador al terminar una
// temporada: sus números y los datos con que se jugó (club, posición y edad de
// entonces), porque después pueden cambiar o el jugador retirarse.
type EstadisticaTemporada struct {
	Temporada int
	Jugador   int // ID
	Nombre    string
	Equipo    string
	Posicion  modelo.Posicion
	Edad      int
	modelo.Estadisticas
}

// EstadisticaEquipo son los números de un equipo en la temporada, con su
// máximo goleador y su máximo asistente.
type EstadisticaEquipo struct {
	liga.EstadisticaEquipo
	Goleador      string
	GolesGoleador int
	Asistente     string
	Asistencias   int
}

// Criterio es una clasificación de jugadores.
type Criterio int

const (
	// PorGoles ordena a los goleadores.
	PorGoles Criterio = iota
	// PorAsistencias ordena a los asistentes.
	PorAsistencias
	// PorTarjetas ordena por puntos de disciplina: una amarilla vale 1 y una
	// roja 3.
	PorTarjetas
	// PorImbatidas ordena a los porteros por porterías imbatidas.
	PorImbatidas
	// PorPocosGolesEncajados ordena a los porteros por goles encajados por
	// partido, de menos a más (el "Zamora").
	PorPocosGolesEncajados
	// PorValoracion ordena por valoración media.
	PorValoracion
)

// MinPartidosClasificacion es la cantidad de partidos que hace falta para entrar
// en las clasificaciones por promedio (el "Zamora" y la valoración), o las
// jornadas jugadas si todavía son menos.
const MinPartidosClasificacion = 5

// puntosDeDisciplina da el peso de las tarjetas de un jugador.
func puntosDeDisciplina(e modelo.Estadisticas) int { return e.Amarillas + 3*e.Rojas }

// EstadisticasJugadores devuelve a todos los jugadores de la liga con sus
// números de la temporada en curso (ceros para quien aún no jugó), por equipo y
// en el orden de su plantilla.
func (c *Carrera) EstadisticasJugadores() ([]EstadisticaJugador, error) {
	est, err := c.Temporada.Estadisticas()
	if err != nil {
		return nil, err
	}
	var out []EstadisticaJugador
	for _, e := range c.Temporada.Equipos {
		for _, j := range e.Plantilla {
			out = append(out, EstadisticaJugador{Jugador: j, Equipo: e.Nombre, Estadisticas: est[j.ID]})
		}
	}
	return out, nil
}

// EstadisticasEquipos devuelve los números de cada equipo en la temporada en
// curso, en el orden de la liga.
func (c *Carrera) EstadisticasEquipos() ([]EstadisticaEquipo, error) {
	jugadores, err := c.EstadisticasJugadores()
	if err != nil {
		return nil, err
	}
	base := c.Temporada.EstadisticasEquipos()
	out := make([]EstadisticaEquipo, len(base))
	for i, b := range base {
		out[i].EstadisticaEquipo = b
	}
	porNombre := map[string]int{}
	for i, b := range base {
		porNombre[b.Equipo] = i
	}
	for _, j := range jugadores {
		e := &out[porNombre[j.Equipo]]
		if j.Goles > e.GolesGoleador || (j.Goles == e.GolesGoleador && j.Goles > 0 && j.Jugador.Nombre < e.Goleador) {
			e.Goleador, e.GolesGoleador = j.Jugador.Nombre, j.Goles
		}
		if j.Asistencias > e.Asistencias || (j.Asistencias == e.Asistencias && j.Asistencias > 0 && j.Jugador.Nombre < e.Asistente) {
			e.Asistente, e.Asistencias = j.Jugador.Nombre, j.Asistencias
		}
	}
	return out, nil
}

// Clasificacion devuelve los mejores jugadores de la liga según el criterio,
// hasta limite (todos si limite <= 0). Solo entran quienes tienen algo que
// mostrar: goles, asistencias o tarjetas, o los partidos mínimos en los
// promedios. Los empates se deciden con criterios propios de cada clasificación
// y, al final, por nombre, de modo que el orden es siempre el mismo.
func (c *Carrera) Clasificacion(criterio Criterio, limite int) ([]EstadisticaJugador, error) {
	todos, err := c.EstadisticasJugadores()
	if err != nil {
		return nil, err
	}
	minimo := min(MinPartidosClasificacion, c.Jornada())

	var filas []EstadisticaJugador
	for _, j := range todos {
		var entra bool
		switch criterio {
		case PorGoles:
			entra = j.Goles > 0
		case PorAsistencias:
			entra = j.Asistencias > 0
		case PorTarjetas:
			entra = puntosDeDisciplina(j.Estadisticas) > 0
		case PorImbatidas:
			entra = j.Jugador.Posicion == modelo.Portero && j.PorteriasImbatidas > 0
		case PorPocosGolesEncajados:
			entra = j.Jugador.Posicion == modelo.Portero && j.Partidos >= max(minimo, 1)
		case PorValoracion:
			entra = j.Partidos >= max(minimo, 1)
		default:
			return nil, fmt.Errorf("criterio de clasificacion desconocido: %d", int(criterio))
		}
		if entra {
			filas = append(filas, j)
		}
	}

	// Cada comparación devuelve un valor negativo si x va antes que y. mayor
	// compara números de los que más es mejor.
	mayor := func(a, b int) int { return cmpInt(b, a) }
	sort.SliceStable(filas, func(a, b int) bool {
		x, y := filas[a], filas[b]
		var d int
		switch criterio {
		case PorGoles:
			d = firstNonZero(mayor(x.Goles, y.Goles), mayor(x.Asistencias, y.Asistencias), cmpInt(x.Minutos, y.Minutos))
		case PorAsistencias:
			d = firstNonZero(mayor(x.Asistencias, y.Asistencias), mayor(x.Goles, y.Goles), cmpInt(x.Minutos, y.Minutos))
		case PorTarjetas:
			d = firstNonZero(mayor(puntosDeDisciplina(x.Estadisticas), puntosDeDisciplina(y.Estadisticas)),
				mayor(x.Rojas, y.Rojas), cmpInt(x.Minutos, y.Minutos))
		case PorImbatidas:
			d = firstNonZero(mayor(x.PorteriasImbatidas, y.PorteriasImbatidas), cmpInt(x.GolesEncajados, y.GolesEncajados))
		case PorPocosGolesEncajados:
			// goles por partido, comparado sin dividir: x.GE/x.PJ < y.GE/y.PJ
			d = firstNonZero(cmpInt(x.GolesEncajados*y.Partidos, y.GolesEncajados*x.Partidos), mayor(x.Partidos, y.Partidos))
		case PorValoracion:
			d = firstNonZero(mayor(x.SumaValoracion*y.Partidos, y.SumaValoracion*x.Partidos), mayor(x.Partidos, y.Partidos))
		}
		if d != 0 {
			return d < 0
		}
		return x.Jugador.Nombre < y.Jugador.Nombre
	})
	if limite > 0 && len(filas) > limite {
		filas = filas[:limite]
	}
	return filas, nil
}

func cmpInt(a, b int) int {
	switch {
	case a < b:
		return -1
	case a > b:
		return 1
	}
	return 0
}

func firstNonZero(valores ...int) int {
	for _, v := range valores {
		if v != 0 {
			return v
		}
	}
	return 0
}

// EstadisticasDeCarrera devuelve los números acumulados de un jugador a lo
// largo de toda la carrera: las temporadas archivadas más la que está en curso.
func (c *Carrera) EstadisticasDeCarrera(jugador int) (modelo.Estadisticas, error) {
	total := modelo.Estadisticas{}
	for _, a := range c.Archivo {
		if a.Jugador == jugador {
			total = total.Sumar(a.Estadisticas)
		}
	}
	actual, err := c.Temporada.Estadisticas()
	if err != nil {
		return modelo.Estadisticas{}, err
	}
	return total.Sumar(actual[jugador]), nil
}

// estadisticasParaArchivar arma las filas que se archivan al terminar la
// temporada en curso: una por jugador que jugó, por equipo y en el orden de su
// plantilla.
func (c *Carrera) estadisticasParaArchivar() ([]EstadisticaTemporada, error) {
	jugadores, err := c.EstadisticasJugadores()
	if err != nil {
		return nil, err
	}
	var filas []EstadisticaTemporada
	for _, j := range jugadores {
		if j.Partidos == 0 {
			continue
		}
		filas = append(filas, EstadisticaTemporada{
			Temporada: c.Numero, Jugador: j.Jugador.ID, Nombre: j.Jugador.Nombre, Equipo: j.Equipo,
			Posicion: j.Jugador.Posicion, Edad: j.Jugador.Edad, Estadisticas: j.Estadisticas,
		})
	}
	return filas, nil
}

// EstadisticasDeEquipo devuelve a los jugadores de un equipo con sus números de
// la temporada en curso, ordenados por posición (portero a delantero), luego por
// valoración descendente y por ID. Devuelve nil si el equipo no existe.
func (c *Carrera) EstadisticasDeEquipo(nombre string) ([]EstadisticaJugador, error) {
	todos, err := c.EstadisticasJugadores()
	if err != nil {
		return nil, err
	}
	var filas []EstadisticaJugador
	for _, j := range todos {
		if j.Equipo == nombre {
			filas = append(filas, j)
		}
	}
	sort.SliceStable(filas, func(a, b int) bool {
		x, y := filas[a].Jugador, filas[b].Jugador
		switch {
		case x.Posicion != y.Posicion:
			return x.Posicion < y.Posicion
		case x.Valoracion() != y.Valoracion():
			return x.Valoracion() > y.Valoracion()
		default:
			return x.ID < y.ID
		}
	})
	return filas, nil
}

// EstadisticaDeJugador devuelve la fila de un jugador de la liga con sus
// números de la temporada en curso; false si no está en ningún equipo (por
// ejemplo, si ya se retiró).
func (c *Carrera) EstadisticaDeJugador(id int) (EstadisticaJugador, bool, error) {
	todos, err := c.EstadisticasJugadores()
	if err != nil {
		return EstadisticaJugador{}, false, err
	}
	for _, j := range todos {
		if j.Jugador.ID == id {
			return j, true, nil
		}
	}
	return EstadisticaJugador{}, false, nil
}

// Trayectoria devuelve las estadísticas archivadas de un jugador, una por cada
// temporada terminada en que jugó, de la más antigua a la más reciente.
func (c *Carrera) Trayectoria(jugador int) []EstadisticaTemporada {
	var filas []EstadisticaTemporada
	for _, a := range c.Archivo {
		if a.Jugador == jugador {
			filas = append(filas, a)
		}
	}
	sort.SliceStable(filas, func(a, b int) bool { return filas[a].Temporada < filas[b].Temporada })
	return filas
}
