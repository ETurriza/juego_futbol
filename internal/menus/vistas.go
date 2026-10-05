package menus

import (
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

var (
	estiloTitulo  = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("212"))
	estiloUsuario = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("86"))
	estiloCursor  = lipgloss.NewStyle().Bold(true)
	estiloAyuda   = lipgloss.NewStyle().Faint(true)
	estiloAviso   = lipgloss.NewStyle().Foreground(lipgloss.Color("203"))
)

// View implementa tea.Model.
func (m Modelo) View() tea.View {
	var cuerpo string
	switch m.pantalla {
	case pantallaPlantilla:
		cuerpo = m.vistaPlantilla()
	case pantallaTabla:
		cuerpo = m.vistaTabla()
	case pantallaJornada:
		cuerpo = m.vistaJornada()
	case pantallaFin:
		cuerpo = m.vistaFin()
	default:
		cuerpo = m.vistaMenu()
	}
	v := tea.NewView(cuerpo)
	v.AltScreen = true
	v.WindowTitle = "Juego de fútbol"
	return v
}

func ayuda(texto string) string { return estiloAyuda.Render(texto) }

// fila antepone el prefijo a una fila y la resalta si es del usuario.
func fila(texto string, delUsuario bool) string {
	if delUsuario {
		return estiloUsuario.Render("> " + texto)
	}
	return "  " + texto
}

func (m Modelo) cabeceraEquipo() string {
	c := m.carrera
	return fmt.Sprintf("%s    Jornada %d / %d", c.NombreEquipo(), c.Jornada(), c.TotalJornadas())
}

func (m Modelo) listaOpciones() string {
	var b strings.Builder
	for i, o := range m.opciones() {
		if i == m.cursor {
			b.WriteString(estiloCursor.Render("> "+o) + "\n")
		} else {
			b.WriteString("  " + o + "\n")
		}
	}
	return b.String()
}

func (m Modelo) vistaMenu() string {
	var b strings.Builder
	b.WriteString(estiloTitulo.Render("JUEGO DE FÚTBOL · Modo carrera") + "\n\n")
	b.WriteString(m.cabeceraEquipo() + "\n")
	fmt.Fprintf(&b, "Valoración del equipo: %d\n\n", m.carrera.ValoracionEquipo())
	b.WriteString(m.listaOpciones())
	if m.aviso != "" {
		b.WriteString("\n" + estiloAviso.Render(m.aviso) + "\n")
	}
	b.WriteString("\n" + ayuda("↑/↓ mover · enter elegir · q salir"))
	return b.String()
}

func (m Modelo) vistaPlantilla() string {
	var b strings.Builder
	b.WriteString(estiloTitulo.Render("PLANTILLA · "+m.carrera.NombreEquipo()) + "\n\n")
	fmt.Fprintf(&b, "  %-14s %-22s %3s %4s  %3s %3s %3s %3s %3s %3s %3s\n",
		"Posición", "Nombre", "Ed", "Val", "RIT", "TIR", "PAS", "REG", "DEF", "FIS", "REF")

	plantilla := m.carrera.Plantilla()
	desde := m.limitarScroll(m.scroll)
	hasta := min(desde+m.filasVisibles(), len(plantilla))
	for _, j := range plantilla[desde:hasta] {
		a := j.Atributos
		b.WriteString(fila(fmt.Sprintf("%-14s %-22s %3d %4d  %3d %3d %3d %3d %3d %3d %3d",
			j.Posicion, j.Nombre, j.Edad, j.Valoracion(),
			a.Ritmo, a.Tiro, a.Pase, a.Regate, a.Defensa, a.Fisico, a.Reflejos), false) + "\n")
	}
	b.WriteString("\n" + ayuda(m.ayudaLista(desde, hasta, len(plantilla))))
	return b.String()
}

func (m Modelo) vistaTabla() string {
	var b strings.Builder
	c := m.carrera
	b.WriteString(estiloTitulo.Render(fmt.Sprintf("TABLA · Jornada %d / %d", c.Jornada(), c.TotalJornadas())) + "\n\n")
	fmt.Fprintf(&b, "  %3s  %-26s %3s %3s %3s %3s %4s %4s %4s %4s\n",
		"#", "Equipo", "PJ", "G", "E", "P", "GF", "GC", "DG", "Pts")

	tabla := c.Tabla()
	desde := m.limitarScroll(m.scroll)
	hasta := min(desde+m.filasVisibles(), len(tabla))
	for _, f := range tabla[desde:hasta] {
		b.WriteString(fila(fmt.Sprintf("%3d  %-26s %3d %3d %3d %3d %4d %4d %+4d %4d",
			f.Posicion, f.Equipo, f.PJ, f.G, f.E, f.P, f.GF, f.GC, f.DG, f.Pts), f.EsDelUsuario) + "\n")
	}
	b.WriteString("\n" + ayuda(m.ayudaLista(desde, hasta, len(tabla))))
	return b.String()
}

// ayudaLista es la línea de ayuda de las listas; indica el tramo mostrado si no
// caben todas las filas.
func (m Modelo) ayudaLista(desde, hasta, total int) string {
	if hasta-desde >= total {
		return "esc volver"
	}
	return fmt.Sprintf("↑/↓ desplazar (%d-%d de %d) · esc volver", desde+1, hasta, total)
}

func (m Modelo) vistaJornada() string {
	c := m.carrera
	var b strings.Builder
	b.WriteString(estiloTitulo.Render(fmt.Sprintf("JORNADA %d / %d", c.Jornada(), c.TotalJornadas())) + "\n\n")

	juegaUsuario := false
	for _, r := range c.UltimaJornada() {
		juegaUsuario = juegaUsuario || r.EsDelUsuario
		b.WriteString(fila(fmt.Sprintf("%-26s %2d - %-2d  %s",
			r.Local, r.GolesLocal, r.GolesVisitante, r.Visitante), r.EsDelUsuario) + "\n")
	}
	if !juegaUsuario {
		b.WriteString("\nTu equipo descansa esta jornada.\n")
	}
	b.WriteString("\n" + ayuda("enter continuar"))
	return b.String()
}

func (m Modelo) vistaFin() string {
	c := m.carrera
	var b strings.Builder
	b.WriteString(estiloTitulo.Render("TEMPORADA TERMINADA") + "\n\n")

	campeon, _ := c.Campeon()
	fmt.Fprintf(&b, "Campeón: %s\n", campeon)
	tabla := c.Tabla()
	for _, f := range tabla {
		if !f.EsDelUsuario {
			continue
		}
		fmt.Fprintf(&b, "Tu equipo: %s, puesto %d de %d, %d puntos\n",
			f.Equipo, f.Posicion, len(tabla), f.Pts)
		if f.Posicion == 1 {
			b.WriteString(estiloUsuario.Render("¡Felicidades, eres el campeón!") + "\n")
		}
	}
	b.WriteString("\n" + m.listaOpciones())
	if m.aviso != "" {
		b.WriteString("\n" + estiloAviso.Render(m.aviso) + "\n")
	}
	b.WriteString("\n" + ayuda("↑/↓ mover · enter elegir · q salir"))
	return b.String()
}
