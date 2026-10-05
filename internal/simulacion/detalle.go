package simulacion

import (
	"math"
	"math/rand"
	"sort"

	"github.com/ETurriza/juego_futbol/internal/modelo"
)

// Constantes de calibración de los sucesos de un partido.
const (
	// Probabilidad de que un gol tenga asistencia.
	probAsistencia = 0.70

	// Tarjetas por equipo y partido: amarillas esperadas y probabilidad de una
	// roja directa. Las segundas amarillas se suman a las rojas.
	amarillasPorEquipo = 1.7
	probRojaDirecta    = 0.05
	// Un jugador ya amonestado tiene este peso relativo para otra amarilla.
	pesoYaAmonestado = 0.12

	// Sustituciones por equipo: de 3 a 5, entre estos minutos.
	cambiosMinutoMin = 46
	cambiosMinutoMax = 85
	// probabilidad de usar 3, 4 y 5 cambios (suman 1).
	probTresCambios  = 0.20
	probCuatroCambio = 0.30
	// Con esta probabilidad el que entra es de la misma posición del que sale.
	probMismaPosicion = 0.80
)

// Pesos por posición. Las proporciones buscadas: de cada diez goles, unos seis
// son de delanteros, casi tres de mediocampistas y menos de uno de defensas;
// las amarillas se reparten sobre todo entre defensas y medios.
var (
	pesoGol = map[modelo.Posicion]float64{
		modelo.Delantero: 1.0, modelo.Mediocampista: 0.6, modelo.Defensa: 0.2, modelo.Portero: 0.002,
	}
	pesoAsistencia = map[modelo.Posicion]float64{
		modelo.Mediocampista: 1.0, modelo.Delantero: 0.7, modelo.Defensa: 0.35, modelo.Portero: 0.01,
	}
	pesoTarjeta = map[modelo.Posicion]float64{
		modelo.Defensa: 1.0, modelo.Mediocampista: 1.0, modelo.Delantero: 0.6, modelo.Portero: 0.15,
	}
)

// enCampo es un jugador durante el partido: está en el campo en los minutos
// (desde, hasta]. Juega en un puesto, que no tiene por qué ser su posición natural.
type enCampo struct {
	j          modelo.Jugador
	puesto     modelo.Posicion
	desde      int
	hasta      int
	titular    bool
	amonestado bool
	// fuera indica que ya salió (cambio) o fue expulsado; ya no puede recibir
	// tarjetas ni ser cambiado, aunque los sucesos de su último minuto sigan
	// siendo válidos.
	fuera bool
}

func (c *enCampo) jugando(minuto int) bool { return minuto > c.desde && minuto <= c.hasta }

type tipoCandidato int

const (
	candCambio tipoCandidato = iota
	candRoja
	candAmarilla
)

type candidato struct {
	minuto int
	tipo   tipoCandidato
}

// tiempo es la línea de tiempo de un equipo en un partido: quién está en el campo
// y los cambios y tarjetas que ocurren.
type tiempo struct {
	local   bool
	campo   []*enCampo
	eventos []modelo.Evento
}

// fraccionConDiez es la fracción del partido que el equipo juega con un hombre
// menos por culpa de sus expulsiones (acotada a 1).
func (t *tiempo) fraccionConDiez() float64 {
	f := 0.0
	for _, e := range t.eventos {
		if e.Tipo == modelo.Roja {
			f += float64(modelo.MinutosPartido-e.Minuto) / modelo.MinutosPartido
		}
	}
	return math.Min(f, 1)
}

// banquilloOrdenado devuelve a los suplentes en orden de preferencia: primero los
// del banquillo de la alineación, en su orden, y después el resto de la
// plantilla que no es titular.
func banquilloOrdenado(e modelo.Equipo, al modelo.Alineacion) []modelo.Jugador {
	plantilla := map[int]modelo.Jugador{}
	for _, j := range e.Plantilla {
		plantilla[j.ID] = j
	}
	usados := map[int]bool{}
	for _, id := range al.Titulares {
		usados[id] = true
	}
	var banca []modelo.Jugador
	for _, id := range al.Banquillo {
		if j, ok := plantilla[id]; ok && !usados[id] {
			banca = append(banca, j)
			usados[id] = true
		}
	}
	for _, j := range e.Plantilla {
		if !usados[j.ID] {
			banca = append(banca, j)
		}
	}
	return banca
}

