package modelo

import (
	"strings"
	"testing"
)

func atributosUniformes(v int) Atributos {
	return Atributos{Ritmo: v, Tiro: v, Pase: v, Regate: v, Defensa: v, Fisico: v, Reflejos: v}
}

func TestPesosSumanCien(t *testing.T) {
	for _, p := range Posiciones {
		w := pesos[p]
		suma := w.Ritmo + w.Tiro + w.Pase + w.Regate + w.Defensa + w.Fisico + w.Reflejos
		if suma != 100 {
			t.Errorf("pesos de %v suman %d, se esperaba 100", p, suma)
		}
	}
}

func TestMediaUniforme(t *testing.T) {
	for _, p := range Posiciones {
		if got := atributosUniformes(70).Media(p); got != 70 {
			t.Errorf("Media(%v) con atributos 70 = %d, se esperaba 70", p, got)
		}
	}
	if got := atributosUniformes(70).Media(Posicion(99)); got != 0 {
		t.Errorf("Media con posicion invalida = %d, se esperaba 0", got)
	}
}

func TestMediaPesaPorPosicion(t *testing.T) {
	a := atributosUniformes(50)
	a.Reflejos = 99
	if a.Media(Portero) <= a.Media(Delantero) {
		t.Error("los reflejos deberian valer mas para un portero que para un delantero")
	}
	b := atributosUniformes(50)
	b.Tiro = 99
	if b.Media(Delantero) <= b.Media(Defensa) {
		t.Error("el tiro deberia valer mas para un delantero que para un defensa")
	}
}

func TestAtributosValidar(t *testing.T) {
	if err := atributosUniformes(50).Validar(); err != nil {
		t.Fatalf("atributos validos rechazados: %v", err)
	}
	for _, v := range []int{AtributoMin - 1, AtributoMax + 1} {
		a := atributosUniformes(50)
		a.Pase = v
		if a.Validar() == nil {
			t.Errorf("pase=%d deberia ser invalido", v)
		}
	}
}

func TestPosicionString(t *testing.T) {
	if Portero.String() != "Portero" || Delantero.String() != "Delantero" {
		t.Error("String de posiciones inesperado")
	}
	if Posicion(42).Valida() {
		t.Error("Posicion(42) no deberia ser valida")
	}
}

func jugadorValido(id int, p Posicion) Jugador {
	return Jugador{ID: id, Nombre: "Jugador", Edad: 25, Posicion: p, Atributos: atributosUniformes(60)}
}

func TestJugadorValidar(t *testing.T) {
	if err := jugadorValido(1, Defensa).Validar(); err != nil {
		t.Fatalf("jugador valido rechazado: %v", err)
	}
	casos := map[string]func(*Jugador){
		"sin nombre":        func(j *Jugador) { j.Nombre = "" },
		"edad baja":         func(j *Jugador) { j.Edad = EdadMin - 1 },
		"edad alta":         func(j *Jugador) { j.Edad = EdadMax + 1 },
		"posicion invalida": func(j *Jugador) { j.Posicion = Posicion(9) },
		"atributo invalido": func(j *Jugador) { j.Atributos.Tiro = 0 },
	}
	for nombre, mutar := range casos {
		j := jugadorValido(1, Defensa)
		mutar(&j)
		if j.Validar() == nil {
			t.Errorf("%s: deberia ser invalido", nombre)
		}
	}
}

func equipoBase() Equipo {
	e := Equipo{Nombre: "Equipo"}
	id := 1
	for _, c := range []struct {
		p Posicion
		n int
	}{{Portero, 2}, {Defensa, 6}, {Mediocampista, 6}, {Delantero, 4}} {
		for i := 0; i < c.n; i++ {
			e.Plantilla = append(e.Plantilla, jugadorValido(id, c.p))
			id++
		}
	}
	return e
}

