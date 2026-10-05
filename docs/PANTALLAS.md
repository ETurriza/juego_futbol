# Galería de pantallas

Cómo se ve el juego en una terminal de 80 por 30, con la semilla 21 (cuatro
temporadas de carrera). Este archivo se **genera** a partir de las pruebas de
`internal/menus`; no se edita a mano. Para actualizarlo tras un cambio visual:

```
go test ./internal/menus -update
```

En las filas con cursor, `>` marca la seleccionada y `*` las de tu club.

## Menú principal

```
JUEGO DE FÚTBOL · Modo carrera

Academia Fontalva    Temporada 4 · Jornada 11 / 18
Valoración del equipo: 81

> Avanzar jornada
  Tabla de posiciones
  Plantilla
  Estadísticas
  Historial
  Salir

↑/↓ mover · enter elegir · q salir
```

## Resultados de la jornada

```
JORNADA 12 / 18 · Temporada 4

  Juventud Granvela           0 - 0   Club Brisanda
  Club Altamira del Sur       1 - 2   Unión Pinarel
> Academia Fontalva           1 - 4   Academia Valmora
  Club Maldoria               0 - 2   Deportivo Valmora
  Unión Granvela              2 - 0   Unión Castelmar

enter continuar
```

## Tabla de posiciones

```
TABLA · Temporada 4 · Jornada 11 / 18

    #  Equipo                      PJ   G   E   P   GF   GC   DG  Pts
    1  Unión Castelmar             11   8   2   1   21    7  +14   26
    2  Academia Valmora            11   8   0   3   18    7  +11   24
    3  Unión Pinarel               11   6   2   3   12    7   +5   20
    4  Unión Granvela              11   5   2   4   17   15   +2   17
    5  Deportivo Valmora           11   5   1   5   16   13   +3   16
    6  Juventud Granvela           11   3   3   5   11   19   -8   12
>   7  Academia Fontalva           11   3   2   6    8   14   -6   11
    8  Club Maldoria               11   3   2   6    5   11   -6   11
    9  Club Brisanda               11   3   2   6   11   19   -8   11
   10  Club Altamira del Sur       11   2   2   7    8   15   -7    8

esc volver
```

## Plantilla: atributos

```
PLANTILLA · Academia Fontalva

  Posición       Nombre                  Ed  Val  RIT TIR PAS REG DEF FIS REF
  Portero        Orlandi Pedrosel        24   75   63  48  68  53  70  74  92
  Portero        Marvel Illanes          32   72   59  46  82  52  71  64  83
  Defensa        Valdor Herbasco         33   87   86  71  84  68  94  86   1
  Defensa        Pelayo Jarvale          30   80   65  52  64  69  90  86   1
  Defensa        Jaldin Bocanegro        29   79   69  69  69  67  84  84   3
  Defensa        Selmo Escalona          24   74   64  44  67  55  85  74   1
  Defensa        Galdo Rivasel           28   74   68  58  66  60  76  85   1
  Defensa        Lorvin Duelmo           19   52   41  25  37  30  64  57   1
  Defensa        Eskardo Pedrosel        16   49   39  19  31  22  65  46   1
  Mediocampista  Wenceo Zaldumar         36   82   63  74  95  92  66  67   1
  Mediocampista  Ivaro Herbasco          29   80   75  85  81  78  83  75   1
  Mediocampista  Lisandor Orbaneja       25   78   75  72  87  82  64  66   1
  Mediocampista  Nervio Torvelo          26   76   56  79  81  74  79  73   1
  Mediocampista  Nolven Jarvale          32   61   64  48  68  70  55  44   1
  Mediocampista  Hilaro Fuentalba        17   55   45  51  63  58  42  53   1
  Mediocampista  Nervio Zaldumar         16   51   40  47  61  50  45  41   1
  Delantero      Nolven Montecal         31   90   91  98  83  90  79  77   1
  Delantero      Joldan Valmeron         25   83   83  88  71  85  62  80   1
  Delantero      Ivaro Jarvale           28   77   75  82  76  85  42  65   1
  Delantero      Tavio Norvalde          36   76   65  86  68  86  64  57   1
  Delantero      Lorvin Pedrosel         32   73   81  73  68  74  66  65   1
  Delantero      Marvel Valmeron         21   73   78  82  61  75  38  58   1

tab estadísticas · esc volver
```

