package menus

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/x/ansi"

	"github.com/ETurriza/juego_futbol/internal/aplicacion"
)

func carreraDePrueba(t *testing.T, semilla int64, equipos int) *aplicacion.Carrera {
	t.Helper()
	c, err := aplicacion.NuevaCarrera(semilla, equipos)
	if err != nil {
		t.Fatal(err)
	}
	return c
}

func modeloDePrueba(t *testing.T, semilla int64, equipos int) Modelo {
	t.Helper()
	siguiente := semilla
	return Nuevo(carreraDePrueba(t, semilla, equipos), func() (*aplicacion.Carrera, error) {
		siguiente++
		return aplicacion.NuevaCarrera(siguiente, equipos)
	})
}

// tecla traduce el nombre de una tecla a su mensaje de Bubble Tea.
func tecla(nombre string) tea.KeyPressMsg {
	switch nombre {
	case "enter":
		return tea.KeyPressMsg{Code: tea.KeyEnter}
	case "esc":
		return tea.KeyPressMsg{Code: tea.KeyEscape}
	case "up":
		return tea.KeyPressMsg{Code: tea.KeyUp}
	case "down":
		return tea.KeyPressMsg{Code: tea.KeyDown}
	case "pgup":
		return tea.KeyPressMsg{Code: tea.KeyPgUp}
	case "pgdown":
		return tea.KeyPressMsg{Code: tea.KeyPgDown}
	case "ctrl+c":
		return tea.KeyPressMsg{Code: 'c', Mod: tea.ModCtrl}
	}
	return tea.KeyPressMsg{Code: rune(nombre[0]), Text: nombre}
}

// pulsar envía teclas al modelo y devuelve el modelo resultante y el último
// comando.
func pulsar(m Modelo, nombres ...string) (Modelo, tea.Cmd) {
	var cmd tea.Cmd
	for _, n := range nombres {
		var siguiente tea.Model
		siguiente, cmd = m.Update(tecla(n))
		m = siguiente.(Modelo)
	}
	return m, cmd
}

func esSalir(cmd tea.Cmd) bool {
	if cmd == nil {
		return false
	}
	_, ok := cmd().(tea.QuitMsg)
	return ok
}

// texto es la pantalla actual sin códigos de color.
func texto(m Modelo) string { return ansi.Strip(m.View().Content) }

func contiene(t *testing.T, m Modelo, fragmentos ...string) {
	t.Helper()
	pantalla := texto(m)
	for _, f := range fragmentos {
		if !strings.Contains(pantalla, f) {
			t.Errorf("la pantalla deberia contener %q:\n%s", f, pantalla)
		}
	}
}

func noContiene(t *testing.T, m Modelo, fragmentos ...string) {
	t.Helper()
	pantalla := texto(m)
	for _, f := range fragmentos {
		if strings.Contains(pantalla, f) {
			t.Errorf("la pantalla no deberia contener %q:\n%s", f, pantalla)
		}
	}
}

func redimensionar(m Modelo, ancho, alto int) Modelo {
	siguiente, _ := m.Update(tea.WindowSizeMsg{Width: ancho, Height: alto})
	return siguiente.(Modelo)
}

func TestMenuPrincipal(t *testing.T) {
	m := modeloDePrueba(t, 1, 10)
	contiene(t, m,
		"JUEGO DE FÚTBOL", m.carrera.NombreEquipo(), "Jornada 0 / 18", "Valoración del equipo",
		"> Avanzar jornada", "  Tabla de posiciones", "  Plantilla", "  Salir", "q salir")
	if !m.View().AltScreen {
		t.Error("la vista deberia usar la pantalla alterna")
	}
}

func TestNavegacionDelMenu(t *testing.T) {
	m := modeloDePrueba(t, 1, 10)
	m, _ = pulsar(m, "up") // no pasa del principio
	contiene(t, m, "> Avanzar jornada")
	m, _ = pulsar(m, "down")
	contiene(t, m, "> Tabla de posiciones")
	m, _ = pulsar(m, "j", "j", "j", "j") // no pasa del final
	contiene(t, m, "> Salir")
	m, _ = pulsar(m, "k")
	contiene(t, m, "> Plantilla")
}

