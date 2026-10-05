package menus

import (
	"fmt"
	"sort"
	"strings"

	"github.com/ETurriza/juego_futbol/internal/aplicacion"
	"github.com/ETurriza/juego_futbol/internal/modelo"
)

// maxFilasClasificacion es cuántos jugadores muestra cada clasificación.
const maxFilasClasificacion = 50

// Opciones del menú de estadísticas, en orden: las seis clasificaciones, la
// comparativa de equipos y volver.
const (
	opEstEquipos = 6
	opEstVolver  = 7
)

var opcionesEstadisticas = []string{
	"Goleadores", "Asistentes", "Tarjetas", "Porteros: porterías imbatidas",
	"Porteros: menos goles encajados", "Mejores valoraciones", "Equipos", "Volver",
}

// clasificaciones: las seis primeras opciones del menú de estadísticas.
var clasificaciones = []struct {
	criterio aplicacion.Criterio
	titulo   string
}{
	{aplicacion.PorGoles, "GOLEADORES"},
	{aplicacion.PorAsistencias, "ASISTENTES"},
	{aplicacion.PorTarjetas, "TARJETAS"},
	{aplicacion.PorImbatidas, "PORTEROS · PORTERÍAS IMBATIDAS"},
	{aplicacion.PorPocosGolesEncajados, "PORTEROS · MENOS GOLES ENCAJADOS"},
	{aplicacion.PorValoracion, "MEJORES VALORACIONES"},
}

// destino es una pantalla de estadísticas con su posición, para poder volver a
// ella con esc tal como estaba.
type destino struct {
	pantalla    pantalla
	cursor      int
	scroll      int
	criterio    aplicacion.Criterio
	equipoSel   string
	fichaID     int
	verAsistent bool
}

func (m Modelo) actual() destino {
	return destino{m.pantalla, m.cursor, m.scroll, m.criterio, m.equipoSel, m.fichaID, m.verAsistent}
}

// irA baja a otra pantalla guardando la actual en la pila.
func (m Modelo) irA(p pantalla) Modelo {
	m.pila = append(append([]destino(nil), m.pila...), m.actual())
	m.pantalla, m.cursor, m.scroll = p, 0, 0
	return m
}

// volver regresa a la pantalla anterior de la pila, tal como estaba.
func (m Modelo) volver() Modelo {
	if len(m.pila) == 0 {
		m.pantalla, m.cursor, m.scroll = pantallaMenu, 0, 0
		return m
	}
	d := m.pila[len(m.pila)-1]
	m.pila = m.pila[:len(m.pila)-1]
	m.pantalla, m.cursor, m.scroll = d.pantalla, d.cursor, d.scroll
	m.criterio, m.equipoSel, m.fichaID, m.verAsistent = d.criterio, d.equipoSel, d.fichaID, d.verAsistent
	return m
}

func (m Modelo) teclaEstadisticas(k string) Modelo {
	switch k {
	case "up", "k":
		m.cursor = max(m.cursor-1, 0)
	case "down", "j":
		m.cursor = min(m.cursor+1, len(opcionesEstadisticas)-1)
	case "esc", "backspace":
		return m.volver()
	case "enter":
		switch {
		case m.cursor < len(clasificaciones):
			m.criterio = clasificaciones[m.cursor].criterio
			return m.irA(pantallaClasificacion)
		case m.cursor == opEstEquipos:
			return m.irA(pantallaEquipos)
		default:
			return m.volver()
		}
	}
	return m
}

// filaEquipo es una línea de la comparativa de equipos: su puesto en la tabla y
// sus números.
type filaEquipo struct {
	Posicion int
	aplicacion.EstadisticaEquipo
	EsDelUsuario bool
}

// filasEquipos devuelve los equipos con sus estadísticas, en el orden de la
// tabla de posiciones.
func (m Modelo) filasEquipos() ([]filaEquipo, error) {
	equipos, err := m.carrera.EstadisticasEquipos()
	if err != nil {
		return nil, err
	}
	puesto := map[string]aplicacion.FilaTabla{}
	for _, f := range m.carrera.Tabla() {
		puesto[f.Equipo] = f
	}
	filas := make([]filaEquipo, len(equipos))
	for i, e := range equipos {
		f := puesto[e.Equipo]
		filas[i] = filaEquipo{Posicion: f.Posicion, EstadisticaEquipo: e, EsDelUsuario: f.EsDelUsuario}
	}
	sort.SliceStable(filas, func(a, b int) bool { return filas[a].Posicion < filas[b].Posicion })
	return filas, nil
}

