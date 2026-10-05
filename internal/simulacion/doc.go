// Package simulacion simula partidos de forma estadística a partir de los
// atributos de los jugadores.
//
// Cada equipo tiene una fuerza de ataque y una de defensa; los goles de cada
// lado se muestrean con una distribución de Poisson cuya esperanza depende de
// la razón entre el ataque propio y la defensa rival. Toda la aleatoriedad se
// recibe como *rand.Rand. Solo depende de modelo.
package simulacion
