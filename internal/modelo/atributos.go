package modelo

import "fmt"

// Límites de cada atributo.
const (
	AtributoMin = 1
	AtributoMax = 99
)

// Atributos son las cualidades de un jugador, cada una entre AtributoMin y
// AtributoMax. Reflejos solo pesa en la valoración de los porteros.
type Atributos struct {
	Ritmo    int
	Tiro     int
	Pase     int
	Regate   int
	Defensa  int
	Fisico   int
	Reflejos int
}

// pesos de cada atributo por posición; cada fila suma 100.
var pesos = map[Posicion]Atributos{
	Portero:       {Ritmo: 10, Tiro: 5, Pase: 15, Regate: 5, Defensa: 15, Fisico: 15, Reflejos: 35},
	Defensa:       {Ritmo: 15, Tiro: 5, Pase: 10, Regate: 5, Defensa: 40, Fisico: 25, Reflejos: 0},
	Mediocampista: {Ritmo: 10, Tiro: 15, Pase: 35, Regate: 20, Defensa: 10, Fisico: 10, Reflejos: 0},
	Delantero:     {Ritmo: 20, Tiro: 35, Pase: 10, Regate: 20, Defensa: 5, Fisico: 10, Reflejos: 0},
}

// Media devuelve la valoración global (0-99) de los atributos para la posición
// dada, ponderando cada atributo según lo que importa en ese puesto. Una
// posición inválida devuelve 0.
func (a Atributos) Media(p Posicion) int {
	w, ok := pesos[p]
	if !ok {
		return 0
	}
	suma := a.Ritmo*w.Ritmo + a.Tiro*w.Tiro + a.Pase*w.Pase + a.Regate*w.Regate +
		a.Defensa*w.Defensa + a.Fisico*w.Fisico + a.Reflejos*w.Reflejos
	return (suma + 50) / 100 // redondeo al entero más cercano
}

// Validar comprueba que todos los atributos estén dentro de rango.
func (a Atributos) Validar() error {
	campos := []struct {
		nombre string
		valor  int
	}{
		{"ritmo", a.Ritmo}, {"tiro", a.Tiro}, {"pase", a.Pase},
		{"regate", a.Regate}, {"defensa", a.Defensa},
		{"fisico", a.Fisico}, {"reflejos", a.Reflejos},
	}
	for _, c := range campos {
		if c.valor < AtributoMin || c.valor > AtributoMax {
			return fmt.Errorf("atributo %s fuera de rango [%d,%d]: %d",
				c.nombre, AtributoMin, AtributoMax, c.valor)
		}
	}
	return nil
}
