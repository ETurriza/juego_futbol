package menus

import (
	"fmt"
	"strings"
	"testing"

	"charm.land/lipgloss/v2"

	"github.com/ETurriza/juego_futbol/internal/aplicacion"
	"github.com/ETurriza/juego_futbol/internal/modelo"
)

// abajo devuelve n pulsaciones de la flecha hacia abajo.
func abajo(n int) []string {
	teclas := make([]string, n)
	for i := range teclas {
		teclas[i] = "down"
	}
	return teclas
}

// conJornadas devuelve un modelo con una carrera de 10 equipos y n jornadas ya
// jugadas.
func conJornadas(t *testing.T, semilla int64, n int) Modelo {
	t.Helper()
	c := carreraDePrueba(t, semilla, 10)
	for i := 0; i < n; i++ {
		if _, err := c.AvanzarJornada(); err != nil {
			t.Fatal(err)
		}
	}
	return Nuevo(c, nil)
}

// abrirEstadisticas abre el menú de estadísticas desde el menú principal, o
// desde el de fin de temporada si la carrera ya terminó.
func abrirEstadisticas(m Modelo) Modelo {
	opcion := opEstadisticas
	if m.pantalla == pantallaFin {
		opcion = opFinEstadisticas
	}
	m, _ = pulsar(m, append(abajo(opcion), "enter")...)
	return m
}

// abrirOpcion elige una opción del menú de estadísticas.
func abrirOpcion(m Modelo, opcion int) Modelo {
	m, _ = pulsar(m, append(abajo(opcion), "enter")...)
	return m
}

func lineasDe(m Modelo, prefijo string) []string {
	var out []string
	for _, l := range strings.Split(texto(m), "\n") {
		if strings.HasPrefix(l, prefijo) {
			out = append(out, l)
		}
	}
	return out
}

func TestMenuDeEstadisticas(t *testing.T) {
	m := conJornadas(t, 1, 4)
	m = abrirEstadisticas(m)
	contiene(t, m, "ESTADÍSTICAS · Temporada 1 · Jornada 4 / 18", "> Goleadores", "  Asistentes", "  Tarjetas",
		"  Porteros: porterías imbatidas", "  Porteros: menos goles encajados", "  Mejores valoraciones",
		"  Equipos", "  Volver", "esc volver")
	m, _ = pulsar(m, "down", "down")
	contiene(t, m, "> Tarjetas")
	m, _ = pulsar(m, "up", "k", "k", "k") // no pasa del principio
	contiene(t, m, "> Goleadores")
	m, _ = pulsar(m, "j", "j", "j", "j", "j", "j", "j", "j", "j") // ni del final
	contiene(t, m, "> Volver")

	// q no cierra el juego aquí; ctrl+c sí.
	if _, cmd := pulsar(m, "q"); esSalir(cmd) {
		t.Error("q no deberia salir desde las estadisticas")
	}
	if _, cmd := pulsar(m, "ctrl+c"); !esSalir(cmd) {
		t.Error("ctrl+c deberia salir")
	}

	// Volver deja el menú principal sobre Estadísticas.
	m, _ = pulsar(m, "enter")
	contiene(t, m, "JUEGO DE FÚTBOL", "> Estadísticas")
	m, _ = pulsar(m, "enter") // el cursor sigue en Estadísticas: se vuelve a abrir
	contiene(t, m, "ESTADÍSTICAS")
	m, _ = pulsar(m, "esc")
	contiene(t, m, "JUEGO DE FÚTBOL", "> Estadísticas")
}

func TestEstadisticasDesdeElFinDeTemporada(t *testing.T) {
	m := jugarTemporada(t, modeloDePrueba(t, 1, 10))
	m, _ = pulsar(m, append(abajo(opFinEstadisticas), "enter")...)
	contiene(t, m, "ESTADÍSTICAS · Temporada 1 · Jornada 18 / 18", "> Goleadores")
	m = abrirOpcion(m, 0)
	contiene(t, m, "GOLEADORES")
	m, _ = pulsar(m, "esc", "esc")
	contiene(t, m, "TEMPORADA 1 TERMINADA", "> Estadísticas")
}

