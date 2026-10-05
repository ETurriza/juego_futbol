package menus

import (
	tea "charm.land/bubbletea/v2"

	"github.com/ETurriza/juego_futbol/internal/aplicacion"
)

type pantalla int

const (
	pantallaMenu pantalla = iota
	pantallaPlantilla
	pantallaTabla
	pantallaJornada
	pantallaFin
	pantallaInicioTemporada
	pantallaHistorial
	pantallaConfirmar
	pantallaEstadisticas
	pantallaClasificacion
	pantallaEquipos
	pantallaPlantillaEquipo
	pantallaFicha
)

// Opciones del menú principal, en orden.
const (
	opAvanzar = iota
	opTabla
	opPlantilla
	opEstadisticas
	opHistorial
	opSalir
)

// Opciones del menú de fin de temporada, en orden.
const (
	opFinSiguiente = iota
	opFinTabla
	opFinEstadisticas
	opFinHistorial
	opFinNueva
	opFinSalir
)

// Opciones de cada menú, en el orden de las constantes de arriba.
var (
	opcionesMenu = []string{"Avanzar jornada", "Tabla de posiciones", "Plantilla", "Estadísticas", "Historial", "Salir"}
	opcionesFin  = []string{"Siguiente temporada", "Ver tabla final", "Estadísticas", "Historial", "Nueva carrera", "Salir"}
	// Confirmación de "Nueva carrera": la opción segura va primero.
	opcionesConfirmar = []string{"No, volver", "Sí, empezar de cero"}
)

// NuevaCarreraFunc crea una carrera nueva; el menú de fin de temporada la usa
// para la opción "Nueva carrera".
type NuevaCarreraFunc func() (*aplicacion.Carrera, error)

// Modelo es el modelo raíz de Bubble Tea del modo carrera.
type Modelo struct {
	carrera  *aplicacion.Carrera
	nueva    NuevaCarreraFunc
	pantalla pantalla
	// origen es la pantalla a la que vuelven la plantilla y la tabla con esc.
	origen      pantalla
	cursor      int
	scroll      int
	ancho, alto int
	aviso       string
	// cambios y valoracionAntes alimentan la pantalla de inicio de temporada.
	cambios         aplicacion.CambiosTemporada
	valoracionAntes int

	// Estado de las pantallas de estadísticas (ver estadisticas.go).
	pila        []destino // pantallas por las que se bajó, para volver con esc
	criterio    aplicacion.Criterio
	equipoSel   string // equipo de la plantilla que se está viendo
	fichaID     int    // jugador de la ficha
	verAsistent bool   // en Equipos, mostrar asistente en vez de goleador
	verStats    bool   // en Plantilla, mostrar estadísticas en vez de atributos
}

// Nuevo crea el modelo para una carrera ya iniciada. nueva se invoca al elegir
// "Nueva carrera" al terminar la temporada.
func Nuevo(carrera *aplicacion.Carrera, nueva NuevaCarreraFunc) Modelo {
	m := Modelo{carrera: carrera, nueva: nueva}
	if carrera.Terminada() {
		m.pantalla = pantallaFin
	}
	return m
}

// Carrera devuelve la carrera en curso.
func (m Modelo) Carrera() *aplicacion.Carrera { return m.carrera }

// Init implementa tea.Model.
func (m Modelo) Init() tea.Cmd { return nil }

// Update implementa tea.Model.
func (m Modelo) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.ancho, m.alto = msg.Width, msg.Height
		m = m.ajustarVentana()
		return m, nil
	case tea.KeyPressMsg:
		return m.tecla(msg.String())
	}
	return m, nil
}

func (m Modelo) tecla(k string) (tea.Model, tea.Cmd) {
	if k == "ctrl+c" {
		return m, tea.Quit
	}
	switch m.pantalla {
	case pantallaMenu, pantallaFin, pantallaConfirmar:
		return m.teclaMenu(k)
	case pantallaPlantilla, pantallaTabla, pantallaHistorial:
		return m.teclaLista(k), nil
	case pantallaJornada:
		if k == "enter" {
			m.cursor, m.aviso = 0, ""
			if m.carrera.Terminada() {
				m.pantalla = pantallaFin
			} else {
				m.pantalla = pantallaMenu
			}
		}
	case pantallaInicioTemporada:
		if k == "enter" {
			m.cursor, m.aviso = 0, ""
			m.pantalla = pantallaMenu
		}
	case pantallaEstadisticas:
		return m.teclaEstadisticas(k), nil
	case pantallaClasificacion, pantallaEquipos, pantallaPlantillaEquipo:
		return m.teclaSeleccion(k), nil
	case pantallaFicha:
		return m.teclaFicha(k), nil
	}
	return m, nil
}

func (m Modelo) opciones() []string {
	switch m.pantalla {
	case pantallaFin:
		return opcionesFin
	case pantallaConfirmar:
		return opcionesConfirmar
	}
	return opcionesMenu
}