## Plantilla: estadísticas (tab)

```
PLANTILLA · Academia Fontalva · estadísticas

  Posición      Jugador                 PJ Tit   Min   G   A  Am  Ro Imb  Val
  Portero       Orlandi Pedrosel        11  11   990   0   0   0   0   2  5.6
  Portero       Marvel Illanes           0   0     0   0   0   0   0   0    -
  Defensa       Valdor Herbasco         11  11   850   0   0   2   0   2  5.8
  Defensa       Pelayo Jarvale          11  11   955   0   1   2   0   2  5.8
  Defensa       Jaldin Bocanegro        11  11   926   0   0   5   0   2  5.7
  Defensa       Selmo Escalona          11  11   845   0   1   0   0   2  5.9
  Defensa       Galdo Rivasel            6   0   182   0   0   0   0   0  5.9
  Defensa       Lorvin Duelmo            4   0    67   0   0   0   0   0  6.0
  Defensa       Eskardo Pedrosel         6   0   147   0   0   0   0   0  5.9
  Mediocampista Wenceo Zaldumar         11  11   895   0   1   2   0   0  5.9
  Mediocampista Ivaro Herbasco          11  11   827   2   1   2   0   0  6.1
  Mediocampista Lisandor Orbaneja       11  11   725   0   0   1   0   0  5.9
  Mediocampista Nervio Torvelo           8   0   178   0   0   1   0   0  6.0
  Mediocampista Nolven Jarvale           5   0   113   0   0   1   0   0  5.9
  Mediocampista Hilaro Fuentalba         5   0   145   0   0   0   0   0  6.0
  Mediocampista Nervio Zaldumar          4   0   109   0   0   0   0   0  6.0
  Delantero     Nolven Montecal         11  11   940   1   3   1   0   0  6.1
  Delantero     Joldan Valmeron         11  11   889   2   0   1   0   0  6.1
  Delantero     Ivaro Jarvale           11  11   758   1   0   2   0   0  6.0
  Delantero     Tavio Norvalde           6   0   103   0   0   0   0   0  6.0
  Delantero     Lorvin Pedrosel          7   0   195   0   0   0   0   0  6.0
  Delantero     Marvel Valmeron          2   0    51   2   0   0   0   0  6.8

tab atributos · esc volver
```

## Historial de temporadas

```
HISTORIAL · Academia Fontalva

  Temp  Campeón                     Tu puesto   Pts
     1  Club Maldoria                 7 de 10    22
     2  Unión Pinarel                 4 de 10    29
     3  Unión Castelmar               7 de 10    22

esc volver
```

## Menú de estadísticas

```
ESTADÍSTICAS · Temporada 4 · Jornada 11 / 18

> Goleadores
  Asistentes
  Tarjetas
  Porteros: porterías imbatidas
  Porteros: menos goles encajados
  Mejores valoraciones
  Equipos
  Volver

↑/↓ mover · enter elegir · esc volver
```

## Clasificación: GOLEADORES

```
GOLEADORES · Temporada 4 · Jornada 11 / 18

   #  Jugador                Equipo                     PJ Goles Asist   Min
>  1  Galdo Duelmo           Unión Castelmar            11     7     3   954
   2  Pelayo Pardovan        Unión Granvela             11     5     0   889
   3  Galdo Lacerna          Unión Castelmar            11     4     3   795
   4  Xandro Abrelda         Club Brisanda              11     4     1   801
   5  Elmiro Duelmo          Deportivo Valmora          11     4     1   955
   6  Kaelo Herbasco         Academia Valmora           11     4     0   829
   7  Baltor Orbaneja        Juventud Granvela          11     3     2   904
   8  Dorvan Jurado          Unión Granvela             11     3     2   921
   9  Hemiro Frondosa        Academia Valmora           11     3     2   974
  10  Lorvin Orbaneja        Deportivo Valmora          11     3     1   806
  11  Nolven Abrelda         Unión Pinarel              11     3     1   825
  12  Joldan Norvalde        Unión Castelmar            11     3     1   897
  13  Kervo Maldovar         Juventud Granvela          11     3     0   795
  14  Lorvin Bastaro         Unión Granvela             11     3     0   820
  15  Kervo Montecal         Club Altamira del Sur      11     3     0   856
  16  Marvel Zaldumar        Academia Valmora           11     2     3   829
  17  Nervio Bastaro         Deportivo Valmora          11     2     2   835
  18  Isidoro Frondosa       Academia Valmora           11     2     2   931
  19  Daveo Garmendo         Club Brisanda              11     2     1   795
* 20  Ivaro Herbasco         Academia Fontalva          11     2     1   827
  21  Corvino Zaldumar       Unión Castelmar            11     2     1   909
* 22  Marvel Valmeron        Academia Fontalva           2     2     0    51
  23  Kervo Frondosa         Unión Granvela              6     2     0   161
  24  Belron Cordival        Club Altamira del Sur       7     2     0   206
  25  Isidoro Dunavar        Unión Granvela             11     2     0   801

↑/↓ mover (1-25 de 50) · enter ver ficha · * tu club · esc volver
```

