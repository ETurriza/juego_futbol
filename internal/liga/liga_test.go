package liga

import (
	"errors"
	"fmt"
	"math/rand"
	"reflect"
	"testing"

	"github.com/ETurriza/juego_futbol/internal/modelo"
	"github.com/ETurriza/juego_futbol/internal/simulacion"
)

func nuevoRand(semilla int64) *rand.Rand {
	return rand.New(rand.NewSource(semilla))
}

// equipoUniforme arma una plantilla de 22 con todos los atributos iguales a v.
func equipoUniforme(nombre string, v int) modelo.Equipo {
	e := modelo.Equipo{Nombre: nombre}
	id := 1
	for _, c := range []struct {
		p modelo.Posicion
		n int
	}{
		{modelo.Portero, 2}, {modelo.Defensa, 7},
		{modelo.Mediocampista, 7}, {modelo.Delantero, 6},
	} {
		for i := 0; i < c.n; i++ {
			e.Plantilla = append(e.Plantilla, modelo.Jugador{
				ID: id, Nombre: "J", Edad: 25, Posicion: c.p,
				Atributos: modelo.Atributos{
					Ritmo: v, Tiro: v, Pase: v, Regate: v, Defensa: v, Fisico: v, Reflejos: v,
				},
			})
			id++
		}
	}
	return e
}

func equiposUniformes(n, v int) []modelo.Equipo {
	equipos := make([]modelo.Equipo, n)
	for i := range equipos {
		equipos[i] = equipoUniforme(fmt.Sprintf("Equipo %02d", i), v)
	}
	return equipos
}

func TestCalendarioEstructura(t *testing.T) {
	for _, n := range []int{2, 3, 4, 5, 6, 7, 8, 11, 20} {
		cal, err := Calendario(n)
		if err != nil {
			t.Fatalf("n=%d: %v", n, err)
		}

		m := n + n%2
		if want := 2 * (m - 1); len(cal) != want {
			t.Errorf("n=%d: %d jornadas, se esperaban %d", n, len(cal), want)
		}

		pares := map[Partido]int{}
		jugados := make([]int, n)
		locales := make([]int, n)
		for i, jornada := range cal {
			enJornada := map[int]bool{}
			for _, p := range jornada {
				if p.Local == p.Visitante {
					t.Errorf("n=%d jornada %d: %d juega contra si mismo", n, i, p.Local)
				}
				if p.Local < 0 || p.Local >= n || p.Visitante < 0 || p.Visitante >= n {
					t.Fatalf("n=%d jornada %d: indice fuera de rango: %+v", n, i, p)
				}
				for _, e := range []int{p.Local, p.Visitante} {
					if enJornada[e] {
						t.Errorf("n=%d jornada %d: equipo %d juega dos veces", n, i, e)
					}
					enJornada[e] = true
					jugados[e]++
				}
				locales[p.Local]++
				pares[p]++
			}
			// Con n par juegan todos; con n impar descansa exactamente uno.
			if want := n / 2; len(jornada) != want {
				t.Errorf("n=%d jornada %d: %d partidos, se esperaban %d", n, i, len(jornada), want)
			}
		}

		// Cada par ordenado (local, visitante) aparece exactamente una vez.
		if len(pares) != n*(n-1) {
			t.Errorf("n=%d: %d pares distintos, se esperaban %d", n, len(pares), n*(n-1))
		}
		for p, veces := range pares {
			if veces != 1 {
				t.Errorf("n=%d: partido %+v aparece %d veces", n, p, veces)
			}
		}
		for e := 0; e < n; e++ {
			if jugados[e] != 2*(n-1) {
				t.Errorf("n=%d: equipo %d juega %d partidos, se esperaban %d", n, e, jugados[e], 2*(n-1))
			}
			if locales[e] != n-1 {
				t.Errorf("n=%d: equipo %d es local %d veces, se esperaban %d", n, e, locales[e], n-1)
			}
		}
	}
}

func TestCalendarioSegundaVueltaInvierteLocalias(t *testing.T) {
	cal, _ := Calendario(6)
	mitad := len(cal) / 2
	for i := 0; i < mitad; i++ {
		for k, p := range cal[i] {
			v := cal[mitad+i][k]
			if v.Local != p.Visitante || v.Visitante != p.Local {
				t.Errorf("jornada %d partido %d: la vuelta %+v no invierte %+v", i, k, v, p)
			}
		}
	}
}