// lineaDeTiempo genera los cambios y las tarjetas de un equipo. Primero decide
// cuándo ocurre cada cosa y luego las resuelve en orden de minuto: quién sale,
// quién entra (según el orden del banquillo), quién es amonestado o expulsado.
// Todo es determinista dada la semilla.
func lineaDeTiempo(r *rand.Rand, local bool, e modelo.Equipo, al modelo.Alineacion) *tiempo {
	t := &tiempo{local: local}
	plantilla := map[int]modelo.Jugador{}
	for _, j := range e.Plantilla {
		plantilla[j.ID] = j
	}
	puestos := al.Formacion.Puestos()
	for i, id := range al.Titulares {
		if j, ok := plantilla[id]; ok {
			t.campo = append(t.campo, &enCampo{j: j, puesto: puestos[i], hasta: modelo.MinutosPartido, titular: true})
		}
	}
	disponibles := banquilloOrdenado(e, al)

	// 1. Cuándo ocurre cada cosa.
	var cands []candidato
	for _, m := range minutosDeCambio(r) {
		cands = append(cands, candidato{m, candCambio})
	}
	if r.Float64() < probRojaDirecta {
		cands = append(cands, candidato{10 + r.Intn(81), candRoja})
	}
	for i, n := 0, poisson(r, amarillasPorEquipo); i < n; i++ {
		cands = append(cands, candidato{1 + r.Intn(modelo.MinutosPartido), candAmarilla})
	}
	sort.SliceStable(cands, func(a, b int) bool {
		if cands[a].minuto != cands[b].minuto {
			return cands[a].minuto < cands[b].minuto
		}
		return cands[a].tipo < cands[b].tipo
	})

	// 2. Resolverlo en orden.
	for _, c := range cands {
		switch c.tipo {
		case candCambio:
			sale := elegirQueSale(r, t.campo, c.minuto)
			if sale == nil {
				continue
			}
			k := elegirQueEntra(r, disponibles, sale.puesto)
			if k < 0 {
				continue
			}
			entra := disponibles[k]
			disponibles = append(disponibles[:k], disponibles[k+1:]...)
			sale.hasta, sale.fuera = c.minuto, true
			t.campo = append(t.campo, &enCampo{j: entra, puesto: sale.puesto, desde: c.minuto, hasta: modelo.MinutosPartido})
			t.eventos = append(t.eventos, modelo.Evento{
				Minuto: c.minuto, Tipo: modelo.Sustitucion, Local: local, Jugador: sale.j.ID, Otro: entra.ID,
			})

		case candRoja:
			p := elegirTarjeta(r, t.campo, c.minuto, true)
			if p == nil {
				continue
			}
			p.hasta, p.fuera = c.minuto, true
			t.eventos = append(t.eventos, modelo.Evento{Minuto: c.minuto, Tipo: modelo.Roja, Local: local, Jugador: p.j.ID})

		case candAmarilla:
			p := elegirTarjeta(r, t.campo, c.minuto, false)
			if p == nil {
				continue
			}
			t.eventos = append(t.eventos, modelo.Evento{Minuto: c.minuto, Tipo: modelo.Amarilla, Local: local, Jugador: p.j.ID})
			if p.amonestado { // segunda amarilla: expulsión
				p.hasta, p.fuera = c.minuto, true
				t.eventos = append(t.eventos, modelo.Evento{Minuto: c.minuto, Tipo: modelo.Roja, Local: local, Jugador: p.j.ID})
			}
			p.amonestado = true
		}
	}
	return t
}

// goles reparte los goles del equipo (y sus asistencias) entre los jugadores que
// estaban en el campo en el minuto de cada gol, según el puesto en que jugaban.
func (t *tiempo) goles(r *rand.Rand, goles int) []modelo.Evento {
	var eventos []modelo.Evento
	for i := 0; i < goles; i++ {
		minuto := 1 + r.Intn(modelo.MinutosPartido)
		jugando := jugandoEn(t.campo, minuto)
		k := elegirPonderado(r, len(jugando), func(i int) float64 {
			c := jugando[i]
			return pesoGol[c.puesto] * float64(c.j.Atributos.Tiro)
		})
		if k < 0 {
			continue
		}
		goleador := jugando[k]
		ev := modelo.Evento{Minuto: minuto, Tipo: modelo.Gol, Local: t.local, Jugador: goleador.j.ID}
		if r.Float64() < probAsistencia {
			a := elegirPonderado(r, len(jugando), func(i int) float64 {
				c := jugando[i]
				if c.j.ID == goleador.j.ID {
					return 0
				}
				return pesoAsistencia[c.puesto] * float64(c.j.Atributos.Pase+c.j.Atributos.Regate) / 2
			})
			if a >= 0 {
				ev.Otro = jugando[a].j.ID
			}
		}
		eventos = append(eventos, ev)
	}
	return eventos
}

