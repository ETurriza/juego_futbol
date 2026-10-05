package simulacion

import (
	"math"
	"math/rand"
	"sort"
	"testing"

	"github.com/ETurriza/juego_futbol/internal/modelo"
)

func formacionPtr(f modelo.Formacion) *modelo.Formacion { return &f }

func TestLaAlineacionAutomaticaEsValidaEnCadaFormacion(t *testing.T) {
	e := equipoRealista("A", 1, 62)
	for _, f := range modelo.Formaciones {
		al := AlineacionAutomatica(e, Criterios{Formacion: formacionPtr(f)})
		if al.Formacion != f {
			t.Errorf("%v: formacion %v", f, al.Formacion)
		}
		if err := al.Validar(e, nil); err != nil {
			t.Errorf("%v: la alineacion automatica no es valida: %v", f, err)
		}
		// Con una plantilla completa, cada puesto lo ocupa un jugador de su posición.
		pos := map[int]modelo.Posicion{}
		for _, j := range e.Plantilla {
			pos[j.ID] = j.Posicion
		}
		for i, id := range al.Titulares {
			if id == 0 || pos[id] != f.Puestos()[i] {
				t.Errorf("%v: el puesto %d lo ocupa %d (%v), se esperaba un %v", f, i, id, pos[id], f.Puestos()[i])
			}
		}
		// El banquillo es el resto de la plantilla, de mayor a menor valoración.
		if len(al.Banquillo) != len(e.Plantilla)-modelo.TitularesPorEquipo {
			t.Errorf("%v: banquillo de %d", f, len(al.Banquillo))
		}
		porID := map[int]modelo.Jugador{}
		for _, j := range e.Plantilla {
			porID[j.ID] = j
		}
		for i := 1; i < len(al.Banquillo); i++ {
			if porID[al.Banquillo[i]].Valoracion() > porID[al.Banquillo[i-1]].Valoracion() {
				t.Errorf("%v: el banquillo no esta ordenado por valoracion", f)
				break
			}
		}
	}
}

func TestLaAlineacionAutomaticaEligeLosMejoresDeCadaLinea(t *testing.T) {
	e := equipoRealista("A", 1, 62)
	// Un defensa claramente mejor que el resto, y un delantero claramente peor.
	var mejorDef, peorDel int
	for i, j := range e.Plantilla {
		if j.Posicion == modelo.Defensa && mejorDef == 0 {
			e.Plantilla[i].Atributos = atributosUniformesSim(90)
			mejorDef = j.ID
		}
	}
	for i := len(e.Plantilla) - 1; i >= 0; i-- {
		if e.Plantilla[i].Posicion == modelo.Delantero {
			e.Plantilla[i].Atributos = atributosUniformesSim(20)
			peorDel = e.Plantilla[i].ID
			break
		}
	}
	al := AlineacionAutomatica(e, Criterios{Formacion: formacionPtr(modelo.F433)})
	titulares := map[int]bool{}
	for _, id := range al.Titulares {
		titulares[id] = true
	}
	if !titulares[mejorDef] {
		t.Error("el mejor defensa deberia ser titular")
	}
	if titulares[peorDel] {
		t.Error("el peor delantero no deberia ser titular habiendo otros")
	}
}

func TestLaFormacionAutomaticaDependeDeLaPlantilla(t *testing.T) {
	// Una plantilla fuerte en defensa y floja en ataque prefiere más defensas; una
	// fuerte en ataque y floja atrás, más delanteros.
	nivelar := func(defensa, ataque int) modelo.Equipo {
		e := equipoRealista("E", 1, 62)
		for i, j := range e.Plantilla {
			switch j.Posicion {
			case modelo.Defensa:
				e.Plantilla[i].Atributos = atributosUniformesSim(defensa)
			case modelo.Delantero, modelo.Mediocampista:
				e.Plantilla[i].Atributos = atributosUniformesSim(ataque)
			}
		}
		return e
	}
	solida := AlineacionAutomatica(nivelar(90, 55), Criterios{})
	ofensiva := AlineacionAutomatica(nivelar(55, 90), Criterios{})
	dS, _, delS := solida.Formacion.Lineas()
	dO, _, delO := ofensiva.Formacion.Lineas()
	t.Logf("plantilla defensiva -> %v; plantilla ofensiva -> %v", solida.Formacion, ofensiva.Formacion)
	if dS < dO || delO < delS || (dS == dO && delO == delS) {
		t.Errorf("la plantilla defensiva (%v) deberia usar mas defensas que la ofensiva (%v), y menos delanteros", solida.Formacion, ofensiva.Formacion)
	}
}

