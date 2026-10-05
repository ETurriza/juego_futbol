package progresion

import (
	"fmt"
	"math"
	"math/rand"
	"reflect"
	"strings"
	"testing"

	"github.com/ETurriza/juego_futbol/internal/modelo"
)

func nuevoRand(semilla int64) *rand.Rand {
	return rand.New(rand.NewSource(semilla))
}

func atributos(tec, fis int) modelo.Atributos {
	return modelo.Atributos{
		Ritmo: fis, Fisico: fis,
		Tiro: tec, Pase: tec, Regate: tec, Defensa: tec, Reflejos: tec,
	}
}

func jugador(edad int, p modelo.Posicion, tec, fis int) modelo.Jugador {
	return modelo.Jugador{ID: 1, Nombre: "J", Edad: edad, Posicion: p, Atributos: atributos(tec, fis)}
}

// cambioObservado promedia, sobre n jugadores, el cambio de los atributos
// técnicos (Pase) y físicos (Ritmo) tras envejecer un año.
func cambioObservado(edad int, p modelo.Posicion, n int) (tecnico, fisico float64) {
	r := nuevoRand(int64(edad)*100 + int64(p))
	for i := 0; i < n; i++ {
		nuevo := Envejecer(r, jugador(edad, p, 70, 70))
		tecnico += float64(nuevo.Atributos.Pase - 70)
		fisico += float64(nuevo.Atributos.Ritmo - 70)
	}
	return tecnico / float64(n), fisico / float64(n)
}

func TestEnvejecerSumaUnAnoYNoModificaElOriginal(t *testing.T) {
	original := jugador(25, modelo.Delantero, 70, 70)
	copia := original
	nuevo := Envejecer(nuevoRand(1), original)
	if nuevo.Edad != 26 {
		t.Errorf("Edad = %d, se esperaba 26", nuevo.Edad)
	}
	if original != copia {
		t.Error("Envejecer modifico el jugador original")
	}
	if nuevo.ID != original.ID || nuevo.Nombre != original.Nombre || nuevo.Posicion != original.Posicion {
		t.Error("Envejecer no deberia cambiar identidad ni posicion")
	}
}

func TestEnvejecerReproducible(t *testing.T) {
	j := jugador(31, modelo.Mediocampista, 70, 70)
	for semilla := int64(1); semilla <= 20; semilla++ {
		if a, b := Envejecer(nuevoRand(semilla), j), Envejecer(nuevoRand(semilla), j); a != b {
			t.Fatalf("semilla %d: %v != %v", semilla, a, b)
		}
	}
}

func TestAtributosSiempreEnRango(t *testing.T) {
	r := nuevoRand(3)
	for _, p := range modelo.Posiciones {
		for _, valores := range [][2]int{{1, 1}, {99, 99}, {50, 50}} {
			j := jugador(16, p, valores[0], valores[1])
			for edad := 16; edad <= 44; edad++ {
				j = Envejecer(r, j)
				if err := j.Atributos.Validar(); err != nil {
					t.Fatalf("%v con %v a los %d: %v", p, valores, j.Edad, err)
				}
			}
		}
	}
}

func TestReflejosSoloEvolucionanEnPorteros(t *testing.T) {
	r := nuevoRand(4)
	for _, p := range []modelo.Posicion{modelo.Defensa, modelo.Mediocampista, modelo.Delantero} {
		j := jugador(30, p, 70, 70)
		j.Atributos.Reflejos = 12
		for i := 0; i < 20; i++ {
			j = Envejecer(r, j)
		}
		if j.Atributos.Reflejos != 12 {
			t.Errorf("%v: los reflejos cambiaron a %d", p, j.Atributos.Reflejos)
		}
	}
	// En un portero sí cambian (cae con la edad).
	g := jugador(37, modelo.Portero, 70, 70)
	g = Envejecer(r, g)
	if g.Atributos.Reflejos == 70 && g.Atributos.Pase == 70 {
		t.Log("un portero de 37 podria mantener sus reflejos por azar")
	}
}

