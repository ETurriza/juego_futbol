package aplicacion

import (
	"reflect"
	"sort"
	"strings"
	"testing"

	"github.com/ETurriza/juego_futbol/internal/modelo"
)

func carreraConJornadas(t *testing.T, semilla int64, jornadas int) *Carrera {
	t.Helper()
	c := nuevaCarrera(t, semilla, 10)
	for i := 0; i < jornadas; i++ {
		if _, err := c.AvanzarJornada(); err != nil {
			t.Fatal(err)
		}
	}
	return c
}

func golesDeLaLiga(c *Carrera) int {
	total := 0
	for _, jornada := range c.Temporada.Resultados {
		for _, r := range jornada {
			total += r.GolesLocal + r.GolesVisitante
		}
	}
	return total
}

func TestEstadisticasJugadoresCubreATodaLaLiga(t *testing.T) {
	c := carreraConJornadas(t, 5, 0)
	filas, err := c.EstadisticasJugadores()
	if err != nil {
		t.Fatal(err)
	}
	if len(filas) != 10*22 {
		t.Fatalf("%d filas, se esperaban 220", len(filas))
	}
	for _, f := range filas {
		if f.Estadisticas != (modelo.Estadisticas{}) {
			t.Fatalf("antes de jugar todo deberia estar en cero: %+v", f)
		}
	}

	c = carreraConJornadas(t, 5, 18)
	filas, err = c.EstadisticasJugadores()
	if err != nil {
		t.Fatal(err)
	}
	goles, partidos := 0, 0
	for _, f := range filas {
		goles += f.Goles
		partidos += f.Partidos
		if f.Equipo == "" || f.Jugador.Nombre == "" {
			t.Fatalf("fila sin equipo o nombre: %+v", f)
		}
	}
	if goles != golesDeLaLiga(c) {
		t.Errorf("los jugadores suman %d goles y la liga %d", goles, golesDeLaLiga(c))
	}
	// 90 partidos con 22 titulares y 6 a 10 suplentes cada uno.
	if partidos < 90*22 || partidos > 90*(22+10) {
		t.Errorf("%d participaciones en 90 partidos", partidos)
	}
}

func TestEstadisticasDeEquiposSumanLosJugadores(t *testing.T) {
	c := carreraConJornadas(t, 6, 18)
	equipos, err := c.EstadisticasEquipos()
	if err != nil {
		t.Fatal(err)
	}
	jugadores, _ := c.EstadisticasJugadores()
	golesPorEquipo := map[string]int{}
	topGoles := map[string]int{}
	topAsist := map[string]int{}
	for _, j := range jugadores {
		golesPorEquipo[j.Equipo] += j.Goles
		topGoles[j.Equipo] = max(topGoles[j.Equipo], j.Goles)
		topAsist[j.Equipo] = max(topAsist[j.Equipo], j.Asistencias)
	}
	tabla := map[string]int{}
	for _, f := range c.Tabla() {
		tabla[f.Equipo] = f.GF
	}
	for _, e := range equipos {
		if e.PJ != 18 || e.G+e.E+e.P != e.PJ {
			t.Errorf("%s: PJ %d, G+E+P %d", e.Equipo, e.PJ, e.G+e.E+e.P)
		}
		if e.GF != golesPorEquipo[e.Equipo] || e.GF != tabla[e.Equipo] {
			t.Errorf("%s: GF %d, goles de jugadores %d, tabla %d", e.Equipo, e.GF, golesPorEquipo[e.Equipo], tabla[e.Equipo])
		}
		if e.GolesGoleador != topGoles[e.Equipo] || e.Asistencias != topAsist[e.Equipo] {
			t.Errorf("%s: goleador %d (esperado %d), asistente %d (esperado %d)",
				e.Equipo, e.GolesGoleador, topGoles[e.Equipo], e.Asistencias, topAsist[e.Equipo])
		}
		if e.GolesGoleador > 0 && e.Goleador == "" || e.Asistencias > 0 && e.Asistente == "" {
			t.Errorf("%s: falta el nombre del goleador o del asistente", e.Equipo)
		}
	}
}