// jugadoresSeleccion devuelve los jugadores de la pantalla de clasificación o
// de plantilla de equipo actual.
func (m Modelo) jugadoresSeleccion() ([]aplicacion.EstadisticaJugador, error) {
	if m.pantalla == pantallaClasificacion {
		return m.carrera.Clasificacion(m.criterio, maxFilasClasificacion)
	}
	return m.carrera.EstadisticasDeEquipo(m.equipoSel)
}

// totalSeleccion es la cantidad de filas de la pantalla seleccionable actual.
func (m Modelo) totalSeleccion() int {
	if m.pantalla == pantallaEquipos {
		filas, _ := m.filasEquipos()
		return len(filas)
	}
	filas, _ := m.jugadoresSeleccion()
	return len(filas)
}

func (m Modelo) esSeleccionable() bool {
	switch m.pantalla {
	case pantallaClasificacion, pantallaEquipos, pantallaPlantillaEquipo:
		return true
	}
	return false
}

// visiblesSeleccion es cuántas filas caben; sin tamaño conocido, todas.
func (m Modelo) visiblesSeleccion(total int) int {
	if m.alto <= 0 {
		return total
	}
	return min(total, max(m.alto-lineasFijas, 3))
}

// ajustarVentana mantiene el cursor y el scroll dentro de rango y con el
// cursor a la vista tras un cambio de tamaño.
func (m Modelo) ajustarVentana() Modelo {
	if !m.esSeleccionable() {
		if m.pantalla == pantallaFicha {
			m.scroll = min(max(m.scroll, 0), max(len(m.carrera.Trayectoria(m.fichaID))-m.visiblesFicha(), 0))
			return m
		}
		m.scroll = m.limitarScroll(m.scroll)
		return m
	}
	total := m.totalSeleccion()
	m.cursor = min(max(m.cursor, 0), max(total-1, 0))
	visibles := m.visiblesSeleccion(total)
	if m.cursor < m.scroll {
		m.scroll = m.cursor
	}
	if m.cursor >= m.scroll+visibles {
		m.scroll = m.cursor - visibles + 1
	}
	m.scroll = min(max(m.scroll, 0), max(total-visibles, 0))
	return m
}

func (m Modelo) teclaSeleccion(k string) Modelo {
	total := m.totalSeleccion()
	visibles := m.visiblesSeleccion(total)
	switch k {
	case "up", "k":
		m.cursor--
	case "down", "j":
		m.cursor++
	case "pgup":
		m.cursor -= visibles
	case "pgdown":
		m.cursor += visibles
	case "home":
		m.cursor = 0
	case "end":
		m.cursor = total - 1
	case "tab":
		if m.pantalla == pantallaEquipos {
			m.verAsistent = !m.verAsistent
		}
	case "esc", "backspace":
		return m.volver()
	case "enter":
		return m.abrirSeleccion()
	}
	return m.ajustarVentana()
}

// abrirSeleccion baja a la ficha del jugador o a la plantilla del equipo
// elegido.
func (m Modelo) abrirSeleccion() Modelo {
	if m.pantalla == pantallaEquipos {
		filas, err := m.filasEquipos()
		if err != nil || m.cursor >= len(filas) {
			return m
		}
		m.equipoSel = filas[m.cursor].Equipo
		return m.irA(pantallaPlantillaEquipo)
	}
	filas, err := m.jugadoresSeleccion()
	if err != nil || m.cursor >= len(filas) {
		return m
	}
	m.fichaID = filas[m.cursor].Jugador.ID
	return m.irA(pantallaFicha)
}

// Líneas de la ficha que no son filas de la trayectoria.
const lineasFicha = 16

// edadMaxProyeccion es la edad hasta la que la ficha muestra la proyección del
// jugador; después ya se ve en su nivel.
const edadMaxProyeccion = 24

func (m Modelo) visiblesFicha() int {
	total := len(m.carrera.Trayectoria(m.fichaID))
	if m.alto <= 0 {
		return total
	}
	return min(total, max(m.alto-lineasFicha, 3))
}

