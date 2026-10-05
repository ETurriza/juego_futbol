package liga

import (
	"reflect"
	"strings"
	"testing"

	"github.com/ETurriza/juego_futbol/internal/modelo"
)

// posicionDePrueba asigna las posiciones de los ids de los partidos de mano:
// 1 portero; 2 a 5 y 13 defensas; 6 a 8 y 14 medios; 9 a 12 delanteros. El
// equipo visitante usa los mismos números más 100.
func posicionDePrueba(id int) modelo.Posicion {
	switch id % 100 {
	case 1:
		return modelo.Portero
	case 2, 3, 4, 5, 13:
		return modelo.Defensa
	case 6, 7, 8, 14:
		return modelo.Mediocampista
	}
	return modelo.Delantero
}

// detalleDePrueba tiene los titulares 1 a 11 (local) y 101 a 111 (visitante).
func detalleDePrueba(eventos ...modelo.Evento) modelo.DetallePartido {
	var d modelo.DetallePartido
	for i := 0; i < modelo.TitularesPorEquipo; i++ {
		d.TitularesLocal[i] = 1 + i
		d.TitularesVisitante[i] = 101 + i
	}
	d.Eventos = eventos
	return d
}

func TestEstadisticasDeUnPartidoConResultadoYSucesos(t *testing.T) {
	d := detalleDePrueba(
		modelo.Evento{Minuto: 10, Tipo: modelo.Gol, Local: true, Jugador: 9, Otro: 6},
		modelo.Evento{Minuto: 30, Tipo: modelo.Amarilla, Local: true, Jugador: 3},
		modelo.Evento{Minuto: 50, Tipo: modelo.Gol, Local: false, Jugador: 111},
		modelo.Evento{Minuto: 60, Tipo: modelo.Sustitucion, Local: true, Jugador: 9, Otro: 12},
		modelo.Evento{Minuto: 70, Tipo: modelo.Gol, Local: true, Jugador: 12, Otro: 8},
		modelo.Evento{Minuto: 80, Tipo: modelo.Roja, Local: false, Jugador: 103},
	)
	// El marcador 2-1 incluye un gol local mas (minuto 10 y 70) y uno visitante.
	est, err := estadisticasDePartido(d, 2, 1)
	if err != nil {
		t.Fatal(err)
	}
	casos := map[string]struct {
		id   int
		want modelo.Estadisticas
	}{
		"goleador que sale":   {9, modelo.Estadisticas{Partidos: 1, Titularidades: 1, Minutos: 60, Goles: 1, GolesEncajados: 1, SumaValoracion: 73}},
		"asistente":           {6, modelo.Estadisticas{Partidos: 1, Titularidades: 1, Minutos: 90, Asistencias: 1, GolesEncajados: 1, SumaValoracion: 69}},
		"suplente goleador":   {12, modelo.Estadisticas{Partidos: 1, Minutos: 30, Goles: 1, SumaValoracion: 66}},
		"defensa amonestado":  {3, modelo.Estadisticas{Partidos: 1, Titularidades: 1, Minutos: 90, Amarillas: 1, GolesEncajados: 1, SumaValoracion: 58}},
		"defensa sin sucesos": {2, modelo.Estadisticas{Partidos: 1, Titularidades: 1, Minutos: 90, GolesEncajados: 1, SumaValoracion: 61}},
		"portero ganador":     {1, modelo.Estadisticas{Partidos: 1, Titularidades: 1, Minutos: 90, GolesEncajados: 1, SumaValoracion: 60}},
		"portero perdedor":    {101, modelo.Estadisticas{Partidos: 1, Titularidades: 1, Minutos: 90, GolesEncajados: 2, SumaValoracion: 51}},
		"goleador perdedor":   {111, modelo.Estadisticas{Partidos: 1, Titularidades: 1, Minutos: 90, Goles: 1, GolesEncajados: 2, SumaValoracion: 67}},
		"expulsado":           {103, modelo.Estadisticas{Partidos: 1, Titularidades: 1, Minutos: 80, Rojas: 1, GolesEncajados: 2, SumaValoracion: 38}},
	}
	for nombre, c := range casos {
		if got := est[c.id]; got != c.want {
			t.Errorf("%s (jugador %d):\n  obtenido %+v\n  esperado %+v", nombre, c.id, got, c.want)
		}
	}
	if len(est) != 23 {
		t.Errorf("%d jugadores con estadisticas, se esperaban 23 (22 titulares y 1 suplente)", len(est))
	}
}