func TestLaAlineacionAutomaticaNoUsaANoDisponibles(t *testing.T) {
	e := equipoRealista("A", 1, 62)
	no := map[int]bool{}
	for _, f := range modelo.Formaciones {
		al := AlineacionAutomatica(e, Criterios{Formacion: formacionPtr(f)})
		// Declarar no disponibles a los tres primeros titulares y a un suplente.
		no = map[int]bool{al.Titulares[0]: true, al.Titulares[1]: true, al.Titulares[5]: true, al.Banquillo[0]: true}
		nueva := AlineacionAutomatica(e, Criterios{Formacion: formacionPtr(f), NoDisponibles: no})
		for _, id := range nueva.Titulares {
			if no[id] {
				t.Fatalf("%v: %d no esta disponible y esta alineado", f, id)
			}
		}
		for _, id := range nueva.Banquillo {
			if no[id] {
				t.Fatalf("%v: %d no esta disponible y esta en el banquillo", f, id)
			}
		}
		if err := nueva.Validar(e, no); err != nil {
			t.Errorf("%v: %v", f, err)
		}
	}
}

func TestSinPorterosDisponiblesJuegaUnJugadorDeCampo(t *testing.T) {
	e := equipoRealista("A", 1, 62)
	no := map[int]bool{}
	for _, j := range e.Plantilla {
		if j.Posicion == modelo.Portero {
			no[j.ID] = true
		}
	}
	al := AlineacionAutomatica(e, Criterios{NoDisponibles: no})
	var enLaPorteria modelo.Jugador
	for _, j := range e.Plantilla {
		if j.ID == al.Titulares[0] {
			enLaPorteria = j
		}
	}
	if enLaPorteria.Posicion == modelo.Portero || enLaPorteria.ID == 0 {
		t.Fatalf("en la porteria deberia jugar un jugador de campo: %+v", enLaPorteria)
	}
	if err := al.Validar(e, no); err != nil {
		t.Errorf("la alineacion de emergencia deberia ser valida: %v", err)
	}
	// Y nunca un portero en el campo.
	for i, id := range al.Titulares[1:] {
		for _, j := range e.Plantilla {
			if j.ID == id && j.Posicion == modelo.Portero {
				t.Errorf("un portero juega de %v", al.Puesto(i+1))
			}
		}
	}
}

func TestSiFaltanJugadoresEnUnaLineaSeCubreConLasContiguas(t *testing.T) {
	e := equipoRealista("A", 1, 62)
	no := map[int]bool{}
	defensas := 0
	for _, j := range e.Plantilla {
		if j.Posicion == modelo.Defensa {
			defensas++
			if defensas > 2 { // quedan solo dos defensas disponibles
				no[j.ID] = true
			}
		}
	}
	al := AlineacionAutomatica(e, Criterios{Formacion: formacionPtr(modelo.F433), NoDisponibles: no})
	if err := al.Validar(e, no); err != nil {
		t.Fatal(err)
	}
	pos := map[int]modelo.Posicion{}
	for _, j := range e.Plantilla {
		pos[j.ID] = j.Posicion
	}
	cuenta := map[modelo.Posicion]int{}
	for _, id := range al.Titulares {
		if id == 0 {
			t.Fatal("quedo un puesto sin cubrir habiendo jugadores")
		}
		cuenta[pos[id]]++
	}
	if cuenta[modelo.Defensa] != 2 || cuenta[modelo.Portero] != 1 {
		t.Errorf("deberia haber 2 defensas naturales y el resto de otras lineas: %v", cuenta)
	}
}

func TestUnaPlantillaCortaNoRompeLaSeleccion(t *testing.T) {
	e := equipoRealista("A", 1, 62)
	e.Plantilla = e.Plantilla[:7]
	// Solo un portero puede jugar (en la porteria); los demas no ocupan puestos de campo.
	utiles, hayPortero := 0, false
	for _, j := range e.Plantilla {
		if j.Posicion != modelo.Portero {
			utiles++
		} else if !hayPortero {
			hayPortero = true
			utiles++
		}
	}
	for _, f := range modelo.Formaciones {
		al := AlineacionAutomatica(e, Criterios{Formacion: formacionPtr(f)})
		huecos := 0
		for _, id := range al.Titulares {
			if id == 0 {
				huecos++
			}
		}
		if huecos != modelo.TitularesPorEquipo-utiles {
			t.Errorf("%v: %d puestos vacios con %d jugadores utiles", f, huecos, utiles)
		}
	}
	vacio := AlineacionAutomatica(modelo.Equipo{Nombre: "V"}, Criterios{})
	if !vacio.Vacia() {
		t.Error("un equipo sin jugadores no tiene alineacion")
	}
}

