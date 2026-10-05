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

  Juventud Granvela           3 - 1   Club Brisanda
  Club Altamira del Sur       0 - 0   Unión Pinarel
> Academia Fontalva           2 - 1   Academia Valmora
  Club Maldoria               0 - 2   Deportivo Valmora
  Unión Granvela              0 - 0   Unión Castelmar

enter continuar
```

## Tabla de posiciones

```
TABLA · Temporada 4 · Jornada 11 / 18

    #  Equipo                      PJ   G   E   P   GF   GC   DG  Pts
    1  Unión Castelmar             11   7   3   1   18    8  +10   24
    2  Unión Pinarel               11   6   2   3   15   11   +4   20
    3  Academia Valmora            11   6   1   4   18   11   +7   19
>   4  Academia Fontalva           11   6   1   4   18   15   +3   19
    5  Club Altamira del Sur       11   5   1   5   15   20   -5   16
    6  Club Brisanda               11   4   2   5   12   15   -3   14
    7  Club Maldoria               11   4   1   6   13   18   -5   13
    8  Unión Granvela              11   3   3   5   16   22   -6   12
    9  Deportivo Valmora           11   3   1   7   13   15   -2   10
   10  Juventud Granvela           11   2   3   6   11   14   -3    9

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
  Portero       Orlandi Pedrosel        11  11   990   0   0   1   0   5  5.8
  Portero       Marvel Illanes           0   0     0   0   0   0   0   0    -
  Defensa       Valdor Herbasco         11  11   945   0   0   2   0   5  6.0
  Defensa       Pelayo Jarvale          11  11   905   0   0   3   0   5  6.0
  Defensa       Jaldin Bocanegro        11  11   890   1   0   2   1   5  5.9
  Defensa       Selmo Escalona          11  11   871   1   0   6   0   4  5.9
  Defensa       Galdo Rivasel            9   0   202   0   0   0   0   0  6.0
  Defensa       Lorvin Duelmo            4   0    55   0   0   0   0   0  6.0
  Defensa       Eskardo Pedrosel         0   0     0   0   0   0   0   0    -
  Mediocampista Wenceo Zaldumar         11  11   812   2   0   1   0   0  6.2
  Mediocampista Ivaro Herbasco          11  11   860   1   1   6   1   0  5.9
  Mediocampista Lisandor Orbaneja       11  11   851   2   2   3   0   0  6.3
  Mediocampista Nervio Torvelo          11  11   920   2   2   2   0   0  6.3
  Mediocampista Nolven Jarvale          10   0   290   1   1   1   0   0  6.1
  Mediocampista Hilaro Fuentalba         5   0   113   0   1   0   0   0  6.1
  Mediocampista Nervio Zaldumar          0   0     0   0   0   0   0   0    -
  Delantero     Nolven Montecal         11  11   821   5   0   1   0   0  6.5
  Delantero     Joldan Valmeron         11  11   816   1   1   1   0   0  6.2
  Delantero     Ivaro Jarvale           10   0   359   1   2   3   0   0  6.1
  Delantero     Tavio Norvalde           7   0   127   1   0   1   0   0  6.0
  Delantero     Lorvin Pedrosel          1   0    13   0   0   0   0   0  6.0
  Delantero     Marvel Valmeron          0   0     0   0   0   0   0   0    -

tab atributos · esc volver
```

## Historial de temporadas

```
HISTORIAL · Academia Fontalva

  Temp  Campeón                     Tu puesto   Pts
     1  Club Altamira del Sur         5 de 10    25
     2  Unión Pinarel                 3 de 10    31
     3  Club Brisanda                 3 de 10    27

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
>  1  Elmiro Duelmo          Deportivo Valmora          11     6     1   928
   2  Orlandi Cordival       Unión Pinarel              11     5     1   868
   3  Dorvan Jurado          Unión Granvela             11     5     0   820
