// Package contrato ofrece las pruebas de contrato de aplicacion.RepositorioPartidas.
// Cada implementación del puerto (memoria, SQLite...) las ejecuta con su propia
// fábrica, de modo que todas cumplan exactamente el mismo comportamiento.
package contrato

import (
	"context"
	"errors"
	"reflect"
	"strings"
	"testing"

	"github.com/ETurriza/juego_futbol/internal/aplicacion"
)

// Fabrica crea un repositorio vacío y aislado para una prueba.
type Fabrica func(t *testing.T) aplicacion.RepositorioPartidas

// carrera devuelve el Guardado de una carrera de 10 equipos con n jornadas
// jugadas.
func carrera(t *testing.T, semilla int64, jornadas int) aplicacion.Guardado {
	t.Helper()
	c, err := aplicacion.NuevaCarrera(semilla, 10)
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < jornadas; i++ {
		if _, err := c.AvanzarJornada(); err != nil {
			t.Fatal(err)
		}
	}
	return c.Exportar()
}

// Ejecutar corre todas las pruebas de contrato contra repositorios de nuevo.
func Ejecutar(t *testing.T, nuevo Fabrica) {
	pruebas := []struct {
		nombre string
		fn     func(*testing.T, aplicacion.RepositorioPartidas)
	}{
		{"CargarInexistente", cargarInexistente},
		{"GuardarYCargar", guardarYCargar},
		{"GuardarSinJornadas", guardarSinJornadas},
		{"GuardarReemplaza", guardarReemplaza},
		{"RanurasIndependientes", ranurasIndependientes},
		{"Listar", listar},
		{"Borrar", borrar},
		{"RanuraInvalida", ranuraInvalida},
		{"ContextoCancelado", contextoCancelado},
		{"NoComparteEstado", noComparteEstado},
	}
	for _, p := range pruebas {
		t.Run(p.nombre, func(t *testing.T) { p.fn(t, nuevo(t)) })
	}
}

func cargarInexistente(t *testing.T, repo aplicacion.RepositorioPartidas) {
	_, err := repo.Cargar(context.Background(), "nada")
	if !errors.Is(err, aplicacion.ErrPartidaNoExiste) {
		t.Errorf("err = %v, se esperaba ErrPartidaNoExiste", err)
	}
}

func guardarYCargar(t *testing.T, repo aplicacion.RepositorioPartidas) {
	ctx := context.Background()
	original := carrera(t, 42, 5)
	if err := repo.Guardar(ctx, "principal", original); err != nil {
		t.Fatal(err)
	}
	cargado, err := repo.Cargar(ctx, "principal")
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(original, cargado) {
		t.Fatal("lo cargado no es identico a lo guardado")
	}
	// Y se puede reconstruir una carrera válida a partir de ello.
	if _, err := aplicacion.Importar(cargado); err != nil {
		t.Errorf("lo cargado no se puede importar: %v", err)
	}
}

func guardarSinJornadas(t *testing.T, repo aplicacion.RepositorioPartidas) {
	ctx := context.Background()
	original := carrera(t, 7, 0)
	if err := repo.Guardar(ctx, "nueva", original); err != nil {
		t.Fatal(err)
	}
	cargado, err := repo.Cargar(ctx, "nueva")
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(original, cargado) {
		t.Fatal("una carrera sin jornadas jugadas no se recupera identica")
	}
}

func guardarReemplaza(t *testing.T, repo aplicacion.RepositorioPartidas) {
	ctx := context.Background()
	if err := repo.Guardar(ctx, "p", carrera(t, 1, 3)); err != nil {
		t.Fatal(err)
	}
	segunda := carrera(t, 2, 8) // otra liga, con mas jornadas
	if err := repo.Guardar(ctx, "p", segunda); err != nil {
		t.Fatal(err)
	}
	cargado, err := repo.Cargar(ctx, "p")
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(segunda, cargado) {
		t.Error("guardar sobre una ranura existente deberia reemplazarla por completo")
	}
	lista, _ := repo.Listar(ctx)
	if len(lista) != 1 {
		t.Errorf("hay %d partidas, se esperaba 1", len(lista))
	}
}

func ranurasIndependientes(t *testing.T, repo aplicacion.RepositorioPartidas) {
	ctx := context.Background()
	a, b := carrera(t, 1, 2), carrera(t, 2, 6)
	if err := repo.Guardar(ctx, "a", a); err != nil {
		t.Fatal(err)
	}
	if err := repo.Guardar(ctx, "b", b); err != nil {
		t.Fatal(err)
	}
	if err := repo.Borrar(ctx, "a"); err != nil {
		t.Fatal(err)
	}
	cargadoB, err := repo.Cargar(ctx, "b")
	if err != nil || !reflect.DeepEqual(b, cargadoB) {
		t.Errorf("borrar la ranura a afecto a la b: %v", err)
	}
}