func TestJugarFueraDePosicionCuesta(t *testing.T) {
	e := equipoUniforme("U", 70)
	natural := alineacion433(e)
	aN, dN := Fuerzas(e, natural)

	// Un delantero de suplente ocupa el lugar de un defensa titular.
	fuera := natural
	var delantero int
	for _, j := range e.Plantilla {
		if j.Posicion == modelo.Delantero && j.ID != natural.Titulares[8] && j.ID != natural.Titulares[9] && j.ID != natural.Titulares[10] {
			delantero = j.ID
			break
		}
	}
	fuera.Titulares[1] = delantero
	aF, dF := Fuerzas(e, fuera)
	if dF >= dN {
		t.Errorf("un delantero de defensa deberia bajar la defensa: %.2f frente a %.2f", dF, dN)
	}
	if aF != aN {
		t.Errorf("el ataque no deberia cambiar: %.2f frente a %.2f", aF, aN)
	}
	// El efecto es moderado: un jugador de once (la familiaridad de saltar dos lineas).
	if caida := (dN - dF) / dN; caida < 0.003 || caida > 0.03 {
		t.Errorf("la caida de la defensa es del %.2f%%", 100*caida)
	}
	// Un defensa de portero cuesta muchísimo más que un medio de defensa.
	porteria := natural
	for _, j := range e.Plantilla {
		if j.Posicion == modelo.Defensa && j.ID != natural.Titulares[1] && j.ID != natural.Titulares[2] &&
			j.ID != natural.Titulares[3] && j.ID != natural.Titulares[4] {
			porteria.Titulares[0] = j.ID
			break
		}
	}
	_, dP := Fuerzas(e, porteria)
	if (dN-dP)/dN < 3*(dN-dF)/dN {
		t.Errorf("un jugador de campo de portero (%.2f) deberia costar mucho mas que uno fuera de linea (%.2f)", dP, dF)
	}
}

func TestUnJugadorFueraDePosicionNoSuperaAUnEspecialistaIgual(t *testing.T) {
	// A igualdad de atributos, jugar en su puesto natural nunca es peor que jugar
	// fuera de él: la familiaridad solo resta.
	for _, natural := range modelo.Posiciones {
		for _, puesto := range modelo.Posiciones {
			j := modelo.Jugador{ID: 1, Nombre: "J", Edad: 25, Posicion: natural, Atributos: atributosUniformesSim(75)}
			especialista := j
			especialista.Posicion = puesto
			if j.ValoracionEn(puesto) > especialista.ValoracionEn(puesto) {
				t.Errorf("un %v en el puesto de %v (%d) supera a un especialista igual (%d)", natural, puesto, j.ValoracionEn(puesto), especialista.ValoracionEn(puesto))
			}
		}
	}
}

func TestLasAlineacionesElegidasSeUsanEnElPartido(t *testing.T) {
	local, visitante := equipoRealista("L", 1, 62), equipoRealista("V", 101, 62)
	alL := AlineacionAutomatica(local, Criterios{Formacion: formacionPtr(modelo.F352)})
	alV := AlineacionAutomatica(visitante, Criterios{Formacion: formacionPtr(modelo.F532)})
	r := rand.New(rand.NewSource(3))
	for i := 0; i < 100; i++ {
		res := SimularConAlineaciones(r, local, alL, visitante, alV)
		d := res.Detalle
		if d.FormacionLocal != modelo.F352 || d.FormacionVisitante != modelo.F532 {
			t.Fatalf("formaciones %v y %v, se esperaban 3-5-2 y 5-3-2", d.FormacionLocal, d.FormacionVisitante)
		}
		if d.TitularesLocal != alL.Titulares || d.TitularesVisitante != alV.Titulares {
			t.Fatal("los titulares del detalle no son los de las alineaciones")
		}
		if _, err := d.Participaciones(); err != nil {
			t.Fatalf("detalle incoherente: %v", err)
		}
	}
}

