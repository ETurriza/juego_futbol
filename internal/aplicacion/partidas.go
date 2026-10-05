package aplicacion

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"
	"unicode"
)

// MaxLongitudRanura es el largo máximo del nombre de una ranura.
const MaxLongitudRanura = 64

// ErrPartidaNoExiste se devuelve cuando no hay una partida en la ranura pedida.
var ErrPartidaNoExiste = errors.New("no existe una partida en esa ranura")

// ResumenPartida describe una partida guardada, para listarlas.
type ResumenPartida struct {
	Ranura        string
	Equipo        string // equipo del usuario
	Jornada       int    // jornadas jugadas
	TotalJornadas int
	Actualizada   time.Time
}

// RepositorioPartidas es el puerto de persistencia: guarda y recupera partidas
// por ranura. En local hay una ranura; en el servidor SSH habrá una por
// usuario. Las implementaciones deben ser seguras para uso concurrente y
// devolver ErrPartidaNoExiste (envuelto) cuando corresponda.
type RepositorioPartidas interface {
	// Guardar crea o reemplaza por completo la partida de la ranura.
	Guardar(ctx context.Context, ranura string, g Guardado) error
	// Cargar devuelve la partida de la ranura.
	Cargar(ctx context.Context, ranura string) (Guardado, error)
	// Listar devuelve las partidas guardadas, la más reciente primero.
	Listar(ctx context.Context) ([]ResumenPartida, error)
	// Borrar elimina la partida de la ranura.
	Borrar(ctx context.Context, ranura string) error
}

// ValidarRanura comprueba que el nombre de una ranura sea utilizable: no
// vacío, de largo acotado y sin caracteres de control ni espacios en los
// extremos.
func ValidarRanura(ranura string) error {
	switch {
	case strings.TrimSpace(ranura) == "":
		return errors.New("la ranura no puede estar vacia")
	case ranura != strings.TrimSpace(ranura):
		return errors.New("la ranura no puede empezar ni terminar con espacios")
	case len([]rune(ranura)) > MaxLongitudRanura:
		return fmt.Errorf("la ranura supera %d caracteres", MaxLongitudRanura)
	}
	for _, r := range ranura {
		if unicode.IsControl(r) {
			return errors.New("la ranura no puede tener caracteres de control")
		}
	}
	return nil
}

// GuardarCarrera guarda la carrera en la ranura dada.
func GuardarCarrera(ctx context.Context, repo RepositorioPartidas, ranura string, c *Carrera) error {
	if err := ValidarRanura(ranura); err != nil {
		return err
	}
	return repo.Guardar(ctx, ranura, c.Exportar())
}

// CargarCarrera recupera la carrera de la ranura dada y valida que sus datos
// sean coherentes.
func CargarCarrera(ctx context.Context, repo RepositorioPartidas, ranura string) (*Carrera, error) {
	if err := ValidarRanura(ranura); err != nil {
		return nil, err
	}
	g, err := repo.Cargar(ctx, ranura)
	if err != nil {
		return nil, err
	}
	return Importar(g)
}
