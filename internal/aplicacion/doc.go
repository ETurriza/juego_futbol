// Package aplicacion contiene los casos de uso del juego: crear una carrera,
// consultar la plantilla y la tabla, y avanzar jornada. Es la capa que usan los
// entrypoints (menus y red); ellos no hablan directamente con liga, simulacion
// ni generador.
//
// Cada jornada usa su propio *rand.Rand, derivado de la semilla de la carrera y
// del número de jornada. Así una carrera queda determinada por su semilla y su
// avance, sin necesidad de guardar el estado de un generador aleatorio.
package aplicacion