func TestPorteriaImbatidaSoloConSesentaMinutos(t *testing.T) {
	d := detalleDePrueba(
		modelo.Evento{Minuto: 20, Tipo: modelo.Gol, Local: true, Jugador: 9},
		modelo.Evento{Minuto: 50, Tipo: modelo.Sustitucion, Local: true, Jugador: 4, Otro: 13},
		modelo.Evento{Minuto: 70, Tipo: modelo.Sustitucion, Local: true, Jugador: 5, Otro: 14},
	)
	est, err := estadisticasDePartido(d, 1, 0)
	if err != nil {
		t.Fatal(err)
	}
	for id, imbatida := range map[int]int{
		1: 1, 2: 1, 3: 1, // titulares todo el partido
		5: 1,        // sale en el 70: juego 70 minutos
		4: 0, 13: 0, // 50 y 40 minutos: no llegan a 60
		6: 0, 9: 0, // medios y delanteros nunca
		14: 0,
	} {
		if got := est[id].PorteriasImbatidas; got != imbatida {
			t.Errorf("jugador %d: %d imbatidas, se esperaba %d", id, got, imbatida)
		}
	}
	// El portero y los defensas valoran la imbatida, y quien juega menos de una
	// hora ve su valoracion acercada al 6,0.
	if got := est[1].SumaValoracion; got != 68 { // 60 + 5 + 3
		t.Errorf("portero imbatido: %d, se esperaba 68", got)
	}
	if got := est[4].SumaValoracion; got != 62 { // 60 + 3 * 50/60 = 62
		t.Errorf("defensa que sale en el 50: %d, se esperaba 62", got)
	}
	// El visitante no marco y no recibio: su portero tambien se queda sin imbatida.
	if got := est[101]; got.PorteriasImbatidas != 0 || got.GolesEncajados != 1 {
		t.Errorf("portero visitante: %+v", got)
	}
}

func TestLaValoracionDeUnPartidoQuedaEntre1y10(t *testing.T) {
	// Un delantero con muchos goles y un portero con muchos encajados.
	var eventos []modelo.Evento
	for m := 1; m <= 12; m++ {
		eventos = append(eventos, modelo.Evento{Minuto: m, Tipo: modelo.Gol, Local: true, Jugador: 9})
	}
	for m := 20; m <= 30; m++ {
		eventos = append(eventos, modelo.Evento{Minuto: m, Tipo: modelo.Gol, Local: false, Jugador: 111})
	}
	est, err := estadisticasDePartido(detalleDePrueba(eventos...), 12, 11)
	if err != nil {
		t.Fatal(err)
	}
	if got := est[9].SumaValoracion; got != valoracionMax {
		t.Errorf("el goleador (%d) deberia topar en %d", got, valoracionMax)
	}
	// El portero local recibio 11 goles pero gano el partido: 60 - 33 + 3 = 30, dentro del rango.
	if got := est[1].SumaValoracion; got < valoracionMin || got > valoracionMax {
		t.Errorf("valoracion fuera de rango: %d", got)
	}
	var recibe []modelo.Evento
	for m := 1; m <= 20; m++ {
		recibe = append(recibe, modelo.Evento{Minuto: m, Tipo: modelo.Gol, Local: false, Jugador: 111})
	}
	est, _ = estadisticasDePartido(detalleDePrueba(recibe...), 0, 20)
	if got := est[101+0].SumaValoracion; got < valoracionMin {
		t.Errorf("la valoracion no puede bajar de %d: %d", valoracionMin, got)
	}
	if got := est[1].SumaValoracion; got != valoracionMin {
		t.Errorf("un portero que recibe 20 goles deberia tocar el piso (%d), tiene %d", valoracionMin, got)
	}
}

