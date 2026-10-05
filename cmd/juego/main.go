package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	tea "charm.land/bubbletea/v2"

	"github.com/ETurriza/juego_futbol/internal/aplicacion"
	"github.com/ETurriza/juego_futbol/internal/menus"
)

func main() {
	semilla := flag.Int64("semilla", time.Now().UnixNano(), "semilla de la carrera (misma semilla, misma liga)")
	equipos := flag.Int("equipos", 10, "número de equipos de la liga")
	flag.Parse()

	if err := ejecutar(*semilla, *equipos); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}

func ejecutar(semilla int64, equipos int) error {
	carrera, err := aplicacion.NuevaCarrera(semilla, equipos)
	if err != nil {
		return err
	}
	// Cada "Nueva carrera" usa la semilla siguiente.
	siguiente := semilla
	nueva := func() (*aplicacion.Carrera, error) {
		siguiente++
		return aplicacion.NuevaCarrera(siguiente, equipos)
	}
	_, err = tea.NewProgram(menus.Nuevo(carrera, nueva)).Run()
	return err
}
