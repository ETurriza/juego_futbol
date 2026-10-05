package aplicacion

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/ETurriza/juego_futbol/internal/generador"
	"github.com/ETurriza/juego_futbol/internal/modelo"
)

// terminar juega todas las jornadas que falten de la temporada en curso.
func terminar(t *testing.T, c *Carrera) {
	t.Helper()
	for !c.Terminada() {
		if _, err := c.AvanzarJornada(); err != nil {
			t.Fatal(err)
		}
	}
}

// jugarTemporadas termina n temporadas y pasa a la siguiente tras cada una; al
// final queda al inicio de la temporada n+1.
func jugarTemporadas(t *testing.T, c *Carrera, n int) {
	t.Helper()
	for i := 0; i < n; i++ {
		terminar(t, c)
		if _, err := c.SiguienteTemporada(); err != nil {
			t.Fatal(err)
		}
	}
}

func idsDe(c *Carrera) map[int]modelo.Jugador {
	ids := map[int]modelo.Jugador{}
	for _, e := range c.Temporada.Equipos {
		for _, j := range e.Plantilla {
			ids[j.ID] = j
		}
	}
	return ids
}

func TestNuevaCarreraEmpiezaEnLaTemporadaUno(t *testing.T) {
	c := nuevaCarrera(t, 5, 10)
	if c.Numero != 1 || len(c.Historial) != 0 {
		t.Errorf("Numero = %d, historial = %d; se esperaba 1 y vacio", c.Numero, len(c.Historial))
	}
	if want := 10*generador.TamanoPlantilla + 1; c.ProximoID != want {
		t.Errorf("ProximoID = %d, se esperaba %d", c.ProximoID, want)
	}
	for id := range idsDe(c) {
		if id >= c.ProximoID {
			t.Fatalf("la ID %d no es menor que ProximoID %d", id, c.ProximoID)
		}
	}
}

func TestSiguienteTemporadaExigeQueLaActualTermine(t *testing.T) {
	c := nuevaCarrera(t, 5, 10)
	c.AvanzarJornada()
	antes := c.Exportar()
	if _, err := c.SiguienteTemporada(); !errors.Is(err, ErrTemporadaEnCurso) {
		t.Errorf("err = %v, se esperaba ErrTemporadaEnCurso", err)
	}
	if !reflect.DeepEqual(antes, c.Exportar()) {
		t.Error("un intento fallido no deberia cambiar la carrera")
	}
}

func TestSiguienteTemporadaReiniciaLaLiga(t *testing.T) {
	c := nuevaCarrera(t, 5, 10)
	terminar(t, c)
	tabla := c.Tabla()
	var puesto, puntos int
	for _, f := range tabla {
		if f.EsDelUsuario {
			puesto, puntos = f.Posicion, f.Pts
		}
	}
	club, usuario := c.NombreEquipo(), c.Usuario
	nombres := nombresDeEquipos(c)

	cambios, err := c.SiguienteTemporada()
	if err != nil {
		t.Fatal(err)
	}
	if cambios.Numero != 2 || c.Numero != 2 {
		t.Errorf("Numero = %d / %d, se esperaba 2", cambios.Numero, c.Numero)
	}
	want := []ResumenTemporada{{Numero: 1, Campeon: tabla[0].Equipo, PuestoUsuario: puesto, PuntosUsuario: puntos}}
	if !reflect.DeepEqual(c.Historial, want) {
		t.Errorf("Historial = %+v, se esperaba %+v", c.Historial, want)
	}
	if c.Terminada() || c.Jornada() != 0 || len(c.Temporada.Resultados) != 0 {
		t.Error("la temporada nueva deberia empezar sin jornadas jugadas")
	}
	if c.NombreEquipo() != club || c.Usuario != usuario {
		t.Error("el usuario deberia seguir en el mismo club")
	}
	if !reflect.DeepEqual(nombres, nombresDeEquipos(c)) {
		t.Error("los clubes de la liga deberian ser los mismos")
	}
	for _, f := range c.Tabla() {
		if f.PJ != 0 || f.Pts != 0 {
			t.Errorf("la tabla deberia empezar en cero: %+v", f)
		}
	}
	// Y se puede jugar la temporada nueva entera.
	terminar(t, c)
	if !c.Terminada() {
		t.Error("la temporada nueva deberia poder jugarse completa")
	}
}