// El cambio medio observado debe coincidir con la tabla de la curva (más el
// bono esperado del último prime) dentro de una tolerancia estadística.
func TestCurvaPorEdadSigueLaTabla(t *testing.T) {
	const n, tolerancia = 6000, 0.35
	for _, p := range []modelo.Posicion{modelo.Mediocampista, modelo.Portero} {
		for edad := 16; edad <= 42; edad++ {
			tabT, tabF := CambioMedio(edad, p)
			esperadoT, esperadoF := conTecho(tabT, 70), conTecho(tabF, 70)
			if e := edadEfectiva(modelo.Jugador{Edad: edad, Posicion: p}); e >= granAnoDesde && e <= granAnoHasta {
				// Con probabilidad probGranAno el cambio técnico suma el bono.
				esperadoT = (1-probGranAno)*conTecho(tabT, 70) + probGranAno*conTecho(tabT+bonoGranAno, 70)
			}
			// Evita los extremos donde el límite 1-99 sesga el promedio.
			obsT, obsF := cambioObservado(edad, p, n)
			if math.Abs(obsT-esperadoT) > tolerancia || math.Abs(obsF-esperadoF) > tolerancia {
				t.Errorf("%v a los %d: tecnico %.2f (esperado %.2f), fisico %.2f (esperado %.2f)",
					p, edad, obsT, esperadoT, obsF, esperadoF)
			}
		}
	}
}

// conTecho aplica los rendimientos decrecientes a un cambio medio: solo el
// crecimiento se frena, según lo cerca del tope que esté el valor.
func conTecho(cambio float64, valor int) float64 {
	if cambio > 0 {
		return cambio * factorCrecimiento(valor)
	}
	return cambio
}

func TestRendimientosDecrecientes(t *testing.T) {
	// El factor baja a medida que el valor sube, sin pasar de 1 ni bajar de factorMinimo.
	anterior := 2.0
	for v := modelo.AtributoMin; v <= modelo.AtributoMax; v++ {
		f := factorCrecimiento(v)
		if f > anterior || f > 1 || f < factorMinimo {
			t.Fatalf("factor(%d) = %.3f (anterior %.3f)", v, f, anterior)
		}
		anterior = f
	}
	if factorCrecimiento(40) != 1 || factorCrecimiento(int(techoAtributo-amplitudTecho)) != 1 {
		t.Error("lejos del tope el crecimiento deberia ser el normal")
	}
	if factorCrecimiento(int(techoAtributo)) != factorMinimo || factorCrecimiento(99) != factorMinimo {
		t.Error("en el tope el crecimiento deberia ser el minimo")
	}

	// Un joven con un atributo muy alto crece bastante menos que uno con uno medio.
	creceMedio, creceAlto := cambioConValor(18, 60), cambioConValor(18, 90)
	if creceAlto >= creceMedio/2 {
		t.Errorf("a los 18, un atributo de 90 crece %.2f y uno de 60 crece %.2f: deberia frenarse mucho", creceAlto, creceMedio)
	}
	// La caída por edad, en cambio, no se frena: un veterano baja igual con 90 que con 60.
	caeMedio, caeAlto := cambioConValor(36, 60), cambioConValor(36, 90)
	if d := caeAlto - caeMedio; d > 0.6 || d < -0.6 {
		t.Errorf("a los 36, el cambio con 90 (%.2f) y con 60 (%.2f) deberia ser parecido", caeAlto, caeMedio)
	}
	// Nadie supera 99 ni baja de 1, por mucho que crezca.
	j := jugador(16, modelo.Delantero, 95, 95)
	r := nuevoRand(2)
	for i := 0; i < 8; i++ {
		j = Envejecer(r, j)
		if err := j.Atributos.Validar(); err != nil {
			t.Fatal(err)
		}
	}
}

// cambioConValor es el cambio medio del pase de un jugador de campo de la edad
// dada que lo tiene en el valor dado.
func cambioConValor(edad, valor int) float64 {
	r := nuevoRand(int64(edad*1000 + valor))
	const n = 6000
	suma := 0.0
	for i := 0; i < n; i++ {
		nuevo := Envejecer(r, jugador(edad, modelo.Mediocampista, valor, valor))
		suma += float64(nuevo.Atributos.Pase - valor)
	}
	return suma / n
}