## Clasificación: ASISTENTES

```
ASISTENTES · Temporada 4 · Jornada 11 / 18

   #  Jugador                Equipo                     PJ Asist Goles   Min
>  1  Quirino Zaldumar       Deportivo Valmora          11     4     1   947
   2  Nervio Jurado          Unión Castelmar            11     4     1   967
   3  Galdo Duelmo           Unión Castelmar            11     3     7   954
   4  Galdo Lacerna          Unión Castelmar            11     3     4   795
   5  Marvel Zaldumar        Academia Valmora           11     3     2   829
   6  Kaelo Yrigaldo         Unión Granvela             11     3     1   815
   7  Tavio Korvena          Club Brisanda              11     3     1   891
*  8  Nolven Montecal        Academia Fontalva          11     3     1   940
   9  Baltor Orbaneja        Juventud Granvela          11     2     3   904
  10  Dorvan Jurado          Unión Granvela             11     2     3   921
  11  Hemiro Frondosa        Academia Valmora           11     2     3   974
  12  Nervio Bastaro         Deportivo Valmora          11     2     2   835
  13  Isidoro Frondosa       Academia Valmora           11     2     2   931
  14  Xandro Zaldumar        Club Altamira del Sur      11     2     1   767
  15  Hilaro Roblemar        Academia Valmora           11     2     1   803
  16  Dorvan Aguirel         Unión Granvela             11     2     1   818
  17  Hemiro Norvalde        Unión Pinarel              11     2     1   818
  18  Elmiro Jarvale         Club Brisanda              11     2     1   841
  19  Corvino Landrosa       Deportivo Valmora          11     2     1   874
  20  Rodvan Maldovar        Academia Valmora           11     2     1   883
  21  Selmo Fuentalba        Club Altamira del Sur      11     2     1   902
  22  Corvino Maldovar       Club Brisanda              11     2     1   949
  23  Zenel Herbasco         Unión Pinarel              11     2     0   827
  24  Ivaro Bocanegro        Unión Granvela             11     2     0   879
  25  Isidoro Irueta         Unión Pinarel              11     2     0   928

↑/↓ mover (1-25 de 50) · enter ver ficha · * tu club · esc volver
```

## Clasificación: TARJETAS

```
TARJETAS · Temporada 4 · Jornada 11 / 18

   #  Jugador                Equipo                     PJ  Amar  Roja   Pts
>  1  Isidoro Pardovan       Club Maldoria              11     5     1     8
   2  Daveo Estoral          Juventud Granvela          11     4     1     7
   3  Kaelo Bocanegro        Club Altamira del Sur      11     4     1     7
   4  Rodvan Illanes         Juventud Granvela          11     3     1     6
   5  Lorvin Bastaro         Unión Granvela             11     2     1     5
   6  Quirino Zaldumar       Deportivo Valmora          11     2     1     5
   7  Isidoro Dunavar        Unión Granvela             11     5     0     5
   8  Ivaro Bocanegro        Unión Granvela             11     5     0     5
   9  Selmo Fuentalba        Club Altamira del Sur      11     5     0     5
* 10  Jaldin Bocanegro       Academia Fontalva          11     5     0     5
  11  Lorvin Orbaneja        Deportivo Valmora          11     1     1     4
  12  Nolven Orbaneja        Club Altamira del Sur      11     4     0     4
  13  Tavio Korvena          Club Brisanda              11     4     0     4
  14  Zenel Hontoria         Club Maldoria              11     4     0     4
  15  Isidoro Norvalde       Unión Granvela             11     4     0     4
  16  Isidoro Sandoral       Unión Castelmar            11     0     1     3
  17  Xandro Abrelda         Club Brisanda              11     0     1     3
  18  Jaldin Cordival        Club Maldoria              11     3     0     3
  19  Xandro Zaldumar        Club Altamira del Sur      11     3     0     3
  20  Kaelo Yrigaldo         Unión Granvela             11     3     0     3
  21  Hemiro Norvalde        Unión Pinarel              11     3     0     3
  22  Dorvan Sandoral        Academia Valmora           11     3     0     3
  23  Daveo Norvalde         Club Maldoria              11     3     0     3
  24  Elmiro Torvelo         Club Maldoria              11     3     0     3
  25  Nervio Garmendo        Deportivo Valmora          11     3     0     3

↑/↓ mover (1-25 de 50) · enter ver ficha · * tu club · esc volver
```

