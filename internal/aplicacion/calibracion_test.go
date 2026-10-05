package aplicacion

import (
	"fmt"
	"math"
	"math/rand"
	"sort"
	"strings"
	"testing"

	"github.com/ETurriza/juego_futbol/internal/liga"
	"github.com/ETurriza/juego_futbol/internal/modelo"
	"github.com/ETurriza/juego_futbol/internal/simulacion"
)

// Las pruebas de calibración miden la distribución de la valoración por posición
// y por edad de una liga a lo largo de muchas temporadas. La media global no
// basta: una liga puede mantener el promedio y, aun así, llenarse de estrellas o
// inflar a los porteros. Con go test -v imprimen las tablas medidas.

const ligasDeCalibracion = 30

// ligaEnvejecida crea una liga de 10 equipos y la hace pasar n temporadas de
// renovación (retiros, envejecimiento y juveniles) sin jugar los partidos, que
// no influyen en las plantillas.
func ligaEnvejecida(t *testing.T, semilla int64, n int) []modelo.Equipo {
	t.Helper()
	c := nuevaCarrera(t, semilla, 10)
	equipos := clonarEquipos(c.Temporada.Equipos)
	proximoID := c.ProximoID
	r := rand.New(rand.NewSource(semilla ^ 0x5eed))
	for i := 0; i < n; i++ {
		for k, e := range equipos {
			equipos[k], _, _ = renovarEquipo(r, e, &proximoID)
		}
	}
	return equipos
}

// distribucion acumula la valoración de un grupo de jugadores.
type distribucion struct {
	n, altos, bajos, topes int
	suma, suma2            float64
}

func (d *distribucion) add(j modelo.Jugador) {
	v := j.Valoracion()
	d.n++
	d.suma += float64(v)
	d.suma2 += float64(v * v)
	if v >= 85 {
		d.altos++
	}
	if v <= 50 {
		d.bajos++
	}
	a := j.Atributos
	for _, x := range []int{a.Ritmo, a.Tiro, a.Pase, a.Regate, a.Defensa, a.Fisico} {
		if x == modelo.AtributoMax {
			d.topes++
		}
	}
	if j.Posicion == modelo.Portero && a.Reflejos == modelo.AtributoMax {
		d.topes++
	}
}

func (d distribucion) media() float64 { return d.suma / float64(d.n) }
func (d distribucion) desv() float64 {
	m := d.media()
	return math.Sqrt(d.suma2/float64(d.n) - m*m)
}
func (d distribucion) pctAltos() float64 { return 100 * float64(d.altos) / float64(d.n) }
func (d distribucion) pctBajos() float64 { return 100 * float64(d.bajos) / float64(d.n) }

// medir reúne, tras n temporadas, la distribución de los porteros y la de los
// jugadores de campo de varias ligas.
func medir(t *testing.T, temporadas int) (porteros, campo distribucion) {
	t.Helper()
	for s := int64(1); s <= ligasDeCalibracion; s++ {
		for _, e := range ligaEnvejecida(t, s*7, temporadas) {
			for _, j := range e.Plantilla {
				if j.Posicion == modelo.Portero {
					porteros.add(j)
				} else {
					campo.add(j)
				}
			}
		}
	}
	return porteros, campo
}