func TestSinPartidosNoHayClasificaciones(t *testing.T) {
	m := abrirEstadisticas(conJornadas(t, 1, 0))
	for i := range clasificaciones {
		pantalla := abrirOpcion(m, i)
		contiene(t, pantalla, clasificaciones[i].titulo, "Todavía no hay estadísticas: juega una jornada.", "esc volver")
		noContiene(t, pantalla, "Jugador")
	}
	// Los equipos sí aparecen, con ceros.
	m = abrirOpcion(m, opEstEquipos)
	contiene(t, m, "EQUIPOS", "Equipo", "Goleador")
	if filas := lineasDe(m, "> "); len(filas) != 1 {
		t.Errorf("deberia haber una fila con cursor: %v", filas)
	}
}

func TestCadaClasificacionTieneSusColumnas(t *testing.T) {
	m := abrirEstadisticas(conJornadas(t, 2, 18))
	columnas := []string{"Goles", "Asist", "Amar", "Imb", "Enc/PJ", "Val"}
	for i, c := range clasificaciones {
		pantalla := abrirOpcion(m, i)
		contiene(t, pantalla, c.titulo+" · Temporada 1 · Jornada 18 / 18", "Jugador", "Equipo", "PJ", columnas[i])
		if filas := lineasDe(pantalla, "> "); len(filas) != 1 {
			t.Errorf("%s: %d filas con cursor", c.titulo, len(filas))
		}
	}
}

func TestLaClasificacionMuestraLosMejoresEnOrden(t *testing.T) {
	m := conJornadas(t, 3, 18)
	esperado, err := m.carrera.Clasificacion(aplicacion.PorGoles, maxFilasClasificacion)
	if err != nil {
		t.Fatal(err)
	}
	m = abrirOpcion(abrirEstadisticas(m), 0)
	// Cada puesto aparece con su número y su jugador, en orden (dos jugadores
	// pueden tener el mismo nombre, así que se busca por puesto y nombre).
	lineas := strings.Split(texto(m), "\n")
	siguiente := 0
	for i, f := range esperado {
		buscado := fmt.Sprintf("%2d  %s", i+1, f.Jugador.Nombre)
		for siguiente < len(lineas) && !strings.Contains(lineas[siguiente], buscado) {
			siguiente++
		}
		if siguiente == len(lineas) {
			t.Fatalf("falta %q o esta fuera de orden:\n%s", buscado, texto(m))
		}
		siguiente++
	}
	// El primero lleva el cursor y sus goles.
	primero := esperado[0]
	if linea := lineasDe(m, "> "); len(linea) != 1 || !strings.Contains(linea[0], primero.Jugador.Nombre) ||
		!strings.Contains(linea[0], fmt.Sprint(primero.Goles)) {
		t.Errorf("la fila del cursor deberia ser la del goleador: %v", linea)
	}
}

func TestLasFilasDeTuClubVanMarcadas(t *testing.T) {
	m := conJornadas(t, 4, 18)
	club := m.carrera.NombreEquipo()
	filas, _ := m.carrera.Clasificacion(aplicacion.PorValoracion, maxFilasClasificacion)
	delClub := 0
	for _, f := range filas {
		if f.Equipo == club {
			delClub++
		}
	}
	m = abrirOpcion(abrirEstadisticas(m), 5) // mejores valoraciones
	marcadas := lineasDe(m, "* ")
	// La fila con el cursor lleva ">", no "*": si es del club, no se cuenta.
	if filas[0].Equipo == club {
		delClub--
	}
	if len(marcadas) != delClub {
		t.Errorf("%d filas marcadas, se esperaban %d", len(marcadas), delClub)
	}
	for _, l := range marcadas {
		if !strings.Contains(l, club) {
			t.Errorf("fila marcada que no es de tu club: %q", l)
		}
	}
	contiene(t, m, "* tu club")
}

