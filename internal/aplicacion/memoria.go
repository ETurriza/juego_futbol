package aplicacion

import (
	"context"
	"fmt"
	"sort"
	"sync"
	"time"
)

// RepositorioMemoria es un RepositorioPartidas en memoria, para pruebas y para
// usar la aplicación sin disco.
type RepositorioMemoria struct {
	mu       sync.Mutex
	partidas map[string]entradaMemoria
	// ahora da la hora de cada guardado; las pruebas pueden reemplazarla.
	ahora func() time.Time
}

type entradaMemoria struct {
	guardado    Guardado
	actualizada time.Time
}

// NuevoRepositorioMemoria crea un repositorio vacío.
func NuevoRepositorioMemoria() *RepositorioMemoria {
	return &RepositorioMemoria{partidas: map[string]entradaMemoria{}, ahora: time.Now}
}

func (r *RepositorioMemoria) Guardar(ctx context.Context, ranura string, g Guardado) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if err := ValidarRanura(ranura); err != nil {
		return err
	}
	if _, err := g.Resumen(ranura); err != nil {
		return err
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	r.partidas[ranura] = entradaMemoria{guardado: g.clonar(), actualizada: r.ahora()}
	return nil
}

func (r *RepositorioMemoria) Cargar(ctx context.Context, ranura string) (Guardado, error) {
	if err := ctx.Err(); err != nil {
		return Guardado{}, err
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	e, ok := r.partidas[ranura]
	if !ok {
		return Guardado{}, fmt.Errorf("%w: %q", ErrPartidaNoExiste, ranura)
	}
	return e.guardado.clonar(), nil
}

func (r *RepositorioMemoria) Listar(ctx context.Context) ([]ResumenPartida, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	var out []ResumenPartida
	for ranura, e := range r.partidas {
		resumen, err := e.guardado.Resumen(ranura)
		if err != nil {
			return nil, err
		}
		resumen.Actualizada = e.actualizada
		out = append(out, resumen)
	}
	sort.Slice(out, func(a, b int) bool {
		if !out[a].Actualizada.Equal(out[b].Actualizada) {
			return out[a].Actualizada.After(out[b].Actualizada)
		}
		return out[a].Ranura < out[b].Ranura
	})
	return out, nil
}

func (r *RepositorioMemoria) Borrar(ctx context.Context, ranura string) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.partidas[ranura]; !ok {
		return fmt.Errorf("%w: %q", ErrPartidaNoExiste, ranura)
	}
	delete(r.partidas, ranura)
	return nil
}