// TestLaLigaNoSeInflaNiSeDegrada comprueba, por posición, que la liga de la
// temporada 5, 10 o 20 se parece a la inicial: mismo nivel, misma dispersión y
// sin llenarse de estrellas.
func TestLaLigaNoSeInflaNiSeDegrada(t *testing.T) {
	inicialP, inicialC := medir(t, 0)
	var tabla strings.Builder
	fmt.Fprintf(&tabla, "\n%-9s %-8s %7s %6s %7s %7s %7s\n", "temporada", "grupo", "media", "desv", ">=85 %", "<=50 %", "a 99")
	linea := func(temp int, grupo string, d distribucion) {
		fmt.Fprintf(&tabla, "%-9d %-8s %7.1f %6.1f %7.1f %7.1f %7d\n", temp, grupo, d.media(), d.desv(), d.pctAltos(), d.pctBajos(), d.topes)
	}
	linea(0, "porteros", inicialP)
	linea(0, "campo", inicialC)

	for _, temporadas := range []int{5, 10, 20} {
		p, c := medir(t, temporadas)
		linea(temporadas, "porteros", p)
		linea(temporadas, "campo", c)

		for nombre, par := range map[string][2]distribucion{"porteros": {inicialP, p}, "campo": {inicialC, c}} {
			ini, fin := par[0], par[1]
			if d := fin.media() - ini.media(); d > 1.5 || d < -1.5 {
				t.Errorf("%s tras %d temporadas: la media pasa de %.1f a %.1f", nombre, temporadas, ini.media(), fin.media())
			}
			if d := fin.desv() - ini.desv(); d > 2 || d < -2 {
				t.Errorf("%s tras %d temporadas: la dispersion pasa de %.1f a %.1f", nombre, temporadas, ini.desv(), fin.desv())
			}
		}
		// Limites absolutos: pocas estrellas, ningún grupo inflado.
		if c.pctAltos() > 8 || p.pctAltos() > 15 {
			t.Errorf("tras %d temporadas hay demasiadas estrellas (>=85): campo %.1f%%, porteros %.1f%%", temporadas, c.pctAltos(), p.pctAltos())
		}
		if c.pctBajos() > 8 || p.pctBajos() > 11 {
			t.Errorf("tras %d temporadas hay demasiados jugadores flojos (<=50): campo %.1f%%, porteros %.1f%%", temporadas, c.pctBajos(), p.pctBajos())
		}
		if d := p.media() - c.media(); d > 4 || d < -4 {
			t.Errorf("tras %d temporadas los porteros (%.1f) y el campo (%.1f) difieren demasiado", temporadas, p.media(), c.media())
		}
		if c.desv() > 12 || p.desv() > 13 {
			t.Errorf("tras %d temporadas la dispersion es excesiva: campo %.1f, porteros %.1f", temporadas, c.desv(), p.desv())
		}
		// Los atributos no se apilan en el tope de 99.
		if float64(c.topes)/float64(c.n) > 0.02 || float64(p.topes)/float64(p.n) > 0.06 {
			t.Errorf("tras %d temporadas hay demasiados atributos en 99: campo %d, porteros %d", temporadas, c.topes, p.topes)
		}
	}
	t.Log(tabla.String())
}

// TestLaCalidadPorEdadEsEstacionaria comprueba que un jugador de cada edad vale
// lo mismo en la temporada 1 que en la 20: la liga inicial y la que resulta de
// la progresión siguen la misma curva.
func TestLaCalidadPorEdadEsEstacionaria(t *testing.T) {
	tramos := []struct {
		nombre   string
		min, max int
	}{{"20-25", 20, 25}, {"26-31", 26, 31}, {"32-35", 32, 35}}
	medias := func(temporadas int) map[string]float64 {
		suma, n := map[string]float64{}, map[string]float64{}
		for s := int64(1); s <= ligasDeCalibracion; s++ {
			for _, e := range ligaEnvejecida(t, s*11, temporadas) {
				for _, j := range e.Plantilla {
					if j.Posicion == modelo.Portero {
						continue
					}
					for _, tr := range tramos {
						if j.Edad >= tr.min && j.Edad <= tr.max {
							suma[tr.nombre] += float64(j.Valoracion())
							n[tr.nombre]++
						}
					}
				}
			}
		}
		out := map[string]float64{}
		for k, v := range suma {
			out[k] = v / n[k]
		}
		return out
	}
	inicial, final := medias(0), medias(20)
	t.Logf("valoracion media del campo por edad: inicio %v, tras 20 temporadas %v", inicial, final)
	for _, tr := range tramos {
		if d := final[tr.nombre] - inicial[tr.nombre]; d > 2.5 || d < -2.5 {
			t.Errorf("edad %s: %.1f al inicio y %.1f tras 20 temporadas", tr.nombre, inicial[tr.nombre], final[tr.nombre])
		}
	}
	// Y la curva tiene la forma esperada: crece hasta la madurez.
	if !(final["20-25"] < final["26-31"]) {
		t.Errorf("los de 26-31 (%.1f) deberian valer mas que los de 20-25 (%.1f)", final["26-31"], final["20-25"])
	}
}