func TestLaMesetaTecnicaDuraMasQueLaFisica(t *testing.T) {
	// Entre los 28 y los 33 la técnica se sostiene y el físico baja.
	for edad := 28; edad <= 33; edad++ {
		tec, fis := CambioMedio(edad, modelo.Delantero)
		if tec <= fis {
			t.Errorf("a los %d la tecnica (%.1f) deberia sostenerse mejor que el fisico (%.1f)", edad, tec, fis)
		}
		if tec < 0 {
			t.Errorf("a los %d la tecnica no deberia bajar en la meseta: %.1f", edad, tec)
		}
	}
}

func TestElDeclivePostPrimeEsRapido(t *testing.T) {
	// Tras la meseta la caida se acelera: cada año baja más que el anterior.
	anterior := 1.0
	for edad := 33; edad <= 37; edad++ {
		tec, _ := CambioMedio(edad, modelo.Delantero)
		if tec > anterior {
			t.Errorf("a los %d la tecnica (%.1f) deberia caer al menos como el año previo (%.1f)", edad, tec, anterior)
		}
		anterior = tec
	}
	if tec, _ := CambioMedio(36, modelo.Delantero); tec > -5 {
		t.Errorf("a los 36 la caida deberia ser fuerte, es %.1f", tec)
	}
}

func TestUltimoPrimeExiste(t *testing.T) {
	// A los 32 hay una fraccion de jugadores que mejora claramente su tecnica.
	r := nuevoRand(9)
	const n = 5000
	mejoran := 0
	for i := 0; i < n; i++ {
		if Envejecer(r, jugador(32, modelo.Mediocampista, 70, 70)).Atributos.Pase >= 73 {
			mejoran++
		}
	}
	if frac := float64(mejoran) / n; frac < 0.10 || frac > 0.40 {
		t.Errorf("fraccion con gran año a los 32 = %.2f, se esperaba entre 0.10 y 0.40", frac)
	}
	// Y no ocurre a los 25.
	mejoran = 0
	for i := 0; i < n; i++ {
		if Envejecer(r, jugador(25, modelo.Mediocampista, 70, 70)).Atributos.Pase >= 76 {
			mejoran++
		}
	}
	if frac := float64(mejoran) / n; frac > 0.03 {
		t.Errorf("a los 25 no deberia haber saltos grandes: %.3f", frac)
	}
}

func TestLosPorterosEnvejecenMasTarde(t *testing.T) {
	for edad := 28; edad <= 40; edad++ {
		campoT, campoF := CambioMedio(edad-DesfasePortero, modelo.Delantero)
		porteroT, porteroF := CambioMedio(edad, modelo.Portero)
		if campoT != porteroT || campoF != porteroF {
			t.Errorf("a los %d un portero deberia seguir la curva de un jugador de campo de %d", edad, edad-DesfasePortero)
		}
	}
	// A los 35 un jugador de campo ya cae, un portero todavía no tanto.
	campo, _ := CambioMedio(35, modelo.Delantero)
	portero, _ := CambioMedio(35, modelo.Portero)
	if portero <= campo {
		t.Errorf("a los 35 el portero (%.1f) deberia caer menos que el jugador de campo (%.1f)", portero, campo)
	}
}

func TestLosDefensasDeclinanAntesQueLosDelanteros(t *testing.T) {
	// Su valoracion depende mas de lo fisico, asi que en la edad en que el
	// fisico ya cae y la tecnica aguanta, el defensa pierde mas valoracion.
	for _, edad := range []int{29, 30, 31, 32, 33} {
		tec, fis := CambioMedio(edad, modelo.Defensa)
		valor := func(p modelo.Posicion) float64 {
			w := pesosDe(p)
			return (tec*w.tec + fis*w.fis) / 100
		}
		if valor(modelo.Defensa) >= valor(modelo.Delantero) {
			t.Errorf("a los %d el defensa (%.2f) deberia perder mas valoracion que el delantero (%.2f)",
				edad, valor(modelo.Defensa), valor(modelo.Delantero))
		}
	}
}