func TestElOrdenDelBanquilloSeRespeta(t *testing.T) {
	e := equipoRealista("L", 1, 62)
	rival := equipoRealista("V", 101, 62)
	al := alineacion433(e)
	// Poner primero en el banquillo a un delantero concreto.
	var preferido int
	for _, id := range al.Banquillo {
		for _, j := range e.Plantilla {
			if j.ID == id && j.Posicion == modelo.Delantero {
				preferido = id
			}
		}
		if preferido != 0 {
			break
		}
	}
	if preferido == 0 {
		t.Fatal("no hay delantero en el banquillo")
	}
	resto := []int{preferido}
	for _, id := range al.Banquillo {
		if id != preferido {
			resto = append(resto, id)
		}
	}
	al.Banquillo = resto
	pos := map[int]modelo.Posicion{}
	for _, j := range e.Plantilla {
		pos[j.ID] = j.Posicion
	}
	puestos := al.Formacion.Puestos()
	puestoDe := map[int]modelo.Posicion{}
	for i, id := range al.Titulares {
		puestoDe[id] = puestos[i]
	}

	r := rand.New(rand.NewSource(4))
	alR := alineacion433(rival)
	reemplazosDelanteros, entroElPreferido := 0, 0
	for i := 0; i < 600; i++ {
		res := SimularConAlineaciones(r, e, al, rival, alR)
		for _, ev := range res.Detalle.Eventos {
			if ev.Tipo != modelo.Sustitucion || !ev.Local {
				continue
			}
			// El primer cambio del equipo: si sale un delantero, entra el preferido.
			if puestoDe[ev.Jugador] == modelo.Delantero {
				reemplazosDelanteros++
				if ev.Otro == preferido {
					entroElPreferido++
				}
			}
			break
		}
	}
	if reemplazosDelanteros < 50 {
		t.Fatalf("solo %d primeros cambios de delantero", reemplazosDelanteros)
	}
	if entroElPreferido != reemplazosDelanteros {
		t.Errorf("de %d cambios de delantero, el preferido entro en %d", reemplazosDelanteros, entroElPreferido)
	}
}

// puntosEsperados calcula, con la distribución de Poisson, los puntos por partido
// de un equipo cuyas esperanzas de goles son a favor y en contra.
func puntosEsperados(aFavor, enContra float64) float64 {
	pmf := func(k int, l float64) float64 {
		return math.Exp(-l) * math.Pow(l, float64(k)) / math.Gamma(float64(k+1))
	}
	var gana, empata float64
	for i := 0; i <= 14; i++ {
		for j := 0; j <= 14; j++ {
			p := pmf(i, aFavor) * pmf(j, enContra)
			switch {
			case i > j:
				gana += p
			case i == j:
				empata += p
			}
		}
	}
	return 3*gana + empata
}

func TestLasFormacionesSonEquilibradas(t *testing.T) {
	// Con equipos de composición variada, ninguna formación es mejor que las demás
	// en promedio: lo que decide es la plantilla. Se calcula con los puntos
	// esperados exactos contra un rival igual de fuerte.
	suma := map[modelo.Formacion]float64{}
	mejor := map[modelo.Formacion]int{}
	const equipos = 120
	r := rand.New(rand.NewSource(5))
	for n := 0; n < equipos; n++ {
		e := equipoRealista("E", 1, 55+r.Intn(15))
		for i, j := range e.Plantilla { // calidades dispares por jugador
			v := 45 + r.Intn(45)
			e.Plantilla[i].Atributos = modelo.Atributos{Ritmo: v, Tiro: v, Pase: v, Regate: v, Defensa: v, Fisico: v, Reflejos: v}
			_ = j
		}
		f433 := modelo.F433
		a0, d0 := Fuerzas(e, AlineacionAutomatica(e, Criterios{Formacion: &f433}))
		rival := (a0 + d0) / 2
		mejorF, mejorP := modelo.F433, -1.0
		for _, f := range modelo.Formaciones {
			f := f
			a, d := Fuerzas(e, AlineacionAutomatica(e, Criterios{Formacion: &f}))
			p := puntosEsperados(golesBase*math.Pow(a/rival, exponenteFuerza), golesBase*math.Pow(rival/d, exponenteFuerza))
			suma[f] += p / equipos
			if p > mejorP {
				mejorF, mejorP = f, p
			}
		}
		mejor[mejorF]++
	}
	var medias []float64
	for _, f := range modelo.Formaciones {
		t.Logf("%v: %.3f puntos esperados, la mejor en %.0f%% de los equipos", f, suma[f], 100*float64(mejor[f])/equipos)
		medias = append(medias, suma[f])
	}
	sort.Float64s(medias)
	if dif := medias[len(medias)-1] - medias[0]; dif > 0.04 {
		t.Errorf("las formaciones difieren en %.3f puntos por partido en promedio; ninguna deberia dominar", dif)
	}
	// Y todas son útiles para alguna plantilla, sin que una sola lo sea para casi todas.
	for _, f := range modelo.Formaciones {
		if p := 100 * float64(mejor[f]) / equipos; p > 50 {
			t.Errorf("%v es la mejor formacion para el %.0f%% de los equipos", f, p)
		}
	}
}