*  4  Nolven Montecal        Academia Fontalva          11     5     0   821
   5  Hemiro Norvalde        Unión Pinarel              11     5     0   849
   6  Lisandor Bastaro       Club Altamira del Sur      11     4     1   812
   7  Isidoro Frondosa       Academia Valmora           11     4     1   932
   8  Corvino Illanes        Club Maldoria              11     4     0   929
   9  Joldan Norvalde        Unión Castelmar            11     3     3   923
  10  Baltor Frondosa        Unión Castelmar            11     3     3   948
  11  Ivaro Bocanegro        Unión Granvela             11     3     2   821
  12  Isidoro Dunavar        Unión Granvela             11     3     2   832
  13  Isidoro Sandoral       Unión Castelmar            11     3     0   735
  14  Marvel Zaldumar        Academia Valmora           11     3     0   856
  15  Pelayo Torvelo         Academia Valmora           11     3     0   923
  16  Baltor Orbaneja        Juventud Granvela          11     2     2   813
* 17  Lisandor Orbaneja      Academia Fontalva          11     2     2   851
* 18  Nervio Torvelo         Academia Fontalva          11     2     2   920
  19  Xandro Zaldumar        Club Altamira del Sur      11     2     2   944
  20  Kervo Montecal         Club Altamira del Sur      11     2     1   781
  21  Galdo Duelmo           Unión Castelmar            11     2     1   817
  22  Nolven Orbaneja        Club Altamira del Sur      11     2     1   851
  23  Tavio Korvena          Club Brisanda              11     2     1   858
  24  Corvino Landrosa       Deportivo Valmora          11     2     1   867
  25  Dorvan Sandoral        Academia Valmora           11     2     1   903

↑/↓ mover (1-25 de 50) · enter ver ficha · * tu club · esc volver
```

## Clasificación: ASISTENTES

```
ASISTENTES · Temporada 4 · Jornada 11 / 18

   #  Jugador                Equipo                     PJ Asist Goles   Min
>  1  Florian Quintaleo      Deportivo Valmora          11     4     0   971
   2  Joldan Norvalde        Unión Castelmar            11     3     3   923
   3  Baltor Frondosa        Unión Castelmar            11     3     3   948
   4  Mendo Olmedrar         Juventud Granvela          11     3     1   866
   5  Nolven Abrelda         Unión Pinarel              11     3     1   883
   6  Valdor Cordival        Academia Valmora           11     3     1   923
   7  Pelayo Pardovan        Unión Granvela             11     3     0   907
   8  Ivaro Bocanegro        Unión Granvela             11     2     3   821
   9  Isidoro Dunavar        Unión Granvela             11     2     3   832
  10  Baltor Orbaneja        Juventud Granvela          11     2     2   813
* 11  Lisandor Orbaneja      Academia Fontalva          11     2     2   851
* 12  Nervio Torvelo         Academia Fontalva          11     2     2   920
  13  Xandro Zaldumar        Club Altamira del Sur      11     2     2   944
* 14  Ivaro Jarvale          Academia Fontalva          10     2     1   359
  15  Arlen Cordival         Unión Pinarel              11     2     1   764
  16  Ivaro Pardovan         Unión Castelmar            11     2     1   767
  17  Fabrel Olmedrar        Club Brisanda              11     2     1   890
  18  Elmiro Torvelo         Club Maldoria              11     2     1   901
  19  Lorvin Pardovan        Club Maldoria              11     2     0   316
  20  Dorvan Aguirel         Unión Granvela             11     2     0   871
  21  Hemiro Aguirel         Unión Castelmar            11     2     0   933
  22  Elmiro Duelmo          Deportivo Valmora          11     1     6   928
  23  Orlandi Cordival       Unión Pinarel              11     1     5   868
  24  Lisandor Bastaro       Club Altamira del Sur      11     1     4   812
  25  Isidoro Frondosa       Academia Valmora           11     1     4   932

↑/↓ mover (1-25 de 50) · enter ver ficha · * tu club · esc volver
```

## Clasificación: TARJETAS

```
TARJETAS · Temporada 4 · Jornada 11 / 18

   #  Jugador                Equipo                     PJ  Amar  Roja   Pts
>  1  Ivaro Herbasco         Academia Fontalva          11     6     1     9
   2  Valdor Cordival        Academia Valmora           11     3     1     6
