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
	contiene(t, m, "> Historial")
}

func TestSalirDelMenu(t *testing.T) {
	m := modeloDePrueba(t, 1, 10)
	if _, cmd := pulsar(m, "q"); !esSalir(cmd) {
		t.Error("q deberia salir desde el menu")
	}
	if _, cmd := pulsar(m, "down", "down", "down", "down", "down", "enter"); !esSalir(cmd) {
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
	contiene(t, m, "TABLA · Temporada 1 · Jornada 1 / 18", "Equipo", "Pts", "esc volver")
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
	contiene(t, m, "TEMPORADA 1 TERMINADA", "Campeón: "+campeon, "Tu equipo: "+m.carrera.NombreEquipo(),
		"de 10", "puntos", "> Siguiente temporada", "  Ver tabla final", "  Historial", "  Nueva carrera", "  Salir")
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
	contiene(t, m, "TEMPORADA 1 TERMINADA")
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
	contiene(t, m, "TABLA · Temporada 1 · Jornada 18 / 18")
	m, _ = pulsar(m, "esc")
	contiene(t, m, "TEMPORADA 1 TERMINADA", "> Ver tabla final")
}

func TestSalirDesdeElFin(t *testing.T) {
	m := jugarTemporada(t, modeloDePrueba(t, 1, 10))
	if _, cmd := pulsar(m, "q"); !esSalir(cmd) {
		t.Error("q deberia salir desde el fin de temporada")
	}
	if _, cmd := pulsar(m, "down", "down", "down", "down", "down", "enter"); !esSalir(cmd) {
		t.Error("la opcion Salir deberia salir")
	}
}

func TestNuevaCarreraPideConfirmacion(t *testing.T) {
	m := jugarTemporada(t, modeloDePrueba(t, 1, 10))
	anterior := m.carrera
	m, _ = pulsar(m, "down", "down", "down", "down", "enter")
	contiene(t, m, "NUEVA CARRERA", "pierdes la actual", "Temporada 1 con "+anterior.NombreEquipo(),
		"> No, volver", "  Sí, empezar de cero")
	if m.carrera != anterior {
		t.Fatal("pedir la confirmacion no deberia cambiar la carrera")
	}

	// q no cierra el juego en la confirmacion; esc y "No" vuelven al fin.
	if _, cmd := pulsar(m, "q"); esSalir(cmd) {
		t.Error("q no deberia salir desde la confirmacion")
	}
	vuelta, _ := pulsar(m, "esc")
	contiene(t, vuelta, "TEMPORADA 1 TERMINADA", "> Nueva carrera")
	vuelta, _ = pulsar(m, "enter")
	contiene(t, vuelta, "TEMPORADA 1 TERMINADA", "> Nueva carrera")
	if vuelta.carrera != anterior {
		t.Error("cancelar no deberia cambiar la carrera")
	}

	// ctrl+c sí cierra desde cualquier pantalla.
	if _, cmd := pulsar(m, "ctrl+c"); !esSalir(cmd) {
		t.Error("ctrl+c deberia salir")
	}
}

func TestNuevaCarreraDesdeElFin(t *testing.T) {
	m := jugarTemporada(t, modeloDePrueba(t, 1, 10))
	anterior := m.carrera
	m, _ = pulsar(m, "down", "down", "down", "down", "enter", "down", "enter")
	if m.carrera == anterior || m.carrera.Jornada() != 0 || m.carrera.Numero != 1 {
		t.Fatal("deberia haber una carrera nueva sin jornadas jugadas")
	}
	contiene(t, m, "JUEGO DE FÚTBOL", "Temporada 1 · Jornada 0 / 18", "> Avanzar jornada")
}

func TestNuevaCarreraConErrorMuestraAviso(t *testing.T) {
	m := jugarTemporada(t, modeloDePrueba(t, 1, 10))
	m.nueva = func() (*aplicacion.Carrera, error) { return nil, errors.New("sin espacio") }
	anterior := m.carrera
	m, _ = pulsar(m, "down", "down", "down", "down", "enter", "down", "enter")
	if m.carrera != anterior {
		t.Error("la carrera no deberia cambiar si falla")
	}
	contiene(t, m, "TEMPORADA 1 TERMINADA", "sin espacio")
}

func TestIniciaEnElFinSiLaTemporadaYaTermino(t *testing.T) {
	c := carreraDePrueba(t, 1, 4)
	for !c.Terminada() {
		if _, err := c.AvanzarJornada(); err != nil {
			t.Fatal(err)
		}
	}
	m := Nuevo(c, nil)
	contiene(t, m, "TEMPORADA 1 TERMINADA")
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

// avanzarCarrera termina n temporadas directamente sobre la carrera, sin pasar
// por la interfaz, y la deja al inicio de la temporada n+1.
func avanzarCarrera(t *testing.T, c *aplicacion.Carrera, n int) {
	t.Helper()
	for i := 0; i < n; i++ {
		for !c.Terminada() {
			if _, err := c.AvanzarJornada(); err != nil {
				t.Fatal(err)
			}
		}
		if _, err := c.SiguienteTemporada(); err != nil {
			t.Fatal(err)
		}
	}
}

func TestSiguienteTemporadaMuestraElInicio(t *testing.T) {
	m := jugarTemporada(t, modeloDePrueba(t, 1, 10))
	antes := m.carrera.ValoracionEquipo()

	m, _ = pulsar(m, "enter") // Siguiente temporada
	if m.carrera.Numero != 2 || m.carrera.Terminada() {
		t.Fatalf("deberia haber empezado la temporada 2: numero %d", m.carrera.Numero)
	}
	contiene(t, m, "TEMPORADA 2 · "+m.carrera.NombreEquipo(),
		fmt.Sprintf("Valoración del equipo: %d → %d", antes, m.carrera.ValoracionEquipo()),
		fmt.Sprintf("En toda la liga: %d retiros y %d juveniles.", m.cambios.RetiradosLiga, m.cambios.JuvenilesLiga),
		"enter continuar")
	// Lista exactamente a los retirados y a los juveniles de su club.
	for _, j := range m.cambios.Retirados {
		contiene(t, m, j.Nombre)
	}
	for _, j := range m.cambios.Juveniles {
		contiene(t, m, j.Nombre, "años")
	}
	if len(m.cambios.Retirados) > 0 {
		contiene(t, m, fmt.Sprintf("Se retiran (%d)", len(m.cambios.Retirados)))
	}
	if len(m.cambios.Juveniles) > 0 {
		contiene(t, m, fmt.Sprintf("Llegan de la cantera (%d)", len(m.cambios.Juveniles)))
	}

	m, _ = pulsar(m, "enter")
	contiene(t, m, "JUEGO DE FÚTBOL", "Temporada 2 · Jornada 0 / 18", "> Avanzar jornada")
}

func TestInicioDeTemporadaSinBajasNiAltas(t *testing.T) {
	m := modeloDePrueba(t, 1, 10)
	m.pantalla = pantallaInicioTemporada
	contiene(t, m, "Nadie se retira este año.", "No llega nadie de la cantera.", "En toda la liga: 0 retiros")
	noContiene(t, m, "Se retiran", "Llegan de la cantera")
}

func TestSiguienteTemporadaSoloEnElFin(t *testing.T) {
	m := modeloDePrueba(t, 1, 10)
	// Forzar la pantalla de fin con una temporada sin terminar: el error se
	// muestra como aviso y el juego no se cierra ni cambia de pantalla.
	m.pantalla = pantallaFin
	m, cmd := pulsar(m, "enter")
	if esSalir(cmd) || m.carrera.Numero != 1 {
		t.Fatal("no deberia cerrarse ni cambiar de temporada")
	}
	contiene(t, m, "TEMPORADA 1 TERMINADA", "todavia no termino")
}

func TestHistorialSinTemporadas(t *testing.T) {
	m := modeloDePrueba(t, 1, 10)
	m, _ = pulsar(m, "down", "down", "down", "down", "enter")
	contiene(t, m, "HISTORIAL · "+m.carrera.NombreEquipo(), "Todavía no has terminado ninguna temporada.", "esc volver")
	noContiene(t, m, "Campeón")
	m, _ = pulsar(m, "esc")
	contiene(t, m, "JUEGO DE FÚTBOL", "> Historial")
}

func TestHistorialConTemporadas(t *testing.T) {
	c := carreraDePrueba(t, 4, 10)
	avanzarCarrera(t, c, 3)
	m := Nuevo(c, nil)
	m, _ = pulsar(m, "down", "down", "down", "down", "enter")
	contiene(t, m, "HISTORIAL", "Temp", "Campeón", "Tu puesto", "Pts", "esc volver")
	for _, h := range c.Historial {
		contiene(t, m, h.Campeon, fmt.Sprintf("%d de 10", h.PuestoUsuario))
	}

	// Solo se marca (con ">") la fila de una temporada en la que el usuario fue
	// campeón.
	c.Historial[1].PuestoUsuario = 1
	marcadas := 0
	for _, linea := range strings.Split(texto(m), "\n") {
		if strings.HasPrefix(linea, ">") {
			marcadas++
			if !strings.Contains(linea, c.Historial[1].Campeon) {
				t.Errorf("fila marcada que no es la del campeonato: %q", linea)
			}
		}
	}
	if marcadas != 1 {
		t.Errorf("%d filas marcadas, se esperaba 1", marcadas)
	}
}

func TestHistorialConScroll(t *testing.T) {
	c := carreraDePrueba(t, 4, 10)
	avanzarCarrera(t, c, 12)
	m := redimensionar(Nuevo(c, nil), 80, 10) // caben 5 filas de 12
	m, _ = pulsar(m, "down", "down", "down", "down", "enter")
	contiene(t, m, "1-5 de 12")
	noContiene(t, m, "  12  ")
	m, _ = pulsar(m, "pgdown", "pgdown", "pgdown")
	contiene(t, m, "8-12 de 12", c.Historial[11].Campeon)
	m, _ = pulsar(m, "esc")
	contiene(t, m, "JUEGO DE FÚTBOL")
}

func TestHistorialDesdeElFinVuelveAlFin(t *testing.T) {
	m := jugarTemporada(t, modeloDePrueba(t, 1, 10))
	m, _ = pulsar(m, "down", "down", "down", "enter")
	contiene(t, m, "HISTORIAL", "Todavía no has terminado")
	m, _ = pulsar(m, "esc")
	contiene(t, m, "TEMPORADA 1 TERMINADA", "> Historial")
}

func TestLaTemporadaApareceEnTodasLasPantallas(t *testing.T) {
	c := carreraDePrueba(t, 6, 10)
	avanzarCarrera(t, c, 1)
	m := Nuevo(c, nil)
	contiene(t, m, "Temporada 2 · Jornada 0 / 18")

	m, _ = pulsar(m, "down", "enter") // tabla
	contiene(t, m, "TABLA · Temporada 2 · Jornada 0 / 18")
	m, _ = pulsar(m, "esc", "up", "enter") // avanzar jornada
	contiene(t, m, "JORNADA 1 / 18 · Temporada 2")
}

func TestSegundaTemporadaSePuedeJugarCompleta(t *testing.T) {
	m := jugarTemporada(t, modeloDePrueba(t, 1, 10))
	m, _ = pulsar(m, "enter", "enter") // siguiente temporada y continuar
	m = jugarTemporada(t, m)
	contiene(t, m, "TEMPORADA 2 TERMINADA", "> Siguiente temporada")
	if len(m.carrera.Historial) != 1 || m.carrera.Numero != 2 {
		t.Errorf("historial %d, temporada %d", len(m.carrera.Historial), m.carrera.Numero)
	}
	m, _ = pulsar(m, "down", "down", "down", "enter") // historial
	contiene(t, m, "HISTORIAL", m.carrera.Historial[0].Campeon)
}