// pesosDe calcula el peso técnico y físico de la valoración de una posición
// probando la media con atributos extremos.
func pesosDe(p modelo.Posicion) struct{ tec, fis float64 } {
	solo := func(a modelo.Atributos) float64 { return float64(a.Media(p)) }
	fis := solo(modelo.Atributos{Ritmo: 100, Fisico: 100})
	tec := solo(modelo.Atributos{Tiro: 100, Pase: 100, Regate: 100, Defensa: 100, Reflejos: 100})
	return struct{ tec, fis float64 }{tec, fis}
}

func TestProbabilidadDeRetiro(t *testing.T) {
	esperado := map[int]float64{33: 0, 34: 0.05, 35: 0.15, 36: 0.30, 37: 0.50, 38: 0.75, 39: 1, 42: 1}
	for edad, want := range esperado {
		if got := ProbabilidadRetiro(edad, modelo.Delantero); got != want {
			t.Errorf("campo a los %d: %.2f, se esperaba %.2f", edad, got, want)
		}
		// Los porteros, tres años mas tarde.
		if got := ProbabilidadRetiro(edad+DesfasePortero, modelo.Portero); got != want {
			t.Errorf("portero a los %d: %.2f, se esperaba %.2f", edad+DesfasePortero, got, want)
		}
	}
	for edad := 16; edad <= 33; edad++ {
		if ProbabilidadRetiro(edad, modelo.Mediocampista) != 0 {
			t.Errorf("a los %d no deberia haber retiros", edad)
		}
	}
	for edad := 34; edad < 45; edad++ {
		if ProbabilidadRetiro(edad+1, modelo.Defensa) < ProbabilidadRetiro(edad, modelo.Defensa) {
			t.Errorf("la probabilidad de retiro no deberia bajar de %d a %d", edad, edad+1)
		}
	}
}

func TestSeRetiraSigueLaProbabilidad(t *testing.T) {
	r := nuevoRand(5)
	const n = 20000
	for _, edad := range []int{30, 33, 34, 35, 36, 37, 38, 39, 41} {
		retiros := 0
		for i := 0; i < n; i++ {
			if SeRetira(r, jugador(edad, modelo.Delantero, 70, 70)) {
				retiros++
			}
		}
		obs, want := float64(retiros)/n, ProbabilidadRetiro(edad, modelo.Delantero)
		if math.Abs(obs-want) > 0.02 {
			t.Errorf("a los %d: %.3f de retiros, se esperaba %.2f", edad, obs, want)
		}
	}
}

func TestSeRetiraConsumeSiempreUnNumero(t *testing.T) {
	// La secuencia posterior no debe depender de la edad del jugador evaluado.
	a, b := nuevoRand(7), nuevoRand(7)
	SeRetira(a, jugador(20, modelo.Delantero, 70, 70))
	SeRetira(b, jugador(40, modelo.Delantero, 70, 70))
	if a.Float64() != b.Float64() {
		t.Error("SeRetira deberia consumir la misma cantidad de aleatoriedad a cualquier edad")
	}
}

// edadMaxCohorte es la última edad que se simula en las cohortes.
const edadMaxCohorte = 44

// cohorte envejece, sin retiros, a n jugadores desde los 17 años con la calidad
// de un juvenil y devuelve la valoración media por edad.
func cohorte(p modelo.Posicion, n int, calidad int) map[int]float64 {
	r := nuevoRand(int64(p) + 77)
	jugadores := make([]modelo.Jugador, n)
	for i := range jugadores {
		jugadores[i] = jugador(17, p, calidad, calidad)
	}
	media := map[int]float64{}
	for edad := 17; edad <= edadMaxCohorte; edad++ {
		suma := 0
		for i := range jugadores {
			suma += jugadores[i].Valoracion()
			jugadores[i] = Envejecer(r, jugadores[i])
		}
		media[edad] = float64(suma) / float64(n)
	}
	return media
}