func TestEstadisticasRechazaUnDetalleQueNoCuadra(t *testing.T) {
	buen := detalleDePrueba(modelo.Evento{Minuto: 10, Tipo: modelo.Gol, Local: true, Jugador: 9})
	if _, err := estadisticasDePartido(buen, 1, 0); err != nil {
		t.Fatal(err)
	}
	if _, err := estadisticasDePartido(buen, 2, 0); err == nil || !strings.Contains(err.Error(), "marcador") {
		t.Errorf("un marcador distinto de los sucesos deberia dar error: %v", err)
	}
	mal := detalleDePrueba(modelo.Evento{Minuto: 10, Tipo: modelo.Gol, Local: true, Jugador: 999})
	if _, err := estadisticasDePartido(mal, 1, 0); err == nil {
		t.Error("un goleador que no participa deberia dar error")
	}
}

func temporadaDePrueba(resultados ...Resultado) *Temporada {
	equipo := func(nombre string, base int) modelo.Equipo {
		e := modelo.Equipo{Nombre: nombre}
		for i := 1; i <= 14; i++ {
			e.Plantilla = append(e.Plantilla, modelo.Jugador{ID: base + i, Posicion: posicionDePrueba(base + i)})
		}
		return e
	}
	return &Temporada{
		Equipos:    []modelo.Equipo{equipo("Local", 0), equipo("Visitante", 100)},
		Resultados: [][]Resultado{resultados},
	}
}

func resultadoConDetalle(local, visitante int, d modelo.DetallePartido) Resultado {
	gl, gv := d.Goles()
	return Resultado{Partido: Partido{Local: local, Visitante: visitante}, GolesLocal: gl, GolesVisitante: gv, Detalle: d}
}

func TestEstadisticasDeLaTemporadaSumanLosPartidos(t *testing.T) {
	uno := resultadoConDetalle(0, 1, detalleDePrueba(
		modelo.Evento{Minuto: 10, Tipo: modelo.Gol, Local: true, Jugador: 9, Otro: 6},
		modelo.Evento{Minuto: 40, Tipo: modelo.Amarilla, Local: false, Jugador: 102},
	))
	dos := resultadoConDetalle(1, 0, detalleDePrueba( // ahora el local es "Visitante"
		modelo.Evento{Minuto: 15, Tipo: modelo.Gol, Local: false, Jugador: 9},
		modelo.Evento{Minuto: 25, Tipo: modelo.Gol, Local: true, Jugador: 109},
	))
	t0 := temporadaDePrueba(uno, dos)
	// El segundo partido usa los mismos ids para ambos lados; se rearma con ids
	// coherentes: el equipo 1 es local con ids 101.., el 0 visitante con 1...
	dos.Detalle = modelo.DetallePartido{}
	for i := 0; i < modelo.TitularesPorEquipo; i++ {
		dos.Detalle.TitularesLocal[i] = 101 + i
		dos.Detalle.TitularesVisitante[i] = 1 + i
	}
	dos.Detalle.Eventos = []modelo.Evento{
		{Minuto: 15, Tipo: modelo.Gol, Local: false, Jugador: 9},
		{Minuto: 25, Tipo: modelo.Gol, Local: true, Jugador: 109},
	}
	dos.GolesLocal, dos.GolesVisitante = dos.Detalle.Goles()
	t0.Resultados = [][]Resultado{{uno, dos}}

	est, err := t0.Estadisticas()
	if err != nil {
		t.Fatal(err)
	}
	if got := est[9]; got.Partidos != 2 || got.Goles != 2 || got.Minutos != 180 || got.Titularidades != 2 {
		t.Errorf("el jugador 9 jugo dos partidos y marco dos goles: %+v", got)
	}
	if got := est[6].Asistencias; got != 1 {
		t.Errorf("asistencias del jugador 6 = %d, se esperaba 1", got)
	}
	if got := est[102].Amarillas; got != 1 {
		t.Errorf("amarillas del jugador 102 = %d, se esperaba 1", got)
	}
	if got := est[109]; got.Goles != 1 || got.Partidos != 2 {
		t.Errorf("jugador 109: %+v", got)
	}
}

