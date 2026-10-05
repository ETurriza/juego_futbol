package arquitectura

import (
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"testing"
)

const prefijoModulo = "github.com/ETurriza/juego_futbol/internal/"

// permitidos lista, por paquete de internal/, los paquetes del proyecto que
// puede importar. Refleja las reglas de dependencias de CLAUDE.md; un paquete
// nuevo debe agregarse aquí al crearse.
var permitidos = map[string][]string{
	"modelo":       {},
	"generador":    {"modelo"},
	"simulacion":   {"modelo", "liga", "mercado"},
	"liga":         {"modelo", "simulacion", "mercado"},
	"mercado":      {"modelo", "simulacion", "liga"},
	"partido":      {"modelo", "simulacion"},
	"aplicacion":   {"modelo", "generador", "simulacion", "liga", "mercado", "partido", "progresion"},
	"persistencia": {"modelo", "aplicacion"},
	"progresion":   {"modelo"},
	"menus":        {"aplicacion", "modelo"},
	"red":          {"aplicacion", "menus", "modelo"},
}

// violaciones devuelve los imports de pkg que no estan permitidos.
func violaciones(pkg string, imports []string, reglas map[string][]string) []string {
	ok := map[string]bool{}
	for _, p := range reglas[pkg] {
		ok[p] = true
	}
	var malos []string
	for _, imp := range imports {
		if imp != pkg && !ok[imp] {
			malos = append(malos, imp)
		}
	}
	sort.Strings(malos)
	return malos
}

// importsDelProyecto devuelve los paquetes internos importados por el codigo no
// de prueba del directorio dir.
func importsDelProyecto(t *testing.T, dir string) []string {
	t.Helper()
	archivos, err := filepath.Glob(filepath.Join(dir, "*.go"))
	if err != nil {
		t.Fatal(err)
	}
	vistos := map[string]bool{}
	for _, f := range archivos {
		if strings.HasSuffix(f, "_test.go") {
			continue
		}
		ast, err := parser.ParseFile(token.NewFileSet(), f, nil, parser.ImportsOnly)
		if err != nil {
			t.Fatalf("%s: %v", f, err)
		}
		for _, spec := range ast.Imports {
			ruta, err := strconv.Unquote(spec.Path.Value)
			if err != nil {
				t.Fatalf("%s: %v", f, err)
			}
			if resto, ok := strings.CutPrefix(ruta, prefijoModulo); ok {
				vistos[strings.SplitN(resto, "/", 2)[0]] = true
			}
		}
	}
	var imports []string
	for p := range vistos {
		imports = append(imports, p)
	}
	sort.Strings(imports)
	return imports
}

func TestReglasDeDependencias(t *testing.T) {
	entradas, err := os.ReadDir("..")
	if err != nil {
		t.Fatal(err)
	}
	for _, entrada := range entradas {
		pkg := entrada.Name()
		if !entrada.IsDir() || pkg == "arquitectura" {
			continue
		}
		if _, ok := permitidos[pkg]; !ok {
			t.Errorf("internal/%s no tiene reglas de dependencias; agregarlo a permitidos", pkg)
			continue
		}
		imports := importsDelProyecto(t, filepath.Join("..", pkg))
		if malos := violaciones(pkg, imports, permitidos); len(malos) > 0 {
			t.Errorf("internal/%s importa paquetes no permitidos: %v (permitidos: %v)",
				pkg, malos, permitidos[pkg])
		}
	}
}

func TestPermitidosSoloNombraPaquetesConocidos(t *testing.T) {
	for pkg, deps := range permitidos {
		for _, d := range deps {
			if _, ok := permitidos[d]; !ok {
				t.Errorf("%s permite a %s, que no tiene reglas propias", pkg, d)
			}
		}
	}
}

func TestModeloNoImportaNada(t *testing.T) {
	if len(permitidos["modelo"]) != 0 {
		t.Error("modelo no debe poder importar ningun paquete del proyecto")
	}
}

func TestPaquetesInternosNoImportanCapasExternas(t *testing.T) {
	externas := []string{"aplicacion", "menus", "persistencia", "red"}
	for _, pkg := range []string{"modelo", "simulacion", "liga", "mercado"} {
		for _, d := range permitidos[pkg] {
			for _, e := range externas {
				if d == e {
					t.Errorf("%s no puede depender de la capa externa %s", pkg, e)
				}
			}
		}
	}
}

func TestViolacionesDetectaImportsProhibidos(t *testing.T) {
	reglas := map[string][]string{"liga": {"modelo"}}
	malos := violaciones("liga", []string{"modelo", "persistencia", "red"}, reglas)
	if len(malos) != 2 || malos[0] != "persistencia" || malos[1] != "red" {
		t.Errorf("violaciones = %v, se esperaba [persistencia red]", malos)
	}
	if malos := violaciones("liga", []string{"modelo"}, reglas); len(malos) != 0 {
		t.Errorf("no deberia haber violaciones, hay %v", malos)
	}
}