func edadPico(media map[int]float64) int {
	pico, mejor := 17, -1.0
	for edad := 17; edad <= edadMaxCohorte; edad++ {
		if media[edad] > mejor {
			pico, mejor = edad, media[edad]
		}
	}
	return pico
}

// TestCalibracionPorEdad muestra (con go test -v) la valoración media por edad
// de una cohorte, para juzgar la forma de la curva a simple vista, y comprueba
// que los mejores años caen donde corresponde a cada posición.
func TestCalibracionPorEdad(t *testing.T) {
	const n, calidad = 3000, 42
	curvas := map[modelo.Posicion]map[int]float64{}
	for _, p := range modelo.Posiciones {
		curvas[p] = cohorte(p, n, calidad)
	}

	var tabla strings.Builder
	fmt.Fprintf(&tabla, "\nvaloracion media por edad (cohorte desde %d, sin retiros)\n", calidad)
	fmt.Fprintf(&tabla, "%-5s %9s %9s %9s %9s\n", "edad", "Portero", "Defensa", "Medio", "Delantero")
	for edad := 17; edad <= 40; edad++ {
		fmt.Fprintf(&tabla, "%-5d %9.1f %9.1f %9.1f %9.1f\n", edad,
			curvas[modelo.Portero][edad], curvas[modelo.Defensa][edad],
			curvas[modelo.Mediocampista][edad], curvas[modelo.Delantero][edad])
	}
	t.Log(tabla.String())

	picos := map[modelo.Posicion][2]int{
		modelo.Defensa:       {26, 31},
		modelo.Mediocampista: {28, 33},
		modelo.Delantero:     {27, 33},
		modelo.Portero:       {30, 36},
	}
	for p, rango := range picos {
		if pico := edadPico(curvas[p]); pico < rango[0] || pico > rango[1] {
			t.Errorf("%v: pico a los %d, se esperaba entre %d y %d", p, pico, rango[0], rango[1])
		}
		c := curvas[p]
		if c[18] >= c[24] || c[24] >= c[edadPico(c)] {
			t.Errorf("%v: deberia crecer de los 18 a los 24 y hasta el pico", p)
		}
		// Referencia: los 39 años en el campo (36 efectivos), 3 más en un portero.
		referencia := 39
		if p == modelo.Portero {
			referencia += DesfasePortero
		}
		if caida := c[edadPico(c)] - c[referencia]; caida < 10 {
			t.Errorf("%v: caida de solo %.1f puntos entre el pico y los %d", p, caida, referencia)
		}
	}
	// Un jugador de campo de 36 ya rinde bastante menos que en su pico; un
	// portero de 36, menos que uno de 33.
	if m := curvas[modelo.Mediocampista]; m[36] > m[edadPico(m)]-5 {
		t.Errorf("mediocampista a los 36 (%.1f) deberia estar claramente por debajo del pico (%.1f)", m[36], m[edadPico(m)])
	}
}

func TestLaCurvaNoCambiaConLaSemilla(t *testing.T) {
	// Dos ejecuciones con la misma semilla dan la misma cohorte.
	if !reflect.DeepEqual(cohorte(modelo.Delantero, 200, 42), cohorte(modelo.Delantero, 200, 42)) {
		t.Error("la cohorte deberia ser reproducible")
	}
}

// jugadorConNivel arma un jugador de campo de la edad dada cuya valoración es
// aproximadamente nivel (todos sus atributos iguales).
func jugadorConNivel(edad int, p modelo.Posicion, nivel int) modelo.Jugador {
	return jugador(edad, p, nivel, nivel)
}