func (m Modelo) teclaFicha(k string) Modelo {
	total := len(m.carrera.Trayectoria(m.fichaID))
	visibles := m.visiblesFicha()
	switch k {
	case "up", "k":
		m.scroll--
	case "down", "j":
		m.scroll++
	case "pgup":
		m.scroll -= visibles
	case "pgdown":
		m.scroll += visibles
	case "esc", "backspace":
		return m.volver()
	}
	m.scroll = min(max(m.scroll, 0), max(total-visibles, 0))
	return m
}

// ---- Vistas ----

func (m Modelo) subtitulo() string {
	c := m.carrera
	return fmt.Sprintf("Temporada %d · Jornada %d / %d", c.Numero, c.Jornada(), c.TotalJornadas())
}

func (m Modelo) vistaEstadisticas() string {
	var b strings.Builder
	b.WriteString(estiloTitulo.Render("ESTADÍSTICAS · "+m.subtitulo()) + "\n\n")
	for i, o := range opcionesEstadisticas {
		if i == m.cursor {
			b.WriteString(estiloCursor.Render("> "+o) + "\n")
		} else {
			b.WriteString("  " + o + "\n")
		}
	}
	b.WriteString("\n" + ayuda("↑/↓ mover · enter elegir · esc volver"))
	return b.String()
}

// marcador devuelve el prefijo de una fila seleccionable: ">" en el cursor y
// "*" en las filas del club del usuario.
func marcador(cursor, delUsuario bool) string {
	switch {
	case cursor:
		return "> "
	case delUsuario:
		return "* "
	}
	return "  "
}

// pintar da estilo a una fila: el cursor en negrita y el club del usuario con
// su color.
func pintar(texto string, cursor, delUsuario bool) string {
	switch {
	case delUsuario:
		return estiloUsuario.Render(texto)
	case cursor:
		return estiloCursor.Render(texto)
	}
	return texto
}

func recortar(s string, n int) string {
	r := []rune(s)
	if len(r) <= n {
		return s
	}
	return string(r[:n-1]) + "…"
}

func sinDatos(m Modelo) string {
	return "Todavía no hay estadísticas: juega una jornada.\n"
}

type columna struct {
	cabecera string
	ancho    int
	valor    func(aplicacion.EstadisticaJugador) string
}

func entero(f func(aplicacion.EstadisticaJugador) int) func(aplicacion.EstadisticaJugador) string {
	return func(e aplicacion.EstadisticaJugador) string { return fmt.Sprint(f(e)) }
}

func valoracion(e modelo.Estadisticas) string {
	if e.Partidos == 0 {
		return "-"
	}
	return fmt.Sprintf("%.1f", e.ValoracionMedia())
}

// columnasDe son las columnas de estadísticas de cada clasificación.
func columnasDe(c aplicacion.Criterio) []columna {
	pj := columna{"PJ", 3, entero(func(e aplicacion.EstadisticaJugador) int { return e.Partidos })}
	switch c {
	case aplicacion.PorGoles:
		return []columna{pj,
			{"Goles", 5, entero(func(e aplicacion.EstadisticaJugador) int { return e.Goles })},
			{"Asist", 5, entero(func(e aplicacion.EstadisticaJugador) int { return e.Asistencias })},
			{"Min", 5, entero(func(e aplicacion.EstadisticaJugador) int { return e.Minutos })}}
	case aplicacion.PorAsistencias:
		return []columna{pj,
			{"Asist", 5, entero(func(e aplicacion.EstadisticaJugador) int { return e.Asistencias })},
			{"Goles", 5, entero(func(e aplicacion.EstadisticaJugador) int { return e.Goles })},
			{"Min", 5, entero(func(e aplicacion.EstadisticaJugador) int { return e.Minutos })}}
	case aplicacion.PorTarjetas:
		return []columna{pj,
			{"Amar", 5, entero(func(e aplicacion.EstadisticaJugador) int { return e.Amarillas })},
			{"Roja", 5, entero(func(e aplicacion.EstadisticaJugador) int { return e.Rojas })},
			{"Pts", 5, entero(func(e aplicacion.EstadisticaJugador) int { return e.Amarillas + 3*e.Rojas })}}
	case aplicacion.PorImbatidas:
		return []columna{pj,
			{"Imb", 5, entero(func(e aplicacion.EstadisticaJugador) int { return e.PorteriasImbatidas })},
			{"Enc", 5, entero(func(e aplicacion.EstadisticaJugador) int { return e.GolesEncajados })}}
	case aplicacion.PorPocosGolesEncajados:
		return []columna{pj,
			{"Enc", 5, entero(func(e aplicacion.EstadisticaJugador) int { return e.GolesEncajados })},
			{"Enc/PJ", 6, func(e aplicacion.EstadisticaJugador) string {
				if e.Partidos == 0 {
					return "-"
				}
				return fmt.Sprintf("%.2f", float64(e.GolesEncajados)/float64(e.Partidos))
			}}}
	default: // valoración
		return []columna{pj,
			{"Val", 5, func(e aplicacion.EstadisticaJugador) string { return valoracion(e.Estadisticas) }},
			{"Goles", 5, entero(func(e aplicacion.EstadisticaJugador) int { return e.Goles })},
			{"Asist", 5, entero(func(e aplicacion.EstadisticaJugador) int { return e.Asistencias })}}
	}
}