func TestCalendarioDeterminista(t *testing.T) {
	a, _ := Calendario(7)
	b, _ := Calendario(7)
	if !reflect.DeepEqual(a, b) {
		t.Error("el calendario deberia ser determinista")
	}
}

func TestCalendarioMinimo(t *testing.T) {
	for _, n := range []int{-1, 0, 1} {
		if _, err := Calendario(n); err == nil {
			t.Errorf("n=%d deberia dar error", n)
		}
	}
}

func TestNuevaValidaEquipos(t *testing.T) {
	if _, err := Nueva(nil); err == nil {
		t.Error("sin equipos deberia dar error")
	}
	if _, err := Nueva([]modelo.Equipo{{Nombre: "Solo"}}); err == nil {
		t.Error("con un equipo deberia dar error")
	}
	if _, err := Nueva([]modelo.Equipo{{Nombre: "A"}, {Nombre: "A"}}); err == nil {
		t.Error("nombres repetidos deberian dar error")
	}
	if _, err := Nueva([]modelo.Equipo{{Nombre: "A"}, {Nombre: ""}}); err == nil {
		t.Error("nombre vacio deberia dar error")
	}
	if _, err := Nueva([]modelo.Equipo{{Nombre: "A"}, {Nombre: "B"}}); err != nil {
		t.Errorf("dos equipos validos rechazados: %v", err)
	}
}

func TestNuevaCopiaLaListaDeEquipos(t *testing.T) {
	equipos := []modelo.Equipo{{Nombre: "A"}, {Nombre: "B"}}
	temp, _ := Nueva(equipos)
	equipos[0].Nombre = "Z"
	if temp.Equipos[0].Nombre != "A" {
		t.Error("la temporada no deberia verse afectada por cambios en la lista original")
	}
}

func tablaDe(nombres []string, resultados ...Resultado) []Fila {
	equipos := make([]modelo.Equipo, len(nombres))
	for i, n := range nombres {
		equipos[i] = modelo.Equipo{Nombre: n}
	}
	return (&Temporada{Equipos: equipos, Resultados: [][]Resultado{resultados}}).Tabla()
}

func res(local, visitante, gl, gv int) Resultado {
	return Resultado{Partido: Partido{Local: local, Visitante: visitante}, GolesLocal: gl, GolesVisitante: gv}
}

func TestTablaConResultadosConocidos(t *testing.T) {
	// A 2-0 B, B 1-1 C, C 0-3 A
	tabla := tablaDe([]string{"A", "B", "C"}, res(0, 1, 2, 0), res(1, 2, 1, 1), res(2, 0, 0, 3))
	want := []Fila{
		{Equipo: "A", PJ: 2, G: 2, GF: 5, GC: 0, DG: 5, Pts: 6},
		{Equipo: "B", PJ: 2, E: 1, P: 1, GF: 1, GC: 3, DG: -2, Pts: 1},
		{Equipo: "C", PJ: 2, E: 1, P: 1, GF: 1, GC: 4, DG: -3, Pts: 1},
	}
	if !reflect.DeepEqual(tabla, want) {
		t.Errorf("tabla = %+v\nesperada = %+v", tabla, want)
	}
}

func TestTablaDesempates(t *testing.T) {
	// Mismos puntos: gana la mayor DG.
	tabla := tablaDe([]string{"A", "B", "C", "D"}, res(0, 2, 1, 0), res(1, 3, 3, 0))
	if tabla[0].Equipo != "B" || tabla[1].Equipo != "A" {
		t.Errorf("deberia desempatar por DG: %+v", tabla)
	}
	// Mismos puntos y DG: gana el que mas goles hizo.
	tabla = tablaDe([]string{"A", "B", "C", "D"}, res(0, 2, 1, 0), res(1, 3, 3, 2))
	if tabla[0].Equipo != "B" || tabla[1].Equipo != "A" {
		t.Errorf("deberia desempatar por GF: %+v", tabla)
	}
	// Todo igual: orden alfabetico.
	tabla = tablaDe([]string{"Zeta", "Alfa"}, res(0, 1, 1, 1))
	if tabla[0].Equipo != "Alfa" || tabla[1].Equipo != "Zeta" {
		t.Errorf("deberia desempatar por nombre: %+v", tabla)
	}
}

