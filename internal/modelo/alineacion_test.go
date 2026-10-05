package modelo

import (
	"strings"
	"testing"
)

func TestFormaciones(t *testing.T) {
	if len(Formaciones) != 6 {
		t.Fatalf("%d formaciones, se esperaban 6", len(Formaciones))
	}
	if F433 != 0 {
		t.Error("el valor cero deberia ser el 4-3-3 (los partidos anteriores)")
	}
	nombres := map[Formacion]string{F433: "4-3-3", F442: "4-4-2", F352: "3-5-2", F451: "4-5-1", F532: "5-3-2", F343: "3-4-3"}
	for _, f := range Formaciones {
		if !f.Valida() || f.String() != nombres[f] {
			t.Errorf("%d: valida %v, nombre %q (se esperaba %q)", int(f), f.Valida(), f.String(), nombres[f])
		}
		d, m, del := f.Lineas()
		if d+m+del != TitularesPorEquipo-1 {
			t.Errorf("%v: %d jugadores de campo, se esperaban 10", f, d+m+del)
		}
		if d < 3 || d > 5 || m < 3 || m > 5 || del < 1 || del > 3 {
			t.Errorf("%v: lineas poco realistas %d-%d-%d", f, d, m, del)
		}
		// Los puestos van portero, defensas, medios, delanteros.
		puestos := f.Puestos()
		if puestos[0] != Portero {
			t.Errorf("%v: el primer puesto deberia ser el portero", f)
		}
		cuenta := map[Posicion]int{}
		for i, p := range puestos {
			cuenta[p]++
			if i > 0 && p < puestos[i-1] {
				t.Errorf("%v: los puestos deberian ir de portero a delantero: %v", f, puestos)
			}
		}
		if cuenta[Portero] != 1 || cuenta[Defensa] != d || cuenta[Mediocampista] != m || cuenta[Delantero] != del {
			t.Errorf("%v: puestos %v, se esperaban 1-%d-%d-%d", f, cuenta, d, m, del)
		}
	}
	if Formacion(9).Valida() || Formacion(-1).Valida() || !strings.Contains(Formacion(9).String(), "9") {
		t.Error("una formacion desconocida no deberia ser valida")
	}
}

func TestFamiliaridadYValoracionEnOtroPuesto(t *testing.T) {
	// En su puesto natural, nada cambia.
	for _, p := range Posiciones {
		if Familiaridad(p, p) != 1 {
			t.Errorf("%v en su puesto deberia tener familiaridad 1", p)
		}
	}
	// Línea contigua mejor que saltar una línea; el portero, muy poco.
	if !(Familiaridad(Defensa, Mediocampista) > Familiaridad(Defensa, Delantero)) {
		t.Error("de defensa a medio deberia costar menos que de defensa a delantero")
	}
	if Familiaridad(Mediocampista, Defensa) != Familiaridad(Defensa, Mediocampista) ||
		Familiaridad(Mediocampista, Delantero) != Familiaridad(Delantero, Mediocampista) ||
		Familiaridad(Defensa, Delantero) != Familiaridad(Delantero, Defensa) {
		t.Error("la familiaridad deberia ser simetrica")
	}
	for _, p := range []Posicion{Defensa, Mediocampista, Delantero} {
		if Familiaridad(Portero, p) >= Familiaridad(Defensa, Delantero) || Familiaridad(p, Portero) >= Familiaridad(Defensa, Delantero) {
			t.Errorf("el portero fuera de su sitio (%v) deberia costar mucho mas", p)
		}
	}

	// Un mediocampista de atributos parejos rinde menos de defensa que de medio.
	j := Jugador{ID: 1, Nombre: "J", Edad: 25, Posicion: Mediocampista, Atributos: atributosUniformes(70)}
	if j.ValoracionEn(Mediocampista) != j.Valoracion() || j.Valoracion() != 70 {
		t.Errorf("en su puesto: %d / %d", j.ValoracionEn(Mediocampista), j.Valoracion())
	}
	if got := j.ValoracionEn(Defensa); got != 67 { // 70 * 0,95 = 66,5 -> 67
		t.Errorf("de defensa: %d, se esperaba 67", got)
	}
	if got := j.ValoracionEn(Portero); got != 35 {
		t.Errorf("de portero: %d, se esperaba 35", got)
	}
	// Los atributos mandan: un medio con muy buen tiro y ritmo rinde mejor de
	// delantero que uno sin ellos.
	rematador := Jugador{ID: 2, Nombre: "R", Edad: 25, Posicion: Mediocampista,
		Atributos: Atributos{Ritmo: 90, Tiro: 90, Pase: 60, Regate: 80, Defensa: 40, Fisico: 60, Reflejos: 1}}
	pasador := Jugador{ID: 3, Nombre: "P", Edad: 25, Posicion: Mediocampista,
		Atributos: Atributos{Ritmo: 50, Tiro: 50, Pase: 90, Regate: 70, Defensa: 60, Fisico: 60, Reflejos: 1}}
	if rematador.ValoracionEn(Delantero) <= pasador.ValoracionEn(Delantero) {
		t.Error("un medio rematador deberia rendir mejor de delantero que uno pasador")
	}
}