// TestPocosVeteranosFlojosEnLasPlantillas comprueba el retiro por nivel: los
// veteranos que ya no dan el nivel no se quedan años en las plantillas.
func TestPocosVeteranosFlojosEnLasPlantillas(t *testing.T) {
	var equipos, conFlojo, de36, flojos36, de38, flojos38 int
	for s := int64(1); s <= ligasDeCalibracion; s++ {
		for _, e := range ligaEnvejecida(t, s*13, 15) {
			equipos++
			hay := false
			for _, j := range e.Plantilla {
				if j.Posicion == modelo.Portero {
					continue
				}
				flojo := j.Valoracion() <= 50
				if j.Edad >= 35 && j.Valoracion() <= 55 {
					hay = true
				}
				if j.Edad >= 36 {
					de36++
					if flojo {
						flojos36++
					}
				}
				if j.Edad >= 38 {
					de38++
					if flojo {
						flojos38++
					}
				}
			}
			if hay {
				conFlojo++
			}
		}
	}
	pct := func(a, b int) float64 { return 100 * float64(a) / float64(max(b, 1)) }
	t.Logf("plantillas con un veterano (>=35) de valoracion <=55: %.0f%% (%d de %d); de los de 36+, con <=50: %.1f%%; de los de 38+: %.1f%%",
		pct(conFlojo, equipos), conFlojo, equipos, pct(flojos36, de36), pct(flojos38, de38))
	if p := pct(conFlojo, equipos); p > 18 {
		t.Errorf("%.0f%% de las plantillas con un veterano flojo; deberia ser poco comun", p)
	}
	if p := pct(flojos36, de36); p > 4 {
		t.Errorf("%.1f%% de los veteranos de 36+ con valoracion <= 50", p)
	}
	if p := pct(flojos38, de38); p > 6 {
		t.Errorf("%.1f%% de los de 38+ con valoracion <= 50", p)
	}
}

// TestElRetiroPorNivelNoCambiaLaEdadMedia comprueba que adelantar los retiros no
// rejuvenece ni envejece la liga: la edad media sigue siendo la de siempre.
func TestElRetiroPorNivelNoCambiaLaEdadMedia(t *testing.T) {
	media := func(temporadas int) float64 {
		suma, n := 0.0, 0.0
		for s := int64(1); s <= ligasDeCalibracion; s++ {
			for _, e := range ligaEnvejecida(t, s*19, temporadas) {
				for _, j := range e.Plantilla {
					suma += float64(j.Edad)
					n++
				}
			}
		}
		return suma / n
	}
	inicial, final := media(0), media(15)
	t.Logf("edad media: inicial %.2f, tras 15 temporadas %.2f", inicial, final)
	if final < 25.5 || final > 28.5 || math.Abs(final-inicial) > 1 {
		t.Errorf("edad media inicial %.2f y final %.2f", inicial, final)
	}
}

