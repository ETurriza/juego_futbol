package menus

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"testing"

	"charm.land/lipgloss/v2"

	"github.com/ETurriza/juego_futbol/internal/aplicacion"
)

// actualizar reescribe las pantallas de referencia: go test ./internal/menus -update
var actualizar = flag.Bool("update", false,
	"reescribe las pantallas de referencia (docs/PANTALLAS.md)")

const (
	archivoDocs  = "../../docs/PANTALLAS.md"
	anchoGaleria = 80
	altoGaleria  = 30
)

// captura es una pantalla de la galería.
type captura struct {
	archivo   string // nombre del archivo de referencia, sin extensión
	titulo    string
	contenido string // texto de la pantalla, tomado en el momento de la captura
}

// galeria arma las pantallas principales del juego con una carrera de semilla
// fija, tal como se verían en una terminal de 80 por 30. Son la referencia visual
// del juego: cualquier cambio en lo que se ve aparece como diferencia en el
// revisor de código.
func galeria(t *testing.T) []captura {
	t.Helper()

	// Una carrera veterana: tres temporadas terminadas y once jornadas de la cuarta.
	c := carreraDePrueba(t, 21, 10)
	avanzarCarrera(t, c, 3)
	for i := 0; i < 11; i++ {
		if _, err := c.AvanzarJornada(); err != nil {
			t.Fatal(err)
		}
	}
	veterana := redimensionar(Nuevo(c, nil), anchoGaleria, altoGaleria)
	filas, _ := c.EstadisticasDeEquipo(c.NombreEquipo())

	// Las capturas se toman en el momento: la carrera es mutable (avanzar jornada
	// la cambia), así que cada pantalla guarda su texto y las que avanzan la
	// carrera usan una copia.
	var out []captura
	desde := func(archivo, titulo string, base Modelo, teclas ...string) {
		m, _ := pulsar(base, teclas...)
		out = append(out, captura{archivo, titulo, normalizar(texto(m))})
	}
	copia := func() Modelo {
		clon, err := aplicacion.Importar(c.Exportar())
		if err != nil {
			t.Fatal(err)
		}
		return redimensionar(Nuevo(clon, nil), anchoGaleria, altoGaleria)
	}
	est := append(abajo(opEstadisticas), "enter")
	alEstadistica := func(extra ...string) []string { return append(append([]string(nil), est...), extra...) }

	desde("01-menu-principal", "Menú principal", veterana)
	desde("02-resultados-de-la-jornada", "Resultados de la jornada", copia(), "enter")
	desde("03-tabla-de-posiciones", "Tabla de posiciones", veterana, "down", "enter")
	desde("04-plantilla-atributos", "Plantilla: atributos", veterana, "down", "down", "enter")
	desde("05-plantilla-estadisticas", "Plantilla: estadísticas (tab)", veterana, "down", "down", "enter", "tab")
	desde("06-historial", "Historial de temporadas", veterana, append(abajo(opHistorial), "enter")...)
	desde("07-menu-de-estadisticas", "Menú de estadísticas", veterana, est...)
	for i, cl := range clasificaciones {
		slug := []string{"goleadores", "asistentes", "tarjetas", "porteros-imbatidas", "porteros-menos-goles", "valoraciones"}[i]
		desde(fmt.Sprintf("%02d-clasificacion-%s", 8+i, slug), "Clasificación: "+cl.titulo, veterana,
			alEstadistica(append(abajo(i), "enter")...)...)
	}
	desde("14-equipos-goleador", "Comparativa de equipos", veterana, alEstadistica(append(abajo(opEstEquipos), "enter")...)...)
	desde("15-equipos-asistente", "Comparativa de equipos (tab)", veterana, alEstadistica(append(abajo(opEstEquipos), "enter", "tab")...)...)
	desde("16-plantilla-de-un-equipo", "Plantilla de un equipo con estadísticas", veterana,
		alEstadistica(append(abajo(opEstEquipos), "enter", "down", "enter")...)...)

	ficha := veterana
	ficha.fichaID, ficha.pantalla = filas[0].Jugador.ID, pantallaFicha
	out = append(out, captura{"17-ficha-portero", "Ficha de un portero", normalizar(texto(ficha))})
	ficha.fichaID = filas[len(filas)-1].Jugador.ID
	out = append(out, captura{"18-ficha-delantero", "Ficha de un delantero", normalizar(texto(ficha))})

	// El cambio de temporada, con otra carrera que acaba de terminar la primera.
	fin := redimensionar(jugarTemporada(t, modeloDePrueba(t, 22, 10)), anchoGaleria, altoGaleria)
	desde("19-fin-de-temporada", "Fin de temporada", fin)
	desde("20-inicio-de-temporada", "Inicio de temporada", fin, "enter")
	desde("21-confirmar-nueva-carrera", "Confirmar una carrera nueva", fin, append(abajo(opFinNueva), "enter")...)
	return out
}