*  3  Selmo Escalona         Academia Fontalva          11     6     0     6
   4  Daveo Gallardel        Unión Pinarel              11     6     0     6
   5  Ivaro Bocanegro        Unión Granvela             11     2     1     5
*  6  Jaldin Bocanegro       Academia Fontalva          11     2     1     5
   7  Rodvan Illanes         Juventud Granvela          11     5     0     5
   8  Corvino Garmendo       Club Altamira del Sur      11     4     0     4
   9  Lorvin Bastaro         Unión Granvela             11     4     0     4
  10  Nervio Jurado          Unión Castelmar            11     4     0     4
  11  Tavio Pedrosel         Juventud Granvela          11     4     0     4
  12  Quirino Zaldumar       Deportivo Valmora          11     4     0     4
  13  Isidoro Frondosa       Academia Valmora           11     4     0     4
  14  Hemiro Aguirel         Unión Castelmar            11     4     0     4
  15  Jaldin Roblemar        Unión Castelmar            11     4     0     4
  16  Kervo Montecal         Club Altamira del Sur      11     0     1     3
  17  Lorvin Orbaneja        Deportivo Valmora          11     0     1     3
* 18  Ivaro Jarvale          Academia Fontalva          10     3     0     3
  19  Hilaro Roblemar        Academia Valmora           11     3     0     3
  20  Kaelo Roblemar         Club Maldoria              11     3     0     3
  21  Hemiro Norvalde        Unión Pinarel              11     3     0     3
  22  Florian Urdangal       Academia Valmora           11     3     0     3
* 23  Lisandor Orbaneja      Academia Fontalva          11     3     0     3
  24  Nolven Orbaneja        Club Altamira del Sur      11     3     0     3
  25  Tavio Korvena          Club Brisanda              11     3     0     3

↑/↓ mover (1-25 de 50) · enter ver ficha · * tu club · esc volver
```

## Clasificación: PORTEROS · PORTERÍAS IMBATIDAS

```
PORTEROS · PORTERÍAS IMBATIDAS · Temporada 4 · Jornada 11 / 18

   #  Jugador                Equipo                     PJ   Imb   Enc
>  1  Lisandor Landrosa      Unión Castelmar            11     6     8
   2  Xandro Abrelda         Unión Pinarel              11     6    11
*  3  Orlandi Pedrosel       Academia Fontalva          11     5    15
   4  Quirino Pardovan       Academia Valmora           11     4    11
   5  Pelayo Cordival        Juventud Granvela          11     3    14
   6  Anselmo Urdangal       Club Brisanda              11     3    15
   7  Lorvin Illanes         Deportivo Valmora          11     3    15
   8  Lisandor Escalona      Unión Granvela             11     3    22
   9  Selmo Frondosa         Club Maldoria              11     2    18
  10  Fabrel Jurado          Club Altamira del Sur      11     2    20

↑/↓ mover · enter ver ficha · * tu club · esc volver
```

## Clasificación: PORTEROS · MENOS GOLES ENCAJADOS

```
PORTEROS · MENOS GOLES ENCAJADOS · Temporada 4 · Jornada 11 / 18

   #  Jugador                Equipo                     PJ   Enc Enc/PJ
>  1  Lisandor Landrosa      Unión Castelmar            11     8   0.73
   2  Quirino Pardovan       Academia Valmora           11    11   1.00
   3  Xandro Abrelda         Unión Pinarel              11    11   1.00
   4  Pelayo Cordival        Juventud Granvela          11    14   1.27
   5  Anselmo Urdangal       Club Brisanda              11    15   1.36
   6  Lorvin Illanes         Deportivo Valmora          11    15   1.36
*  7  Orlandi Pedrosel       Academia Fontalva          11    15   1.36
   8  Selmo Frondosa         Club Maldoria              11    18   1.64
   9  Fabrel Jurado          Club Altamira del Sur      11    20   1.82
  10  Lisandor Escalona      Unión Granvela             11    22   2.00

