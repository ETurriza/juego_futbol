package generador

import (
	"math/rand"
	"reflect"
	"sort"
	"testing"

	"github.com/ETurriza/juego_futbol/internal/modelo"
)

func nuevoRand(semilla int64) *rand.Rand {
	return rand.New(rand.NewSource(semilla))
}

func TestEquipoReproducible(t *testing.T) {
	a := Equipo(nuevoRand(42), "Deportivo Prueba", 1)
	b := Equipo(nuevoRand(42), "Deportivo Prueba", 1)
	if !reflect.DeepEqual(a, b) {
		t.Fatal("la misma semilla deberia dar el mismo equipo")
	}
	c := Equipo(nuevoRand(43), "Deportivo Prueba", 1)
	if reflect.DeepEqual(a, c) {
		t.Fatal("semillas distintas deberian dar equipos distintos")
	}
}

func TestEquipoComposicionYValidez(t *testing.T) {
	for semilla := int64(1); semilla <= 50; semilla++ {
		e := Equipo(nuevoRand(semilla), NombreEquipo(nuevoRand(semilla)), 100)
		if err := e.Validar(); err != nil {
			t.Fatalf("semilla %d: equipo invalido: %v", semilla, err)
		}
		if len(e.Plantilla) != TamanoPlantilla {
			t.Fatalf("semilla %d: plantilla de %d", semilla, len(e.Plantilla))
		}
		for _, c := range composicion {
			if got := e.Contar(c.posicion); got != c.cantidad {
				t.Errorf("semilla %d: %d de %v, se esperaban %d", semilla, got, c.posicion, c.cantidad)
			}
		}
	}
}

func TestEquipoIDsConsecutivasYNombresUnicos(t *testing.T) {
	e := Equipo(nuevoRand(7), "Club Prueba", 500)
	nombres := map[string]bool{}
	for i, j := range e.Plantilla {
		if j.ID != 500+i {
			t.Errorf("jugador %d tiene ID %d, se esperaba %d", i, j.ID, 500+i)
		}
		if nombres[j.Nombre] {
			t.Errorf("nombre repetido en el equipo: %s", j.Nombre)
		}
		nombres[j.Nombre] = true
	}
}

func TestAtributosEnRango(t *testing.T) {
	r := nuevoRand(1)
	for i := 0; i < 2000; i++ {
		for _, p := range modelo.Posiciones {
			if err := Atributos(r, p).Validar(); err != nil {
				t.Fatalf("%v: %v", p, err)
			}
		}
	}
}

func TestPerfilesPorPosicion(t *testing.T) {
	r := nuevoRand(2)
	const n = 1000
	var reflejosPortero, reflejosDefensa, tiroDelantero, tiroDefensa int
	for i := 0; i < n; i++ {
		reflejosPortero += Atributos(r, modelo.Portero).Reflejos
		reflejosDefensa += Atributos(r, modelo.Defensa).Reflejos
		tiroDelantero += Atributos(r, modelo.Delantero).Tiro
		tiroDefensa += Atributos(r, modelo.Defensa).Tiro
	}
	if reflejosPortero <= reflejosDefensa {
		t.Error("los porteros deberian tener mas reflejos que los defensas")
	}
	if tiroDelantero <= tiroDefensa {
		t.Error("los delanteros deberian tener mas tiro que los defensas")
	}
}

func TestJugadorValido(t *testing.T) {
	r := nuevoRand(3)
	for i := 0; i < 500; i++ {
		for _, p := range modelo.Posiciones {
			if err := Jugador(r, i, p).Validar(); err != nil {
				t.Fatalf("%v: %v", p, err)
			}
		}
	}
}

func TestListasSinDuplicados(t *testing.T) {
	for nombre, lista := range map[string][]string{
		"nombres": nombres, "apellidos": apellidos, "ciudades": ciudades, "prefijos": prefijos,
	} {
		vistos := map[string]bool{}
		for _, s := range lista {
			if s == "" {
				t.Errorf("%s: entrada vacia", nombre)
			}
			if vistos[s] {
				t.Errorf("%s: duplicado %q", nombre, s)
			}
			vistos[s] = true
		}
	}
}

func TestNombreEquipoReproducible(t *testing.T) {
	if NombreEquipo(nuevoRand(9)) != NombreEquipo(nuevoRand(9)) {
		t.Error("NombreEquipo deberia ser reproducible con la misma semilla")
	}
	if NombreJugador(nuevoRand(9)) != NombreJugador(nuevoRand(9)) {
		t.Error("NombreJugador deberia ser reproducible con la misma semilla")
	}
}