// equipoParaAlineacion tiene 2 porteros (IDs 1-2), 7 defensas (3-9), 7 medios
// (10-16) y 6 delanteros (17-22).
func equipoParaAlineacion() Equipo {
	e := Equipo{Nombre: "E"}
	id := 1
	for _, c := range []struct {
		p Posicion
		n int
	}{{Portero, 2}, {Defensa, 7}, {Mediocampista, 7}, {Delantero, 6}} {
		for i := 0; i < c.n; i++ {
			e.Plantilla = append(e.Plantilla, Jugador{ID: id, Nombre: "J", Edad: 25, Posicion: c.p, Atributos: atributosUniformes(60)})
			id++
		}
	}
	return e
}

// alineacion433 es un 4-3-3 válido: portero 1, defensas 3-6, medios 10-12 y
// delanteros 17-19, con el resto en el banquillo.
func alineacion433() Alineacion {
	return Alineacion{
		Formacion: F433,
		Titulares: [TitularesPorEquipo]int{1, 3, 4, 5, 6, 10, 11, 12, 17, 18, 19},
		Banquillo: []int{2, 7, 13, 20},
	}
}

func TestAlineacionValida(t *testing.T) {
	e := equipoParaAlineacion()
	if err := alineacion433().Validar(e, nil); err != nil {
		t.Fatalf("una alineacion correcta se rechazo: %v", err)
	}
	a := alineacion433()
	if a.Vacia() || (Alineacion{}).Vacia() == false {
		t.Error("Vacia deberia ser verdadero solo sin titulares")
	}
	if a.Puesto(0) != Portero || a.Puesto(1) != Defensa || a.Puesto(5) != Mediocampista || a.Puesto(10) != Delantero {
		t.Error("Puesto no sigue el orden de la formacion")
	}
	// Una formacion distinta cambia los puestos.
	a.Formacion = F352
	if a.Puesto(4) != Mediocampista || a.Puesto(9) != Delantero {
		t.Error("en un 3-5-2 el puesto 5 es un medio")
	}
}

func TestAlineacionRechazaLoIncoherente(t *testing.T) {
	e := equipoParaAlineacion()
	casos := map[string]struct {
		mutar     func(*Alineacion)
		noDisp    map[int]bool
		fragmento string
	}{
		"formacion invalida":   {func(a *Alineacion) { a.Formacion = Formacion(9) }, nil, "formacion invalida"},
		"titular ajeno":        {func(a *Alineacion) { a.Titulares[3] = 999 }, nil, "no es de la plantilla"},
		"titular repetido":     {func(a *Alineacion) { a.Titulares[4] = a.Titulares[3] }, nil, "dos veces"},
		"titular sin puesto":   {func(a *Alineacion) { a.Titulares[10] = 0 }, nil, "no es de la plantilla"},
		"no disponible":        {func(a *Alineacion) {}, map[int]bool{4: true}, "no esta disponible"},
		"portero que no lo es": {func(a *Alineacion) { a.Titulares[0] = 10 }, nil, "no es portero"},
		"portero en el campo":  {func(a *Alineacion) { a.Titulares[1] = 2; a.Banquillo = []int{7} }, nil, "no puede jugar de"},
		"suplente ajeno":       {func(a *Alineacion) { a.Banquillo = []int{999} }, nil, "no es de la plantilla"},
		"suplente titular":     {func(a *Alineacion) { a.Banquillo = []int{1} }, nil, "titular y esta en el banquillo"},
		"suplente repetido":    {func(a *Alineacion) { a.Banquillo = []int{2, 2} }, nil, "dos veces en el banquillo"},
	}
	for nombre, c := range casos {
		a := alineacion433()
		c.mutar(&a)
		err := a.Validar(e, c.noDisp)
		if err == nil {
			t.Errorf("%s: deberia dar error", nombre)
		} else if !strings.Contains(err.Error(), c.fragmento) {
			t.Errorf("%s: error %q, se esperaba que contenga %q", nombre, err, c.fragmento)
		}
	}
}

