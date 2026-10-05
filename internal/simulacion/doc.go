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
// Cada equipo juega con una alineación (modelo.Alineacion): una formación y once
// titulares con su banquillo. Quien juega fuera de su posición rinde según sus
// atributos en ese puesto por una familiaridad menor que 1. Las formaciones
// tienen perfiles de ataque y defensa distintos pero de efecto neto parecido: lo
// que decide es la plantilla. AlineacionAutomatica elige la formación y el once
// que mejor le van a un equipo (con criterios del entrenador y sin los jugadores
// no disponibles), y es la que usan los rivales y el modo jugador. Una roja deja
// a un equipo con diez y cambia la esperanza de goles el resto del partido.
//
// SimularConAlineaciones juega un partido con las alineaciones dadas y Simular
// deja que cada equipo elija las suyas; ambos devuelven el marcador con su
// detalle (alineaciones y sucesos). Marcador devuelve solo el resultado, con un
// 4-3-3 automático y sin tarjetas, mucho más rápido. Toda la aleatoriedad se
// recibe como *rand.Rand. Solo depende de modelo.
package simulacion