func listar(t *testing.T, repo aplicacion.RepositorioPartidas) {
	ctx := context.Background()
	vacia, err := repo.Listar(ctx)
	if err != nil || len(vacia) != 0 {
		t.Fatalf("Listar en vacio = %v, %v", vacia, err)
	}

	a, b := carrera(t, 1, 4), carrera(t, 2, 0)
	if err := repo.Guardar(ctx, "antigua", a); err != nil {
		t.Fatal(err)
	}
	if err := repo.Guardar(ctx, "reciente", b); err != nil {
		t.Fatal(err)
	}
	lista, err := repo.Listar(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(lista) != 2 {
		t.Fatalf("%d partidas, se esperaban 2", len(lista))
	}
	for i := 1; i < len(lista); i++ {
		if lista[i].Actualizada.After(lista[i-1].Actualizada) {
			t.Error("Listar deberia ordenar de la mas reciente a la mas antigua")
		}
	}
	porRanura := map[string]aplicacion.ResumenPartida{}
	for _, r := range lista {
		porRanura[r.Ranura] = r
		if r.Actualizada.IsZero() {
			t.Errorf("%s: sin fecha de actualizacion", r.Ranura)
		}
	}
	if r := porRanura["antigua"]; r.Equipo != a.Equipos[a.Usuario].Nombre || r.Jornada != 4 || r.TotalJornadas != 18 {
		t.Errorf("resumen de antigua incorrecto: %+v", r)
	}
	if r := porRanura["reciente"]; r.Equipo != b.Equipos[b.Usuario].Nombre || r.Jornada != 0 || r.TotalJornadas != 18 {
		t.Errorf("resumen de reciente incorrecto: %+v", r)
	}
}

func borrar(t *testing.T, repo aplicacion.RepositorioPartidas) {
	ctx := context.Background()
	if err := repo.Borrar(ctx, "nada"); !errors.Is(err, aplicacion.ErrPartidaNoExiste) {
		t.Errorf("borrar inexistente: err = %v, se esperaba ErrPartidaNoExiste", err)
	}
	if err := repo.Guardar(ctx, "p", carrera(t, 1, 2)); err != nil {
		t.Fatal(err)
	}
	if err := repo.Borrar(ctx, "p"); err != nil {
		t.Fatal(err)
	}
	if _, err := repo.Cargar(ctx, "p"); !errors.Is(err, aplicacion.ErrPartidaNoExiste) {
		t.Errorf("tras borrar, Cargar: err = %v", err)
	}
	if err := repo.Borrar(ctx, "p"); !errors.Is(err, aplicacion.ErrPartidaNoExiste) {
		t.Errorf("borrar dos veces: err = %v", err)
	}
}

func ranuraInvalida(t *testing.T, repo aplicacion.RepositorioPartidas) {
	g := carrera(t, 1, 1)
	for _, ranura := range []string{"", "   ", " x", "x ", "a\nb", strings.Repeat("x", aplicacion.MaxLongitudRanura+1)} {
		if err := repo.Guardar(context.Background(), ranura, g); err == nil {
			t.Errorf("la ranura %q deberia rechazarse", ranura)
		}
	}
	// Los nombres con acentos, espacios internos y largo máximo sí valen.
	for _, ranura := range []string{"Carrera de Ñandú", strings.Repeat("ñ", aplicacion.MaxLongitudRanura)} {
		if err := repo.Guardar(context.Background(), ranura, g); err != nil {
			t.Errorf("la ranura %q deberia aceptarse: %v", ranura, err)
		}
	}
}

func contextoCancelado(t *testing.T, repo aplicacion.RepositorioPartidas) {
	ctx, cancelar := context.WithCancel(context.Background())
	cancelar()
	if err := repo.Guardar(ctx, "p", carrera(t, 1, 1)); err == nil {
		t.Error("Guardar con contexto cancelado deberia fallar")
	}
	if _, err := repo.Cargar(ctx, "p"); err == nil {
		t.Error("Cargar con contexto cancelado deberia fallar")
	}
	if _, err := repo.Listar(ctx); err == nil {
		t.Error("Listar con contexto cancelado deberia fallar")
	}
	if err := repo.Borrar(ctx, "p"); err == nil {
		t.Error("Borrar con contexto cancelado deberia fallar")
	}
	// Una guardada cancelada no deja nada a medias.
	if _, err := repo.Cargar(context.Background(), "p"); !errors.Is(err, aplicacion.ErrPartidaNoExiste) {
		t.Errorf("tras un guardado cancelado no deberia haber partida: %v", err)
	}
}

func noComparteEstado(t *testing.T, repo aplicacion.RepositorioPartidas) {
	ctx := context.Background()
	original := carrera(t, 3, 3)
	if err := repo.Guardar(ctx, "p", original); err != nil {
		t.Fatal(err)
	}
	// Alterar lo guardado despues de Guardar no cambia lo almacenado...
	original.Equipos[0].Plantilla[0].Nombre = "Alterado"
	original.Resultados[0][0].GolesLocal = 99
	cargado, _ := repo.Cargar(ctx, "p")
	if cargado.Equipos[0].Plantilla[0].Nombre == "Alterado" || cargado.Resultados[0][0].GolesLocal == 99 {
		t.Error("el repositorio comparte memoria con el Guardado que recibio")
	}
	// ...ni alterar lo cargado cambia lo almacenado.
	cargado.Equipos[0].Plantilla[0].Nombre = "Otro"
	otra, _ := repo.Cargar(ctx, "p")
	if otra.Equipos[0].Plantilla[0].Nombre == "Otro" {
		t.Error("el repositorio devuelve memoria compartida")
	}
}