func TestEquipoValidar(t *testing.T) {
	e := equipoBase()
	if err := e.Validar(); err != nil {
		t.Fatalf("equipo valido rechazado: %v", err)
	}
	if e.Contar(Portero) != 2 || e.Contar(Defensa) != 6 {
		t.Error("Contar devuelve cantidades inesperadas")
	}

	sinNombre := equipoBase()
	sinNombre.Nombre = ""
	if sinNombre.Validar() == nil {
		t.Error("equipo sin nombre deberia ser invalido")
	}

	corto := equipoBase()
	corto.Plantilla = corto.Plantilla[:MinPlantilla-1]
	if corto.Validar() == nil {
		t.Error("plantilla corta deberia ser invalida")
	}

	sinPorteros := equipoBase()
	sinPorteros.Plantilla[0].Posicion = Defensa
	if sinPorteros.Validar() == nil {
		t.Error("equipo con un solo portero deberia ser invalido")
	}

	repetido := equipoBase()
	repetido.Plantilla[3].ID = repetido.Plantilla[2].ID
	if repetido.Validar() == nil {
		t.Error("IDs repetidas deberian ser invalidas")
	}

	malo := equipoBase()
	malo.Plantilla[5].Edad = 0
	if malo.Validar() == nil {
		t.Error("jugador invalido deberia invalidar el equipo")
	}
}

func TestEquipoValoracion(t *testing.T) {
	if (Equipo{}).Valoracion() != 0 {
		t.Error("equipo vacio deberia valorar 0")
	}
	e := equipoBase()
	if got := e.Valoracion(); got != 60 {
		t.Errorf("Valoracion = %d, se esperaba 60", got)
	}
	// Un suplente muy bueno entra en el once y sube la valoracion.
	antes := e.Valoracion()
	crack := jugadorValido(100, Delantero)
	crack.Atributos = atributosUniformes(99)
	e.Plantilla = append(e.Plantilla, crack)
	if e.Valoracion() <= antes {
		t.Error("un jugador mejor deberia subir la valoracion")
	}
}

func TestTalentoEfectivoYValidacion(t *testing.T) {
	j := jugadorValido(1, Delantero)
	if j.Talento != 0 || j.TalentoEfectivo() != TalentoNeutro {
		t.Errorf("un jugador sin talento definido deberia tener el neutro: %d / %d", j.Talento, j.TalentoEfectivo())
	}
	j.Talento = 130
	if j.TalentoEfectivo() != 130 || j.Validar() != nil {
		t.Errorf("un talento de 130 deberia ser valido y respetarse: %v", j.Validar())
	}
	for _, t0 := range []int{TalentoMin, TalentoMax, TalentoNeutro} {
		j.Talento = t0
		if err := j.Validar(); err != nil {
			t.Errorf("el talento %d deberia ser valido: %v", t0, err)
		}
	}
	for _, t0 := range []int{1, TalentoMin - 1, TalentoMax + 1, -5} {
		j.Talento = t0
		if j.Validar() == nil {
			t.Errorf("el talento %d deberia ser invalido", t0)
		}
	}
}

func TestProyeccionPorTalento(t *testing.T) {
	casos := []struct {
		talento int
		want    Proyeccion
	}{
		{0, ProyeccionNormal}, // sin definir: neutro
		{TalentoMin, ProyeccionLimitada},
		{talentoProyeccionNormal - 1, ProyeccionLimitada},
		{talentoProyeccionNormal, ProyeccionNormal},
		{100, ProyeccionNormal},
		{talentoProyeccionAlta - 1, ProyeccionNormal},
		{talentoProyeccionAlta, ProyeccionAlta},
		{talentoProyeccionExcepcional - 1, ProyeccionAlta},
		{talentoProyeccionExcepcional, ProyeccionExcepcional},
		{TalentoMax, ProyeccionExcepcional},
	}
	for _, c := range casos {
		j := jugadorValido(1, Mediocampista)
		j.Talento = c.talento
		if got := j.Proyeccion(); got != c.want {
			t.Errorf("talento %d: proyeccion %v, se esperaba %v", c.talento, got, c.want)
		}
	}
	for p, want := range map[Proyeccion]string{ProyeccionLimitada: "limitada", ProyeccionNormal: "normal", ProyeccionAlta: "alta", ProyeccionExcepcional: "excepcional"} {
		if p.String() != want {
			t.Errorf("%d: %q, se esperaba %q", int(p), p.String(), want)
		}
	}
	if !strings.Contains(Proyeccion(9).String(), "9") {
		t.Error("una proyeccion desconocida deberia mostrar su numero")
	}
}