## Clasificación: PORTEROS · PORTERÍAS IMBATIDAS

```
PORTEROS · PORTERÍAS IMBATIDAS · Temporada 4 · Jornada 11 / 18

   #  Jugador                Equipo                     PJ   Imb   Enc
>  1  Lisandor Landrosa      Unión Castelmar            11     6     7
   2  Quirino Pardovan       Academia Valmora           11     5     7
   3  Xandro Abrelda         Unión Pinarel              11     5     7
   4  Selmo Frondosa         Club Maldoria              11     4    11
   5  Lisandor Escalona      Unión Granvela             11     4    15
   6  Lorvin Illanes         Deportivo Valmora          11     3    13
   7  Fabrel Jurado          Club Altamira del Sur      11     3    15
*  8  Orlandi Pedrosel       Academia Fontalva          11     2    14
   9  Pelayo Cordival        Juventud Granvela          11     2    19
  10  Anselmo Urdangal       Club Brisanda              11     1    19

↑/↓ mover · enter ver ficha · * tu club · esc volver
```

## Clasificación: PORTEROS · MENOS GOLES ENCAJADOS

```
PORTEROS · MENOS GOLES ENCAJADOS · Temporada 4 · Jornada 11 / 18

   #  Jugador                Equipo                     PJ   Enc Enc/PJ
>  1  Lisandor Landrosa      Unión Castelmar            11     7   0.64
   2  Quirino Pardovan       Academia Valmora           11     7   0.64
   3  Xandro Abrelda         Unión Pinarel              11     7   0.64
   4  Selmo Frondosa         Club Maldoria              11    11   1.00
   5  Lorvin Illanes         Deportivo Valmora          11    13   1.18
*  6  Orlandi Pedrosel       Academia Fontalva          11    14   1.27
   7  Fabrel Jurado          Club Altamira del Sur      11    15   1.36
   8  Lisandor Escalona      Unión Granvela             11    15   1.36
   9  Anselmo Urdangal       Club Brisanda              11    19   1.73
  10  Pelayo Cordival        Juventud Granvela          11    19   1.73

↑/↓ mover · enter ver ficha · * tu club · esc volver
```

## Clasificación: MEJORES VALORACIONES