func nombresDeEquipos(c *Carrera) []string {
	var nombres []string
	for _, e := range c.Temporada.Equipos {
		nombres = append(nombres, e.Nombre)
	}
	return nombres
}

func TestSiguienteTemporadaEnvejeceRetiraYRepone(t *testing.T) {
	c := nuevaCarrera(t, 7, 10)
	terminar(t, c)
	antes := idsDe(c)
	proximoAntes := c.ProximoID

	cambios, err := c.SiguienteTemporada()
	if err != nil {
		t.Fatal(err)
	}
	despues := idsDe(c)

	retirados, nuevos := 0, 0
	for id, j := range antes {
		d, sigue := despues[id]
		if !sigue {
			retirados++
			continue
		}
		if d.Edad != j.Edad+1 {
			t.Errorf("jugador %d: edad %d, se esperaba %d", id, d.Edad, j.Edad+1)
		}
		if d.Nombre != j.Nombre || d.Posicion != j.Posicion {
			t.Errorf("jugador %d cambio de identidad o posicion", id)
		}
	}
	for id, j := range despues {
		if _, existia := antes[id]; existia {
			continue
		}
		nuevos++
		if id < proximoAntes {
			t.Errorf("el juvenil %d reutiliza una ID anterior (ProximoID era %d)", id, proximoAntes)
		}
		if j.Edad < generador.EdadJuvenilMin || j.Edad > generador.EdadJuvenilMax {
			t.Errorf("juvenil %d con edad %d", id, j.Edad)
		}
	}
	if retirados != cambios.RetiradosLiga || nuevos != cambios.JuvenilesLiga {
		t.Errorf("cambios (%d retirados, %d juveniles) no coinciden con lo observado (%d, %d)",
			cambios.RetiradosLiga, cambios.JuvenilesLiga, retirados, nuevos)
	}
	if retirados != nuevos {
		t.Errorf("cada retiro se repone con un juvenil: %d retirados, %d juveniles", retirados, nuevos)
	}
	if c.ProximoID != proximoAntes+nuevos {
		t.Errorf("ProximoID = %d, se esperaba %d", c.ProximoID, proximoAntes+nuevos)
	}
	// Los cambios del usuario son los de su equipo.
	if len(cambios.Retirados) != len(cambios.Juveniles) {
		t.Errorf("en el club del usuario: %d retirados y %d juveniles", len(cambios.Retirados), len(cambios.Juveniles))
	}
	for _, j := range cambios.Retirados {
		if _, sigue := despues[j.ID]; sigue {
			t.Errorf("el retirado %d sigue en la liga", j.ID)
		}
	}
	for _, j := range cambios.Juveniles {
		if despues[j.ID].Nombre != j.Nombre {
			t.Errorf("el juvenil %d no esta en la plantilla del usuario", j.ID)
		}
	}
}