func TestSalirDelMenu(t *testing.T) {
	m := modeloDePrueba(t, 1, 10)
	if _, cmd := pulsar(m, "q"); !esSalir(cmd) {
		t.Error("q deberia salir desde el menu")
	}
	if _, cmd := pulsar(m, "down", "down", "down", "enter"); !esSalir(cmd) {
		t.Error("la opcion Salir deberia salir")
	}
	if _, cmd := pulsar(m, "esc"); esSalir(cmd) {
		t.Error("esc no deberia salir desde el menu")
	}
}

func TestCtrlCSaleDesdeCualquierPantalla(t *testing.T) {
	m := modeloDePrueba(t, 1, 10)
	for _, secuencia := range [][]string{
		{}, {"enter"}, {"down", "enter"}, {"down", "down", "enter"},
	} {
		pantalla, _ := pulsar(m, secuencia...)
		if _, cmd := pulsar(pantalla, "ctrl+c"); !esSalir(cmd) {
			t.Errorf("ctrl+c deberia salir tras %v", secuencia)
		}
	}
}

func TestAvanzarJornadaMuestraResultados(t *testing.T) {
	m := modeloDePrueba(t, 1, 10)
	m, _ = pulsar(m, "enter")
	contiene(t, m, "JORNADA 1 / 18", "enter continuar")
	if m.carrera.Jornada() != 1 {
		t.Fatalf("Jornada = %d, se esperaba 1", m.carrera.Jornada())
	}
	partidos := m.carrera.UltimaJornada()
	if len(partidos) != 5 {
		t.Fatalf("%d partidos, se esperaban 5", len(partidos))
	}
	for _, p := range partidos {
		contiene(t, m, p.Local, p.Visitante, fmt.Sprintf("%2d - %-2d", p.GolesLocal, p.GolesVisitante))
		if p.EsDelUsuario {
			contiene(t, m, "> "+p.Local)
		}
	}
	noContiene(t, m, "descansa")

	m, _ = pulsar(m, "enter")
	contiene(t, m, "Jornada 1 / 18", "> Avanzar jornada")
}

func TestTablaDePosiciones(t *testing.T) {
	m := modeloDePrueba(t, 1, 10)
	m, _ = pulsar(m, "enter", "enter") // jugar una jornada y volver al menu
	m, _ = pulsar(m, "down", "enter")
	contiene(t, m, "TABLA · Jornada 1 / 18", "Equipo", "Pts", "esc volver")
	// Solo la fila del equipo del usuario va marcada con ">".
	marcadas := 0
	for _, linea := range strings.Split(texto(m), "\n") {
		if strings.HasPrefix(linea, ">") {
			marcadas++
			if !strings.Contains(linea, m.carrera.NombreEquipo()) {
				t.Errorf("fila marcada que no es del usuario: %q", linea)
			}
		}
	}
	if marcadas != 1 {
		t.Errorf("%d filas marcadas, se esperaba 1", marcadas)
	}
	for _, f := range m.carrera.Tabla() {
		contiene(t, m, f.Equipo)
	}
	m, _ = pulsar(m, "esc")
	contiene(t, m, "JUEGO DE FÚTBOL", "> Tabla de posiciones")
}

func TestPlantilla(t *testing.T) {
	m := modeloDePrueba(t, 1, 10)
	m, _ = pulsar(m, "down", "down", "enter")
	contiene(t, m, "PLANTILLA · "+m.carrera.NombreEquipo(), "Posición", "Portero", "Delantero", "REF", "esc volver")
	for _, j := range m.carrera.Plantilla() {
		contiene(t, m, j.Nombre)
	}
	m, _ = pulsar(m, "esc")
	contiene(t, m, "> Plantilla")
}