func TestTablaInicial(t *testing.T) {
	temp, _ := Nueva(equiposUniformes(4, 60))
	tabla := temp.Tabla()
	if len(tabla) != 4 {
		t.Fatalf("la tabla tiene %d filas", len(tabla))
	}
	for _, f := range tabla {
		if f.PJ != 0 || f.Pts != 0 {
			t.Errorf("fila inicial no vacia: %+v", f)
		}
	}
}

func TestTemporadaCompleta(t *testing.T) {
	for _, n := range []int{4, 5, 20} {
		temp, err := Nueva(equiposUniformes(n, 60))
		if err != nil {
			t.Fatal(err)
		}
		if err := temp.JugarTemporada(nuevoRand(1)); err != nil {
			t.Fatal(err)
		}
		if !temp.Terminada() || temp.JornadaActual() != len(temp.Calendario) {
			t.Errorf("n=%d: la temporada deberia estar terminada", n)
		}

		var gf, gc, g, e, p, pts int
		for _, f := range temp.Tabla() {
			if f.PJ != 2*(n-1) {
				t.Errorf("n=%d: %s jugo %d partidos, se esperaban %d", n, f.Equipo, f.PJ, 2*(n-1))
			}
			if f.G+f.E+f.P != f.PJ {
				t.Errorf("n=%d: %s: G+E+P != PJ: %+v", n, f.Equipo, f)
			}
			if f.Pts != PuntosVictoria*f.G+PuntosEmpate*f.E || f.DG != f.GF-f.GC {
				t.Errorf("n=%d: %s: puntos o DG incoherentes: %+v", n, f.Equipo, f)
			}
			gf, gc, g, e, p, pts = gf+f.GF, gc+f.GC, g+f.G, e+f.E, p+f.P, pts+f.Pts
		}
		if gf != gc {
			t.Errorf("n=%d: goles a favor %d != en contra %d", n, gf, gc)
		}
		if g != p {
			t.Errorf("n=%d: victorias %d != derrotas %d", n, g, p)
		}
		// Cada empate reparte 2 puntos y cada victoria 3.
		if want := PuntosVictoria*g + 2*PuntosEmpate*(e/2); pts != want {
			t.Errorf("n=%d: puntos totales %d, se esperaban %d", n, pts, want)
		}
	}
}

func TestTemporadaReproducible(t *testing.T) {
	jugar := func(semilla int64) []Fila {
		temp, _ := Nueva(equiposUniformes(6, 60))
		if err := temp.JugarTemporada(nuevoRand(semilla)); err != nil {
			t.Fatal(err)
		}
		return temp.Tabla()
	}
	if !reflect.DeepEqual(jugar(7), jugar(7)) {
		t.Error("la misma semilla deberia dar la misma tabla")
	}
	if reflect.DeepEqual(jugar(7), jugar(8)) {
		t.Error("semillas distintas deberian dar tablas distintas")
	}
}

func TestJugarJornadaAvanzaYTermina(t *testing.T) {
	temp, _ := Nueva(equiposUniformes(2, 60))
	r := nuevoRand(1)
	for i := 0; i < 2; i++ {
		resultados, err := temp.JugarJornada(r)
		if err != nil {
			t.Fatalf("jornada %d: %v", i, err)
		}
		if len(resultados) != 1 || temp.JornadaActual() != i+1 {
			t.Errorf("jornada %d: %d resultados, jornada actual %d", i, len(resultados), temp.JornadaActual())
		}
	}
	if _, err := temp.JugarJornada(r); !errors.Is(err, ErrTemporadaTerminada) {
		t.Errorf("err = %v, se esperaba ErrTemporadaTerminada", err)
	}
	// Jugar la temporada ya terminada no hace nada ni falla.
	if err := temp.JugarTemporada(r); err != nil {
		t.Errorf("JugarTemporada sobre temporada terminada: %v", err)
	}
}

func TestResultadosSiguenElCalendario(t *testing.T) {
	temp, _ := Nueva(equiposUniformes(5, 60))
	if err := temp.JugarTemporada(nuevoRand(3)); err != nil {
		t.Fatal(err)
	}
	for i, jornada := range temp.Calendario {
		if len(temp.Resultados[i]) != len(jornada) {
			t.Fatalf("jornada %d: %d resultados para %d partidos", i, len(temp.Resultados[i]), len(jornada))
		}
		for k, p := range jornada {
			if temp.Resultados[i][k].Partido != p {
				t.Errorf("jornada %d partido %d: resultado %+v no corresponde a %+v", i, k, temp.Resultados[i][k].Partido, p)
			}
		}
	}
}

