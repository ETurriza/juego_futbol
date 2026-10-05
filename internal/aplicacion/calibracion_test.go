package aplicacion

import (
	"fmt"
	"math"
	"math/rand"
	"strings"
	"testing"

	"github.com/ETurriza/juego_futbol/internal/modelo"
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
		if c.pctBajos() > 8 || p.pctBajos() > 8 {
			t.Errorf("tras %d temporadas hay demasiados jugadores flojos (<=50): campo %.1f%%, porteros %.1f%%", temporadas, c.pctBajos(), p.pctBajos())
		}
		if d := p.media() - c.media(); d > 4 || d < -4 {
			t.Errorf("tras %d temporadas los porteros (%.1f) y el campo (%.1f) difieren demasiado", temporadas, p.media(), c.media())
		}
		if c.desv() > 12 || p.desv() > 13 {
			t.Errorf("tras %d temporadas la dispersion es excesiva: campo %.1f, porteros %.1f", temporadas, c.desv(), p.desv())
		}
		// Los atributos no se apilan en el tope de 99.
		if float64(c.topes)/float64(c.n) > 0.02 || float64(p.topes)/float64(p.n) > 0.04 {
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