// TestLaLigaInicialTienePocosVeteranosFlojos comprueba que la liga que se crea
// al empezar una carrera no incluye veteranos que, siguiendo la progresión, ya se
// habrían retirado por falta de nivel (sin el filtro de supervivencia del
// generador, el porcentaje casi se duplica).
func TestLaLigaInicialTienePocosVeteranosFlojos(t *testing.T) {
	equipos, con := 0, 0
	for s := int64(1); s <= 300; s++ {
		c := nuevaCarrera(t, s*3, 10)
		for _, e := range c.Temporada.Equipos {
			equipos++
			for _, j := range e.Plantilla {
				if j.Posicion != modelo.Portero && j.Edad >= 35 && j.Valoracion() <= 55 {
					con++
					break
				}
			}
		}
	}
	pct := 100 * float64(con) / float64(equipos)
	t.Logf("liga inicial: plantillas con un veterano (>=35) de valoracion <=55: %.1f%% (%d de %d)", pct, con, equipos)
	if pct > 5.5 {
		t.Errorf("%.1f%% de las plantillas iniciales con un veterano flojo; deberia ser ~4 %%", pct)
	}
}

// ---- Cracks: cuántos hay, cómo se reparten y cuánto pesan ----

// picosDeCarrera devuelve, por posición, la mejor valoración que alcanza cada
// jugador de muchas ligas a lo largo de 25 temporadas de renovación.
func picosDeCarrera(t *testing.T) map[modelo.Posicion][]int {
	t.Helper()
	picos := map[modelo.Posicion][]int{}
	for s := int64(1); s <= 40; s++ {
		c := nuevaCarrera(t, s*17, 10)
		equipos := clonarEquipos(c.Temporada.Equipos)
		proximoID := c.ProximoID
		r := rand.New(rand.NewSource(s ^ 0xabc))
		pico := map[int]int{}
		pos := map[int]modelo.Posicion{}
		for n := 0; n < 25; n++ {
			for _, e := range equipos {
				for _, j := range e.Plantilla {
					if v := j.Valoracion(); v > pico[j.ID] {
						pico[j.ID] = v
					}
					pos[j.ID] = j.Posicion
				}
			}
			for k, e := range equipos {
				equipos[k], _, _ = renovarEquipo(r, e, &proximoID)
			}
		}
		for id, v := range pico {
			picos[pos[id]] = append(picos[pos[id]], v)
		}
	}
	return picos
}

func porcentajeDesde(v []int, minimo int) float64 {
	n := 0
	for _, x := range v {
		if x >= minimo {
			n++
		}
	}
	return 100 * float64(n) / float64(len(v))
}

// TestHayCracksYSeRepartenPorTodasLasPosiciones comprueba el objetivo de diseño:
// de cada cien jugadores de campo, unos tres o cuatro llegan a 90 o más en su
// carrera y uno de cada trescientos a 95 o más, y los cracks no se concentran en
// una posición (en particular, no en la portería).
func TestHayCracksYSeRepartenPorTodasLasPosiciones(t *testing.T) {
	picos := picosDeCarrera(t)
	var campo []int
	var tabla strings.Builder
	fmt.Fprintf(&tabla, "\n%-14s %6s %8s %7s %7s\n", "posicion", "n", "mediana", ">=90", ">=95")
	porPos := map[modelo.Posicion]float64{}
	for _, p := range modelo.Posiciones {
		v := picos[p]
		sort.Ints(v)
		if p != modelo.Portero {
			campo = append(campo, v...)
		}
		porPos[p] = porcentajeDesde(v, 90)
		fmt.Fprintf(&tabla, "%-14s %6d %8d %6.2f%% %6.2f%%\n", p, len(v), v[len(v)/2], porPos[p], porcentajeDesde(v, 95))
	}
	sort.Ints(campo)
	fmt.Fprintf(&tabla, "%-14s %6d %8d %6.2f%% %6.2f%%\n", "campo (todos)", len(campo), campo[len(campo)/2], porcentajeDesde(campo, 90), porcentajeDesde(campo, 95))
	t.Log(tabla.String())

	if p := porcentajeDesde(campo, 90); p < 2.8 || p > 4.6 {
		t.Errorf("%.2f%% de los jugadores de campo llegan a 90 o mas; el objetivo es ~3-4%%", p)
	}
	if p := porcentajeDesde(campo, 95); p < 0.15 || p > 0.6 {
		t.Errorf("%.2f%% llegan a 95 o mas; el objetivo es ~0,33%%", p)
	}
	// Reparto: ninguna posición concentra los cracks, ni los porteros.
	menor, mayor := math.Inf(1), math.Inf(-1)
	for p, v := range porPos {
		menor, mayor = math.Min(menor, v), math.Max(mayor, v)
		if v < 1.8 || v > 6.5 {
			t.Errorf("%v: %.2f%% con 90 o mas, fuera de [1,8, 6,5]", p, v)
		}
	}
	if mayor > 2.8*menor {
		t.Errorf("los cracks se concentran: de %.2f%% a %.2f%% segun la posicion", menor, mayor)
	}
	if porPos[modelo.Portero] > 1.6*porcentajeDesde(campo, 90) {
		t.Errorf("los porteros tienen demasiados cracks (%.2f%%) frente al campo (%.2f%%)", porPos[modelo.Portero], porcentajeDesde(campo, 90))
	}
}

