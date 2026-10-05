// Package simulacion simula partidos de forma estadística a partir de los
// atributos de los jugadores.
//
// Cada equipo tiene una fuerza de ataque y una de defensa; los goles de cada
// lado se muestrean con una distribución de Poisson cuya esperanza depende de
// la razón entre el ataque propio y la defensa rival. La fuerza de cada línea es
// la media ponderada de sus jugadores más un bono por sus cracks (jugadores de
// valoración superior a 85), con un tope de saturación, de modo que un crack se
// nota en el resultado pero muchos cracks no deciden solos el partido.
//
// Simular devuelve el marcador con su detalle (alineaciones y sucesos); Marcador
// devuelve solo el resultado, mucho más rápido. Toda la aleatoriedad se recibe
// como *rand.Rand. Solo depende de modelo.
package simulacion
