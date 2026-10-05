package simulacion

import (
	"fmt"
	"math/rand"
	"reflect"
	"testing"

	"github.com/ETurriza/juego_futbol/internal/modelo"
)

// equipoRealista arma una plantilla de 22 (2/7/7/6) con atributos parecidos a
// los que genera el generador: cada posición destaca en lo suyo. Las IDs
// empiezan en idInicial.
func equipoRealista(nombre string, idInicial, calidad int) modelo.Equipo {
	perfil := map[modelo.Posicion]modelo.Atributos{
		modelo.Portero:       {Ritmo: -15, Tiro: -35, Pase: -10, Regate: -30, Defensa: -10, Fisico: 0, Reflejos: 25},
		modelo.Defensa:       {Ritmo: 0, Tiro: -20, Pase: -5, Regate: -15, Defensa: 20, Fisico: 10, Reflejos: -50},
		modelo.Mediocampista: {Ritmo: 0, Tiro: 0, Pase: 15, Regate: 10, Defensa: -5, Fisico: 0, Reflejos: -50},
		modelo.Delantero:     {Ritmo: 10, Tiro: 20, Pase: 0, Regate: 10, Defensa: -25, Fisico: 0, Reflejos: -50},
	}
	lim := func(v int) int { return min(max(v, 1), 99) }
	e := modelo.Equipo{Nombre: nombre}
	id := idInicial
	for _, c := range []struct {
		p modelo.Posicion
		n int
	}{{modelo.Portero, 2}, {modelo.Defensa, 7}, {modelo.Mediocampista, 7}, {modelo.Delantero, 6}} {
		a := perfil[c.p]
		for i := 0; i < c.n; i++ {
			e.Plantilla = append(e.Plantilla, modelo.Jugador{
				ID: id, Nombre: fmt.Sprintf("%s %d", nombre, id), Edad: 25, Posicion: c.p,
				Atributos: modelo.Atributos{
					Ritmo: lim(calidad + a.Ritmo), Tiro: lim(calidad + a.Tiro), Pase: lim(calidad + a.Pase),
					Regate: lim(calidad + a.Regate), Defensa: lim(calidad + a.Defensa),
					Fisico: lim(calidad + a.Fisico), Reflejos: lim(calidad + a.Reflejos),
				},
			})
			id++
		}
	}
	return e
}

// jugarMuchos simula n partidos entre dos equipos realistas y devuelve los
// resultados.
func jugarMuchos(n int, semilla int64) (local, visitante modelo.Equipo, resultados []Resultado) {
	local, visitante = equipoRealista("L", 1, 62), equipoRealista("V", 101, 62)
	r := rand.New(rand.NewSource(semilla))
	for i := 0; i < n; i++ {
		resultados = append(resultados, Simular(r, local, visitante))
	}
	return local, visitante, resultados
}

func TestElDetalleEsCoherenteConElMarcador(t *testing.T) {
	local, visitante, resultados := jugarMuchos(3000, 1)
	idsLocal, idsVis := map[int]bool{}, map[int]bool{}
	for _, j := range local.Plantilla {
		idsLocal[j.ID] = true
	}
	for _, j := range visitante.Plantilla {
		idsVis[j.ID] = true
	}
	for i, res := range resultados {
		gl, gv := res.Detalle.Goles()
		if gl != res.GolesLocal || gv != res.GolesVisitante {
			t.Fatalf("partido %d: los sucesos suman %d-%d y el marcador es %d-%d", i, gl, gv, res.GolesLocal, res.GolesVisitante)
		}
		partes, err := res.Detalle.Participaciones()
		if err != nil {
			t.Fatalf("partido %d: detalle incoherente: %v", i, err)
		}
		titulares := 0
		for _, p := range partes {
			if p.Titular() {
				titulares++
			}
			if p.Minutos() <= 0 || p.Minutos() > modelo.MinutosPartido {
				t.Fatalf("partido %d: %d minutos para el jugador %d", i, p.Minutos(), p.Jugador)
			}
			if p.Local && !idsLocal[p.Jugador] || !p.Local && !idsVis[p.Jugador] {
				t.Fatalf("partido %d: el jugador %d juega en el equipo equivocado", i, p.Jugador)
			}
		}
		if titulares != 22 {
			t.Fatalf("partido %d: %d titulares, se esperaban 22", i, titulares)
		}
		// Los goles con asistente: el asistente es de otro jugador del mismo equipo (ya
		// lo comprueba Participaciones) y nunca hay mas asistencias que goles.
		goles, asistencias := 0, 0
		for _, e := range res.Detalle.Eventos {
			if e.Tipo == modelo.Gol {
				goles++
				if e.Otro != 0 {
					asistencias++
				}
			}
		}
		if asistencias > goles {
			t.Fatalf("partido %d: %d asistencias para %d goles", i, asistencias, goles)
		}
	}
}