// ventana devuelve el tramo [desde, hasta) de n filas que se muestra.
func (m Modelo) ventana(total int) (desde, hasta int) {
	visibles := m.visiblesSeleccion(total)
	desde = min(max(m.scroll, 0), max(total-visibles, 0))
	return desde, min(desde+visibles, total)
}

func (m Modelo) vistaClasificacion() string {
	titulo := ""
	for _, c := range clasificaciones {
		if c.criterio == m.criterio {
			titulo = c.titulo
		}
	}
	var b strings.Builder
	b.WriteString(estiloTitulo.Render(titulo+" · "+m.subtitulo()) + "\n\n")

	filas, err := m.jugadoresSeleccion()
	switch {
	case err != nil:
		b.WriteString(estiloAviso.Render("No se pudieron calcular las estadísticas: "+err.Error()) + "\n\n")
		b.WriteString(ayuda("esc volver"))
		return b.String()
	case len(filas) == 0:
		b.WriteString(sinDatos(m) + "\n" + ayuda("esc volver"))
		return b.String()
	}

	cols := columnasDe(m.criterio)
	cab := fmt.Sprintf("  %2s  %-22s %-25s", "#", "Jugador", "Equipo")
	for _, c := range cols {
		cab += fmt.Sprintf(" %*s", c.ancho, c.cabecera)
	}
	b.WriteString(cab + "\n")

	club := m.carrera.NombreEquipo()
	desde, hasta := m.ventana(len(filas))
	for i := desde; i < hasta; i++ {
		f := filas[i]
		delUsuario := f.Equipo == club
		linea := fmt.Sprintf("%2d  %-22s %-25s", i+1, recortar(f.Jugador.Nombre, 22), recortar(f.Equipo, 25))
		for _, c := range cols {
			linea += fmt.Sprintf(" %*s", c.ancho, c.valor(f))
		}
		b.WriteString(pintar(marcador(i == m.cursor, delUsuario)+linea, i == m.cursor, delUsuario) + "\n")
	}
	b.WriteString("\n" + ayuda(m.ayudaSeleccion(desde, hasta, len(filas), "enter ver ficha", true)))
	return b.String()
}

// ayudaSeleccion es la línea de ayuda de las listas con cursor.
func (m Modelo) ayudaSeleccion(desde, hasta, total int, accion string, club bool) string {
	partes := []string{"↑/↓ mover", accion}
	if hasta-desde < total {
		partes[0] = fmt.Sprintf("↑/↓ mover (%d-%d de %d)", desde+1, hasta, total)
	}
	if club {
		partes = append(partes, "* tu club")
	}
	partes = append(partes, "esc volver")
	return strings.Join(partes, " · ")
}