func TestCursorYFichaDelJugador(t *testing.T) {
	m := conJornadas(t, 5, 12)
	filas, _ := m.carrera.Clasificacion(aplicacion.PorGoles, maxFilasClasificacion)
	m = abrirOpcion(abrirEstadisticas(m), 0)
	m, _ = pulsar(m, "down", "down", "j")
	elegido := filas[3]
	if linea := lineasDe(m, "> "); len(linea) != 1 || !strings.Contains(linea[0], elegido.Jugador.Nombre) {
		t.Fatalf("el cursor deberia estar sobre %s: %v", elegido.Jugador.Nombre, linea)
	}

	m, _ = pulsar(m, "enter")
	a := elegido.Jugador.Atributos
	contiene(t, m, "FICHA · "+elegido.Jugador.Nombre, elegido.Equipo, elegido.Jugador.Posicion.String(),
		fmt.Sprintf("%d años", elegido.Jugador.Edad), fmt.Sprintf("Valoración %d", elegido.Jugador.Valoracion()),
		fmt.Sprintf("RIT %d", a.Ritmo), fmt.Sprintf("REF %d", a.Reflejos),
		"Temporada 1 (en curso)", fmt.Sprintf("PJ %d", elegido.Partidos), fmt.Sprintf("G %d", elegido.Goles),
		"Trayectoria", "Todavía no ha terminado ninguna temporada", "esc volver")
	noContiene(t, m, "Carrera:")

	// Volver deja el cursor donde estaba.
	m, _ = pulsar(m, "esc")
	contiene(t, m, "GOLEADORES")
	if linea := lineasDe(m, "> "); len(linea) != 1 || !strings.Contains(linea[0], elegido.Jugador.Nombre) {
		t.Errorf("tras volver, el cursor deberia seguir en %s: %v", elegido.Jugador.Nombre, linea)
	}
	m, _ = pulsar(m, "backspace")
	contiene(t, m, "ESTADÍSTICAS", "> Goleadores")
}

func TestScrollDeLaClasificacion(t *testing.T) {
	m := conJornadas(t, 6, 18)
	todos, _ := m.carrera.Clasificacion(aplicacion.PorValoracion, maxFilasClasificacion)
	if len(todos) < 30 {
		t.Fatalf("solo %d jugadores en la clasificacion", len(todos))
	}
	m = redimensionar(m, 90, 12) // caben 7 filas
	m = abrirOpcion(abrirEstadisticas(m), 5)
	contiene(t, m, todos[0].Jugador.Nombre, fmt.Sprintf("1-7 de %d", len(todos)))
	noContiene(t, m, todos[10].Jugador.Nombre)

	// Bajar más allá de la ventana la desplaza para mantener el cursor a la vista.
	m, _ = pulsar(m, append(abajo(8), "j")...)
	contiene(t, m, todos[9].Jugador.Nombre)
	if linea := lineasDe(m, "> "); len(linea) != 1 || !strings.Contains(linea[0], todos[9].Jugador.Nombre) {
		t.Errorf("el cursor deberia seguir visible en el puesto 10: %v", linea)
	}
	noContiene(t, m, todos[0].Jugador.Nombre)

	m, _ = pulsar(m, "end")
	contiene(t, m, todos[len(todos)-1].Jugador.Nombre, fmt.Sprintf("de %d", len(todos)))
	m, _ = pulsar(m, "home")
	contiene(t, m, todos[0].Jugador.Nombre)
	m, _ = pulsar(m, "pgdown", "pgdown")
	if linea := lineasDe(m, "> "); len(linea) != 1 || !strings.Contains(linea[0], todos[14].Jugador.Nombre) {
		t.Errorf("dos paginas de 7 deberian llevar al puesto 15: %v", linea)
	}
	m, _ = pulsar(m, "pgup")
	if linea := lineasDe(m, "> "); len(linea) != 1 || !strings.Contains(linea[0], todos[7].Jugador.Nombre) {
		t.Errorf("subir una pagina deberia llevar al puesto 8: %v", linea)
	}

	// Agrandar la ventana muestra todo sin romper el cursor.
	m = redimensionar(m, 90, 80)
	noContiene(t, m, " de "+fmt.Sprint(len(todos)))
	if linea := lineasDe(m, "> "); len(linea) != 1 {
		t.Errorf("tras redimensionar deberia haber un cursor: %v", linea)
	}
}