// victoriasLiga juega partidos de ida y vuelta alternados entre dos equipos de una
// liga real y devuelve el porcentaje de victorias del primero.
func victoriasLiga(mio, rival modelo.Equipo, partidos int, semilla int64) float64 {
	r := rand.New(rand.NewSource(semilla))
	g := 0
	for i := 0; i < partidos; i++ {
		var gm, gr int
		if i%2 == 0 {
			gm, gr = simulacion.Marcador(r, mio, rival)
		} else {
			gr, gm = simulacion.Marcador(r, rival, mio)
		}
		if gm > gr {
			g++
		}
	}
	return 100 * float64(g) / float64(partidos)
}

// mejorDeLaPlantilla devuelve el índice del k-ésimo mejor de los n titulares de
// la posición, por valoración.
func mejorDeLaPlantilla(e modelo.Equipo, p modelo.Posicion, n, k int) int {
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
	return idx[min(k, n-1)]
}

func crack(e modelo.Equipo, i int) {
	e.Plantilla[i].Atributos = modelo.Atributos{Ritmo: 95, Tiro: 95, Pase: 95, Regate: 95, Defensa: 95, Fisico: 95, Reflejos: 95}
}

// TestUnCrackPesaOchoPuntosEnUnaLigaReal mide, con equipos de ligas generadas de
// verdad (que ya tienen sus propios cracks), cuánto suma un crack de 95: el
// objetivo de diseño son unos +8 puntos de victoria en cada posición.
func TestUnCrackPesaOchoPuntosEnUnaLigaReal(t *testing.T) {
	type puesto struct {
		nombre string
		p      modelo.Posicion
		n      int
	}
	puestos := []puesto{
		{"portero", modelo.Portero, 1}, {"defensa", modelo.Defensa, 4},
		{"medio", modelo.Mediocampista, 3}, {"delantero", modelo.Delantero, 3},
	}
	const ligas, partidos = 6, 4000
	efecto := map[string]float64{}
	var tres, once, control float64
	for s := int64(1); s <= ligas; s++ {
		c := nuevaCarrera(t, s*29, 10)
		base := clonarEquipos(c.Temporada.Equipos)[(s*3)%10]
		rival := clonarEquipos([]modelo.Equipo{base})[0]
		sin := victoriasLiga(clonarEquipos([]modelo.Equipo{base})[0], rival, partidos, s)
		control += sin / ligas
		for _, pu := range puestos {
			mio := clonarEquipos([]modelo.Equipo{base})[0]
			crack(mio, mejorDeLaPlantilla(mio, pu.p, pu.n, 0))
			efecto[pu.nombre] += (victoriasLiga(mio, rival, partidos, s) - sin) / ligas
		}
		m3 := clonarEquipos([]modelo.Equipo{base})[0]
		crack(m3, mejorDeLaPlantilla(m3, modelo.Portero, 1, 0))
		crack(m3, mejorDeLaPlantilla(m3, modelo.Mediocampista, 3, 0))
		crack(m3, mejorDeLaPlantilla(m3, modelo.Delantero, 3, 0))
		tres += (victoriasLiga(m3, rival, partidos, s) - sin) / ligas
		m11 := clonarEquipos([]modelo.Equipo{base})[0]
		for _, pu := range puestos {
			for k := 0; k < pu.n; k++ {
				crack(m11, mejorDeLaPlantilla(m11, pu.p, pu.n, k))
			}
		}
		once += victoriasLiga(m11, rival, partidos, s) / ligas
	}
	t.Logf("equipo sin cracks extra: %.1f%% de victorias contra su copia; un crack de 95 suma (puntos de victoria): %v", control, efecto)
	t.Logf("tres cracks suman %.1f; un once de cracks gana el %.1f%%", tres, once)
	menor, mayor := math.Inf(1), math.Inf(-1)
	for pos, e := range efecto {
		if e < 5 || e > 11.5 {
			t.Errorf("un crack de %s suma %.1f puntos de victoria; el objetivo es ~8 (entre 5 y 11,5)", pos, e)
		}
		menor, mayor = math.Min(menor, e), math.Max(mayor, e)
	}
	if mayor > 1.6*menor {
		t.Errorf("el efecto de un crack varia demasiado entre posiciones: de %.1f a %.1f", menor, mayor)
	}
	if tres < 14 || tres > 27 {
		t.Errorf("tres cracks suman %.1f; el objetivo es ~20 (entre 14 y 27)", tres)
	}
	if once < 60 || once > 82 {
		t.Errorf("un once de cracks gana el %.1f%%; deberia estar entre 60%% y 82%%", once)
	}
}