func (m Modelo) vistaEquipos() string {
	var b strings.Builder
	b.WriteString(estiloTitulo.Render("EQUIPOS · "+m.subtitulo()) + "\n\n")
	filas, err := m.filasEquipos()
	if err != nil {
		b.WriteString(estiloAviso.Render("No se pudieron calcular las estadísticas: "+err.Error()) + "\n\n")
		b.WriteString(ayuda("esc volver"))
		return b.String()
	}
	ultima := "Goleador"
	if m.verAsistent {
		ultima = "Asistente"
	}
	fmt.Fprintf(&b, "  %2s  %-25s %2s %3s %3s %3s %2s %3s %2s  %s\n",
		"#", "Equipo", "PJ", "GF", "GC", "Imb", "SM", "Am", "Ro", ultima)
	desde, hasta := m.ventana(len(filas))
	for i := desde; i < hasta; i++ {
		f := filas[i]
		nombre, n := f.Goleador, f.GolesGoleador
		if m.verAsistent {
			nombre, n = f.Asistente, f.Asistencias
		}
		destacado := "-"
		if n > 0 {
			destacado = fmt.Sprintf("%s (%d)", recortar(nombre, 17), n)
		}
		linea := fmt.Sprintf("%2d  %-25s %2d %3d %3d %3d %2d %3d %2d  %s",
			f.Posicion, recortar(f.Equipo, 25), f.PJ, f.GF, f.GC, f.Imbatidas, f.SinMarcar, f.Amarillas, f.Rojas, destacado)
		b.WriteString(pintar(marcador(i == m.cursor, f.EsDelUsuario)+linea, i == m.cursor, f.EsDelUsuario) + "\n")
	}
	b.WriteString("\n" + ayuda(strings.Replace(
		m.ayudaSeleccion(desde, hasta, len(filas), "enter plantilla", true),
		"* tu club", "tab goleador/asistente · * tu club", 1)))
	return b.String()
}

// cabeceraPlantillaStats y filaPlantillaStats dan el formato de la plantilla con
// estadísticas.
const formatoPlantillaStats = "%-13s %-22s %3s %3s %5s %3s %3s %3s %3s %3s %4s"

func cabeceraPlantillaStats() string {
	return fmt.Sprintf(formatoPlantillaStats, "Posición", "Jugador", "PJ", "Tit", "Min", "G", "A", "Am", "Ro", "Imb", "Val")
}

func filaPlantillaStats(f aplicacion.EstadisticaJugador) string {
	return fmt.Sprintf(formatoPlantillaStats, f.Jugador.Posicion, recortar(f.Jugador.Nombre, 22),
		fmt.Sprint(f.Partidos), fmt.Sprint(f.Titularidades), fmt.Sprint(f.Minutos), fmt.Sprint(f.Goles),
		fmt.Sprint(f.Asistencias), fmt.Sprint(f.Amarillas), fmt.Sprint(f.Rojas),
		fmt.Sprint(f.PorteriasImbatidas), valoracion(f.Estadisticas))
}

func (m Modelo) vistaPlantillaEquipo() string {
	var b strings.Builder
	b.WriteString(estiloTitulo.Render(strings.ToUpper(m.equipoSel)+" · estadísticas") + "\n\n")
	filas, err := m.jugadoresSeleccion()
	if err != nil {
		b.WriteString(estiloAviso.Render("No se pudieron calcular las estadísticas: "+err.Error()) + "\n\n")
		b.WriteString(ayuda("esc volver"))
		return b.String()
	}
	b.WriteString("  " + cabeceraPlantillaStats() + "\n")
	desde, hasta := m.ventana(len(filas))
	for i := desde; i < hasta; i++ {
		b.WriteString(pintar(marcador(i == m.cursor, false)+filaPlantillaStats(filas[i]), i == m.cursor, false) + "\n")
	}
	b.WriteString("\n" + ayuda(m.ayudaSeleccion(desde, hasta, len(filas), "enter ver ficha", false)))
	return b.String()
}

// vistaPlantillaStats es la plantilla del usuario con estadísticas (se alterna
// con los atributos con tab).
func (m Modelo) vistaPlantillaStats() string {
	var b strings.Builder
	b.WriteString(estiloTitulo.Render("PLANTILLA · "+m.carrera.NombreEquipo()+" · estadísticas") + "\n\n")
	filas, err := m.carrera.EstadisticasDeEquipo(m.carrera.NombreEquipo())
	if err != nil {
		b.WriteString(estiloAviso.Render("No se pudieron calcular las estadísticas: "+err.Error()) + "\n\n")
		b.WriteString(ayuda("esc volver"))
		return b.String()
	}
	b.WriteString("  " + cabeceraPlantillaStats() + "\n")
	desde := m.limitarScroll(m.scroll)
	hasta := min(desde+m.filasVisibles(), len(filas))
	for _, f := range filas[desde:hasta] {
		b.WriteString("  " + filaPlantillaStats(f) + "\n")
	}
	b.WriteString("\n" + ayuda(strings.Replace(m.ayudaLista(desde, hasta, len(filas)), "esc volver", "tab atributos · esc volver", 1)))
	return b.String()
}