func TestComparativaDeEquipos(t *testing.T) {
	m := conJornadas(t, 7, 18)
	club := m.carrera.NombreEquipo()
	tabla := m.carrera.Tabla()
	equipos, _ := m.carrera.EstadisticasEquipos()
	porNombre := map[string]aplicacion.EstadisticaEquipo{}
	for _, e := range equipos {
		porNombre[e.Equipo] = e
	}

	m = abrirOpcion(abrirEstadisticas(m), opEstEquipos)
	contiene(t, m, "EQUIPOS · Temporada 1 · Jornada 18 / 18", "Equipo", "GF", "GC", "Imb", "SM", "Am", "Ro", "Goleador",
		"tab goleador/asistente", "* tu club")
	// En el orden de la tabla, con sus datos.
	pantalla := texto(m)
	ultimo := -1
	for _, f := range tabla {
		pos := strings.Index(pantalla, f.Equipo)
		if pos < 0 || pos < ultimo {
			t.Fatalf("%s no aparece en el orden de la tabla", f.Equipo)
		}
		ultimo = pos
		e := porNombre[f.Equipo]
		if e.GolesGoleador > 0 {
			contiene(t, m, fmt.Sprintf("(%d)", e.GolesGoleador))
		}
	}
	// El cursor está en el líder; la fila de tu club, si no es el líder, lleva "*".
	if linea := lineasDe(m, "> "); len(linea) != 1 || !strings.Contains(linea[0], tabla[0].Equipo) {
		t.Errorf("el cursor deberia estar en el lider: %v", linea)
	}
	if tabla[0].Equipo != club {
		if linea := lineasDe(m, "* "); len(linea) != 1 || !strings.Contains(linea[0], club) {
			t.Errorf("la fila de tu club deberia ir marcada: %v", linea)
		}
	}

	// tab alterna goleador y asistente.
	m, _ = pulsar(m, "tab")
	contiene(t, m, "Asistente")
	noContiene(t, m, "Goleador")
	m, _ = pulsar(m, "tab")
	contiene(t, m, "Goleador")
}

func TestDeEquiposALaPlantillaYALaFicha(t *testing.T) {
	m := conJornadas(t, 8, 18)
	tabla := m.carrera.Tabla()
	m = abrirOpcion(abrirEstadisticas(m), opEstEquipos)
	m, _ = pulsar(m, "down", "down", "enter") // el tercero de la tabla
	elegido := tabla[2].Equipo
	contiene(t, m, strings.ToUpper(elegido)+" · estadísticas", "Posición", "Jugador", "PJ", "Tit", "Min", "Imb", "Val", "enter ver ficha")
	filas, _ := m.carrera.EstadisticasDeEquipo(elegido)
	for _, f := range filas {
		contiene(t, m, f.Jugador.Nombre)
	}

	// A la ficha de un jugador de ese equipo, y de vuelta sin perder el sitio.
	m, _ = pulsar(m, "down", "enter")
	contiene(t, m, "FICHA · "+filas[1].Jugador.Nombre, elegido)
	m, _ = pulsar(m, "esc")
	contiene(t, m, strings.ToUpper(elegido)+" · estadísticas")
	if linea := lineasDe(m, "> "); len(linea) != 1 || !strings.Contains(linea[0], filas[1].Jugador.Nombre) {
		t.Errorf("el cursor deberia seguir en el segundo jugador: %v", linea)
	}
	m, _ = pulsar(m, "esc")
	contiene(t, m, "EQUIPOS")
	if linea := lineasDe(m, "> "); len(linea) != 1 || !strings.Contains(linea[0], tabla[2].Equipo) {
		t.Errorf("el cursor deberia seguir en el tercer equipo: %v", linea)
	}
	// Con la temporada terminada, se vuelve al menú de fin de temporada.
	m, _ = pulsar(m, "esc", "esc")
	contiene(t, m, "TEMPORADA 1 TERMINADA", "> Estadísticas")
}

