package simulacion

import (
	"math"
	"math/rand"
	"sort"
	"sync"
	"testing"

	"github.com/ETurriza/juego_futbol/internal/modelo"
)

// victorias simula partidos de ida y vuelta alternados entre mio y rival y
// devuelve el porcentaje de victorias de mio.
func victorias(mio, rival modelo.Equipo, partidos int, semilla int64) float64 {
	r := rand.New(rand.NewSource(semilla))
	g := 0
	for i := 0; i < partidos; i++ {
		var gm, gr int
		if i%2 == 0 {
			gm, gr = Marcador(r, mio, rival)
		} else {
			gr, gm = Marcador(r, rival, mio)
		}
		if gm > gr {
			g++
		}
	}
	return 100 * float64(g) / float64(partidos)
}

// titularIdx devuelve los índices de plantilla de los n mejores jugadores de la
// posición, como los elige la simulación.
func titularIdx(e modelo.Equipo, p modelo.Posicion, n int) []int {
	var idx []int
	for i, j := range e.Plantilla {
		if j.Posicion == p {
			idx = append(idx, i)
		}
	}
	sort.SliceStable(idx, func(a, b int) bool {
		x, y := e.Plantilla[idx[a]], e.Plantilla[idx[b]]
		if x.Valoracion() != y.Valoracion() {
			return x.Valoracion() > y.Valoracion()
		}
		return x.ID < y.ID
	})
	return idx[:min(n, len(idx))]
}

// plaza es un puesto del once: la posición y cuántos titulares tiene.
type plaza struct {
	p    modelo.Posicion
	n, k int // el k-ésimo mejor de los n titulares de la posición
}

var (
	plazaPortero   = plaza{modelo.Portero, 1, 0}
	plazaDefensa   = plaza{modelo.Defensa, 4, 0}
	plazaMedio     = plaza{modelo.Mediocampista, 3, 0}
	plazaDelantero = plaza{modelo.Delantero, 3, 0}
	onceCompleto   []plaza
)

func init() {
	onceCompleto = []plaza{plazaPortero}
	for k := 0; k < 4; k++ {
		onceCompleto = append(onceCompleto, plaza{modelo.Defensa, 4, k})
	}
	for k := 0; k < 3; k++ {
		onceCompleto = append(onceCompleto, plaza{modelo.Mediocampista, 3, k}, plaza{modelo.Delantero, 3, k})
	}
}

// conCracks devuelve una copia del equipo con un crack (todos los atributos en
// 95, valoración 95) en cada plaza dada.
func conCracks(e modelo.Equipo, plazas ...plaza) modelo.Equipo {
	copia := modelo.Equipo{Nombre: e.Nombre, Plantilla: append([]modelo.Jugador(nil), e.Plantilla...)}
	for _, pl := range plazas {
		i := titularIdx(e, pl.p, pl.n)[pl.k]
		copia.Plantilla[i].Atributos = atributosUniformesSim(95)
	}
	return copia
}

func atributosUniformesSim(v int) modelo.Atributos {
	return modelo.Atributos{Ritmo: v, Tiro: v, Pase: v, Regate: v, Defensa: v, Fisico: v, Reflejos: v}
}

const partidosCracks = 20000

var (
	controlUnaVez    sync.Once
	controlVictorias float64
)

// efecto devuelve cuántos puntos de victoria suma tener cracks en las plazas dadas
// frente al mismo equipo sin ellos. El equipo de control se simula una sola vez.
func efecto(plazas ...plaza) float64 {
	base := equipoRealista("Base", 1, 62)
	rival := equipoRealista("Rival", 101, 62)
	controlUnaVez.Do(func() { controlVictorias = victorias(base, rival, partidosCracks, 7) })
	return victorias(conCracks(base, plazas...), rival, partidosCracks, 7) - controlVictorias
}