func TestLosTitularesSonElOnceEsperado(t *testing.T) {
	local, _, resultados := jugarMuchos(50, 2)
	pos := map[int]modelo.Posicion{}
	for _, j := range local.Plantilla {
		pos[j.ID] = j.Posicion
	}
	for _, res := range resultados {
		cuenta := map[modelo.Posicion]int{}
		for _, id := range res.Detalle.TitularesLocal {
			cuenta[pos[id]]++
		}
		want := map[modelo.Posicion]int{modelo.Portero: 1, modelo.Defensa: 4, modelo.Mediocampista: 3, modelo.Delantero: 3}
		if !reflect.DeepEqual(cuenta, want) {
			t.Fatalf("titulares por posicion = %v, se esperaba %v", cuenta, want)
		}
	}
}

func TestLosEventosVanOrdenadosPorMinuto(t *testing.T) {
	_, _, resultados := jugarMuchos(500, 3)
	for i, res := range resultados {
		ultimo := 0
		for _, e := range res.Detalle.Eventos {
			if e.Minuto < ultimo {
				t.Fatalf("partido %d: sucesos fuera de orden", i)
			}
			ultimo = e.Minuto
		}
	}
}

// El reparto de goles y asistencias por posición debe parecerse al fútbol: la
// mayoría de los goles son de delanteros y casi ninguno del portero.
func TestRepartoDeGolesYAsistenciasPorPosicion(t *testing.T) {
	local, visitante, resultados := jugarMuchos(6000, 4)
	pos := map[int]modelo.Posicion{}
	for _, e := range []modelo.Equipo{local, visitante} {
		for _, j := range e.Plantilla {
			pos[j.ID] = j.Posicion
		}
	}
	goles := map[modelo.Posicion]int{}
	asist := map[modelo.Posicion]int{}
	totalGoles, totalAsist := 0, 0
	for _, res := range resultados {
		for _, e := range res.Detalle.Eventos {
			if e.Tipo != modelo.Gol {
				continue
			}
			goles[pos[e.Jugador]]++
			totalGoles++
			if e.Otro != 0 {
				asist[pos[e.Otro]]++
				totalAsist++
			}
		}
	}
	pct := func(parte, total int) float64 { return 100 * float64(parte) / float64(total) }
	t.Logf("goles: delanteros %.1f%%, medios %.1f%%, defensas %.1f%%, porteros %.2f%%",
		pct(goles[modelo.Delantero], totalGoles), pct(goles[modelo.Mediocampista], totalGoles),
		pct(goles[modelo.Defensa], totalGoles), pct(goles[modelo.Portero], totalGoles))
	t.Logf("asistencias: delanteros %.1f%%, medios %.1f%%, defensas %.1f%%; con asistencia: %.1f%% de los goles",
		pct(asist[modelo.Delantero], totalAsist), pct(asist[modelo.Mediocampista], totalAsist),
		pct(asist[modelo.Defensa], totalAsist), pct(totalAsist, totalGoles))

	rangos := []struct {
		nombre string
		valor  float64
		min    float64
		max    float64
	}{
		{"goles de delanteros", pct(goles[modelo.Delantero], totalGoles), 52, 70},
		{"goles de mediocampistas", pct(goles[modelo.Mediocampista], totalGoles), 22, 36},
		{"goles de defensas", pct(goles[modelo.Defensa], totalGoles), 4, 16},
		{"goles de porteros", pct(goles[modelo.Portero], totalGoles), 0, 1},
		{"asistencias de mediocampistas", pct(asist[modelo.Mediocampista], totalAsist), 38, 55},
		{"asistencias de delanteros", pct(asist[modelo.Delantero], totalAsist), 22, 40},
		{"asistencias de defensas", pct(asist[modelo.Defensa], totalAsist), 8, 24},
		{"goles con asistencia", pct(totalAsist, totalGoles), 66, 74},
	}
	for _, r := range rangos {
		if r.valor < r.min || r.valor > r.max {
			t.Errorf("%s: %.1f%% fuera de [%.0f, %.0f]", r.nombre, r.valor, r.min, r.max)
		}
	}
}