```
MEJORES VALORACIONES · Temporada 4 · Jornada 11 / 18

   #  Jugador                Equipo                     PJ   Val Goles Asist
>  1  Galdo Duelmo           Unión Castelmar            11   7.0     7     3
   2  Galdo Lacerna          Unión Castelmar            11   6.7     4     3
   3  Nervio Jurado          Unión Castelmar            11   6.6     1     4
   4  Hemiro Aguirel         Unión Castelmar            11   6.5     2     0
   5  Hemiro Frondosa        Academia Valmora           11   6.5     3     2
   6  Pelayo Pardovan        Unión Granvela             11   6.5     5     0
   7  Kaelo Herbasco         Academia Valmora           11   6.5     4     0
   8  Joldan Norvalde        Unión Castelmar            11   6.4     3     1
   9  Marvel Zaldumar        Academia Valmora           11   6.4     2     3
  10  Jaldin Roblemar        Unión Castelmar            11   6.4     1     0
  11  Rodvan Maldovar        Academia Valmora           11   6.4     1     2
  12  Hilaro Roblemar        Academia Valmora           11   6.4     1     2
  13  Dorvan Jurado          Unión Granvela             11   6.4     3     2
  14  Valdor Cordival        Academia Valmora           11   6.4     1     1
  15  Isidoro Frondosa       Academia Valmora           11   6.4     2     2
  16  Corvino Zaldumar       Unión Castelmar            11   6.4     2     1
  17  Elmiro Duelmo          Deportivo Valmora          11   6.4     4     1
  18  Nolven Abrelda         Unión Pinarel              11   6.4     3     1
  19  Baltor Orbaneja        Juventud Granvela          11   6.3     3     2
  20  Ivaro Pardovan         Unión Castelmar            11   6.3     1     1
  21  Arlen Quintaleo        Unión Castelmar            11   6.3     0     0
  22  Florian Urdangal       Academia Valmora           11   6.3     1     1
  23  Lisandor Landrosa      Unión Castelmar            11   6.2     0     0
  24  Lorvin Orbaneja        Deportivo Valmora          11   6.2     3     1
  25  Xandro Abrelda         Club Brisanda              11   6.2     4     1

↑/↓ mover (1-25 de 50) · enter ver ficha · * tu club · esc volver
```

## Comparativa de equipos

```
EQUIPOS · Temporada 4 · Jornada 11 / 18

   #  Equipo                    PJ  GF  GC Imb SM  Am Ro  Goleador
>  1  Unión Castelmar           11  21   7   6  0  13  1  Galdo Duelmo (7)
   2  Academia Valmora          11  18   7   5  3  16  0  Kaelo Herbasco (4)
   3  Unión Pinarel             11  12   7   5  3  16  0  Nolven Abrelda (3)
   4  Unión Granvela            11  17  15   4  1  28  1  Pelayo Pardovan (5)
   5  Deportivo Valmora         11  16  13   3  4  21  2  Elmiro Duelmo (4)
   6  Juventud Granvela         11  11  19   2  5  23  2  Baltor Orbaneja (3)
*  7  Academia Fontalva         11   8  14   2  5  20  0  Ivaro Herbasco (2)
   8  Club Maldoria             11   5  11   4  6  23  1  Lorvin Abrelda (2)
   9  Club Brisanda             11  11  19   1  4  18  1  Xandro Abrelda (4)
  10  Club Altamira del Sur     11   8  15   3  4  26  1  Kervo Montecal (3)

↑/↓ mover · enter plantilla · tab goleador/asistente · * tu club · esc volver
```

## Comparativa de equipos (tab)

```
EQUIPOS · Temporada 4 · Jornada 11 / 18

   #  Equipo                    PJ  GF  GC Imb SM  Am Ro  Asistente
>  1  Unión Castelmar           11  21   7   6  0  13  1  Nervio Jurado (4)
   2  Academia Valmora          11  18   7   5  3  16  0  Marvel Zaldumar (3)
   3  Unión Pinarel             11  12   7   5  3  16  0  Hemiro Norvalde (2)
   4  Unión Granvela            11  17  15   4  1  28  1  Kaelo Yrigaldo (3)
   5  Deportivo Valmora         11  16  13   3  4  21  2  Quirino Zaldumar (4)
   6  Juventud Granvela         11  11  19   2  5  23  2  Baltor Orbaneja (2)
*  7  Academia Fontalva         11   8  14   2  5  20  0  Nolven Montecal (3)
   8  Club Maldoria             11   5  11   4  6  23  1  Baltor Irueta (1)
   9  Club Brisanda             11  11  19   1  4  18  1  Tavio Korvena (3)
  10  Club Altamira del Sur     11   8  15   3  4  26  1  Selmo Fuentalba (2)

↑/↓ mover · enter plantilla · tab goleador/asistente · * tu club · esc volver
```

## Plantilla de un equipo con estadísticas

