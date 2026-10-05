package aplicacion

import (
	"context"
	"errors"
	"reflect"
	"strings"
	"testing"
	"time"
)

func guardadoDePrueba(t *testing.T, jornadas int) Guardado {
	t.Helper()
	c := nuevaCarrera(t, 5, 10)
	for i := 0; i < jornadas; i++ {
		if _, err := c.AvanzarJornada(); err != nil {
			t.Fatal(err)
		}
	}
	return c.Exportar()
}

func TestExportarImportarIdaYVuelta(t *testing.T) {
	for _, jornadas := range []int{0, 1, 7, 18} {
		c := nuevaCarrera(t, 5, 10)
		for i := 0; i < jornadas; i++ {
			c.AvanzarJornada()
		}
		reconstruida, err := Importar(c.Exportar())
		if err != nil {
			t.Fatalf("%d jornadas: %v", jornadas, err)
		}
		if !reflect.DeepEqual(c, reconstruida) {
			t.Errorf("%d jornadas: la carrera reconstruida no es identica", jornadas)
		}
	}
}

func TestCarreraImportadaContinuaIgualQueLaContinua(t *testing.T) {
	continua := nuevaCarrera(t, 8, 10)
	for i := 0; i < 6; i++ {
		continua.AvanzarJornada()
	}
	reanudada, err := Importar(continua.Exportar())
	if err != nil {
		t.Fatal(err)
	}
	for !continua.Terminada() {
		a, _ := continua.AvanzarJornada()
		b, err := reanudada.AvanzarJornada()
		if err != nil || !reflect.DeepEqual(a, b) {
			t.Fatalf("jornada %d difiere al reanudar: %v", continua.Jornada(), err)
		}
	}
	if !reflect.DeepEqual(continua.Tabla(), reanudada.Tabla()) {
		t.Error("las tablas finales deberian coincidir")
	}
}

func TestExportarNoComparteMemoria(t *testing.T) {
	c := nuevaCarrera(t, 5, 10)
	c.AvanzarJornada()
	g := c.Exportar()
	g.Equipos[0].Plantilla[0].Nombre = "Alterado"
	g.Resultados[0][0].GolesLocal = 99
	if c.Temporada.Equipos[0].Plantilla[0].Nombre == "Alterado" || c.Temporada.Resultados[0][0].GolesLocal == 99 {
		t.Error("Exportar comparte memoria con la carrera")
	}

	g = c.Exportar()
	reconstruida, _ := Importar(g)
	g.Equipos[1].Plantilla[0].Nombre = "Alterado"
	if reconstruida.Temporada.Equipos[1].Plantilla[0].Nombre == "Alterado" {
		t.Error("Importar comparte memoria con el Guardado")
	}
}

func TestImportarRechazaDatosIncoherentes(t *testing.T) {
	casos := map[string]func(*Guardado){
		"usuario negativo":     func(g *Guardado) { g.Usuario = -1 },
		"usuario fuera":        func(g *Guardado) { g.Usuario = len(g.Equipos) },
		"sin equipos":          func(g *Guardado) { g.Equipos = nil },
		"un solo equipo":       func(g *Guardado) { g.Equipos = g.Equipos[:1] },
		"nombre repetido":      func(g *Guardado) { g.Equipos[1].Nombre = g.Equipos[0].Nombre },
		"equipo invalido":      func(g *Guardado) { g.Equipos[2].Plantilla[0].Edad = 0 },
		"ID repetida en liga":  func(g *Guardado) { g.Equipos[1].Plantilla[0].ID = g.Equipos[0].Plantilla[0].ID },
		"demasiadas jornadas":  func(g *Guardado) { g.Resultados = append(g.Resultados, g.Resultados...) },
		"jornada incompleta":   func(g *Guardado) { g.Resultados[1] = g.Resultados[1][:2] },
		"partido distinto":     func(g *Guardado) { g.Resultados[0][0].Local, g.Resultados[0][0].Visitante = 0, 0 },
		"goles negativos":      func(g *Guardado) { g.Resultados[2][1].GolesVisitante = -1 },
		"orden de partidos":    func(g *Guardado) { g.Resultados[0][0], g.Resultados[0][1] = g.Resultados[0][1], g.Resultados[0][0] },
		"jornada con relleno":  func(g *Guardado) { g.Resultados[0] = append(g.Resultados[0], ResultadoGuardado{}) },
		"jornada sin partidos": func(g *Guardado) { g.Resultados[0] = nil },
		"resultados de sobra": func(g *Guardado) {
			g.Resultados = append(g.Resultados, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil)
		},
	}
	for nombre, mutar := range casos {
		g := guardadoDePrueba(t, 4)
		mutar(&g)
		if _, err := Importar(g); err == nil {
			t.Errorf("%s: deberia dar error", nombre)
		} else if !strings.Contains(err.Error(), "guardado invalido") {
			t.Errorf("%s: el error deberia indicar un guardado invalido: %v", nombre, err)
		}
	}
	if _, err := Importar(guardadoDePrueba(t, 4)); err != nil {
		t.Errorf("el guardado sin alterar deberia importarse: %v", err)
	}
}