func TestFichaConTrayectoriaYCarrera(t *testing.T) {
	c := carreraDePrueba(t, 9, 10)
	avanzarCarrera(t, c, 3)
	for i := 0; i < 8; i++ {
		c.AvanzarJornada()
	}
	// Un jugador que haya jugado las tres temporadas archivadas.
	cuenta := map[int]int{}
	for _, a := range c.Archivo {
		cuenta[a.Jugador]++
	}
	var id int
	for j, n := range cuenta {
		if n == 3 && (id == 0 || j < id) {
			id = j
		}
	}
	if id == 0 {
		t.Fatal("nadie jugo las tres temporadas")
	}
	fila, _, _ := c.EstadisticaDeJugador(id)
	tray := c.Trayectoria(id)
	total, _ := c.EstadisticasDeCarrera(id)

	m := Nuevo(c, nil)
	m.fichaID, m.pantalla = id, pantallaFicha
	contiene(t, m, "FICHA · "+fila.Jugador.Nombre, "Temporada 4 (en curso)", "Trayectoria", "Temp", "Club", "Val",
		fmt.Sprintf("Carrera: %d partidos · %d goles · %d asistencias · %d amarillas · %d rojas",
			total.Partidos, total.Goles, total.Asistencias, total.Amarillas, total.Rojas))
	for _, a := range tray {
		contiene(t, m, fmt.Sprintf("%4d  %s", a.Temporada, recortar(a.Equipo, 25)))
	}
	noContiene(t, m, "Todavía no ha terminado")
}

func TestFichaDeUnPorteroMuestraGolesEncajados(t *testing.T) {
	m := conJornadas(t, 10, 8)
	filas, _ := m.carrera.EstadisticasDeEquipo(m.carrera.NombreEquipo())
	m.fichaID, m.pantalla = filas[0].Jugador.ID, pantallaFicha // el primero es portero
	if filas[0].Jugador.Posicion != modelo.Portero {
		t.Fatal("el primero deberia ser un portero")
	}
	contiene(t, m, "Enc ")
	m.fichaID = filas[len(filas)-1].Jugador.ID // un delantero
	noContiene(t, m, "Enc ")
}

func TestFichaConScrollEnLaTrayectoria(t *testing.T) {
	c := carreraDePrueba(t, 11, 10)
	avanzarCarrera(t, c, 12)
	cuenta := map[int]int{}
	for _, a := range c.Archivo {
		cuenta[a.Jugador]++
	}
	var id, mejor int
	for j, n := range cuenta {
		if n > mejor || (n == mejor && j < id) {
			id, mejor = j, n
		}
	}
	if mejor < 6 {
		t.Skipf("nadie jugo suficientes temporadas (%d)", mejor)
	}
	tray := c.Trayectoria(id)
	m := redimensionar(Nuevo(c, nil), 90, lineasFicha+4) // caben 4 filas
	m.fichaID, m.pantalla = id, pantallaFicha
	contiene(t, m, fmt.Sprintf("%4d  ", tray[0].Temporada), "↑/↓ desplazar")
	noContiene(t, m, fmt.Sprintf("%4d  %s", tray[len(tray)-1].Temporada, recortar(tray[len(tray)-1].Equipo, 25)))
	m, _ = pulsar(m, "pgdown", "pgdown", "pgdown")
	contiene(t, m, fmt.Sprintf("%4d  %s", tray[len(tray)-1].Temporada, recortar(tray[len(tray)-1].Equipo, 25)))
	m, _ = pulsar(m, "up", "k")
	contiene(t, m, "FICHA")
}

func TestLaPlantillaAlternaAtributosYEstadisticas(t *testing.T) {
	m := conJornadas(t, 12, 9)
	m, _ = pulsar(m, "down", "down", "enter") // plantilla
	contiene(t, m, "PLANTILLA · "+m.carrera.NombreEquipo(), "RIT", "REF", "tab estadísticas")
	noContiene(t, m, "Tit")

	m, _ = pulsar(m, "tab")
	contiene(t, m, "PLANTILLA · "+m.carrera.NombreEquipo()+" · estadísticas", "Posición", "PJ", "Tit", "Min", "Imb", "Val", "tab atributos")
	filas, _ := m.carrera.EstadisticasDeEquipo(m.carrera.NombreEquipo())
	for _, f := range filas {
		contiene(t, m, f.Jugador.Nombre)
	}
	m, _ = pulsar(m, "tab")
	contiene(t, m, "RIT", "tab estadísticas")
	noContiene(t, m, "Tit")

	// Con scroll, igual que la plantilla de atributos.
	m = redimensionar(m, 90, 12)
	m, _ = pulsar(m, "tab", "end")
	contiene(t, m, filas[len(filas)-1].Jugador.Nombre, "de 22")
	noContiene(t, m, filas[0].Jugador.Nombre)
}

