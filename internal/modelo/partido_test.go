package modelo

import (
	"reflect"
	"strings"
	"testing"
)

// detalleBase tiene titulares 1..11 (local) y 101..111 (visitante), y a 12 y
// 112 en la banca.
func detalleBase() DetallePartido {
	var d DetallePartido
	for i := 0; i < TitularesPorEquipo; i++ {
		d.TitularesLocal[i] = 1 + i
		d.TitularesVisitante[i] = 101 + i
	}
	return d
}

func TestParticipacionesSinSucesos(t *testing.T) {
	ps, err := detalleBase().Participaciones()
	if err != nil {
		t.Fatal(err)
	}
	if len(ps) != 22 {
		t.Fatalf("%d participaciones, se esperaban 22", len(ps))
	}
	for i, p := range ps {
		if !p.Titular() || p.Minutos() != MinutosPartido || p.Local != (i < 11) {
			t.Errorf("participacion %d incorrecta: %+v", i, p)
		}
	}
}

func TestParticipacionesConSustitucionYExpulsion(t *testing.T) {
	d := detalleBase()
	d.Eventos = []Evento{
		{Minuto: 20, Tipo: Amarilla, Local: false, Jugador: 105},
		{Minuto: 30, Tipo: Gol, Local: true, Jugador: 10, Otro: 7},
		{Minuto: 60, Tipo: Sustitucion, Local: true, Jugador: 9, Otro: 12},
		{Minuto: 60, Tipo: Gol, Local: true, Jugador: 9}, // el que sale aun puede marcar en ese minuto
		{Minuto: 70, Tipo: Gol, Local: true, Jugador: 12, Otro: 10},
		{Minuto: 75, Tipo: Roja, Local: false, Jugador: 105},
	}
	ps, err := d.Participaciones()
	if err != nil {
		t.Fatal(err)
	}
	por := map[int]Participacion{}
	for _, p := range ps {
		por[p.Jugador] = p
	}
	if len(ps) != 23 {
		t.Errorf("%d participaciones, se esperaban 23", len(ps))
	}
	if p := por[9]; p.Hasta != 60 || p.Minutos() != 60 || !p.Titular() {
		t.Errorf("el que sale: %+v", p)
	}
	if p := por[12]; p.Desde != 60 || p.Hasta != 90 || p.Minutos() != 30 || p.Titular() {
		t.Errorf("el suplente: %+v", p)
	}
	if p := por[105]; p.Hasta != 75 || p.Minutos() != 75 {
		t.Errorf("el expulsado: %+v", p)
	}
	if l, v := d.Goles(); l != 3 || v != 0 {
		t.Errorf("Goles = %d-%d, se esperaba 3-0", l, v)
	}
}

func TestParticipacionesRechazaDetallesIncoherentes(t *testing.T) {
	casos := map[string]struct {
		eventos   []Evento
		fragmento string
	}{
		"tipo invalido":          {[]Evento{{Minuto: 5, Tipo: TipoEvento(9), Local: true, Jugador: 1}}, "tipo invalido"},
		"minuto cero":            {[]Evento{{Minuto: 0, Tipo: Gol, Local: true, Jugador: 1}}, "fuera de"},
		"minuto 91":              {[]Evento{{Minuto: 91, Tipo: Gol, Local: true, Jugador: 1}}, "fuera de"},
		"fuera de orden":         {[]Evento{{Minuto: 50, Tipo: Gol, Local: true, Jugador: 1}, {Minuto: 40, Tipo: Gol, Local: true, Jugador: 2}}, "fuera de orden"},
		"jugador ajeno":          {[]Evento{{Minuto: 5, Tipo: Gol, Local: true, Jugador: 999}}, "no participa"},
		"equipo equivocado":      {[]Evento{{Minuto: 5, Tipo: Gol, Local: false, Jugador: 1}}, "otro equipo"},
		"suplente sin entrar":    {[]Evento{{Minuto: 5, Tipo: Gol, Local: true, Jugador: 12}}, "no participa"},
		"expulsado que marca":    {[]Evento{{Minuto: 10, Tipo: Roja, Local: true, Jugador: 3}, {Minuto: 20, Tipo: Gol, Local: true, Jugador: 3}}, "no estaba en el campo"},
		"sustituido que marca":   {[]Evento{{Minuto: 50, Tipo: Sustitucion, Local: true, Jugador: 3, Otro: 12}, {Minuto: 60, Tipo: Gol, Local: true, Jugador: 3}}, "no estaba en el campo"},
		"entra y marca a la vez": {[]Evento{{Minuto: 50, Tipo: Sustitucion, Local: true, Jugador: 3, Otro: 12}, {Minuto: 50, Tipo: Gol, Local: true, Jugador: 12}}, "no estaba en el campo"},
		"entra dos veces":        {[]Evento{{Minuto: 50, Tipo: Sustitucion, Local: true, Jugador: 3, Otro: 12}, {Minuto: 60, Tipo: Sustitucion, Local: true, Jugador: 4, Otro: 12}}, "dos veces"},
		"entra un titular":       {[]Evento{{Minuto: 50, Tipo: Sustitucion, Local: true, Jugador: 3, Otro: 4}}, "dos veces"},
		"entra el 0":             {[]Evento{{Minuto: 50, Tipo: Sustitucion, Local: true, Jugador: 3, Otro: 0}}, "invalida"},
		"autoasistencia":         {[]Evento{{Minuto: 5, Tipo: Gol, Local: true, Jugador: 3, Otro: 3}}, "asiste a si mismo"},
		"asistente ajeno":        {[]Evento{{Minuto: 5, Tipo: Gol, Local: true, Jugador: 3, Otro: 101}}, "otro equipo"},
		"asistente ausente":      {[]Evento{{Minuto: 5, Tipo: Gol, Local: true, Jugador: 3, Otro: 12}}, "no participa"},
	}
	for nombre, c := range casos {
		d := detalleBase()
		d.Eventos = c.eventos
		_, err := d.Participaciones()
		if err == nil {
			t.Errorf("%s: deberia dar error", nombre)
		} else if !strings.Contains(err.Error(), c.fragmento) {
			t.Errorf("%s: error %q, se esperaba que contenga %q", nombre, err, c.fragmento)
		}
	}

	repetido := detalleBase()
	repetido.TitularesVisitante[0] = repetido.TitularesLocal[0]
	if _, err := repetido.Participaciones(); err == nil {
		t.Error("un titular repetido entre equipos deberia dar error")
	}
}