func resumenStats(e modelo.Estadisticas, porteria bool) string {
	s := fmt.Sprintf("PJ %d · Tit %d · Min %d · G %d · A %d · Am %d · Ro %d · Imb %d",
		e.Partidos, e.Titularidades, e.Minutos, e.Goles, e.Asistencias, e.Amarillas, e.Rojas, e.PorteriasImbatidas)
	if porteria {
		s += fmt.Sprintf(" · Enc %d", e.GolesEncajados)
	}
	return s + " · Val " + valoracion(e)
}

func (m Modelo) vistaFicha() string {
	c := m.carrera
	var b strings.Builder
	fila, ok, err := c.EstadisticaDeJugador(m.fichaID)
	if err != nil || !ok {
		b.WriteString(estiloTitulo.Render("FICHA") + "\n\n")
		b.WriteString("No se encontró al jugador.\n\n" + ayuda("esc volver"))
		return b.String()
	}
	j := fila.Jugador
	a := j.Atributos
	b.WriteString(estiloTitulo.Render("FICHA · "+j.Nombre) + "\n")
	fmt.Fprintf(&b, "%s · %s · %d años · Valoración %d\n", fila.Equipo, j.Posicion, j.Edad, j.Valoracion())
	fmt.Fprintf(&b, "RIT %d  TIR %d  PAS %d  REG %d  DEF %d  FIS %d  REF %d\n",
		a.Ritmo, a.Tiro, a.Pase, a.Regate, a.Defensa, a.Fisico, a.Reflejos)
	// Pista del talento: en los jóvenes todavía no se sabe hasta dónde llegarán.
	if j.Edad <= edadMaxProyeccion {
		fmt.Fprintf(&b, "Proyección: %s\n", j.Proyeccion())
	}
	b.WriteString("\n")

	porteria := j.Posicion == modelo.Portero
	fmt.Fprintf(&b, "Temporada %d (en curso)\n  %s\n\n", c.Numero, resumenStats(fila.Estadisticas, porteria))

	tray := c.Trayectoria(j.ID)
	b.WriteString("Trayectoria\n")
	if len(tray) == 0 {
		b.WriteString("  Todavía no ha terminado ninguna temporada en esta carrera.\n")
	} else {
		fmt.Fprintf(&b, "  %4s  %-25s %3s %5s %3s %3s %3s %3s %3s %4s\n",
			"Temp", "Club", "PJ", "Min", "G", "A", "Am", "Ro", "Imb", "Val")
		visibles := m.visiblesFicha()
		desde := min(max(m.scroll, 0), max(len(tray)-visibles, 0))
		for _, t := range tray[desde:min(desde+visibles, len(tray))] {
			fmt.Fprintf(&b, "  %4d  %-25s %3d %5d %3d %3d %3d %3d %3d %4s\n",
				t.Temporada, recortar(t.Equipo, 25), t.Partidos, t.Minutos, t.Goles, t.Asistencias,
				t.Amarillas, t.Rojas, t.PorteriasImbatidas, valoracion(t.Estadisticas))
		}
	}
	if total, err := c.EstadisticasDeCarrera(j.ID); err == nil && len(tray) > 0 {
		fmt.Fprintf(&b, "\nCarrera: %d partidos · %d goles · %d asistencias · %d amarillas · %d rojas\n",
			total.Partidos, total.Goles, total.Asistencias, total.Amarillas, total.Rojas)
	} else {
		b.WriteString("\n")
	}
	ayudaFicha := "esc volver"
	if len(tray) > m.visiblesFicha() {
		ayudaFicha = "↑/↓ desplazar · esc volver"
	}
	b.WriteString("\n" + ayuda(ayudaFicha))
	return b.String()
}