func TestClasificacionDeGoleadores(t *testing.T) {
	c := carreraConJornadas(t, 7, 18)
	top, err := c.Clasificacion(PorGoles, 10)
	if err != nil {
		t.Fatal(err)
	}
	if len(top) != 10 {
		t.Fatalf("%d goleadores, se esperaban 10", len(top))
	}
	for i, f := range top {
		if f.Goles == 0 {
			t.Errorf("el puesto %d no tiene goles", i+1)
		}
		if i > 0 && f.Goles > top[i-1].Goles {
			t.Errorf("goleadores fuera de orden en el puesto %d", i+1)
		}
	}
	// Es de verdad el mejor de la liga.
	todos, _ := c.EstadisticasJugadores()
	mejor := 0
	for _, j := range todos {
		mejor = max(mejor, j.Goles)
	}
	if top[0].Goles != mejor {
		t.Errorf("el goleador tiene %d goles y el maximo de la liga es %d", top[0].Goles, mejor)
	}

	// Sin limite devuelve a todos los que marcaron; el orden es siempre el mismo.
	todosLosQueMarcaron, _ := c.Clasificacion(PorGoles, 0)
	otra, _ := c.Clasificacion(PorGoles, 0)
	if !reflect.DeepEqual(todosLosQueMarcaron, otra) {
		t.Error("el orden de la clasificacion deberia ser determinista")
	}
	marcaron := 0
	for _, j := range todos {
		if j.Goles > 0 {
			marcaron++
		}
	}
	if len(todosLosQueMarcaron) != marcaron {
		t.Errorf("%d goleadores, %d jugadores marcaron", len(todosLosQueMarcaron), marcaron)
	}
}

func TestClasificacionesVaciasAntesDeJugar(t *testing.T) {
	c := carreraConJornadas(t, 7, 0)
	for _, crit := range []Criterio{PorGoles, PorAsistencias, PorTarjetas, PorImbatidas, PorPocosGolesEncajados, PorValoracion} {
		filas, err := c.Clasificacion(crit, 5)
		if err != nil || len(filas) != 0 {
			t.Errorf("criterio %d: %d filas, %v; se esperaba vacio", crit, len(filas), err)
		}
	}
	if _, err := c.Clasificacion(Criterio(99), 5); err == nil {
		t.Error("un criterio desconocido deberia dar error")
	}
}

func TestClasificacionDeAsistentesYTarjetas(t *testing.T) {
	c := carreraConJornadas(t, 8, 18)
	asist, _ := c.Clasificacion(PorAsistencias, 0)
	for i := 1; i < len(asist); i++ {
		if asist[i].Asistencias > asist[i-1].Asistencias || asist[i].Asistencias == 0 {
			t.Fatalf("asistentes mal ordenados en el puesto %d", i+1)
		}
	}
	tarjetas, _ := c.Clasificacion(PorTarjetas, 0)
	if len(tarjetas) == 0 {
		t.Fatal("tras una temporada deberia haber tarjetas")
	}
	for i := 1; i < len(tarjetas); i++ {
		a, b := puntosDeDisciplina(tarjetas[i-1].Estadisticas), puntosDeDisciplina(tarjetas[i].Estadisticas)
		if b > a || b == 0 {
			t.Fatalf("tarjetas mal ordenadas en el puesto %d: %d tras %d", i+1, b, a)
		}
	}
}

func TestClasificacionesDePorteros(t *testing.T) {
	c := carreraConJornadas(t, 9, 18)
	imbatidas, _ := c.Clasificacion(PorImbatidas, 0)
	if len(imbatidas) == 0 {
		t.Fatal("deberia haber porteros con porterias imbatidas")
	}
	for i, f := range imbatidas {
		if f.Jugador.Posicion != modelo.Portero || f.PorteriasImbatidas == 0 {
			t.Fatalf("puesto %d: %v con %d imbatidas", i+1, f.Jugador.Posicion, f.PorteriasImbatidas)
		}
		if i > 0 && f.PorteriasImbatidas > imbatidas[i-1].PorteriasImbatidas {
			t.Fatalf("imbatidas mal ordenadas en el puesto %d", i+1)
		}
	}

	zamora, _ := c.Clasificacion(PorPocosGolesEncajados, 0)
	if len(zamora) == 0 {
		t.Fatal("deberia haber porteros en la clasificacion de goles encajados")
	}
	for i, f := range zamora {
		if f.Jugador.Posicion != modelo.Portero || f.Partidos < MinPartidosClasificacion {
			t.Fatalf("puesto %d: %v con %d partidos", i+1, f.Jugador.Posicion, f.Partidos)
		}
		if i > 0 {
			ant := zamora[i-1]
			// goles por partido, sin dividir
			if f.GolesEncajados*ant.Partidos < ant.GolesEncajados*f.Partidos {
				t.Fatalf("el puesto %d encaja menos por partido que el anterior", i+1)
			}
		}
	}
}