↑/↓ mover · enter ver ficha · * tu club · esc volver
```

## Clasificación: MEJORES VALORACIONES

```
MEJORES VALORACIONES · Temporada 4 · Jornada 11 / 18

   #  Jugador                Equipo                     PJ   Val Goles Asist
>  1  Orlandi Cordival       Unión Pinarel              11   6.6     5     1
   2  Baltor Frondosa        Unión Castelmar            11   6.5     3     3
   3  Joldan Norvalde        Unión Castelmar            11   6.5     3     3
*  4  Nolven Montecal        Academia Fontalva          11   6.5     5     0
   5  Elmiro Duelmo          Deportivo Valmora          11   6.5     6     1
   6  Hemiro Norvalde        Unión Pinarel              11   6.4     5     0
   7  Galdo Duelmo           Unión Castelmar            11   6.4     2     1
   8  Isidoro Frondosa       Academia Valmora           11   6.4     4     1
   9  Isidoro Sandoral       Unión Castelmar            11   6.4     3     0
  10  Isidoro Dunavar        Unión Granvela             11   6.3     3     2
  11  Ivaro Pardovan         Unión Castelmar            11   6.3     1     2
  12  Nolven Abrelda         Unión Pinarel              11   6.3     1     3
  13  Dorvan Jurado          Unión Granvela             11   6.3     5     0
  14  Arlen Cordival         Unión Pinarel              11   6.3     1     2
  15  Lisandor Bastaro       Club Altamira del Sur      11   6.3     4     1
  16  Nervio Jurado          Unión Castelmar            11   6.3     1     1
  17  Hemiro Aguirel         Unión Castelmar            11   6.3     0     2
* 18  Nervio Torvelo         Academia Fontalva          11   6.3     2     2
  19  Marvel Yrigaldo        Unión Pinarel              11   6.3     2     0
  20  Marvel Zaldumar        Academia Valmora           11   6.3     3     0
  21  Pelayo Torvelo         Academia Valmora           11   6.3     3     0
* 22  Lisandor Orbaneja      Academia Fontalva          11   6.3     2     2
  23  Pelayo Korvena         Unión Castelmar            11   6.3     1     1
  24  Xandro Zaldumar        Club Altamira del Sur      11   6.3     2     2
  25  Dorvan Sandoral        Academia Valmora           11   6.2     2     1

↑/↓ mover (1-25 de 50) · enter ver ficha · * tu club · esc volver
```

## Comparativa de equipos

```
EQUIPOS · Temporada 4 · Jornada 11 / 18

   #  Equipo                    PJ  GF  GC Imb SM  Am Ro  Goleador
>  1  Unión Castelmar           11  18   8   6  3  25  0  Baltor Frondosa (3)
   2  Unión Pinarel             11  15  11   6  3  17  0  Hemiro Norvalde (5)
   3  Academia Valmora          11  18  11   4  2  23  1  Isidoro Frondosa (4)
*  4  Academia Fontalva         11  18  15   5  3  33  2  Nolven Montecal (5)
   5  Club Altamira del Sur     11  15  20   2  4  17  1  Lisandor Bastaro (4)
   6  Club Brisanda             11  12  15   3  4  19  0  Jaldin Duelmo (2)
   7  Club Maldoria             11  13  18   2  3  18  0  Corvino Illanes (4)
   8  Unión Granvela            11  16  22   3  3  24  1  Dorvan Jurado (5)
   9  Deportivo Valmora         11  13  15   3  6  12  1  Elmiro Duelmo (6)
  10  Juventud Granvela         11  11  14   3  6  20  0  Arlen Calderan (2)

↑/↓ mover · enter plantilla · tab goleador/asistente · * tu club · esc volver
```

## Comparativa de equipos (tab)

```
EQUIPOS · Temporada 4 · Jornada 11 / 18

   #  Equipo                    PJ  GF  GC Imb SM  Am Ro  Asistente
>  1  Unión Castelmar           11  18   8   6  3  25  0  Baltor Frondosa (3)
   2  Unión Pinarel             11  15  11   6  3  17  0  Nolven Abrelda (3)
   3  Academia Valmora          11  18  11   4  2  23  1  Valdor Cordival (3)