func TestJuvenilValidoYEnRango(t *testing.T) {
	r := nuevoRand(11)
	for i := 0; i < 500; i++ {
		for _, p := range modelo.Posiciones {
			j := Juvenil(r, i, p)
			if err := j.Validar(); err != nil {
				t.Fatalf("%v: %v", p, err)
			}
			if j.Edad < EdadJuvenilMin || j.Edad > EdadJuvenilMax {
				t.Fatalf("edad %d fuera de [%d, %d]", j.Edad, EdadJuvenilMin, EdadJuvenilMax)
			}
			if j.ID != i || j.Posicion != p {
				t.Fatalf("ID o posicion incorrectas: %+v", j)
			}
		}
	}
}

func TestJuvenilReproducible(t *testing.T) {
	if Juvenil(nuevoRand(5), 1, modelo.Defensa) != Juvenil(nuevoRand(5), 1, modelo.Defensa) {
		t.Error("la misma semilla deberia dar el mismo juvenil")
	}
	if Juvenil(nuevoRand(5), 1, modelo.Defensa) == Juvenil(nuevoRand(6), 1, modelo.Defensa) {
		t.Error("semillas distintas deberian dar juveniles distintos")
	}
}

func TestLosJuvenilesSonMasFlojosQueLosJugadoresHechos(t *testing.T) {
	r := nuevoRand(12)
	const n = 1000
	for _, p := range modelo.Posiciones {
		var juveniles, hechos int
		for i := 0; i < n; i++ {
			juveniles += Juvenil(r, i, p).Valoracion()
			hechos += Jugador(r, i, p).Valoracion()
		}
		if juveniles+10*n > hechos {
			t.Errorf("%v: los juveniles (%.1f) deberian valer al menos 10 puntos menos que los hechos (%.1f)",
				p, float64(juveniles)/n, float64(hechos)/n)
		}
	}
}

func TestMinimoPorPosicionCoincideConLaPlantillaGenerada(t *testing.T) {
	e := Equipo(nuevoRand(1), "Club Prueba", 1)
	total := 0
	for _, p := range modelo.Posiciones {
		if MinimoPorPosicion(p) != e.Contar(p) {
			t.Errorf("%v: minimo %d, plantilla generada %d", p, MinimoPorPosicion(p), e.Contar(p))
		}
		total += MinimoPorPosicion(p)
	}
	if total != TamanoPlantilla {
		t.Errorf("los minimos suman %d, se esperaban %d", total, TamanoPlantilla)
	}
	if MinimoPorPosicion(modelo.Posicion(99)) != 0 {
		t.Error("una posicion invalida deberia tener minimo 0")
	}
}

// mediaPorEdad promedia la valoración de n jugadores creados por f, agrupados por
// tramos de edad.
func mediaPorEdad(n int, f func(i int) modelo.Jugador) map[string]float64 {
	suma, cuenta := map[string]float64{}, map[string]float64{}
	for i := 0; i < n; i++ {
		j := f(i)
		var tramo string
		switch {
		case j.Edad <= 19:
			tramo = "<=19"
		case j.Edad <= 25:
			tramo = "20-25"
		case j.Edad <= 31:
			tramo = "26-31"
		case j.Edad <= 35:
			tramo = "32-35"
		default:
			tramo = ">=36"
		}
		suma[tramo] += float64(j.Valoracion())
		cuenta[tramo]++
	}
	out := map[string]float64{}
	for t, s := range suma {
		out[t] = s / cuenta[t]
	}
	return out
}

func TestLaCalidadDeLosJugadoresSigueLaCurvaDeEdad(t *testing.T) {
	r := nuevoRand(31)
	media := mediaPorEdad(20000, func(i int) modelo.Jugador { return Jugador(r, i, modelo.Delantero) })
	t.Logf("valoracion media de un delantero por edad: %v", media)
	// Los jóvenes valen mucho menos que los jugadores en su mejor momento, y los
	// veteranos de 36 o más han caído por debajo.
	if media["26-31"]-media["<=19"] < 12 {
		t.Errorf("de <=19 (%.1f) a 26-31 (%.1f) deberia haber al menos 12 puntos", media["<=19"], media["26-31"])
	}
	if media["20-25"] <= media["<=19"] || media["26-31"] <= media["20-25"] {
		t.Errorf("la calidad deberia crecer con la edad hasta la madurez: %v", media)
	}
	if media[">=36"] > media["26-31"]-4 {
		t.Errorf("a partir de los 36 (%.1f) deberia haber una caida clara respecto de 26-31 (%.1f)", media[">=36"], media["26-31"])
	}
}