func TestPlantillaConScroll(t *testing.T) {
	m := modeloDePrueba(t, 1, 10)
	plantilla := m.carrera.Plantilla()
	primero, ultimo := plantilla[0].Nombre, plantilla[len(plantilla)-1].Nombre

	m = redimensionar(m, 80, 12) // caben 7 filas de 22
	m, _ = pulsar(m, "down", "down", "enter")
	contiene(t, m, primero, "1-7 de 22")
	noContiene(t, m, ultimo)

	m, _ = pulsar(m, "up") // no pasa del principio
	contiene(t, m, "1-7 de 22")

	m, _ = pulsar(m, "down", "down", "down")
	contiene(t, m, "4-10 de 22")
	noContiene(t, m, primero)

	m, _ = pulsar(m, "pgdown", "pgdown", "pgdown", "pgdown")
	contiene(t, m, ultimo, "16-22 de 22") // no pasa del final
	m, _ = pulsar(m, "pgup")
	contiene(t, m, "9-15 de 22")
}

func TestRedimensionarRecortaElScroll(t *testing.T) {
	m := modeloDePrueba(t, 1, 10)
	m = redimensionar(m, 80, 12)
	m, _ = pulsar(m, "down", "down", "enter", "pgdown", "pgdown", "pgdown", "pgdown")
	m = redimensionar(m, 80, 40) // ahora caben todas
	contiene(t, m, "esc volver")
	noContiene(t, m, " de 22")
}

func TestVentanaMuyChicaNoFalla(t *testing.T) {
	m := modeloDePrueba(t, 1, 10)
	for _, tam := range [][2]int{{0, 0}, {1, 1}, {20, 3}} {
		m = redimensionar(m, tam[0], tam[1])
		for _, secuencia := range [][]string{{}, {"down", "enter"}, {"down", "down", "enter"}} {
			pantalla, _ := pulsar(m, secuencia...)
			if texto(pantalla) == "" {
				t.Errorf("tamaño %v tras %v: pantalla vacia", tam, secuencia)
			}
			pulsar(pantalla, "down", "pgdown", "end", "up", "pgup", "home", "esc")
		}
	}
}

func TestLigaImparMuestraDescanso(t *testing.T) {
	m := modeloDePrueba(t, 3, 5)
	vistoDescanso := false
	for i := 0; i < m.carrera.TotalJornadas(); i++ {
		m, _ = pulsar(m, "enter")
		juega := false
		for _, p := range m.carrera.UltimaJornada() {
			juega = juega || p.EsDelUsuario
		}
		descansa := strings.Contains(texto(m), "Tu equipo descansa esta jornada.")
		if descansa == juega {
			t.Fatalf("jornada %d: el usuario juega=%v pero el mensaje de descanso=%v", i+1, juega, descansa)
		}
		vistoDescanso = vistoDescanso || descansa
		m, _ = pulsar(m, "enter")
	}
	if !vistoDescanso {
		t.Error("en una liga de 5 equipos el usuario deberia descansar alguna jornada")
	}
}

// jugarTemporada avanza jornada a jornada hasta que aparece la pantalla de fin.
func jugarTemporada(t *testing.T, m Modelo) Modelo {
	t.Helper()
	for i := 0; i < m.carrera.TotalJornadas(); i++ {
		m, _ = pulsar(m, "enter", "enter")
	}
	return m
}

func TestFinDeTemporada(t *testing.T) {
	m := modeloDePrueba(t, 1, 10)
	m = jugarTemporada(t, m)
	if !m.carrera.Terminada() {
		t.Fatal("la temporada deberia estar terminada")
	}
	campeon, _ := m.carrera.Campeon()
	contiene(t, m, "TEMPORADA TERMINADA", "Campeón: "+campeon, "Tu equipo: "+m.carrera.NombreEquipo(),
		"de 10", "puntos", "> Nueva carrera", "  Ver tabla final", "  Salir")
	// El menú principal ya no está disponible: no se puede avanzar más.
	noContiene(t, m, "Avanzar jornada")
}

func TestUltimaJornadaLlevaAlFin(t *testing.T) {
	m := modeloDePrueba(t, 1, 10)
	for i := 0; i < m.carrera.TotalJornadas()-1; i++ {
		m, _ = pulsar(m, "enter", "enter")
	}
	m, _ = pulsar(m, "enter") // última jornada
	contiene(t, m, "JORNADA 18 / 18")
	m, _ = pulsar(m, "enter")
	contiene(t, m, "TEMPORADA TERMINADA")
}