func TestUnCrackSeNotaEnElResultadoEnCadaPosicion(t *testing.T) {
	efectos := map[string]float64{
		"portero":   efecto(plazaPortero),
		"defensa":   efecto(plazaDefensa),
		"medio":     efecto(plazaMedio),
		"delantero": efecto(plazaDelantero),
	}
	t.Logf("puntos de victoria que suma un crack de 95: %v", efectos)
	menor, mayor := math.Inf(1), math.Inf(-1)
	for pos, e := range efectos {
		// Un crack aislado (en un equipo sin otros cracks) suma bastante más que
		// los +2 de antes del bono; dentro de una liga real, donde los equipos ya
		// tienen algún crack, suma ~8 (se comprueba en aplicacion).
		if e < 5 || e > 14 {
			t.Errorf("un crack de %s suma %.1f puntos de victoria, se esperaban entre 5 y 14", pos, e)
		}
		menor, mayor = math.Min(menor, e), math.Max(mayor, e)
	}
	// Ninguna posición pesa mucho más que otra.
	if mayor > 2.2*menor {
		t.Errorf("el efecto de un crack varia demasiado entre posiciones: de %.1f a %.1f", menor, mayor)
	}
}

func TestVariosCracksPesanMasPeroConTope(t *testing.T) {
	uno := efecto(plazaDelantero)
	tres := efecto(plazaPortero, plazaMedio, plazaDelantero)
	t.Logf("1 crack: +%.1f; 3 cracks: +%.1f", uno, tres)
	if tres < 14 || tres > 33 {
		t.Errorf("tres cracks suman %.1f puntos de victoria, se esperaban entre 14 y 33", tres)
	}
	if tres < 1.8*uno {
		t.Errorf("tres cracks (%.1f) deberian pesar bastante más que uno (%.1f)", tres, uno)
	}
	// Pero los efectos se saturan: cada crack adicional suma menos que el primero.
	if tres > 3.3*uno {
		t.Errorf("el efecto no se satura: 1 crack %.1f, 3 cracks %.1f", uno, tres)
	}
}

func TestUnEquipoDeCracksNoGanaSiempre(t *testing.T) {
	base := equipoRealista("Base", 1, 62)
	rival := equipoRealista("Rival", 101, 62)
	todos := conCracks(base, onceCompleto...)
	gana := victorias(todos, rival, partidosCracks, 8)
	t.Logf("un once de cracks de 95 gana el %.1f%% frente a un equipo normal", gana)
	if gana < 60 || gana > 88 {
		t.Errorf("un once de cracks gana el %.1f%%; deberia estar entre 60%% y 88%%", gana)
	}
	// Hay mucho fútbol: aun así pierde partidos.
	if perdidos := victorias(rival, todos, partidosCracks, 9); perdidos < 3 {
		t.Errorf("el equipo normal solo gana el %.1f%% contra un once de cracks", perdidos)
	}
}

func TestLaBrechaNormalDeNivelNoCambia(t *testing.T) {
	// Equipos sin cracks (nadie llega a la valoración del umbral): un equipo 10
	// puntos mejor sigue ganando algo más de la mitad, como antes del bono.
	rival := equipoUniforme("R", 72)
	gana := victorias(equipoUniforme("M", 82), rival, 30000, 3)
	t.Logf("un equipo de nivel 82 contra uno de 72 gana el %.1f%%", gana)
	if gana < 51 || gana > 57 {
		t.Errorf("un equipo +10 gana el %.1f%%, se esperaba ~54%%", gana)
	}
	if empate := victorias(equipoUniforme("M", 72), rival, 30000, 4); empate < 35.5 || empate > 39 {
		t.Errorf("dos equipos iguales: %.1f%% de victorias, se esperaba ~37%% (hay empates)", empate)
	}
}

func TestElBonoDeLosCracksEstaAcotado(t *testing.T) {
	// Aun con un once perfecto (todos en 99) contra uno normal, la liga sigue
	// teniendo azar.
	gana := victorias(equipoUniforme("M", 99), equipoUniforme("R", 72), 30000, 5)
	t.Logf("un equipo de 99 contra uno de 72 gana el %.1f%%", gana)
	if gana > 95 {
		t.Errorf("un equipo perfecto gana el %.1f%%; el bono deberia estar acotado", gana)
	}
}