func TestLosJuvenilesCrecenConLaEdad(t *testing.T) {
	// Un juvenil de 19 años ha pasado tres años de crecimiento desde los 16: vale
	// bastante mas que uno de 16, igual que lo haria cualquier jugador de la liga.
	r := nuevoRand(32)
	suma, cuenta := map[int]float64{}, map[int]float64{}
	for i := 0; i < 20000; i++ {
		j := Juvenil(r, i, modelo.Mediocampista)
		suma[j.Edad] += float64(j.Valoracion())
		cuenta[j.Edad]++
	}
	var anterior float64
	for edad := EdadJuvenilMin; edad <= EdadJuvenilMax; edad++ {
		if cuenta[edad] == 0 {
			t.Fatalf("no se genero ningun juvenil de %d anos", edad)
		}
		media := suma[edad] / cuenta[edad]
		t.Logf("juvenil de %d: valoracion media %.1f", edad, media)
		if edad > EdadJuvenilMin && media < anterior+3 {
			t.Errorf("a los %d (%.1f) deberia valer al menos 3 puntos mas que a los %d (%.1f)", edad, media, edad-1, anterior)
		}
		anterior = media
	}
}

func TestUnJuvenilYUnJugadorDeLaMismaEdadSonIguales(t *testing.T) {
	// Son el mismo proceso: la calidad de un jugador de 18 años no depende de
	// cómo se creó.
	r := nuevoRand(33)
	var juv, nJuv, jug, nJug float64
	for i := 0; nJuv < 6000 || nJug < 6000; i++ {
		if j := Juvenil(r, i, modelo.Defensa); j.Edad == 18 {
			juv += float64(j.Valoracion())
			nJuv++
		}
		if j := Jugador(r, i, modelo.Defensa); j.Edad == 18 {
			jug += float64(j.Valoracion())
			nJug++
		}
		if i > 400000 {
			t.Fatal("no salieron suficientes jugadores de 18 años")
		}
	}
	if d := juv/nJuv - jug/nJug; d > 1.5 || d < -1.5 {
		t.Errorf("un juvenil de 18 (%.1f) y un jugador de 18 (%.1f) deberian valer lo mismo", juv/nJuv, jug/nJug)
	}
}

func TestLasEdadesIniciales(t *testing.T) {
	r := nuevoRand(34)
	minC, maxC, minP, maxP := 99, 0, 99, 0
	for i := 0; i < 5000; i++ {
		for _, p := range modelo.Posiciones {
			j := Jugador(r, i, p)
			if p == modelo.Portero {
				minP, maxP = min(minP, j.Edad), max(maxP, j.Edad)
			} else {
				minC, maxC = min(minC, j.Edad), max(maxC, j.Edad)
			}
		}
	}
	// Los jugadores de campo, de 17 a 36; los porteros juegan tres años más.
	if minC != 17 || maxC != 36 || minP != 17 || maxP != 39 {
		t.Errorf("edades de campo %d-%d y de porteros %d-%d; se esperaban 17-36 y 17-39", minC, maxC, minP, maxP)
	}
}

func TestLaLigaInicialNoTieneVeteranosQueYaSeHabrianRetirado(t *testing.T) {
	// Un veterano muy flojo habria dejado el fútbol: pocos de 36 o más con
	// valoración de 50 o menos.
	r := nuevoRand(35)
	veteranos, flojos := 0, 0
	for i := 0; i < 40000; i++ {
		j := Jugador(r, i, modelo.Mediocampista)
		if j.Edad >= 36 {
			veteranos++
			if j.Valoracion() <= 50 {
				flojos++
			}
		}
	}
	if veteranos == 0 {
		t.Fatal("no salio ningun veterano")
	}
	if frac := float64(flojos) / float64(veteranos); frac > 0.04 {
		t.Errorf("%.1f%% de los veteranos de 36 o mas con valoracion <= 50; deberia ser casi ninguno", 100*frac)
	}
}

func TestJugadorSiempreDevuelveUnJugadorValido(t *testing.T) {
	// Aunque haya que reintentar (el jugador se retiraria antes de llegar a la
	// edad), siempre se devuelve uno con la ID pedida, de la edad y posicion
	// correctas y con datos validos.
	r := nuevoRand(36)
	for i := 0; i < 3000; i++ {
		for _, p := range modelo.Posiciones {
			j := Jugador(r, 1000+i, p)
			if err := j.Validar(); err != nil {
				t.Fatalf("%v: %v", p, err)
			}
			if j.ID != 1000+i || j.Posicion != p {
				t.Fatalf("ID o posicion incorrectas: %+v", j)
			}
		}
	}
}

