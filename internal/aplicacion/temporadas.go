package aplicacion

import (
	"errors"
	"math/rand"

	"github.com/ETurriza/juego_futbol/internal/generador"
	"github.com/ETurriza/juego_futbol/internal/liga"
	"github.com/ETurriza/juego_futbol/internal/modelo"
	"github.com/ETurriza/juego_futbol/internal/progresion"
)

// ErrTemporadaEnCurso se devuelve al pedir la siguiente temporada antes de que
// termine la actual.
var ErrTemporadaEnCurso = errors.New("la temporada todavia no termino")

// ResumenTemporada es lo que queda en el historial de una temporada terminada.
type ResumenTemporada struct {
	Numero        int
	Campeon       string
	PuestoUsuario int
	PuntosUsuario int
}

// CambiosTemporada describe qué cambió en las plantillas al pasar de
// temporada.
type CambiosTemporada struct {
	// Numero es el número de la temporada que empieza.
	Numero int
	// Retirados y Juveniles son los del equipo del usuario.
	Retirados []modelo.Jugador
	Juveniles []modelo.Jugador
	// RetiradosLiga y JuvenilesLiga cuentan los de todos los equipos.
	RetiradosLiga int
	JuvenilesLiga int
}

// SiguienteTemporada cierra la temporada terminada y empieza la siguiente en el
// mismo club: guarda su resumen en el historial, hace envejecer a todos los
// jugadores de la liga, retira a los veteranos, repone con juveniles y arma un
// calendario nuevo. Antes archiva las estadísticas de la temporada que termina. Es atómica: si falla, la carrera queda como estaba.
func (c *Carrera) SiguienteTemporada() (CambiosTemporada, error) {
	if !c.Terminada() {
		return CambiosTemporada{}, ErrTemporadaEnCurso
	}

	resumen := c.resumenTemporada()
	archivo, err := c.estadisticasParaArchivar()
	if err != nil {
		return CambiosTemporada{}, err
	}
	r := rand.New(rand.NewSource(semillaEvolucion(c.Semilla, c.Numero)))
	proximoID := c.ProximoID

	cambios := CambiosTemporada{Numero: c.Numero + 1}
	equipos := make([]modelo.Equipo, len(c.Temporada.Equipos))
	for i, e := range c.Temporada.Equipos {
		nuevo, retirados, juveniles := renovarEquipo(r, e, &proximoID)
		equipos[i] = nuevo
		cambios.RetiradosLiga += len(retirados)
		cambios.JuvenilesLiga += len(juveniles)
		if i == c.Usuario {
			cambios.Retirados, cambios.Juveniles = retirados, juveniles
		}
	}

	temporada, err := liga.Nueva(equipos)
	if err != nil {
		return CambiosTemporada{}, err
	}
	c.Historial = append(c.Historial, resumen)
	c.Archivo = append(c.Archivo, archivo...)
	c.Temporada = temporada
	c.Numero++
	c.ProximoID = proximoID
	return cambios, nil
}

// resumenTemporada arma el resumen de la temporada terminada.
func (c *Carrera) resumenTemporada() ResumenTemporada {
	tabla := c.Tabla()
	resumen := ResumenTemporada{Numero: c.Numero, Campeon: tabla[0].Equipo}
	for _, f := range tabla {
		if f.EsDelUsuario {
			resumen.PuestoUsuario, resumen.PuntosUsuario = f.Posicion, f.Pts
		}
	}
	return resumen
}

// renovarEquipo pasa un año: cada jugador se retira o envejece, y las
// posiciones que queden por debajo de la composición mínima se reponen con
// juveniles. Nunca recorta una plantilla que ya supera el mínimo. Devuelve el
// equipo nuevo, los retirados y los juveniles incorporados.
func renovarEquipo(r *rand.Rand, e modelo.Equipo, proximoID *int) (modelo.Equipo, []modelo.Jugador, []modelo.Jugador) {
	nuevo := modelo.Equipo{Nombre: e.Nombre, Plantilla: make([]modelo.Jugador, 0, len(e.Plantilla))}
	var retirados, juveniles []modelo.Jugador

	for _, j := range e.Plantilla {
		if progresion.SeRetira(r, j) {
			retirados = append(retirados, j)
			continue
		}
		nuevo.Plantilla = append(nuevo.Plantilla, progresion.Envejecer(r, j))
	}

	for _, p := range modelo.Posiciones {
		for faltan := generador.MinimoPorPosicion(p) - nuevo.Contar(p); faltan > 0; faltan-- {
			j := generador.Juvenil(r, *proximoID, p)
			// Evita repetir el nombre de otro jugador del equipo; el tope
			// impide un bucle infinito si se achican las listas de nombres.
			for intento := 0; nombreUsado(nuevo, j.Nombre) && intento < 50; intento++ {
				j.Nombre = generador.NombreJugador(r)
			}
			*proximoID++
			nuevo.Plantilla = append(nuevo.Plantilla, j)
			juveniles = append(juveniles, j)
		}
	}
	return nuevo, retirados, juveniles
}

func nombreUsado(e modelo.Equipo, nombre string) bool {
	for _, j := range e.Plantilla {
		if j.Nombre == nombre {
			return true
		}
	}
	return false
}