func TestLosPartidosSinDetalleNoCuentanEnLasEstadisticas(t *testing.T) {
	sinDetalle := Resultado{Partido: Partido{Local: 0, Visitante: 1}, GolesLocal: 3, GolesVisitante: 1}
	con := resultadoConDetalle(0, 1, detalleDePrueba(modelo.Evento{Minuto: 10, Tipo: modelo.Gol, Local: true, Jugador: 9}))
	t0 := temporadaDePrueba(sinDetalle, con)
	est, err := t0.Estadisticas()
	if err != nil {
		t.Fatal(err)
	}
	if got := est[9]; got.Partidos != 1 || got.Goles != 1 {
		t.Errorf("solo debe contar el partido con detalle: %+v", got)
	}
	// Pero los goles del partido sin detalle sí están en las cifras de equipo.
	eq := t0.EstadisticasEquipos()
	if eq[0].PJ != 2 || eq[0].GF != 4 || eq[0].GC != 1 {
		t.Errorf("el equipo local deberia tener PJ 2, GF 4, GC 1: %+v", eq[0])
	}
}

func TestEstadisticasDeLaTemporadaFallanConUnDetalleIncoherente(t *testing.T) {
	malo := resultadoConDetalle(0, 1, detalleDePrueba(modelo.Evento{Minuto: 10, Tipo: modelo.Gol, Local: true, Jugador: 9}))
	malo.GolesLocal = 5
	if _, err := temporadaDePrueba(malo).Estadisticas(); err == nil || !strings.Contains(err.Error(), "jornada 1 partido 1") {
		t.Errorf("el error deberia indicar el partido: %v", err)
	}
}

func TestEstadisticasDeEquipos(t *testing.T) {
	d1 := detalleDePrueba(
		modelo.Evento{Minuto: 10, Tipo: modelo.Gol, Local: true, Jugador: 9},
		modelo.Evento{Minuto: 20, Tipo: modelo.Gol, Local: true, Jugador: 10},
		modelo.Evento{Minuto: 30, Tipo: modelo.Amarilla, Local: false, Jugador: 102},
		modelo.Evento{Minuto: 40, Tipo: modelo.Amarilla, Local: true, Jugador: 3},
		modelo.Evento{Minuto: 80, Tipo: modelo.Roja, Local: false, Jugador: 103},
	) // Local 2-0 Visitante
	d2 := detalleDePrueba(modelo.Evento{Minuto: 50, Tipo: modelo.Gol, Local: false, Jugador: 109}) // 0-1
	eq := temporadaDePrueba(resultadoConDetalle(0, 1, d1), resultadoConDetalle(0, 1, d2)).EstadisticasEquipos()
	want := []EstadisticaEquipo{
		{Equipo: "Local", PJ: 2, G: 1, E: 0, P: 1, GF: 2, GC: 1, Imbatidas: 1, SinMarcar: 1, Amarillas: 1},
		{Equipo: "Visitante", PJ: 2, G: 1, E: 0, P: 1, GF: 1, GC: 2, Imbatidas: 1, SinMarcar: 1, Amarillas: 1, Rojas: 1},
	}
	if !reflect.DeepEqual(eq, want) {
		t.Errorf("EstadisticasEquipos:\n  obtenido %+v\n  esperado %+v", eq, want)
	}
}

// ligaConIDsUnicas arma una liga de n equipos uniformes con ids de jugador
// distintas en toda la liga.
func ligaConIDsUnicas(n, v int) []modelo.Equipo {
	equipos := equiposUniformes(n, v)
	for i := range equipos {
		for k := range equipos[i].Plantilla {
			equipos[i].Plantilla[k].ID = (i+1)*100 + k + 1
		}
	}
	return equipos
}