func TestFelicitaAlCampeonDelUsuario(t *testing.T) {
	// Busca una semilla en la que el usuario salga campeón y otra en que no.
	var campeon, otro bool
	for semilla := int64(1); semilla <= 200 && !(campeon && otro); semilla++ {
		m := jugarTemporada(t, modeloDePrueba(t, semilla, 4))
		nombre, _ := m.carrera.Campeon()
		felicita := strings.Contains(texto(m), "¡Felicidades, eres el campeón!")
		if esUsuario := nombre == m.carrera.NombreEquipo(); felicita != esUsuario {
			t.Fatalf("semilla %d: felicitacion=%v, usuario campeon=%v", semilla, felicita, esUsuario)
		} else if esUsuario {
			campeon = true
		} else {
			otro = true
		}
	}
	if !campeon || !otro {
		t.Fatalf("no se probaron ambos casos: campeon=%v otro=%v", campeon, otro)
	}
}

func TestVerTablaFinalYVolverAlFin(t *testing.T) {
	m := jugarTemporada(t, modeloDePrueba(t, 1, 10))
	m, _ = pulsar(m, "down", "enter")
	contiene(t, m, "TABLA · Jornada 18 / 18")
	m, _ = pulsar(m, "esc")
	contiene(t, m, "TEMPORADA TERMINADA", "> Ver tabla final")
}

func TestSalirDesdeElFin(t *testing.T) {
	m := jugarTemporada(t, modeloDePrueba(t, 1, 10))
	if _, cmd := pulsar(m, "q"); !esSalir(cmd) {
		t.Error("q deberia salir desde el fin de temporada")
	}
	if _, cmd := pulsar(m, "down", "down", "enter"); !esSalir(cmd) {
		t.Error("la opcion Salir deberia salir")
	}
}

func TestNuevaCarreraDesdeElFin(t *testing.T) {
	m := jugarTemporada(t, modeloDePrueba(t, 1, 10))
	anterior := m.carrera
	m, _ = pulsar(m, "enter")
	if m.carrera == anterior || m.carrera.Jornada() != 0 {
		t.Fatal("deberia haber una carrera nueva sin jornadas jugadas")
	}
	contiene(t, m, "JUEGO DE FÚTBOL", "Jornada 0 / 18", "> Avanzar jornada")
}

func TestNuevaCarreraConErrorMuestraAviso(t *testing.T) {
	m := jugarTemporada(t, modeloDePrueba(t, 1, 10))
	m.nueva = func() (*aplicacion.Carrera, error) { return nil, errors.New("sin espacio") }
	anterior := m.carrera
	m, _ = pulsar(m, "enter")
	if m.carrera != anterior {
		t.Error("la carrera no deberia cambiar si falla")
	}
	contiene(t, m, "TEMPORADA TERMINADA", "sin espacio")
}

func TestIniciaEnElFinSiLaTemporadaYaTermino(t *testing.T) {
	c := carreraDePrueba(t, 1, 4)
	for !c.Terminada() {
		if _, err := c.AvanzarJornada(); err != nil {
			t.Fatal(err)
		}
	}
	m := Nuevo(c, nil)
	contiene(t, m, "TEMPORADA TERMINADA")
}

func TestMensajesDesconocidosSeIgnoran(t *testing.T) {
	m := modeloDePrueba(t, 1, 10)
	antes := texto(m)
	siguiente, cmd := m.Update(struct{}{})
	if cmd != nil || texto(siguiente.(Modelo)) != antes {
		t.Error("un mensaje desconocido no deberia cambiar nada")
	}
	// En la pantalla de jornada solo enter continua.
	m, _ = pulsar(m, "enter")
	enJornada := texto(m)
	m, _ = pulsar(m, "down", "esc", "q")
	if texto(m) != enJornada {
		t.Error("solo enter deberia salir de la pantalla de jornada")
	}
}
