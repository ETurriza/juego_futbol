package liga

import "fmt"

// Partido enfrenta a dos equipos, identificados por su índice en la lista de
// equipos de la temporada.
type Partido struct {
	Local     int
	Visitante int
}

// Jornada es el conjunto de partidos que se juegan a la vez. Ningún equipo
// aparece más de una vez en una jornada.
type Jornada []Partido

// Calendario genera el calendario todos contra todos, ida y vuelta, para n
// equipos con el método del círculo. Cada pareja se enfrenta dos veces, una
// como local y otra como visitante; la segunda vuelta repite la primera con las
// localías invertidas.
//
// Con n par hay 2(n-1) jornadas. Con n impar se agrega un descanso: hay 2n
// jornadas y cada equipo descansa en dos de ellas. El resultado es
// determinista.
func Calendario(n int) ([]Jornada, error) {
	if n < 2 {
		return nil, fmt.Errorf("se necesitan al menos 2 equipos, hay %d", n)
	}
	m := n + n%2 // con impares se suma un equipo fantasma (índice n) = descanso
	orden := make([]int, m)
	for i := range orden {
		orden[i] = i
	}

	var ida []Jornada
	for ronda := 0; ronda < m-1; ronda++ {
		var jornada Jornada
		for i := 0; i < m/2; i++ {
			local, visitante := orden[i], orden[m-1-i]
			if (ronda+i)%2 == 1 {
				local, visitante = visitante, local
			}
			if local >= n || visitante >= n {
				continue // descanso
			}
			jornada = append(jornada, Partido{Local: local, Visitante: visitante})
		}
		ida = append(ida, jornada)

		// Rota todos menos el primero una posición a la derecha.
		ultimo := orden[m-1]
		copy(orden[2:], orden[1:m-1])
		orden[1] = ultimo
	}

	calendario := append([]Jornada(nil), ida...)
	for _, j := range ida {
		vuelta := make(Jornada, len(j))
		for i, p := range j {
			vuelta[i] = Partido{Local: p.Visitante, Visitante: p.Local}
		}
		calendario = append(calendario, vuelta)
	}
	return calendario, nil
}