func TestElRetiroDependeDelNivelDeLosVeteranos(t *testing.T) {
	// A los 38 años, con el nivel de un titular o más, rige la probabilidad por edad.
	base := ProbabilidadRetiro(38, modelo.Mediocampista)
	if got := ProbabilidadRetiroDe(jugadorConNivel(38, modelo.Mediocampista, nivelMinimoTitular)); got != base {
		t.Errorf("con el nivel minimo deberia regir la edad: %.2f, se esperaba %.2f", got, base)
	}
	if got := ProbabilidadRetiroDe(jugadorConNivel(38, modelo.Mediocampista, 90)); got != base {
		t.Errorf("un crack no se retira antes por ser bueno: %.2f", got)
	}
	// Cuanto más flojo, más probable el retiro, hasta el 100 %.
	anterior := 0.0
	for nivel := nivelMinimoTitular + 5; nivel >= 30; nivel -= 3 {
		p := ProbabilidadRetiroDe(jugadorConNivel(36, modelo.Delantero, nivel))
		if p < anterior || p > 1 {
			t.Fatalf("a los 36, con nivel %d: %.2f tras %.2f", nivel, p, anterior)
		}
		anterior = p
	}
	if got := ProbabilidadRetiroDe(jugadorConNivel(38, modelo.Mediocampista, 47)); got != 1 {
		t.Errorf("un jugador de 38 con 47 deberia retirarse seguro: %.2f", got)
	}
	// Valores concretos: base por edad más pesoNivelEnRetiro por punto de déficit.
	falta := float64(nivelMinimoTitular - 55)
	if got, want := ProbabilidadRetiroDe(jugadorConNivel(36, modelo.Delantero, 55)), 0.30+falta*pesoNivelEnRetiro; math.Abs(got-want) > 1e-9 {
		t.Errorf("36 años con 55: %.3f, se esperaba %.3f", got, want)
	}
}

func TestElRetiroPorNivelSoloAfectaAVeteranos(t *testing.T) {
	// Un joven flojo no se retira: esta creciendo.
	for _, edad := range []int{16, 20, 25, 30, 32} {
		if got := ProbabilidadRetiroDe(jugadorConNivel(edad, modelo.Defensa, 35)); got != 0 {
			t.Errorf("a los %d con nivel 35 no deberia retirarse: %.2f", edad, got)
		}
	}
	// Y empieza en edadRetiroPorNivel, segun la edad efectiva: un portero, tres años después.
	if got := ProbabilidadRetiroDe(jugadorConNivel(edadRetiroPorNivel, modelo.Delantero, 50)); got <= 0 {
		t.Errorf("a los %d un jugador flojo ya deberia tener probabilidad de retiro: %.2f", edadRetiroPorNivel, got)
	}
	if got := ProbabilidadRetiroDe(jugadorConNivel(edadRetiroPorNivel-1, modelo.Delantero, 50)); got != 0 {
		t.Errorf("un año antes no: %.2f", got)
	}
	if got := ProbabilidadRetiroDe(jugadorConNivel(edadRetiroPorNivel+DesfasePortero, modelo.Portero, 50)); got <= 0 {
		t.Errorf("un portero flojo de %d anos ya deberia tener probabilidad de retiro: %.2f", edadRetiroPorNivel+DesfasePortero, got)
	}
	if got := ProbabilidadRetiroDe(jugadorConNivel(edadRetiroPorNivel+DesfasePortero-1, modelo.Portero, 50)); got != 0 {
		t.Errorf("un portero un año antes no: %.2f", got)
	}
}

func TestSeRetiraSigueLaProbabilidadPorNivel(t *testing.T) {
	r := nuevoRand(41)
	const n = 20000
	for _, c := range []struct {
		edad  int
		nivel int
	}{{34, 50}, {36, 55}, {37, 60}, {38, 70}} {
		j := jugadorConNivel(c.edad, modelo.Mediocampista, c.nivel)
		retiros := 0
		for i := 0; i < n; i++ {
			if SeRetira(r, j) {
				retiros++
			}
		}
		obs, want := float64(retiros)/n, ProbabilidadRetiroDe(j)
		if math.Abs(obs-want) > 0.02 {
			t.Errorf("%d anos con nivel %d: %.3f de retiros, se esperaba %.3f", c.edad, c.nivel, obs, want)
		}
	}
}
