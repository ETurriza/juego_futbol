package modelo

// Estadisticas acumula el rendimiento de un jugador en un periodo (una
// temporada, o toda la carrera).
type Estadisticas struct {
	Partidos           int // partidos con algún minuto jugado
	Titularidades      int
	Minutos            int
	Goles              int
	Asistencias        int
	Amarillas          int
	Rojas              int
	PorteriasImbatidas int // porteros y defensas, en partidos sin goles en contra
	GolesEncajados     int // goles recibidos mientras estaba en el campo
	// SumaValoracion es la suma de las valoraciones de cada partido, en décimas
	// (65 es un 6,5).
	SumaValoracion int
}

// ValoracionMedia es la valoración media por partido, de 1 a 10; 0 si no jugó.
func (e Estadisticas) ValoracionMedia() float64 {
	if e.Partidos == 0 {
		return 0
	}
	return float64(e.SumaValoracion) / 10 / float64(e.Partidos)
}

// Sumar devuelve la suma de dos acumulados.
func (e Estadisticas) Sumar(o Estadisticas) Estadisticas {
	return Estadisticas{
		Partidos:           e.Partidos + o.Partidos,
		Titularidades:      e.Titularidades + o.Titularidades,
		Minutos:            e.Minutos + o.Minutos,
		Goles:              e.Goles + o.Goles,
		Asistencias:        e.Asistencias + o.Asistencias,
		Amarillas:          e.Amarillas + o.Amarillas,
		Rojas:              e.Rojas + o.Rojas,
		PorteriasImbatidas: e.PorteriasImbatidas + o.PorteriasImbatidas,
		GolesEncajados:     e.GolesEncajados + o.GolesEncajados,
		SumaValoracion:     e.SumaValoracion + o.SumaValoracion,
	}
}