func TestLasPantallasDeEstadisticasAguantanVentanasChicas(t *testing.T) {
	base := conJornadas(t, 13, 6)
	for _, tam := range [][2]int{{0, 0}, {1, 1}, {20, 3}, {200, 100}} {
		m := redimensionar(base, tam[0], tam[1])
		recorridos := [][]string{
			append(abajo(opEstadisticas), "enter"),
			append(append(abajo(opEstadisticas), "enter"), "enter"),
			append(append(abajo(opEstadisticas), "enter"), append(abajo(5), "enter")...),
			append(append(abajo(opEstadisticas), "enter"), append(abajo(opEstEquipos), "enter")...),
			append(append(abajo(opEstadisticas), "enter"), append(abajo(opEstEquipos), "enter", "enter")...),
			append(append(abajo(opEstadisticas), "enter"), "enter", "enter"),
			append(append(abajo(opEstadisticas), "enter"), append(abajo(opEstEquipos), "enter", "enter", "enter")...),
		}
		for _, teclas := range recorridos {
			pantalla, _ := pulsar(m, teclas...)
			if texto(pantalla) == "" {
				t.Errorf("tamaño %v tras %v: pantalla vacia", tam, teclas)
			}
			pulsar(pantalla, "down", "pgdown", "end", "up", "pgup", "home", "tab", "esc", "esc", "esc")
		}
	}
}

func TestSiHayUnErrorDeEstadisticasSeMuestraSinCerrar(t *testing.T) {
	m := conJornadas(t, 14, 3)
	// Un partido cuyo detalle no cuadra con el marcador hace fallar el cálculo.
	m.carrera.Temporada.Resultados[0][0].GolesLocal += 7

	for _, opcion := range []int{0, opEstEquipos} {
		pantalla := abrirOpcion(abrirEstadisticas(m), opcion)
		contiene(t, pantalla, "No se pudieron calcular las estadísticas", "esc volver")
		if _, cmd := pulsar(pantalla, "down", "enter", "tab", "esc"); esSalir(cmd) {
			t.Error("un error no deberia cerrar el juego")
		}
	}
	m.pantalla, m.fichaID = pantallaFicha, 1
	contiene(t, m, "No se encontró al jugador")
	m.pantalla = pantallaPlantilla
	m.verStats = true
	contiene(t, m, "No se pudieron calcular las estadísticas")
}

func TestLaPilaDeNavegacionNoSeComparteEntreModelos(t *testing.T) {
	a := abrirEstadisticas(conJornadas(t, 15, 3))
	b := abrirOpcion(a, 0)
	c := abrirOpcion(a, opEstEquipos)
	if b.pantalla != pantallaClasificacion || c.pantalla != pantallaEquipos {
		t.Fatal("cada modelo deberia seguir su propio camino")
	}
	b, _ = pulsar(b, "esc")
	c, _ = pulsar(c, "esc")
	contiene(t, b, "> Goleadores")
	contiene(t, c, "> Equipos")
	if a.pantalla != pantallaEstadisticas || len(a.pila) != 1 {
		t.Errorf("el modelo original no deberia cambiar: pantalla %v, pila %d", a.pantalla, len(a.pila))
	}
}

func TestVolverConLaPilaVaciaVaAlMenu(t *testing.T) {
	m := conJornadas(t, 16, 2)
	m.pantalla = pantallaEstadisticas // sin pasar por el menu: la pila esta vacia
	m, _ = pulsar(m, "esc")
	contiene(t, m, "JUEGO DE FÚTBOL")
}