// normalizar quita los espacios sobrantes al final de cada línea.
func normalizar(s string) string {
	lineas := strings.Split(strings.TrimRight(s, "\n"), "\n")
	for i, l := range lineas {
		lineas[i] = strings.TrimRight(l, " ")
	}
	return strings.Join(lineas, "\n") + "\n"
}

func TestNingunaPantallaPasaDe80Columnas(t *testing.T) {
	for _, c := range galeria(t) {
		ancho := 0
		for _, l := range strings.Split(c.contenido, "\n") {
			ancho = max(ancho, lipgloss.Width(l))
		}
		if ancho > anchoGaleria {
			t.Errorf("%s: %d columnas, deberia caber en %d:\n%s", c.archivo, ancho, anchoGaleria, c.contenido)
		}
	}
}

// primeraDiferencia describe dónde difieren dos textos.
func primeraDiferencia(obtenido, esperado string) string {
	a, b := strings.Split(obtenido, "\n"), strings.Split(esperado, "\n")
	for i := 0; i < max(len(a), len(b)); i++ {
		var x, y string
		if i < len(a) {
			x = a[i]
		}
		if i < len(b) {
			y = b[i]
		}
		if x != y {
			return fmt.Sprintf("linea %d:\n  obtenida:  %q\n  esperada:  %q", i+1, x, y)
		}
	}
	return "sin diferencias visibles"
}

// TestGaleriaDePantallas compara las pantallas con la página de referencia
// docs/PANTALLAS.md. Si el cambio es intencionado:
//
//	go test ./internal/menus -update
func TestGaleriaDePantallas(t *testing.T) {
	capturas := galeria(t)

	var pagina strings.Builder
	pagina.WriteString("# Galería de pantallas\n\n")
	pagina.WriteString("Cómo se ve el juego en una terminal de 80 por 30, con la semilla 21 (cuatro\n")
	pagina.WriteString("temporadas de carrera). Este archivo se **genera** a partir de las pruebas de\n")
	pagina.WriteString("`internal/menus`; no se edita a mano. Para actualizarlo tras un cambio visual:\n\n")
	pagina.WriteString("```\ngo test ./internal/menus -update\n```\n\n")
	pagina.WriteString("En las filas con cursor, `>` marca la seleccionada y `*` las de tu club.\n")

	for _, c := range capturas {
		contenido := c.contenido
		fmt.Fprintf(&pagina, "\n## %s\n\n```\n%s```\n", c.titulo, contenido)
	}

	if *actualizar {
		if err := os.WriteFile(archivoDocs, []byte(pagina.String()), 0o644); err != nil {
			t.Fatal(err)
		}
	} else if previo, err := os.ReadFile(archivoDocs); err != nil {
		t.Errorf("falta docs/PANTALLAS.md (%v); generarlo con -update", err)
	} else if string(previo) != pagina.String() {
		t.Errorf("docs/PANTALLAS.md esta desactualizado; regenerarlo con -update. Primera diferencia, %s",
			primeraDiferencia(pagina.String(), string(previo)))
	}
}
