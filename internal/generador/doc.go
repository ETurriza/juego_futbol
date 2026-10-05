// Package generador crea jugadores y equipos con nombres inventados.
//
// Los jugadores (iniciales y juveniles) nacen con 16 años y crecen con la propia
// progresión, de modo que su calidad por edad es la misma que tendrán durante la
// carrera.
//
// Toda la aleatoriedad se recibe como *rand.Rand, nunca global, para que la
// generación sea reproducible con una semilla fija. Depende de modelo y de
// progresion.
package generador