func TestTarjetasPorPartidoSonRealistas(t *testing.T) {
	_, _, resultados := jugarMuchos(8000, 5)
	var amarillas, rojas, cambios int
	for _, res := range resultados {
		for _, e := range res.Detalle.Eventos {
			switch e.Tipo {
			case modelo.Amarilla:
				amarillas++
			case modelo.Roja:
				rojas++
			case modelo.Sustitucion:
				cambios++
			}
		}
	}
	n := float64(len(resultados))
	t.Logf("por partido: %.2f amarillas, %.3f rojas, %.2f cambios", float64(amarillas)/n, float64(rojas)/n, float64(cambios)/n)
	if v := float64(amarillas) / n; v < 3.0 || v > 4.2 {
		t.Errorf("%.2f amarillas por partido fuera de [3.0, 4.2]", v)
	}
	if v := float64(rojas) / n; v < 0.07 || v > 0.25 {
		t.Errorf("%.3f rojas por partido fuera de [0.07, 0.25]", v)
	}
	if v := float64(cambios) / n; v < 8 || v > 9.6 {
		t.Errorf("%.2f cambios por partido fuera de [8, 9.6] (3 a 5 por equipo)", v)
	}
}

func TestLasRojasNoLlegaranAPorteros(t *testing.T) {
	local, visitante, resultados := jugarMuchos(8000, 6)
	pos := map[int]modelo.Posicion{}
	for _, e := range []modelo.Equipo{local, visitante} {
		for _, j := range e.Plantilla {
			pos[j.ID] = j.Posicion
		}
	}
	for _, res := range resultados {
		amonestados := map[int]int{}
		for _, e := range res.Detalle.Eventos {
			switch e.Tipo {
			case modelo.Amarilla:
				amonestados[e.Jugador]++
			case modelo.Roja:
				// Una roja a un portero solo podria venir de una segunda amarilla.
				if pos[e.Jugador] == modelo.Portero && amonestados[e.Jugador] < 2 {
					t.Fatal("roja directa a un portero")
				}
			}
		}
	}
}

func TestSegundaAmarillaVaSeguidaDeRoja(t *testing.T) {
	_, _, resultados := jugarMuchos(20000, 7)
	vistos := 0
	for _, res := range resultados {
		amarillas := map[int]int{}
		eventos := res.Detalle.Eventos
		for i, e := range eventos {
			if e.Tipo != modelo.Amarilla {
				continue
			}
			amarillas[e.Jugador]++
			if amarillas[e.Jugador] == 2 {
				vistos++
				sigue := false
				for _, s := range eventos[i+1:] {
					if s.Minuto != e.Minuto {
						break
					}
					sigue = sigue || (s.Tipo == modelo.Roja && s.Jugador == e.Jugador)
				}
				if !sigue {
					t.Fatalf("la segunda amarilla del jugador %d no va seguida de su roja", e.Jugador)
				}
			}
			if amarillas[e.Jugador] > 2 {
				t.Fatalf("el jugador %d tiene %d amarillas en un partido", e.Jugador, amarillas[e.Jugador])
			}
		}
	}
	if vistos == 0 {
		t.Error("en 20000 partidos deberia haber alguna doble amarilla")
	}
}

func TestLosSuplentesJueganYSumanMinutos(t *testing.T) {
	_, _, resultados := jugarMuchos(2000, 8)
	golesSuplentes, goles := 0, 0
	for _, res := range resultados {
		partes, _ := res.Detalle.Participaciones()
		suplentes := map[int]bool{}
		porSide := map[bool]int{}
		for _, p := range partes {
			if !p.Titular() {
				suplentes[p.Jugador] = true
				porSide[p.Local]++
				if p.Desde < 46 || p.Desde > 85 {
					t.Fatalf("un suplente entro en el minuto %d, fuera de [46, 85]", p.Desde)
				}
			}
		}
		for _, lado := range []bool{true, false} {
			if n := porSide[lado]; n < 3 || n > 5 {
				t.Fatalf("un equipo uso %d suplentes, se esperaban de 3 a 5", n)
			}
		}
		for _, e := range res.Detalle.Eventos {
			if e.Tipo == modelo.Gol {
				goles++
				if suplentes[e.Jugador] {
					golesSuplentes++
				}
			}
		}
	}
	// Los suplentes juegan ~25 de 90 minutos, así que marcan una parte de los goles.
	if golesSuplentes == 0 || float64(golesSuplentes)/float64(goles) > 0.20 {
		t.Errorf("goles de suplentes: %d de %d", golesSuplentes, goles)
	}
}

func TestLosPorterosNoSeSustituyen(t *testing.T) {
	local, visitante, resultados := jugarMuchos(3000, 9)
	pos := map[int]modelo.Posicion{}
	for _, e := range []modelo.Equipo{local, visitante} {
		for _, j := range e.Plantilla {
			pos[j.ID] = j.Posicion
		}
	}
	for _, res := range resultados {
		for _, e := range res.Detalle.Eventos {
			if e.Tipo == modelo.Sustitucion && (pos[e.Jugador] == modelo.Portero || pos[e.Otro] == modelo.Portero) {
				t.Fatal("un portero entro o salio en una sustitucion")
			}
		}
	}
}