func TestMediaSaturaElAporteDeLosCracks(t *testing.T) {
	var sin media
	sin.sumar(1, 70)
	sin.sumar(1, 80)
	if got := sin.valor(); math.Abs(got-75) > 1e-9 {
		t.Errorf("sin cracks, la media deberia ser simple: %.3f", got)
	}

	// Cada jugador por debajo del umbral no aporta nada.
	var bajo media
	bajo.sumar(1, 70)
	bajo.sumarCrack(modelo.Jugador{Posicion: modelo.Delantero, Atributos: atributosUniformesSim(umbralCrack)}, modelo.Delantero, crackDelantero)
	if bajo.cracks != 0 {
		t.Errorf("un jugador justo en el umbral no deberia aportar: %.3f", bajo.cracks)
	}

	// El aporte crece con los cracks pero nunca pasa del tope.
	anterior := 0.0
	for n := 1; n <= 40; n++ {
		var m media
		m.sumar(1, 70)
		for i := 0; i < n; i++ {
			m.sumarCrack(modelo.Jugador{Posicion: modelo.Defensa, Atributos: atributosUniformesSim(99)}, modelo.Defensa, crackDefensa)
		}
		aporte := m.valor() - 70
		if aporte < anterior-1e-9 || aporte > topeCracks+1e-9 {
			t.Fatalf("con %d cracks el aporte es %.3f (anterior %.3f, tope %.1f)", n, aporte, anterior, topeCracks)
		}
		anterior = aporte
	}
	// Con pocos cracks el aporte es casi lineal.
	var uno media
	uno.sumar(1, 70)
	uno.sumarCrack(modelo.Jugador{Posicion: modelo.Delantero, Atributos: atributosUniformesSim(90)}, modelo.Delantero, crackDelantero)
	esperado := crackDelantero * (90 - umbralCrack)
	if got := uno.valor() - 70; got > esperado || got < 0.8*esperado {
		t.Errorf("con un crack leve el aporte (%.3f) deberia ser casi el lineal (%.3f)", got, esperado)
	}
}

func TestUnCrackSoloCuentaPorSuValoracion(t *testing.T) {
	// Un jugador cuyo valor de rol es altísimo pero su valoración es normal no es
	// un crack: los porteros tienen reflejos altos por naturaleza.
	var m media
	m.sumar(1, 70)
	portero := modelo.Jugador{Posicion: modelo.Portero, Atributos: modelo.Atributos{
		Ritmo: 55, Tiro: 35, Pase: 60, Regate: 40, Defensa: 60, Fisico: 70, Reflejos: 95}}
	if portero.Valoracion() >= umbralCrack {
		t.Fatalf("el caso de prueba esta mal: valoracion %d", portero.Valoracion())
	}
	m.sumarCrack(portero, modelo.Portero, crackPortero)
	if m.cracks != 0 {
		t.Errorf("un portero de valoracion %d con 95 en reflejos no deberia ser crack", portero.Valoracion())
	}
}

func TestMarcadorYSimularDanLosMismosGolesEnPromedio(t *testing.T) {
	// Marcador (4-3-3, sin tarjetas) y Simular con la misma formación (con tarjetas
	// y su efecto sobre los goles) no dan el mismo resultado partido a partido,
	// pero sí un promedio de goles muy parecido: las rojas son raras.
	a, b := equipoRealista("A", 1, 62), equipoRealista("B", 101, 62)
	alA, alB := alineacion433(a), alineacion433(b)
	rm := rand.New(rand.NewSource(1))
	rs := rand.New(rand.NewSource(2))
	const n = 20000
	var gm, gs float64
	for i := 0; i < n; i++ {
		gl, gv := Marcador(rm, a, b)
		gm += float64(gl + gv)
		res := SimularConAlineaciones(rs, a, alA, b, alB)
		gs += float64(res.GolesLocal + res.GolesVisitante)
	}
	if d := math.Abs(gm-gs) / n; d > 0.08 {
		t.Errorf("goles por partido: Marcador %.2f, Simular %.2f", gm/n, gs/n)
	}
}