func TestResumen(t *testing.T) {
	g := guardadoDePrueba(t, 5)
	r, err := g.Resumen("mi ranura")
	if err != nil {
		t.Fatal(err)
	}
	want := ResumenPartida{Ranura: "mi ranura", Equipo: g.Equipos[g.Usuario].Nombre, Jornada: 5, TotalJornadas: 18}
	if r != want {
		t.Errorf("Resumen = %+v, se esperaba %+v", r, want)
	}
	g.Usuario = 99
	if _, err := g.Resumen("x"); err == nil {
		t.Error("un usuario fuera de rango deberia dar error")
	}
}

func TestValidarRanura(t *testing.T) {
	validas := []string{"principal", "a", "Carrera de Ñandú", "con-guion_bajo.1", strings.Repeat("x", MaxLongitudRanura)}
	for _, r := range validas {
		if err := ValidarRanura(r); err != nil {
			t.Errorf("%q deberia ser valida: %v", r, err)
		}
	}
	invalidas := []string{"", " ", "\t", " x", "x ", "a\nb", "a\x00b", strings.Repeat("x", MaxLongitudRanura+1)}
	for _, r := range invalidas {
		if err := ValidarRanura(r); err == nil {
			t.Errorf("%q deberia ser invalida", r)
		}
	}
}

func TestGuardarYCargarCarrera(t *testing.T) {
	ctx := context.Background()
	repo := NuevoRepositorioMemoria()
	c := nuevaCarrera(t, 5, 10)
	for i := 0; i < 4; i++ {
		c.AvanzarJornada()
	}
	if err := GuardarCarrera(ctx, repo, "principal", c); err != nil {
		t.Fatal(err)
	}
	cargada, err := CargarCarrera(ctx, repo, "principal")
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(c, cargada) {
		t.Error("la carrera cargada no es identica a la guardada")
	}

	if _, err := CargarCarrera(ctx, repo, "otra"); !errors.Is(err, ErrPartidaNoExiste) {
		t.Errorf("err = %v, se esperaba ErrPartidaNoExiste", err)
	}
	if err := GuardarCarrera(ctx, repo, "", c); err == nil {
		t.Error("una ranura vacia deberia rechazarse al guardar")
	}
	if _, err := CargarCarrera(ctx, repo, ""); err == nil {
		t.Error("una ranura vacia deberia rechazarse al cargar")
	}
}

func TestCargarCarreraRechazaGuardadoCorrupto(t *testing.T) {
	ctx := context.Background()
	repo := NuevoRepositorioMemoria()
	g := guardadoDePrueba(t, 3)
	g.Resultados[0][0].Local = 7 // no coincide con el calendario
	if err := repo.Guardar(ctx, "mala", g); err != nil {
		t.Fatal(err)
	}
	if _, err := CargarCarrera(ctx, repo, "mala"); err == nil {
		t.Error("un guardado corrupto deberia rechazarse al cargar")
	}
}

func TestRepositorioMemoriaOrdenaPorFecha(t *testing.T) {
	repo := NuevoRepositorioMemoria()
	reloj := time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC)
	repo.ahora = func() time.Time { return reloj }
	g := guardadoDePrueba(t, 1)

	for _, ranura := range []string{"b", "a", "c"} {
		if err := repo.Guardar(context.Background(), ranura, g); err != nil {
			t.Fatal(err)
		}
		reloj = reloj.Add(time.Minute)
	}
	lista, _ := repo.Listar(context.Background())
	var orden []string
	for _, r := range lista {
		orden = append(orden, r.Ranura)
	}
	if want := []string{"c", "a", "b"}; !reflect.DeepEqual(orden, want) {
		t.Errorf("orden = %v, se esperaba %v (mas reciente primero)", orden, want)
	}

	// Con la misma hora, desempata por nombre.
	igual := NuevoRepositorioMemoria()
	igual.ahora = func() time.Time { return reloj }
	for _, ranura := range []string{"z", "y", "x"} {
		if err := igual.Guardar(context.Background(), ranura, g); err != nil {
			t.Fatal(err)
		}
	}
	lista, _ = igual.Listar(context.Background())
	orden = nil
	for _, r := range lista {
		orden = append(orden, r.Ranura)
	}
	if want := []string{"x", "y", "z"}; !reflect.DeepEqual(orden, want) {
		t.Errorf("con la misma hora, orden = %v, se esperaba %v", orden, want)
	}
}
