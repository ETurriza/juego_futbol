// Package progresion modela cómo cambian los jugadores con la edad: su
// evolución de un año al siguiente y su retiro.
//
// Los atributos técnicos (tiro, pase, regate, defensa y los reflejos de los
// porteros) se sostienen más que los físicos (ritmo y físico): la técnica tiene
// una meseta larga, con un posible "último prime" pasados los 30, y el declive
// llega rápido después de los 34. Los porteros envejecen unos años más tarde.
//
// El crecimiento tiene rendimientos decrecientes cerca del tope de los atributos
// (la caída por edad, no), de modo que los atributos no se apilan en 99.
//
// El crecimiento lo escala el talento oculto del jugador (modelo.Talento): los
// de mucho talento crecen más rápido y llegan más lejos, los de poco, menos.
//
// El retiro depende de la edad y también del nivel: un veterano que ya no da el
// nivel de un titular se retira antes.
//
// Toda la aleatoriedad se recibe como *rand.Rand. Solo depende de modelo.
package progresion