func TestLasPlantillasMantienenLaComposicionDurante10Temporadas(t *testing.T) {
	c := nuevaCarrera(t, 3, 10)
	vistas := map[int]bool{}
	for temporada := 1; temporada <= 10; temporada++ {
		terminar(t, c)
		retiradasAntes := idsDe(c)
		if _, err := c.SiguienteTemporada(); err != nil {
			t.Fatal(err)
		}
		despues := idsDe(c)
		for id := range retiradasAntes {
			if _, sigue := despues[id]; !sigue {
				vistas[id] = true // retirado
			}
		}
		for id := range despues {
			if vistas[id] {
				t.Fatalf("temporada %d: la ID %d de un retirado reaparece", temporada, id)
			}
		}
		for _, e := range c.Temporada.Equipos {
			if err := e.Validar(); err != nil {
				t.Fatalf("temporada %d: %v", temporada, err)
			}
			if len(e.Plantilla) != generador.TamanoPlantilla {
				t.Fatalf("temporada %d: %s con %d jugadores", temporada, e.Nombre, len(e.Plantilla))
			}
			for _, p := range modelo.Posiciones {
				if e.Contar(p) < generador.MinimoPorPosicion(p) {
					t.Fatalf("temporada %d: %s con %d %v", temporada, e.Nombre, e.Contar(p), p)
				}
			}
			nombres := map[string]bool{}
			for _, j := range e.Plantilla {
				if nombres[j.Nombre] {
					t.Fatalf("temporada %d: %s repite el nombre %q", temporada, e.Nombre, j.Nombre)
				}
				nombres[j.Nombre] = true
			}
		}
	}
	if c.Numero != 11 || len(c.Historial) != 10 {
		t.Errorf("Numero = %d, historial = %d; se esperaba 11 y 10", c.Numero, len(c.Historial))
	}
	for i, h := range c.Historial {
		if h.Numero != i+1 || h.Campeon == "" || h.PuestoUsuario < 1 || h.PuestoUsuario > 10 {
			t.Errorf("resumen %d incorrecto: %+v", i+1, h)
		}
	}
}

func TestVariasTemporadasSonReproducibles(t *testing.T) {
	a, b := nuevaCarrera(t, 12, 10), nuevaCarrera(t, 12, 10)
	jugarTemporadas(t, a, 4)
	jugarTemporadas(t, b, 4)
	if !reflect.DeepEqual(a, b) {
		t.Error("la misma semilla deberia dar la misma carrera tras varias temporadas")
	}
	c := nuevaCarrera(t, 13, 10)
	jugarTemporadas(t, c, 4)
	if reflect.DeepEqual(a.Historial, c.Historial) && reflect.DeepEqual(a.Temporada.Equipos, c.Temporada.Equipos) {
		t.Error("semillas distintas deberian dar carreras distintas")
	}
}

func TestLaPrimeraTemporadaConservaSusSemillas(t *testing.T) {
	c := nuevaCarrera(t, 42, 10)
	if c.semillaTemporada() != c.Semilla {
		t.Error("la temporada 1 deberia usar la semilla de la carrera tal cual")
	}
	vistas := map[int64]bool{c.Semilla: true}
	for n := 2; n <= 50; n++ {
		c.Numero = n
		s := c.semillaTemporada()
		if vistas[s] {
			t.Fatalf("la temporada %d repite semilla", n)
		}
		vistas[s] = true
	}
	// La semilla de evolución también es distinta por temporada.
	if semillaEvolucion(42, 1) == semillaEvolucion(42, 2) || semillaEvolucion(42, 1) == semillaEvolucion(43, 1) {
		t.Error("las semillas de evolucion deberian variar con la temporada y la carrera")
	}
}

func TestTemporadaGuardadaAMitadSeReanudaIgual(t *testing.T) {
	continua := nuevaCarrera(t, 21, 10)
	jugarTemporadas(t, continua, 2)
	for i := 0; i < 5; i++ {
		continua.AvanzarJornada()
	}

	reanudada, err := Importar(continua.Exportar())
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(continua, reanudada) {
		t.Fatal("la carrera importada deberia ser identica")
	}
	jugarTemporadas(t, continua, 2) // termina la 3 y la 4
	jugarTemporadas(t, reanudada, 2)
	if !reflect.DeepEqual(continua, reanudada) {
		t.Error("una carrera reanudada a mitad de temporada deberia dar lo mismo que una continua")
	}
}