func TestElTalentoSigueLaDistribucionEsperada(t *testing.T) {
	r := nuevoRand(51)
	const n = 100000
	var proy [4]int
	suma := 0
	valores := make([]int, 0, n)
	for i := 0; i < n; i++ {
		j := Juvenil(r, i, modelo.Delantero)
		if j.Talento < modelo.TalentoMin || j.Talento > modelo.TalentoMax {
			t.Fatalf("talento %d fuera de rango", j.Talento)
		}
		valores = append(valores, j.Talento)
		suma += j.Talento
		proy[j.Proyeccion()]++
	}
	sort.Ints(valores)
	t.Logf("talento: mediana %d, p10 %d, p90 %d, p99 %d; proyeccion limitada %.1f%%, normal %.1f%%, alta %.1f%%, excepcional %.1f%%",
		valores[n/2], valores[n/10], valores[n*9/10], valores[n*99/100],
		100*float64(proy[0])/n, 100*float64(proy[1])/n, 100*float64(proy[2])/n, 100*float64(proy[3])/n)
	if m := valores[n/2]; m < 98 || m > 102 {
		t.Errorf("la mediana del talento deberia ser ~100: %d", m)
	}
	// Excepcional: unos pocos; alta: bastantes más; limitada: una minoría.
	pct := func(i int) float64 { return 100 * float64(proy[i]) / n }
	if p := pct(int(modelo.ProyeccionExcepcional)); p < 2.5 || p > 5.5 {
		t.Errorf("proyeccion excepcional: %.1f%%, se esperaba entre 2,5%% y 5,5%%", p)
	}
	if p := pct(int(modelo.ProyeccionAlta)); p < 8 || p > 18 {
		t.Errorf("proyeccion alta: %.1f%%, se esperaba entre 8%% y 18%%", p)
	}
	if p := pct(int(modelo.ProyeccionLimitada)); p < 8 || p > 25 {
		t.Errorf("proyeccion limitada: %.1f%%, se esperaba entre 8%% y 25%%", p)
	}
}

func TestElTalentoEsIndependienteDeLaPosicion(t *testing.T) {
	// La misma distribución en todas las posiciones: los cracks no se concentran.
	r := nuevoRand(52)
	const n = 30000
	medias := map[modelo.Posicion]float64{}
	excepcionales := map[modelo.Posicion]float64{}
	for _, p := range modelo.Posiciones {
		for i := 0; i < n; i++ {
			j := Juvenil(r, i, p)
			medias[p] += float64(j.Talento) / n
			if j.Proyeccion() == modelo.ProyeccionExcepcional {
				excepcionales[p] += 100.0 / n
			}
		}
	}
	for _, p := range modelo.Posiciones {
		if medias[p] < 99 || medias[p] > 104 || excepcionales[p] < 2.5 || excepcionales[p] > 5.5 {
			t.Errorf("%v: talento medio %.1f, excepcionales %.1f%%", p, medias[p], excepcionales[p])
		}
	}
}

func TestElTalentoSeConservaEnLosJugadoresIniciales(t *testing.T) {
	// Jugador crea un juvenil y lo envejece: su talento no se pierde ni se redibuja.
	r := nuevoRand(53)
	altos := 0
	for i := 0; i < 5000; i++ {
		j := Jugador(r, i, modelo.Mediocampista)
		if j.Talento == 0 {
			t.Fatal("un jugador generado debe tener talento")
		}
		if j.Talento >= 120 {
			altos++
		}
	}
	if altos < 300 || altos > 1200 {
		t.Errorf("%d de 5000 jugadores con talento >= 120", altos)
	}
}

func TestElTalentoSeNotaEnLaCalidadAlMadurar(t *testing.T) {
	// A igualdad de edad (26 a 31), los de talento alto valen mucho más que los de
	// talento bajo: es lo que hace a un crack.
	r := nuevoRand(54)
	var altos, bajos, nA, nB float64
	for i := 0; i < 60000; i++ {
		j := Jugador(r, i, modelo.Mediocampista)
		if j.Edad < 26 || j.Edad > 31 {
			continue
		}
		switch {
		case j.Talento >= 125:
			altos += float64(j.Valoracion())
			nA++
		case j.Talento <= 85:
			bajos += float64(j.Valoracion())
			nB++
		}
	}
	if nA < 100 || nB < 100 {
		t.Fatalf("muestras pequeñas: %v y %v", nA, nB)
	}
	dif := altos/nA - bajos/nB
	t.Logf("a los 26-31: talento alto %.1f, bajo %.1f (diferencia %.1f)", altos/nA, bajos/nB, dif)
	if dif < 10 {
		t.Errorf("la diferencia entre talento alto y bajo es de solo %.1f puntos", dif)
	}
}