```
ACADEMIA VALMORA · estadísticas

  Posición      Jugador                 PJ Tit   Min   G   A  Am  Ro Imb  Val
> Portero       Quirino Pardovan        11  11   990   0   0   0   0   5  6.2
  Portero       Arlen Sandoral           0   0     0   0   0   0   0   0    -
  Defensa       Valdor Cordival         11  11   957   1   1   0   0   5  6.4
  Defensa       Valdor Dunavar          11  11   849   0   0   1   0   5  6.2
  Defensa       Hilaro Roblemar         11  11   803   1   2   2   0   4  6.4
  Defensa       Rodvan Maldovar         11  11   883   1   2   1   0   5  6.4
  Defensa       Selmo Pedrosel           8   0   180   0   0   0   0   0  6.1
  Defensa       Zenel Irueta             6   0   175   0   0   2   0   0  5.9
  Defensa       Elmiro Cordival          4   0    99   0   0   0   0   0  6.0
  Mediocampista Isidoro Frondosa        11  11   931   2   2   2   0   0  6.4
  Mediocampista Dorvan Sandoral         11  11   819   1   1   3   0   0  6.2
  Mediocampista Marvel Zaldumar         11  11   829   2   3   1   0   0  6.4
  Mediocampista Quirino Rivasel          2   0    54   0   0   0   0   0  6.1
  Mediocampista Pelayo Torvelo           6   0    64   0   0   0   0   0  6.0
  Mediocampista Baltor Frondosa          3   0   114   0   0   1   0   0  6.1
  Mediocampista Galdo Jarvale            7   0   191   0   0   1   0   0  6.0
  Delantero     Florian Urdangal        11  11   802   1   1   0   0   0  6.3
  Delantero     Hemiro Frondosa         11  11   974   3   2   1   0   0  6.5
  Delantero     Kaelo Herbasco          11  11   829   4   0   1   0   0  6.5
  Delantero     Xandro Torvelo           3   0   100   0   0   0   0   0  6.1
  Delantero     Elmiro Hontoria          5   0   146   1   0   0   0   0  6.2
  Delantero     Marvel Bocanegro         3   0   101   1   0   0   0   0  6.2

↑/↓ mover · enter ver ficha · esc volver
```

## Ficha de un portero

```
FICHA · Orlandi Pedrosel
Academia Fontalva · Portero · 24 años · Valoración 75
RIT 63  TIR 48  PAS 68  REG 53  DEF 70  FIS 74  REF 92
Proyección: normal

Temporada 4 (en curso)
  PJ 11 · Tit 11 · Min 990 · G 0 · A 0 · Am 0 · Ro 0 · Imb 2 · Enc 14 · Val 5.6

Trayectoria
  Todavía no ha terminado ninguna temporada en esta carrera.


esc volver
```

## Ficha de un delantero

```
FICHA · Marvel Valmeron
Academia Fontalva · Delantero · 21 años · Valoración 73
RIT 78  TIR 82  PAS 61  REG 75  DEF 38  FIS 58  REF 1
Proyección: normal

Temporada 4 (en curso)
  PJ 2 · Tit 0 · Min 51 · G 2 · A 0 · Am 0 · Ro 0 · Imb 0 · Val 6.8

Trayectoria
  Temp  Club                       PJ   Min   G   A  Am  Ro Imb  Val
     2  Academia Fontalva           9   232   0   0   0   0   0  6.0
     3  Academia Fontalva           9   170   0   0   0   0   0  6.0

Carrera: 20 partidos · 2 goles · 0 asistencias · 0 amarillas · 0 rojas

esc volver
```

## Fin de temporada

```
TEMPORADA 1 TERMINADA

Campeón: Deportivo Granvela
Tu equipo: Deportivo Granvela, puesto 1 de 10, 38 puntos
¡Felicidades, eres el campeón!

> Siguiente temporada
  Ver tabla final
  Estadísticas
  Historial
  Nueva carrera
  Salir

↑/↓ mover · enter elegir · q salir
```

## Inicio de temporada

```
TEMPORADA 2 · Deportivo Granvela

Valoración del equipo: 80 → 81

Se retiran (1)
  Kaelo Dunavar          Defensa        36 años   Val 65

Llegan de la cantera (1)
  Gervasio Orbaneja      Defensa        18 años   Val 56

En toda la liga: 7 retiros y 7 juveniles.

enter continuar
```

## Confirmar una carrera nueva

```
NUEVA CARRERA

Esto empieza una carrera nueva y pierdes la actual
(Temporada 2 con Deportivo Granvela). Aún no hay guardado automático.
¿Seguro?

> No, volver
  Sí, empezar de cero

↑/↓ mover · enter elegir · esc volver
```