// minutosDeCambio sortea cuántos cambios hace un equipo (3 a 5) y en qué
// minutos distintos.
func minutosDeCambio(r *rand.Rand) []int {
	n := 5
	switch x := r.Float64(); {
	case x < probTresCambios:
		n = 3
	case x < probTresCambios+probCuatroCambio:
		n = 4
	}
	rango := cambiosMinutoMax - cambiosMinutoMin + 1
	minutos := make([]int, 0, n)
	for _, p := range r.Perm(rango)[:n] {
		minutos = append(minutos, cambiosMinutoMin+p)
	}
	return minutos
}

func jugandoEn(campo []*enCampo, minuto int) []*enCampo {
	var out []*enCampo
	for _, c := range campo {
		if c.jugando(minuto) {
			out = append(out, c)
		}
	}
	return out
}

// elegirQueSale elige, entre los titulares de campo que siguen jugando, a quien
// sale; con más probabilidad cuanto menor sea su físico. Los porteros no se
// cambian. Devuelve nil si no hay nadie.
func elegirQueSale(r *rand.Rand, campo []*enCampo, minuto int) *enCampo {
	var cand []*enCampo
	for _, c := range campo {
		if c.titular && !c.fuera && c.puesto != modelo.Portero && c.jugando(minuto) {
			cand = append(cand, c)
		}
	}
	k := elegirPonderado(r, len(cand), func(i int) float64 {
		return float64(modelo.AtributoMax + 2 - cand[i].j.Atributos.Fisico)
	})
	if k < 0 {
		return nil
	}
	return cand[k]
}

// elegirQueEntra elige un suplente de la banca, que está en orden de preferencia
// (nunca un portero): casi siempre el primero que juega la misma posición que el
// puesto que queda libre, y si no, el primero que pueda jugar en el campo.
// Devuelve su índice en la banca, o -1 si no hay. Siempre consume un número
// aleatorio.
func elegirQueEntra(r *rand.Rand, banca []modelo.Jugador, puesto modelo.Posicion) int {
	misma := r.Float64() < probMismaPosicion
	for _, soloMisma := range []bool{misma, false} {
		for k, j := range banca {
			if j.Posicion == modelo.Portero || (soloMisma && j.Posicion != puesto) {
				continue
			}
			return k
		}
	}
	return -1
}

// elegirTarjeta elige, entre quienes juegan en ese minuto, a quien recibe una
// tarjeta, con más probabilidad cuanto más defensivo y fuerte sea. Los porteros
// pueden ver amarilla, pero no se les expulsa con roja directa. Devuelve nil si
// no hay nadie.
func elegirTarjeta(r *rand.Rand, campo []*enCampo, minuto int, roja bool) *enCampo {
	jugando := jugandoEn(campo, minuto)
	k := elegirPonderado(r, len(jugando), func(i int) float64 {
		j := jugando[i].j
		if jugando[i].fuera || (roja && jugando[i].puesto == modelo.Portero) {
			return 0
		}
		w := pesoTarjeta[jugando[i].puesto] * float64(j.Atributos.Defensa+j.Atributos.Fisico) / 2
		if !roja && jugando[i].amonestado {
			w *= pesoYaAmonestado
		}
		return w
	})
	if k < 0 {
		return nil
	}
	return jugando[k]
}

// elegirPonderado devuelve un índice de 0 a n-1 elegido con probabilidad
// proporcional a peso(i), o -1 si n es 0 o ningún peso es positivo. Siempre
// consume un número aleatorio si hay algo que elegir.
func elegirPonderado(r *rand.Rand, n int, peso func(i int) float64) int {
	total := 0.0
	for i := 0; i < n; i++ {
		total += max(peso(i), 0)
	}
	if n == 0 || total <= 0 {
		return -1
	}
	x := r.Float64() * total
	for i := 0; i < n; i++ {
		x -= max(peso(i), 0)
		if x < 0 {
			return i
		}
	}
	return n - 1
}
