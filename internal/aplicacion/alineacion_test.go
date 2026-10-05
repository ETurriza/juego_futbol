package aplicacion

import (
	"reflect"
	"testing"

	"github.com/ETurriza/juego_futbol/internal/modelo"
)

func formacionDe(f modelo.Formacion) *modelo.Formacion { return &f }

func TestSinEleccionLaAlineacionEsLaAutomatica(t *testing.T) {
	c := nuevaCarrera(t, 3, 10)
	al, manual := c.AlineacionVigente()
	if manual {
		t.Error("al empezar no hay alineacion elegida")
	}
	if err := al.Validar(c.equipoUsuario(), nil); err != nil {
		t.Errorf("la alineacion automatica no es valida: %v", err)
	}
	if !reflect.DeepEqual(al, c.AlineacionAutomatica(nil)) {
		t.Error("la vigente deberia ser la automatica")
	}
	if c.AlineacionElegidaInvalida() != nil {
		t.Error("sin eleccion no hay nada invalido")
	}
}

func TestElegirAlineacionValidaLaGuardaYLaUsa(t *testing.T) {
	c := nuevaCarrera(t, 3, 10)
	elegida := c.AlineacionAutomatica(formacionDe(modelo.F343))
	if err := c.ElegirAlineacion(elegida); err != nil {
		t.Fatal(err)
	}
	al, manual := c.AlineacionVigente()
	if !manual || !reflect.DeepEqual(al, elegida) {
		t.Fatalf("la vigente deberia ser la elegida: %v (manual=%v)", al, manual)
	}
	// Modificar el argumento original no cambia la alineacion guardada.
	elegida.Banquillo[0] = -1
	if al2, _ := c.AlineacionVigente(); al2.Banquillo[0] == -1 {
		t.Error("la alineacion elegida comparte memoria con el argumento")
	}

	if _, err := c.AvanzarJornada(); err != nil {
		t.Fatal(err)
	}
	usada := false
	for _, res := range c.Temporada.Resultados[0] {
		switch c.Usuario {
		case res.Partido.Local:
			usada = res.Detalle.FormacionLocal == modelo.F343 && res.Detalle.TitularesLocal == al.Titulares
		case res.Partido.Visitante:
			usada = res.Detalle.FormacionVisitante == modelo.F343 && res.Detalle.TitularesVisitante == al.Titulares
		}
		if usada {
			break
		}
	}
	if !usada {
		t.Error("el equipo del usuario no jugo con la alineacion elegida")
	}

	c.UsarAlineacionAutomatica()
	if _, manual := c.AlineacionVigente(); manual {
		t.Error("tras volver a la automatica no deberia haber eleccion")
	}
}

func TestElegirAlineacionInvalidaNoCambiaNada(t *testing.T) {
	c := nuevaCarrera(t, 3, 10)
	buena := c.AlineacionAutomatica(nil)
	if err := c.ElegirAlineacion(buena); err != nil {
		t.Fatal(err)
	}
	casos := map[string]func(al *modelo.Alineacion){
		"un titular repetido": func(al *modelo.Alineacion) { al.Titulares[5] = al.Titulares[4] },
		"un jugador ajeno":    func(al *modelo.Alineacion) { al.Titulares[5] = 999999 },
		"formacion invalida":  func(al *modelo.Alineacion) { al.Formacion = modelo.Formacion(99) },
		"sin portero":         func(al *modelo.Alineacion) { al.Titulares[0] = al.Titulares[1]; al.Titulares[1] = 0 },
	}
	for nombre, rompe := range casos {
		mala := c.AlineacionAutomatica(nil)
		rompe(&mala)
		if err := c.ElegirAlineacion(mala); err == nil {
			t.Errorf("%s: deberia dar error", nombre)
		}
		if al, manual := c.AlineacionVigente(); !manual || !reflect.DeepEqual(al, buena) {
			t.Errorf("%s: el intento fallido cambio la alineacion elegida", nombre)
		}
	}
}

func TestSiElegidaDejaDeSerValidaSeUsaLaAutomatica(t *testing.T) {
	c := nuevaCarrera(t, 3, 10)
	if err := c.ElegirAlineacion(c.AlineacionAutomatica(nil)); err != nil {
		t.Fatal(err)
	}
	// Un titular se va del equipo (retiro, venta...).
	retirado := c.Alineacion.Titulares[4]
	e := &c.Temporada.Equipos[c.Usuario]
	var resto []modelo.Jugador
	for _, j := range e.Plantilla {
		if j.ID != retirado {
			resto = append(resto, j)
		}
	}
	e.Plantilla = resto

	if c.AlineacionElegidaInvalida() == nil {
		t.Error("la alineacion con un jugador que ya no esta deberia ser invalida")
	}
	al, manual := c.AlineacionVigente()
	if manual {
		t.Error("deberia usarse la automatica")
	}
	if err := al.Validar(c.equipoUsuario(), nil); err != nil {
		t.Errorf("la automatica de reemplazo no es valida: %v", err)
	}
	if _, err := c.AvanzarJornada(); err != nil {
		t.Errorf("una alineacion obsoleta no deberia impedir jugar: %v", err)
	}
}

func TestLaAlineacionElegidaSobreviveAGuardarYReanudar(t *testing.T) {
	c := nuevaCarrera(t, 9, 10)
	if err := c.ElegirAlineacion(c.AlineacionAutomatica(formacionDe(modelo.F532))); err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 4; i++ {
		c.AvanzarJornada()
	}
	g := c.Exportar()
	g.Alineacion.Banquillo[0] = -7 // no debe afectar a la carrera original
	if c.Alineacion.Banquillo[0] == -7 {
		t.Fatal("Exportar comparte la alineacion con la carrera")
	}

	reanudada, err := Importar(c.Exportar())
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(c, reanudada) {
		t.Error("la carrera reanudada no es identica")
	}
	for !c.Terminada() {
		a, _ := c.AvanzarJornada()
		b, err := reanudada.AvanzarJornada()
		if err != nil || !reflect.DeepEqual(a, b) {
			t.Fatalf("jornada %d difiere al reanudar con alineacion elegida", c.Jornada())
		}
	}
}

func TestImportarRechazaUnaFormacionInvalida(t *testing.T) {
	c := nuevaCarrera(t, 9, 10)
	g := c.Exportar()
	g.Alineacion.Formacion = modelo.Formacion(42)
	g.AlineacionManual = true
	if _, err := Importar(g); err == nil {
		t.Error("una formacion invalida en la alineacion deberia rechazarse")
	}
}

func TestAlineacionElegidaNoCambiaLaAutomaticaDeLosRivales(t *testing.T) {
	// Elegir una alineacion afecta solo al equipo del usuario: los demas partidos de
	// la jornada son iguales con o sin ella cuando es la misma que habria elegido solo.
	a := nuevaCarrera(t, 12, 10)
	b := nuevaCarrera(t, 12, 10)
	if err := b.ElegirAlineacion(b.AlineacionAutomatica(nil)); err != nil {
		t.Fatal(err)
	}
	ja, _ := a.AvanzarJornada()
	jb, _ := b.AvanzarJornada()
	if !reflect.DeepEqual(ja, jb) {
		t.Error("elegir la misma alineacion que la automatica no deberia cambiar nada")
	}
}