// TestLaLigaSigueAbiertaConLosCracks comprueba que pesar más los cracks no
// vuelve la liga predecible: el campeón no domina y el mejor equipo no gana
// siempre, pero la calidad sí cuenta.
func TestLaLigaSigueAbiertaConLosCracks(t *testing.T) {
	const temporadas = 90
	var campeon, ultimo float64
	mejorGana := 0
	for s := int64(1); s <= temporadas; s++ {
		c := nuevaCarrera(t, s*37, 10)
		// Una liga ya envejecida (sin jugar los partidos de esas temporadas).
		envejecida, err := liga.Nueva(ligaEnvejecida(t, s*37, 2))
		if err != nil {
			t.Fatal(err)
		}
		c.Temporada = envejecida
		mejor, mejorV := "", -1
		for _, e := range c.Temporada.Equipos {
			if v := e.Valoracion(); v > mejorV {
				mejor, mejorV = e.Nombre, v
			}
		}
		terminar(t, c)
		tabla := c.Tabla()
		campeon += float64(tabla[0].Pts) / temporadas
		ultimo += float64(tabla[len(tabla)-1].Pts) / temporadas
		if tabla[0].Equipo == mejor {
			mejorGana++
		}
	}
	pct := 100 * float64(mejorGana) / temporadas
	t.Logf("campeon %.1f puntos de 54, ultimo %.1f; el equipo de mayor valoracion gana la liga el %.0f%% (azar puro: 10%%)", campeon, ultimo, pct)
	if campeon > 40 {
		t.Errorf("el campeon suma %.1f de 54 puntos: la liga esta dominada", campeon)
	}
	if ultimo < 10 {
		t.Errorf("el ultimo suma solo %.1f puntos", ultimo)
	}
	if pct < 17 || pct > 50 {
		t.Errorf("el equipo de mayor valoracion gana el %.0f%% de las ligas; deberia estar entre 17%% y 50%%", pct)
	}
}