func TestLasFormacionesTienenPerfilesDistintos(t *testing.T) {
	e := equipoRealista("E", 1, 62)
	ataque := map[modelo.Formacion]float64{}
	defensa := map[modelo.Formacion]float64{}
	for _, f := range modelo.Formaciones {
		ataque[f], defensa[f] = Fuerzas(e, AlineacionAutomatica(e, Criterios{Formacion: formacionPtr(f)}))
	}
	if !(ataque[modelo.F343] > ataque[modelo.F433] && ataque[modelo.F433] > ataque[modelo.F532]) {
		t.Errorf("el ataque deberia ir de 3-4-3 a 5-3-2: %v", ataque)
	}
	if !(defensa[modelo.F532] > defensa[modelo.F433] && defensa[modelo.F433] > defensa[modelo.F343]) {
		t.Errorf("la defensa deberia ir de 5-3-2 a 3-4-3: %v", defensa)
	}
}

func TestDiezHombresMarcanMenosYRecibenMas(t *testing.T) {
	local, visitante := equipoRealista("L", 1, 62), equipoRealista("V", 101, 62)
	alL, alV := alineacion433(local), alineacion433(visitante)
	r := rand.New(rand.NewSource(6))

	var base, conRojaLocal [2]float64 // goles del local y del visitante
	var nBase, nRoja int
	for i := 0; i < 24000; i++ {
		res := SimularConAlineaciones(r, local, alL, visitante, alV)
		rojaLocalTemprana, otraRoja := false, false
		for _, e := range res.Detalle.Eventos {
			if e.Tipo != modelo.Roja {
				continue
			}
			if e.Local && e.Minuto <= 30 {
				rojaLocalTemprana = true
			} else {
				otraRoja = true
			}
		}
		switch {
		case rojaLocalTemprana && !otraRoja:
			conRojaLocal[0] += float64(res.GolesLocal)
			conRojaLocal[1] += float64(res.GolesVisitante)
			nRoja++
		case !rojaLocalTemprana && !otraRoja:
			base[0] += float64(res.GolesLocal)
			base[1] += float64(res.GolesVisitante)
			nBase++
		}
	}
	if nRoja < 200 {
		t.Fatalf("solo %d partidos con roja temprana del local", nRoja)
	}
	marcaBase, recibeBase := base[0]/float64(nBase), base[1]/float64(nBase)
	marcaRoja, recibeRoja := conRojaLocal[0]/float64(nRoja), conRojaLocal[1]/float64(nRoja)
	t.Logf("sin rojas: marca %.2f recibe %.2f; con roja temprana del local (%d partidos): marca %.2f recibe %.2f",
		marcaBase, recibeBase, nRoja, marcaRoja, recibeRoja)
	if marcaRoja > marcaBase*0.93 {
		t.Errorf("con diez, el local deberia marcar menos: %.2f frente a %.2f", marcaRoja, marcaBase)
	}
	if recibeRoja < recibeBase*1.07 {
		t.Errorf("con diez, el local deberia recibir mas: %.2f frente a %.2f", recibeRoja, recibeBase)
	}
	// El efecto es moderado: no se hunde ni se vuelve imposible marcar.
	if marcaRoja < marcaBase*0.6 || recibeRoja > recibeBase*1.5 {
		t.Errorf("el efecto de la roja es excesivo: marca %.2f/%.2f, recibe %.2f/%.2f", marcaRoja, marcaBase, recibeRoja, recibeBase)
	}
}