func TestMejorEquipoTerminaArriba(t *testing.T) {
	const semillas = 40
	campeon := 0
	for semilla := int64(1); semilla <= semillas; semilla++ {
		equipos := equiposUniformes(6, 50)
		equipos[3] = equipoUniforme("Gigante", 80)
		temp, _ := Nueva(equipos)
		if err := temp.JugarTemporada(nuevoRand(semilla)); err != nil {
			t.Fatal(err)
		}
		if temp.Tabla()[0].Equipo == "Gigante" {
			campeon++
		}
	}
	if campeon*100 < semillas*90 {
		t.Errorf("el mejor equipo fue campeon %d/%d veces, se esperaba al menos el 90%%", campeon, semillas)
	}
}

func TestJugarJornadaConUsaLaAlineacionElegida(t *testing.T) {
	equipos := ligaConIDsUnicas(4, 62)
	temp, err := Nueva(equipos)
	if err != nil {
		t.Fatal(err)
	}
	// El equipo 0 juega con una formación fija y sus titulares en el orden inverso
	// al automático dentro de cada línea no es válido; basta fijar un 5-3-2.
	f := modelo.F532
	elegida := simulacion.AlineacionAutomatica(equipos[0], simulacion.Criterios{Formacion: &f})
	resultados, err := temp.JugarJornadaCon(nuevoRand(1), map[int]modelo.Alineacion{0: elegida})
	if err != nil {
		t.Fatal(err)
	}
	vistos := 0
	for _, res := range resultados {
		for _, lado := range []struct {
			indice int
			f      modelo.Formacion
			t      [modelo.TitularesPorEquipo]int
		}{
			{res.Partido.Local, res.Detalle.FormacionLocal, res.Detalle.TitularesLocal},
			{res.Partido.Visitante, res.Detalle.FormacionVisitante, res.Detalle.TitularesVisitante},
		} {
			if lado.indice != 0 {
				continue
			}
			vistos++
			if lado.f != modelo.F532 || lado.t != elegida.Titulares {
				t.Errorf("el equipo 0 no jugo con su alineacion: %v %v", lado.f, lado.t)
			}
		}
	}
	if vistos != 1 {
		t.Fatalf("el equipo 0 jugo %d veces en la jornada", vistos)
	}
}

func TestJugarJornadaSinAlineacionesEsComoConLasAutomaticas(t *testing.T) {
	a, _ := Nueva(ligaConIDsUnicas(6, 62))
	b, _ := Nueva(ligaConIDsUnicas(6, 62))
	for !a.Terminada() {
		ra, err := a.JugarJornada(nuevoRand(int64(a.JornadaActual())))
		if err != nil {
			t.Fatal(err)
		}
		rb, err := b.JugarJornadaCon(nuevoRand(int64(b.JornadaActual())), map[int]modelo.Alineacion{})
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(ra, rb) {
			t.Fatal("una jornada sin alineaciones elegidas deberia ser igual a JugarJornada")
		}
	}
}

func TestLasEstadisticasSeAcreditanAlPuestoJugado(t *testing.T) {
	// Un delantero que juega de defensa (puesto 2) cuenta la porteria imbatida del
	// equipo, aunque su posicion natural sea otra; y un defensa de delantero, no.
	d := detalleDePrueba()
	// slot 1 (defensa) <- 9 (delantero) y slot 8 (delantero) <- 2 (defensa); slot 2 y slot 9 se cruzan igual.
	d.TitularesLocal[1], d.TitularesLocal[8] = 9, 2
	d.TitularesLocal[2], d.TitularesLocal[9] = 10, 3
	est, err := estadisticasDePartido(d, 0, 0)
	if err != nil {
		t.Fatal(err)
	}
	delanteroDeDefensa, defensaDeDelantero, portero := est[9], est[3], est[1]
	if delanteroDeDefensa.PorteriasImbatidas != 1 {
		t.Errorf("el delantero que jugo de defensa deberia sumar la porteria imbatida: %+v", delanteroDeDefensa)
	}
	if defensaDeDelantero.PorteriasImbatidas != 0 {
		t.Errorf("el defensa que jugo de delantero no suma porterias imbatidas: %+v", defensaDeDelantero)
	}
	if portero.PorteriasImbatidas != 1 {
		t.Errorf("el portero deberia sumar la porteria imbatida: %+v", portero)
	}
}