func TestLosMejoresDelanterosMarcanMas(t *testing.T) {
	e := equipoRealista("L", 1, 62)
	// Un delantero estrella: tiro altisimo.
	var estrella, comun int
	for i, j := range e.Plantilla {
		if j.Posicion == modelo.Delantero {
			if estrella == 0 {
				estrella = j.ID
				e.Plantilla[i].Atributos.Tiro = 99
			} else if comun == 0 {
				comun = j.ID
			}
		}
	}
	rival := equipoRealista("V", 101, 62)
	r := rand.New(rand.NewSource(10))
	golesEstrella, golesComun := 0, 0
	for i := 0; i < 5000; i++ {
		for _, ev := range Simular(r, e, rival).Detalle.Eventos {
			if ev.Tipo == modelo.Gol && ev.Local {
				if ev.Jugador == estrella {
					golesEstrella++
				} else if ev.Jugador == comun {
					golesComun++
				}
			}
		}
	}
	if golesEstrella <= golesComun {
		t.Errorf("la estrella (%d goles) deberia marcar mas que un delantero comun (%d)", golesEstrella, golesComun)
	}
}

func TestLosDefensasYMediosVenMasAmarillasQueLosDelanteros(t *testing.T) {
	local, visitante, resultados := jugarMuchos(8000, 11)
	pos := map[int]modelo.Posicion{}
	n := map[modelo.Posicion]int{}
	for _, e := range []modelo.Equipo{local, visitante} {
		for _, j := range e.Plantilla {
			pos[j.ID] = j.Posicion
			n[j.Posicion]++
		}
	}
	amarillas := map[modelo.Posicion]int{}
	for _, res := range resultados {
		for _, e := range res.Detalle.Eventos {
			if e.Tipo == modelo.Amarilla {
				amarillas[pos[e.Jugador]]++
			}
		}
	}
	// Por jugador de la plantilla (los suplentes también cuentan).
	por := func(p modelo.Posicion) float64 { return float64(amarillas[p]) / float64(n[p]) }
	if por(modelo.Defensa) <= por(modelo.Delantero) || por(modelo.Mediocampista) <= por(modelo.Delantero) {
		t.Errorf("amarillas por jugador: defensa %.0f, medio %.0f, delantero %.0f",
			por(modelo.Defensa), por(modelo.Mediocampista), por(modelo.Delantero))
	}
	if por(modelo.Portero) >= por(modelo.Defensa) {
		t.Errorf("un portero (%.0f) deberia ver menos amarillas que un defensa (%.0f)", por(modelo.Portero), por(modelo.Defensa))
	}
}

func TestElDetalleEsReproducible(t *testing.T) {
	_, _, a := jugarMuchos(200, 12)
	_, _, b := jugarMuchos(200, 12)
	if !reflect.DeepEqual(a, b) {
		t.Error("la misma semilla deberia dar los mismos sucesos")
	}
	_, _, c := jugarMuchos(200, 13)
	if reflect.DeepEqual(a, c) {
		t.Error("semillas distintas deberian dar sucesos distintos")
	}
}

func TestPlantillasIncompletasNoFallan(t *testing.T) {
	// Equipos con pocos jugadores, sin banca o vacíos: no deben fallar y el
	// detalle debe seguir siendo coherente.
	pocos := equipoRealista("P", 1, 62)
	pocos.Plantilla = pocos.Plantilla[:5]
	sinBanca := equipoRealista("S", 100, 62)
	sinBanca.Plantilla = sinBanca.Plantilla[:11]
	vacio := modelo.Equipo{Nombre: "Vacio"}
	r := rand.New(rand.NewSource(14))
	equipos := []modelo.Equipo{pocos, sinBanca, vacio, equipoRealista("N", 200, 62)}
	for i, a := range equipos {
		for k, b := range equipos {
			if i == k {
				continue // el mismo equipo repetiria las IDs de jugador
			}
			for i := 0; i < 200; i++ {
				res := Simular(r, a, b)
				if _, err := res.Detalle.Participaciones(); err != nil {
					t.Fatalf("%s vs %s: detalle incoherente: %v", a.Nombre, b.Nombre, err)
				}
			}
		}
	}
}

func TestSimularNoModificaLaPlantillaConDetalle(t *testing.T) {
	a, b := equipoRealista("A", 1, 62), equipoRealista("B", 101, 62)
	copiaA, copiaB := equipoRealista("A", 1, 62), equipoRealista("B", 101, 62)
	r := rand.New(rand.NewSource(15))
	for i := 0; i < 100; i++ {
		Simular(r, a, b)
	}
	if !reflect.DeepEqual(a, copiaA) || !reflect.DeepEqual(b, copiaB) {
		t.Error("Simular modifico los equipos")
	}
}