func TestClasificacionPorValoracionExigePartidos(t *testing.T) {
	c := carreraConJornadas(t, 10, 18)
	val, _ := c.Clasificacion(PorValoracion, 0)
	if len(val) == 0 {
		t.Fatal("deberia haber valoraciones")
	}
	for i, f := range val {
		if f.Partidos < MinPartidosClasificacion {
			t.Fatalf("puesto %d con solo %d partidos", i+1, f.Partidos)
		}
		if i > 0 && f.ValoracionMedia() > val[i-1].ValoracionMedia()+1e-9 {
			t.Fatalf("valoraciones mal ordenadas en el puesto %d", i+1)
		}
	}
	// Al principio de la temporada el umbral baja a las jornadas jugadas.
	c = carreraConJornadas(t, 10, 1)
	val, _ = c.Clasificacion(PorValoracion, 0)
	if len(val) == 0 {
		t.Error("tras una jornada deberia haber valoraciones de quienes jugaron")
	}
}

func TestSiguienteTemporadaArchivaLasEstadisticas(t *testing.T) {
	c := carreraConJornadas(t, 11, 18)
	golesTemporada := golesDeLaLiga(c)
	antes, _ := c.EstadisticasJugadores()
	jugaron := 0
	for _, j := range antes {
		if j.Partidos > 0 {
			jugaron++
		}
	}

	if _, err := c.SiguienteTemporada(); err != nil {
		t.Fatal(err)
	}
	if len(c.Archivo) != jugaron {
		t.Errorf("%d filas archivadas, jugaron %d", len(c.Archivo), jugaron)
	}
	goles := 0
	for _, a := range c.Archivo {
		goles += a.Goles
		if a.Temporada != 1 || a.Nombre == "" || a.Equipo == "" || a.Partidos < 1 || a.Edad < modelo.EdadMin {
			t.Fatalf("fila archivada incorrecta: %+v", a)
		}
	}
	if goles != golesTemporada {
		t.Errorf("el archivo suma %d goles y la temporada tuvo %d", goles, golesTemporada)
	}

	// La temporada nueva empieza en cero.
	nuevos, _ := c.EstadisticasJugadores()
	for _, j := range nuevos {
		if j.Estadisticas != (modelo.Estadisticas{}) {
			t.Fatalf("la temporada nueva deberia empezar en cero: %+v", j)
		}
	}
}

func TestEstadisticasDeCarreraSumanTemporadas(t *testing.T) {
	c := carreraConJornadas(t, 12, 18)
	primera, _ := c.EstadisticasJugadores()
	// Un delantero que haya marcado en la primera temporada.
	var elegido modelo.Jugador
	golesPrimera := 0
	for _, j := range primera {
		if j.Goles > golesPrimera && j.Jugador.Edad < 33 {
			elegido, golesPrimera = j.Jugador, j.Goles
		}
	}
	if golesPrimera == 0 {
		t.Fatal("nadie marco en la primera temporada")
	}
	if _, err := c.SiguienteTemporada(); err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 10; i++ {
		c.AvanzarJornada()
	}
	actual, _ := c.EstadisticasJugadores()
	golesActual := 0
	for _, j := range actual {
		if j.Jugador.ID == elegido.ID {
			golesActual = j.Goles
		}
	}
	carrera, err := c.EstadisticasDeCarrera(elegido.ID)
	if err != nil {
		t.Fatal(err)
	}
	if carrera.Goles != golesPrimera+golesActual {
		t.Errorf("goles de carrera %d, se esperaban %d + %d", carrera.Goles, golesPrimera, golesActual)
	}
	if carrera.Partidos < 18 {
		t.Errorf("%d partidos en carrera, deberia haber jugado al menos la mayoria de la primera temporada", carrera.Partidos)
	}
	if ninguno, _ := c.EstadisticasDeCarrera(-5); ninguno != (modelo.Estadisticas{}) {
		t.Error("un jugador que no existe no tiene estadisticas")
	}
}

func TestLasEstadisticasSobrevivenAGuardarYReanudar(t *testing.T) {
	continua := nuevaCarrera(t, 13, 10)
	jugarTemporadas(t, continua, 2)
	for i := 0; i < 7; i++ {
		continua.AvanzarJornada()
	}
	reanudada, err := Importar(continua.Exportar())
	if err != nil {
		t.Fatal(err)
	}
	a, _ := continua.EstadisticasJugadores()
	b, _ := reanudada.EstadisticasJugadores()
	if !reflect.DeepEqual(a, b) {
		t.Error("las estadisticas de la temporada en curso deberian ser iguales tras importar")
	}
	if !reflect.DeepEqual(continua.Archivo, reanudada.Archivo) || len(continua.Archivo) == 0 {
		t.Error("el archivo deberia ser identico y no estar vacio")
	}
	terminar(t, continua)
	terminar(t, reanudada)
	a, _ = continua.EstadisticasJugadores()
	b, _ = reanudada.EstadisticasJugadores()
	if !reflect.DeepEqual(a, b) {
		t.Error("al terminar la temporada las estadisticas deberian coincidir")
	}
}