*  4  Academia Fontalva         11  18  15   5  3  33  2  Ivaro Jarvale (2)
   5  Club Altamira del Sur     11  15  20   2  4  17  1  Xandro Zaldumar (2)
   6  Club Brisanda             11  12  15   3  4  19  0  Fabrel Olmedrar (2)
   7  Club Maldoria             11  13  18   2  3  18  0  Elmiro Torvelo (2)
   8  Unión Granvela            11  16  22   3  3  24  1  Pelayo Pardovan (3)
   9  Deportivo Valmora         11  13  15   3  6  12  1  Florian Quintaleo (4)
  10  Juventud Granvela         11  11  14   3  6  20  0  Mendo Olmedrar (3)

↑/↓ mover · enter plantilla · tab goleador/asistente · * tu club · esc volver
```

## Plantilla de un equipo con estadísticas

```
UNIÓN PINAREL · estadísticas

  Posición      Jugador                 PJ Tit   Min   G   A  Am  Ro Imb  Val
> Portero       Xandro Abrelda          11  11   990   0   0   0   0   6  6.1
  Portero       Daveo Urdangal           0   0     0   0   0   0   0   0    -
  Defensa       Anselmo Bocanegro       11  11   932   0   1   3   0   6  6.2
  Defensa       Daveo Gallardel         11  11   904   0   1   6   0   5  6.0
  Defensa       Arlen Cordival          11  11   764   1   2   0   0   4  6.3
  Defensa       Ivaro Landrosa          11  11   869   0   1   1   0   5  6.2
  Defensa       Valdor Korvena          11  11   784   1   0   0   0   3  6.2
  Defensa       Belron Garmendo         11   0   287   0   1   0   0   0  6.0
  Defensa       Nolven Montecal          7   0   171   0   0   1   0   0  6.0
  Mediocampista Hemiro Norvalde         11  11   849   5   0   3   0   0  6.4
  Mediocampista Zenel Herbasco          11  11   888   0   1   1   0   0  6.1
  Mediocampista Nolven Abrelda          11  11   883   1   3   0   0   0  6.3
  Mediocampista Yeraldo Pedrosel        11   0   328   0   0   0   0   0  6.0
  Mediocampista Lorvin Quintaleo         6   0    91   0   0   0   0   0  6.0
  Mediocampista Wenceo Gallardel         1   0    12   0   0   0   0   0  6.0
  Mediocampista Nolven Rivasel           0   0     0   0   0   0   0   0    -
  Delantero     Orlandi Cordival        11  11   868   5   1   1   0   0  6.6
  Delantero     Marvel Yrigaldo         11  11   905   2   0   0   0   0  6.3
  Delantero     Isidoro Irueta           9   0   239   0   0   1   0   0  6.0
  Delantero     Pelayo Zaldumar          5   0   126   0   0   0   0   0  6.0
  Delantero     Belron Torvelo           0   0     0   0   0   0   0   0    -
  Delantero     Mendo Estoral            0   0     0   0   0   0   0   0    -

↑/↓ mover · enter ver ficha · esc volver
```

## Ficha de un portero

```
FICHA · Orlandi Pedrosel
Academia Fontalva · Portero · 24 años · Valoración 75
RIT 63  TIR 48  PAS 68  REG 53  DEF 70  FIS 74  REF 92
Proyección: normal

Temporada 4 (en curso)
  PJ 11 · Tit 11 · Min 990 · G 0 · A 0 · Am 1 · Ro 0 · Imb 5 · Enc 15 · Val 5.8

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
  PJ 0 · Tit 0 · Min 0 · G 0 · A 0 · Am 0 · Ro 0 · Imb 0 · Val -

Trayectoria
  Todavía no ha terminado ninguna temporada en esta carrera.


esc volver
```

## Fin de temporada

```
TEMPORADA 1 TERMINADA

Campeón: Estrella Altamira del Sur
Tu equipo: Deportivo Granvela, puesto 4 de 10, 27 puntos

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
