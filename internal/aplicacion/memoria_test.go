package aplicacion_test

import (
	"testing"

	"github.com/ETurriza/juego_futbol/internal/aplicacion"
	"github.com/ETurriza/juego_futbol/internal/aplicacion/contrato"
)

func TestRepositorioMemoriaCumpleElContrato(t *testing.T) {
	contrato.Ejecutar(t, func(t *testing.T) aplicacion.RepositorioPartidas {
		return aplicacion.NuevoRepositorioMemoria()
	})
}