func TestImportarRechazaDetallesYArchivosIncoherentes(t *testing.T) {
	casos := map[string]func(*Guardado){
		"sucesos que no suman el marcador": func(g *Guardado) { g.Resultados[0][0].GolesLocal += 3 },
		"jugador de otro equipo": func(g *Guardado) {
			g.Resultados[0][0].Detalle.TitularesLocal[3] = g.Resultados[0][0].Detalle.TitularesVisitante[3]
		},
		"titular que no existe": func(g *Guardado) { g.Resultados[0][0].Detalle.TitularesLocal[2] = 99999 },
		"sucesos fuera de orden": func(g *Guardado) {
			ev := g.Resultados[0][0].Detalle.Eventos
			ev[0], ev[len(ev)-1] = ev[len(ev)-1], ev[0]
		},
		"minuto imposible": func(g *Guardado) { g.Resultados[0][0].Detalle.Eventos[0].Minuto = 0 },
		"archivo de una temporada sin terminar": func(g *Guardado) {
			g.Archivo = append(g.Archivo, EstadisticaTemporada{Temporada: g.Numero, Jugador: 1, Estadisticas: modelo.Estadisticas{Partidos: 1}})
		},
		"archivo sin jugador": func(g *Guardado) {
			g.Numero, g.Historial = 2, []ResumenTemporada{{Numero: 1}}
			g.Archivo = []EstadisticaTemporada{{Temporada: 1, Estadisticas: modelo.Estadisticas{Partidos: 1}}}
		},
		"archivo sin partidos": func(g *Guardado) {
			g.Numero, g.Historial = 2, []ResumenTemporada{{Numero: 1}}
			g.Archivo = []EstadisticaTemporada{{Temporada: 1, Jugador: 1}}
		},
	}
	for nombre, mutar := range casos {
		g := guardadoDePrueba(t, 3)
		mutar(&g)
		if _, err := Importar(g); err == nil {
			t.Errorf("%s: deberia dar error", nombre)
		} else if !strings.Contains(err.Error(), "guardado invalido") {
			t.Errorf("%s: el error deberia indicar un guardado invalido: %v", nombre, err)
		}
	}
}

// Las partidas guardadas antes de las estadísticas no tienen detalle: se
// importan, cuentan en la tabla y se puede seguir jugando; solo faltan sus
// estadísticas.
func TestImportarAceptaPartidosSinDetalle(t *testing.T) {
	g := guardadoDePrueba(t, 4)
	for _, jornada := range g.Resultados {
		for k := range jornada {
			jornada[k].Detalle = modelo.DetallePartido{}
		}
	}
	c, err := Importar(g)
	if err != nil {
		t.Fatalf("un guardado antiguo deberia importarse: %v", err)
	}
	if len(c.Tabla()) != 10 || c.Tabla()[0].PJ == 0 {
		t.Error("los partidos sin detalle deberian contar en la tabla")
	}
	filas, err := c.EstadisticasJugadores()
	if err != nil {
		t.Fatal(err)
	}
	for _, f := range filas {
		if f.Partidos != 0 {
			t.Fatalf("sin detalle no deberia haber estadisticas: %+v", f)
		}
	}
	// Se puede seguir: los partidos nuevos si traen detalle y sus estadisticas cuentan.
	terminar(t, c)
	filas, _ = c.EstadisticasJugadores()
	jugaron := 0
	for _, f := range filas {
		if f.Partidos > 0 {
			jugaron++
		}
	}
	if jugaron == 0 {
		t.Error("tras seguir jugando deberia haber estadisticas de los partidos nuevos")
	}
	if _, err := c.SiguienteTemporada(); err != nil {
		t.Errorf("deberia poder pasar de temporada: %v", err)
	}
}

func TestOrdenDeLasClasificacionesEsTotal(t *testing.T) {
	// Sin importar el orden de entrada, el resultado es el mismo: se ordena por
	// criterios y al final por nombre.
	c := carreraConJornadas(t, 14, 18)
	a, _ := c.Clasificacion(PorGoles, 0)
	nombres := func(f []EstadisticaJugador) []string {
		var out []string
		for _, x := range f {
			out = append(out, x.Jugador.Nombre)
		}
		return out
	}
	// Invertir las plantillas no cambia la clasificacion.
	for i := range c.Temporada.Equipos {
		p := c.Temporada.Equipos[i].Plantilla
		sort.SliceStable(p, func(x, y int) bool { return p[x].ID > p[y].ID })
	}
	b, _ := c.Clasificacion(PorGoles, 0)
	if !reflect.DeepEqual(nombres(a), nombres(b)) {
		t.Error("la clasificacion no deberia depender del orden de las plantillas")
	}
}
