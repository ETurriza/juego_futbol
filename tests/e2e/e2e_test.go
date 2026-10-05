// Package e2e prueba el juego de extremo a extremo: un tea.Program real recibe
// bytes de teclado como los de una terminal y se revisa la pantalla final.
package e2e

import (
	"bytes"
	"context"
	"io"
	"strings"
	"sync"
	"testing"
	"time"

	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/x/ansi"

	"github.com/ETurriza/juego_futbol/internal/aplicacion"
	"github.com/ETurriza/juego_futbol/internal/menus"
)

// Bytes que envía una terminal para cada tecla.
const (
	enter = "\r"
	abajo = "\x1b[B"
	// retroceso vuelve atrás igual que esc, sin la ambigüedad de un escape
	// suelto.
	retroceso = "\x7f"
)

const equipos = 10

// salida es un bytes.Buffer seguro para usar desde varias goroutines.
type salida struct {
	mu  sync.Mutex
	buf bytes.Buffer
}

func (s *salida) Write(p []byte) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.buf.Write(p)
}

func (s *salida) Len() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.buf.Len()
}

// jugar arranca el programa con la semilla dada, envía las teclas en orden y
// devuelve el modelo final cuando el programa termina.
func jugar(t *testing.T, semilla int64, teclas ...string) (menus.Modelo, *salida) {
	t.Helper()
	carrera, err := aplicacion.NuevaCarrera(semilla, equipos)
	if err != nil {
		t.Fatal(err)
	}
	siguiente := semilla
	nueva := func() (*aplicacion.Carrera, error) {
		siguiente++
		return aplicacion.NuevaCarrera(siguiente, equipos)
	}

	entrada, escritor := io.Pipe()
	defer escritor.Close()
	out := &salida{}

	ctx, cancelar := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancelar()
	programa := tea.NewProgram(menus.Nuevo(carrera, nueva),
		tea.WithContext(ctx), tea.WithInput(entrada), tea.WithOutput(out),
		tea.WithWindowSize(80, 30), tea.WithoutSignals())

	fin := make(chan error, 1)
	go func() {
		for _, k := range teclas {
			if _, err := escritor.Write([]byte(k)); err != nil {
				fin <- err
				return
			}
		}
		fin <- nil
	}()

	modelo, err := programa.Run()
	if err != nil {
		t.Fatalf("el programa fallo: %v", err)
	}
	if err := <-fin; err != nil {
		t.Fatalf("no se pudieron enviar las teclas: %v", err)
	}
	return modelo.(menus.Modelo), out
}

func pantalla(m menus.Modelo) string { return ansi.Strip(m.View().Content) }

func TestTemporadaCompletaConTeclado(t *testing.T) {
	var teclas []string
	for i := 0; i < 18; i++ {
		teclas = append(teclas, enter, enter) // avanzar jornada y continuar
	}
	teclas = append(teclas, "q") // salir desde el fin de temporada

	modelo, out := jugar(t, 1, teclas...)

	carrera := modelo.Carrera()
	if !carrera.Terminada() || carrera.Jornada() != 18 {
		t.Fatalf("la temporada deberia estar terminada: jornada %d", carrera.Jornada())
	}
	campeon, _ := carrera.Campeon()
	for _, f := range []string{"TEMPORADA TERMINADA", "Campeón: " + campeon, "Tu equipo: " + carrera.NombreEquipo()} {
		if !strings.Contains(pantalla(modelo), f) {
			t.Errorf("la pantalla final deberia contener %q:\n%s", f, pantalla(modelo))
		}
	}
	if out.Len() == 0 {
		t.Error("el programa no escribio nada en la terminal")
	}
}

func TestRecorridoPorLasPantallas(t *testing.T) {
	modelo, _ := jugar(t, 7,
		enter, enter, // avanzar una jornada y volver al menú
		abajo, enter, // tabla de posiciones
		retroceso,    // volver al menú
		abajo, enter, // plantilla
		retroceso, // volver al menú
		"q",
	)
	if modelo.Carrera().Jornada() != 1 {
		t.Errorf("Jornada = %d, se esperaba 1", modelo.Carrera().Jornada())
	}
	for _, f := range []string{"JUEGO DE FÚTBOL", "Jornada 1 / 18", "> Plantilla"} {
		if !strings.Contains(pantalla(modelo), f) {
			t.Errorf("la pantalla deberia contener %q:\n%s", f, pantalla(modelo))
		}
	}
}

func TestSalirConQEnElMenu(t *testing.T) {
	modelo, _ := jugar(t, 1, "q")
	if modelo.Carrera().Jornada() != 0 {
		t.Error("no deberia haberse jugado ninguna jornada")
	}
}

func TestNuevaCarreraAlTerminar(t *testing.T) {
	var teclas []string
	for i := 0; i < 18; i++ {
		teclas = append(teclas, enter, enter)
	}
	teclas = append(teclas, enter) // "Nueva carrera"
	teclas = append(teclas, "q")   // salir desde el menú de la carrera nueva

	modelo, _ := jugar(t, 1, teclas...)
	carrera := modelo.Carrera()
	if carrera.Jornada() != 0 || carrera.Terminada() {
		t.Fatalf("deberia haber una carrera nueva: jornada %d", carrera.Jornada())
	}
	anterior, _ := aplicacion.NuevaCarrera(1, equipos)
	if carrera.NombreEquipo() == anterior.NombreEquipo() && carrera.Semilla == anterior.Semilla {
		t.Error("la carrera nueva deberia usar otra semilla")
	}
	if carrera.Semilla != 2 {
		t.Errorf("Semilla = %d, se esperaba 2", carrera.Semilla)
	}
}