func TestUnJugadorDeCampoSoloEsPorteroEnUnaEmergencia(t *testing.T) {
	e := equipoParaAlineacion()
	a := alineacion433()
	a.Titulares[0] = 13 // un medio (que no esta en el once) de portero
	a.Banquillo = []int{2, 7, 20}
	if a.Validar(e, nil) == nil {
		t.Fatal("con porteros disponibles, un medio no puede ser portero")
	}
	// Con los dos porteros no disponibles, sí.
	a.Banquillo = []int{7, 20}
	if err := a.Validar(e, map[int]bool{1: true, 2: true}); err != nil {
		t.Errorf("en una emergencia un jugador de campo puede ir a la porteria: %v", err)
	}
	// Pero si queda un portero disponible, no.
	a.Banquillo = []int{2, 7, 20}
	if a.Validar(e, map[int]bool{1: true}) == nil {
		t.Error("con un portero todavia disponible no hay emergencia")
	}
}

func TestLosJugadoresPuedenJugarFueraDeSuLinea(t *testing.T) {
	e := equipoParaAlineacion()
	a := alineacion433()
	a.Titulares[1] = 13 // un medio (suplente) de defensa
	a.Titulares[9] = 7  // un defensa (suplente) de delantero
	a.Banquillo = []int{2, 14, 20}
	if err := a.Validar(e, nil); err != nil {
		t.Errorf("jugar fuera de su linea esta permitido: %v", err)
	}
}

func TestParticipacionesLlevanElPuestoEnQueJugaron(t *testing.T) {
	d := detalleBase()
	d.FormacionLocal, d.FormacionVisitante = F352, F433
	d.Eventos = []Evento{
		{Minuto: 60, Tipo: Sustitucion, Local: true, Jugador: 5, Otro: 12}, // el 5 (puesto 5) es un medio en 3-5-2
	}
	ps, err := d.Participaciones()
	if err != nil {
		t.Fatal(err)
	}
	puestos := map[int]Posicion{}
	for _, p := range ps {
		puestos[p.Jugador] = p.Puesto
	}
	// 3-5-2: portero, 3 defensas (2-4), 5 medios (5-9), 2 delanteros (10-11).
	// El local juega 3-5-2 y el visitante 4-3-3 (portero, 4 defensas, 3 medios, 3
	// delanteros: los puestos 5 a 7 son medios y 8 a 10, delanteros).
	for id, want := range map[int]Posicion{1: Portero, 2: Defensa, 4: Defensa, 5: Mediocampista, 9: Mediocampista, 10: Delantero, 11: Delantero,
		101: Portero, 102: Defensa, 105: Defensa, 106: Mediocampista, 108: Mediocampista, 109: Delantero, 111: Delantero} {
		if puestos[id] != want {
			t.Errorf("jugador %d: puesto %v, se esperaba %v", id, puestos[id], want)
		}
	}
	if puestos[12] != Mediocampista {
		t.Errorf("el suplente hereda el puesto del que sustituye: %v", puestos[12])
	}
	d.FormacionLocal = Formacion(9)
	if _, err := d.Participaciones(); err == nil {
		t.Error("una formacion invalida deberia dar error")
	}
}
