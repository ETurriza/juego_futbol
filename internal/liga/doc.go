// Package liga organiza una competición todos contra todos: el calendario de
// jornadas (ida y vuelta), la simulación jornada a jornada y la tabla de
// posiciones.
//
// Los equipos se identifican por su índice en la lista que recibe la temporada.
// Toda la aleatoriedad se recibe como *rand.Rand. Solo depende de modelo y
// simulacion; los tipos son simples y sin estado oculto para poder
// persistirlos más adelante.
package liga