func TestRecortar(t *testing.T) {
	if recortar("corto", 10) != "corto" || recortar("exacto", 6) != "exacto" {
		t.Error("un texto que cabe no se recorta")
	}
	if got := recortar("Academia Altamira del Sur", 10); got != "Academia …" || len([]rune(got)) != 10 {
		t.Errorf("recortar = %q", got)
	}
	if got := recortar("Ñandú largo", 5); len([]rune(got)) != 5 {
		t.Errorf("recortar con acentos = %q", got)
	}
}

// anchoMaximo devuelve el ancho de la línea más ancha de la pantalla.
func anchoMaximo(m Modelo) int {
	ancho := 0
	for _, l := range strings.Split(texto(m), "\n") {
		ancho = max(ancho, lipgloss.Width(l))
	}
	return ancho
}

// valoresEsperados son las celdas numéricas que debe mostrar cada clasificación
// para un jugador, escritas aparte del código de la pantalla.
func valoresEsperados(c aplicacion.Criterio, e aplicacion.EstadisticaJugador) []string {
	pj := fmt.Sprint(e.Partidos)
	switch c {
	case aplicacion.PorGoles:
		return []string{pj, fmt.Sprint(e.Goles), fmt.Sprint(e.Asistencias), fmt.Sprint(e.Minutos)}
	case aplicacion.PorAsistencias:
		return []string{pj, fmt.Sprint(e.Asistencias), fmt.Sprint(e.Goles), fmt.Sprint(e.Minutos)}
	case aplicacion.PorTarjetas:
		return []string{pj, fmt.Sprint(e.Amarillas), fmt.Sprint(e.Rojas), fmt.Sprint(e.Amarillas + 3*e.Rojas)}
	case aplicacion.PorImbatidas:
		return []string{pj, fmt.Sprint(e.PorteriasImbatidas), fmt.Sprint(e.GolesEncajados)}
	case aplicacion.PorPocosGolesEncajados:
		return []string{pj, fmt.Sprint(e.GolesEncajados), fmt.Sprintf("%.2f", float64(e.GolesEncajados)/float64(e.Partidos))}
	}
	return []string{pj, fmt.Sprintf("%.1f", e.ValoracionMedia()), fmt.Sprint(e.Goles), fmt.Sprint(e.Asistencias)}
}

func TestCadaClasificacionMuestraLosValoresCorrectos(t *testing.T) {
	m := conJornadas(t, 22, 18)
	hub := abrirEstadisticas(m)
	for i, cl := range clasificaciones {
		esperadas, err := m.carrera.Clasificacion(cl.criterio, maxFilasClasificacion)
		if err != nil || len(esperadas) < 3 {
			t.Fatalf("%s: %d filas, %v", cl.titulo, len(esperadas), err)
		}
		pantalla := abrirOpcion(hub, i)
		var filas []string
		for _, l := range strings.Split(texto(pantalla), "\n") {
			// Las filas de datos empiezan con el marcador y el puesto.
			if len(l) > 4 && (l[0] == '>' || l[0] == '*' || l[0] == ' ') && strings.Contains(l[:6], fmt.Sprint(len(filas)+1)) {
				filas = append(filas, l)
			}
		}
		if len(filas) < 3 {
			t.Fatalf("%s: solo se reconocieron %d filas de datos:\n%s", cl.titulo, len(filas), texto(pantalla))
		}
		for k := 0; k < 3; k++ {
			want := valoresEsperados(cl.criterio, esperadas[k])
			celdas := strings.Fields(filas[k])
			got := celdas[len(celdas)-len(want):]
			if strings.Join(got, " ") != strings.Join(want, " ") {
				t.Errorf("%s, puesto %d (%s): celdas %v, se esperaban %v\n%s",
					cl.titulo, k+1, esperadas[k].Jugador.Nombre, got, want, filas[k])
			}
			if !strings.Contains(filas[k], esperadas[k].Jugador.Nombre) {
				t.Errorf("%s, puesto %d: deberia ser %s:\n%s", cl.titulo, k+1, esperadas[k].Jugador.Nombre, filas[k])
			}
		}
	}
}