func (m Modelo) teclaMenu(k string) (tea.Model, tea.Cmd) {
	n := len(m.opciones())
	switch k {
	case "up", "k":
		m.cursor = max(m.cursor-1, 0)
	case "down", "j":
		m.cursor = min(m.cursor+1, n-1)
	case "q":
		// En la confirmación una tecla suelta no debe cerrar el juego.
		if m.pantalla != pantallaConfirmar {
			return m, tea.Quit
		}
	case "esc":
		if m.pantalla == pantallaConfirmar {
			return m.cancelarNueva(), nil
		}
	case "enter":
		switch m.pantalla {
		case pantallaFin:
			return m.elegirFin()
		case pantallaConfirmar:
			return m.elegirConfirmar()
		}
		return m.elegirMenu()
	}
	return m, nil
}

func (m Modelo) elegirMenu() (tea.Model, tea.Cmd) {
	m.aviso = ""
	switch m.cursor {
	case opAvanzar:
		if _, err := m.carrera.AvanzarJornada(); err != nil {
			m.aviso = err.Error()
			return m, nil
		}
		m.pantalla = pantallaJornada
	case opTabla:
		m.abrirLista(pantallaTabla)
	case opPlantilla:
		m.abrirLista(pantallaPlantilla)
	case opEstadisticas:
		return m.irA(pantallaEstadisticas), nil
	case opHistorial:
		m.abrirLista(pantallaHistorial)
	case opSalir:
		return m, tea.Quit
	}
	return m, nil
}

func (m Modelo) elegirFin() (tea.Model, tea.Cmd) {
	m.aviso = ""
	switch m.cursor {
	case opFinSiguiente:
		antes := m.carrera.ValoracionEquipo()
		cambios, err := m.carrera.SiguienteTemporada()
		if err != nil {
			m.aviso = err.Error()
			return m, nil
		}
		m.cambios, m.valoracionAntes = cambios, antes
		m.pantalla, m.cursor = pantallaInicioTemporada, 0
	case opFinTabla:
		m.abrirLista(pantallaTabla)
	case opFinEstadisticas:
		return m.irA(pantallaEstadisticas), nil
	case opFinHistorial:
		m.abrirLista(pantallaHistorial)
	case opFinNueva: // pide confirmación
		m.pantalla, m.cursor = pantallaConfirmar, 0
	case opFinSalir:
		return m, tea.Quit
	}
	return m, nil
}

// cancelarNueva vuelve de la confirmación al fin de temporada, sobre la opción
// "Nueva carrera".
func (m Modelo) cancelarNueva() Modelo {
	m.pantalla, m.cursor = pantallaFin, opFinNueva
	return m
}

func (m Modelo) elegirConfirmar() (tea.Model, tea.Cmd) {
	if m.cursor == 0 {
		return m.cancelarNueva(), nil
	}
	carrera, err := m.nueva()
	if err != nil {
		m.aviso = err.Error()
		return m.cancelarNueva(), nil
	}
	m.carrera, m.pantalla, m.cursor = carrera, pantallaMenu, 0
	return m, nil
}

// abrirLista pasa a la plantilla o la tabla recordando de dónde viene.
func (m *Modelo) abrirLista(p pantalla) {
	m.origen, m.pantalla, m.scroll = m.pantalla, p, 0
}

func (m Modelo) teclaLista(k string) Modelo {
	visibles := m.filasVisibles()
	switch k {
	case "up", "k":
		m.scroll--
	case "down", "j":
		m.scroll++
	case "pgup":
		m.scroll -= visibles
	case "pgdown":
		m.scroll += visibles
	case "home":
		m.scroll = 0
	case "end":
		m.scroll = m.totalFilas()
	case "tab":
		if m.pantalla == pantallaPlantilla {
			m.verStats = !m.verStats
		}
	case "esc", "backspace":
		m.pantalla = m.origen
		return m
	}
	m.scroll = m.limitarScroll(m.scroll)
	return m
}

// Líneas de la pantalla que no son filas de la lista: título, línea en blanco,
// encabezado, línea en blanco y ayuda.
const lineasFijas = 5

func (m Modelo) totalFilas() int {
	switch m.pantalla {
	case pantallaPlantilla:
		return len(m.carrera.Plantilla())
	case pantallaHistorial:
		return len(m.carrera.Historial)
	}
	return len(m.carrera.Tabla())
}

// filasVisibles es cuántas filas de la lista caben en la terminal; sin tamaño
// conocido se muestran todas.
func (m Modelo) filasVisibles() int {
	total := m.totalFilas()
	if m.alto <= 0 {
		return total
	}
	return min(total, max(m.alto-lineasFijas, 3))
}

func (m Modelo) limitarScroll(s int) int {
	return min(max(s, 0), max(m.totalFilas()-m.filasVisibles(), 0))
}