func TestParticipacionesConEquipoIncompleto(t *testing.T) {
	var d DetallePartido
	d.TitularesLocal[0], d.TitularesLocal[1] = 1, 2
	d.TitularesVisitante[0] = 101
	d.Eventos = []Evento{{Minuto: 10, Tipo: Gol, Local: true, Jugador: 2, Otro: 1}}
	ps, err := d.Participaciones()
	if err != nil || len(ps) != 3 {
		t.Fatalf("%d participaciones, %v", len(ps), err)
	}
}

func TestDetalleVacioYGoles(t *testing.T) {
	var d DetallePartido
	if !d.Vacio() {
		t.Error("un detalle sin datos deberia ser vacio")
	}
	if ps, err := d.Participaciones(); err != nil || len(ps) != 0 {
		t.Errorf("un detalle vacio no tiene participaciones: %v, %v", ps, err)
	}
	if detalleBase().Vacio() {
		t.Error("un detalle con titulares no es vacio")
	}
	d.Eventos = []Evento{{Minuto: 1, Tipo: Gol, Local: true, Jugador: 1}}
	if d.Vacio() {
		t.Error("un detalle con sucesos no es vacio")
	}

	d = detalleBase()
	d.Eventos = []Evento{
		{Minuto: 5, Tipo: Gol, Local: true, Jugador: 1},
		{Minuto: 6, Tipo: Gol, Local: false, Jugador: 101},
		{Minuto: 7, Tipo: Amarilla, Local: false, Jugador: 102},
		{Minuto: 8, Tipo: Gol, Local: false, Jugador: 101},
	}
	if l, v := d.Goles(); l != 1 || v != 2 {
		t.Errorf("Goles = %d-%d, se esperaba 1-2", l, v)
	}
}

func TestTipoEventoString(t *testing.T) {
	for tipo, want := range map[TipoEvento]string{Gol: "Gol", Amarilla: "Amarilla", Roja: "Roja", Sustitucion: "Sustitucion"} {
		if tipo.String() != want || !tipo.Valido() {
			t.Errorf("%d: %q, valido %v", int(tipo), tipo.String(), tipo.Valido())
		}
	}
	if TipoEvento(7).Valido() || TipoEvento(-1).Valido() || !strings.Contains(TipoEvento(7).String(), "7") {
		t.Error("un tipo desconocido no deberia ser valido")
	}
}

func TestEstadisticasSumarYValoracionMedia(t *testing.T) {
	a := Estadisticas{Partidos: 2, Titularidades: 1, Minutos: 120, Goles: 3, Asistencias: 1, Amarillas: 1,
		Rojas: 0, PorteriasImbatidas: 1, GolesEncajados: 2, SumaValoracion: 130}
	b := Estadisticas{Partidos: 1, Titularidades: 1, Minutos: 90, Goles: 1, Rojas: 1, SumaValoracion: 45}
	want := Estadisticas{Partidos: 3, Titularidades: 2, Minutos: 210, Goles: 4, Asistencias: 1, Amarillas: 1,
		Rojas: 1, PorteriasImbatidas: 1, GolesEncajados: 2, SumaValoracion: 175}
	if got := a.Sumar(b); !reflect.DeepEqual(got, want) {
		t.Errorf("Sumar = %+v, se esperaba %+v", got, want)
	}
	if a.Sumar(Estadisticas{}) != a {
		t.Error("sumar el cero no deberia cambiar nada")
	}
	if got := a.ValoracionMedia(); got != 6.5 {
		t.Errorf("ValoracionMedia = %.2f, se esperaba 6.5", got)
	}
	if (Estadisticas{}).ValoracionMedia() != 0 {
		t.Error("sin partidos la valoracion media es 0")
	}
}