func TestEstadisticasDeUnaTemporadaSimuladaSonCoherentes(t *testing.T) {
	for _, n := range []int{4, 5, 10} {
		temp, err := Nueva(ligaConIDsUnicas(n, 62))
		if err != nil {
			t.Fatal(err)
		}
		if err := temp.JugarTemporada(nuevoRand(int64(n))); err != nil {
			t.Fatal(err)
		}
		est, err := temp.Estadisticas()
		if err != nil {
			t.Fatalf("n=%d: %v", n, err)
		}
		equipos := temp.EstadisticasEquipos()

		partidos, goles := 0, 0
		minutosEsperados := 0
		var amarillasEquipos, rojasEquipos int
		for _, jornada := range temp.Resultados {
			for _, res := range jornada {
				partidos++
				goles += res.GolesLocal + res.GolesVisitante
				minutosEsperados += 2 * modelo.TitularesPorEquipo * modelo.MinutosPartido
				for _, ev := range res.Detalle.Eventos {
					if ev.Tipo == modelo.Roja {
						minutosEsperados -= modelo.MinutosPartido - ev.Minuto
					}
				}
			}
		}
		var sumaGoles, sumaAsist, sumaMinutos, sumaAmarillas, sumaRojas, sumaPartidos int
		for id, e := range est {
			sumaGoles += e.Goles
			sumaAsist += e.Asistencias
			sumaMinutos += e.Minutos
			sumaAmarillas += e.Amarillas
			sumaRojas += e.Rojas
			sumaPartidos += e.Partidos
			if e.Partidos == 0 || e.Minutos <= 0 || e.Minutos > e.Partidos*modelo.MinutosPartido {
				t.Fatalf("n=%d jugador %d: %d partidos y %d minutos", n, id, e.Partidos, e.Minutos)
			}
			if e.Titularidades > e.Partidos {
				t.Fatalf("n=%d jugador %d: %d titularidades en %d partidos", n, id, e.Titularidades, e.Partidos)
			}
			if media := e.ValoracionMedia(); media < 1 || media > 10 {
				t.Fatalf("n=%d jugador %d: valoracion media %.2f", n, id, media)
			}
			if e.Rojas > e.Partidos {
				t.Fatalf("n=%d jugador %d: mas rojas que partidos", n, id)
			}
		}
		for _, e := range equipos {
			amarillasEquipos += e.Amarillas
			rojasEquipos += e.Rojas
		}
		if sumaGoles != goles {
			t.Errorf("n=%d: los jugadores suman %d goles y los marcadores %d", n, sumaGoles, goles)
		}
		if sumaAsist > goles {
			t.Errorf("n=%d: %d asistencias para %d goles", n, sumaAsist, goles)
		}
		if sumaMinutos != minutosEsperados {
			t.Errorf("n=%d: %d minutos jugados, se esperaban %d", n, sumaMinutos, minutosEsperados)
		}
		if sumaAmarillas != amarillasEquipos || sumaRojas != rojasEquipos {
			t.Errorf("n=%d: tarjetas de jugadores (%d, %d) no coinciden con las de equipos (%d, %d)",
				n, sumaAmarillas, sumaRojas, amarillasEquipos, rojasEquipos)
		}
		// Cada partido tiene al menos 22 participantes.
		if sumaPartidos < 22*partidos {
			t.Errorf("n=%d: %d participaciones para %d partidos", n, sumaPartidos, partidos)
		}
		// Goles encajados: los de los porteros titulares suman los goles en contra
		// de sus equipos mientras estuvieron en el campo (todo el partido).
		encajadosPorteros, golesEnContra := 0, 0
		for _, e := range temp.Equipos {
			for _, j := range e.Plantilla {
				if j.Posicion == modelo.Portero {
					encajadosPorteros += est[j.ID].GolesEncajados
				}
			}
		}
		for _, e := range equipos {
			golesEnContra += e.GC
		}
		if encajadosPorteros != golesEnContra {
			t.Errorf("n=%d: los porteros encajaron %d goles y los equipos recibieron %d", n, encajadosPorteros, golesEnContra)
		}
	}
}