func TestExportarImportarVariasTemporadas(t *testing.T) {
	c := nuevaCarrera(t, 8, 10)
	jugarTemporadas(t, c, 3)
	for i := 0; i < 4; i++ {
		c.AvanzarJornada()
	}
	g := c.Exportar()
	if g.Numero != 4 || len(g.Historial) != 3 || g.ProximoID != c.ProximoID {
		t.Fatalf("Guardado = temporada %d, historial %d, ProximoID %d", g.Numero, len(g.Historial), g.ProximoID)
	}
	reconstruida, err := Importar(g)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(c, reconstruida) {
		t.Error("la carrera reconstruida no es identica")
	}
	// No comparten memoria.
	g.Historial[0].Campeon = "Alterado"
	if c.Historial[0].Campeon == "Alterado" {
		t.Error("Exportar comparte el historial con la carrera")
	}
	resumen, _ := c.Exportar().Resumen("x")
	if resumen.Temporada != 4 || resumen.Jornada != 4 {
		t.Errorf("Resumen = %+v", resumen)
	}
}

// TestEstabilidadDeLaLigaAlLargoPlazo comprueba que la liga no se infla ni se
// degrada con los años: la valoración y la edad medias se estabilizan. Con -v
// imprime la evolución para poder calibrar las constantes.
func TestEstabilidadDeLaLigaAlLargoPlazo(t *testing.T) {
	const temporadas = 14
	for _, semilla := range []int64{1, 2, 3} {
		c := nuevaCarrera(t, semilla, 10)
		var valoracion, edad []float64
		var tabla strings.Builder
		fmt.Fprintf(&tabla, "\nsemilla %d: evolucion de la liga\n%-10s %12s %10s %10s\n", semilla, "temporada", "valoracion", "edad", ">=34 anos")
		for n := 1; n <= temporadas; n++ {
			v, e, veteranos := estadisticasDeLiga(c)
			valoracion, edad = append(valoracion, v), append(edad, e)
			fmt.Fprintf(&tabla, "%-10d %12.1f %10.1f %10d\n", n, v, e, veteranos)
			terminar(t, c)
			if _, err := c.SiguienteTemporada(); err != nil {
				t.Fatal(err)
			}
		}
		t.Log(tabla.String())

		// Una vez estabilizada (de la temporada 6 en adelante), casi no se mueve.
		estable := valoracion[5:]
		minV, maxV := estable[0], estable[0]
		for _, v := range estable {
			minV, maxV = min(minV, v), max(maxV, v)
		}
		if maxV-minV > 4 {
			t.Errorf("semilla %d: la valoracion media oscila %.1f entre las temporadas 6 y %d", semilla, maxV-minV, temporadas)
		}
		for n, v := range valoracion {
			if v < 55 || v > 75 {
				t.Errorf("semilla %d temporada %d: valoracion media %.1f fuera de [55, 75]", semilla, n+1, v)
			}
		}
		// No debe haber un salto brusco respecto de la liga inicial.
		if d := valoracion[temporadas-1] - valoracion[0]; d > 5 || d < -5 {
			t.Errorf("semilla %d: la valoracion media cambio %.1f puntos en %d temporadas", semilla, d, temporadas)
		}
		for n, e := range edad[5:] {
			if e < 24 || e > 29.5 {
				t.Errorf("semilla %d temporada %d: edad media %.1f fuera de [24, 29.5]", semilla, n+6, e)
			}
		}
	}
}

// estadisticasDeLiga devuelve la valoración media y la edad media de todos los
// jugadores de la liga, y cuántos tienen 34 años o más.
func estadisticasDeLiga(c *Carrera) (valoracion, edad float64, veteranos int) {
	n := 0
	for _, e := range c.Temporada.Equipos {
		for _, j := range e.Plantilla {
			valoracion += float64(j.Valoracion())
			edad += float64(j.Edad)
			if j.Edad >= 34 {
				veteranos++
			}
			n++
		}
	}
	return valoracion / float64(n), edad / float64(n), veteranos
}
