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

Unión Dorada    Temporada 4 · Jornada 11 / 18
Valoración del equipo: 77

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

  Unión Embalse Nuevo         0 - 0   Club Brisanda
> Academia Dorada             1 - 2   Unión Dorada
  Juventud Peñaviva           1 - 3   Academia Isleta Roja
  Club Granvela               0 - 1   Academia Nuevaluna
  Deportivo Sanvedra          2 - 3   Unión Pinarel

enter continuar
```

## Tabla de posiciones

```
TABLA · Temporada 4 · Jornada 11 / 18

    #  Equipo                      PJ   G   E   P   GF   GC   DG  Pts
    1  Unión Pinarel               11   7   2   2   20    9  +11   23
    2  Club Granvela               11   7   1   3   11   10   +1   22
    3  Juventud Peñaviva           11   6   2   3   10    8   +2   20
>   4  Unión Dorada                11   6   1   4   10   10   +0   19
    5  Deportivo Sanvedra          11   4   4   3   21   11  +10   16
    6  Club Brisanda               11   4   3   4   16   15   +1   15
    7  Academia Isleta Roja        11   4   2   5   13   13   +0   14
    8  Academia Nuevaluna          11   4   0   7    9   15   -6   12
    9  Unión Embalse Nuevo         11   3   1   7   15   23   -8   10
   10  Academia Dorada             11   2   0   9    6   17  -11    6

esc volver
```

## Plantilla: atributos

```
PLANTILLA · Unión Dorada

  Posición       Nombre                  Ed  Val  RIT TIR PAS REG DEF FIS REF
  Portero        Nervio Illanes          29   85   75  63  81  63  80  82  89
  Portero        Isidoro Calderan        16   47   27  10  33   7  34  46  57
  Defensa        Florian Abrelda         25   79   65  55  59  67  91  82   1
  Defensa        Arlen Orbaneja          30   77   67  48  78  45  88  78   1
  Defensa        Joldan Norvalde         24   72   71  46  59  61  78  74   1
  Defensa        Eskardo Hontoria        35   71   49  67  75  65  85  61   1
  Defensa        Daveo Zaldumar          20   64   49  35  63  50  73  68   1
  Defensa        Dorvan Fuentalba        21   60   53  39  52  35  71  58   1
  Defensa        Zenel Norvalde          21   59   39  34  46  36  76  57   1
  Mediocampista  Rodvan Olmedrar         31   79   78  77  88  80  67  57   1
  Mediocampista  Kervo Yrigaldo          23   77   75  71  83  75  71  76   1
  Mediocampista  Hemiro Calderan         34   74   57  76  81  73  76  59   1
  Mediocampista  Cedrio Pardovan         28   72   56  68  80  76  59  71   1
  Mediocampista  Tavio Escalona          35   65   59  63  69  70  62  54   1
  Mediocampista  Orlandi Korvena         20   55   58  44  65  55  47  45   1
  Mediocampista  Zenel Valmeron          38   50   33  52  58  49  54  35   1
  Delantero      Xandro Herbasco         29   80   78  89  78  87  45  62   1
  Delantero      Hemiro Jarvale          31   77   71  84  72  81  61  66   1
  Delantero      Anselmo Sandoral        32   76   77  79  78  79  55  62   1
  Delantero      Ivaro Urdangal          23   75   76  80  75  78  43  66   1
  Delantero      Belron Abrelda          19   67   60  79  56  71  25  62   1
  Delantero      Baltor Orbaneja         18   55   56  67  43  53  17  43   1

tab estadísticas · esc volver
```

## Plantilla: estadísticas (tab)

```
PLANTILLA · Unión Dorada · estadísticas

  Posición      Jugador                 PJ Tit   Min   G   A  Am  Ro Imb  Val
  Portero       Nervio Illanes          11  11   990   0   0   0   0   5  6.0
  Portero       Isidoro Calderan         0   0     0   0   0   0   0   0    -
  Defensa       Florian Abrelda         11  11   899   1   0   2   0   5  6.2
  Defensa       Arlen Orbaneja          11  11   775   0   0   2   1   3  5.9
  Defensa       Joldan Norvalde         11  11   951   0   0   1   0   5  6.1
  Defensa       Eskardo Hontoria        11  11   873   0   0   2   0   2  5.9
  Defensa       Daveo Zaldumar           4   0   127   0   0   0   0   0  6.0
  Defensa       Dorvan Fuentalba         3   0   110   0   0   1   0   0  6.0
  Defensa       Zenel Norvalde           5   0   137   0   0   0   0   0  6.0
  Mediocampista Rodvan Olmedrar         11  11   896   3   0   2   0   0  6.3
  Mediocampista Kervo Yrigaldo          11  11   903   0   1   0   0   0  6.1
  Mediocampista Hemiro Calderan         11  11   815   1   2   0   0   0  6.2
  Mediocampista Cedrio Pardovan          6   0   136   0   0   0   0   0  6.0
  Mediocampista Tavio Escalona           5   0    76   0   0   0   0   0  6.0
  Mediocampista Orlandi Korvena          5   0   123   0   0   0   0   0  6.1
  Mediocampista Zenel Valmeron           5   0    68   0   0   0   0   0  6.0
  Delantero     Xandro Herbasco         11  11   864   2   1   1   0   0  6.2
  Delantero     Hemiro Jarvale          11  11   826   0   2   2   0   0  6.1
  Delantero     Anselmo Sandoral        11  11   950   2   1   1   0   0  6.3
  Delantero     Ivaro Urdangal           5   0   168   1   0   0   0   0  6.2
  Delantero     Belron Abrelda           4   0    72   0   0   0   0   0  6.0
  Delantero     Baltor Orbaneja          5   0    95   0   0   0   0   0  6.0

tab atributos · esc volver
```

## Historial de temporadas

```
HISTORIAL · Unión Dorada

  Temp  Campeón                     Tu puesto   Pts
     1  Club Granvela                 7 de 10    25
     2  Deportivo Sanvedra            7 de 10    22
     3  Juventud Peñaviva             2 de 10    30

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
>  1  Rodvan Cordival        Unión Pinarel              11     7     2   857
   2  Mendo Orbaneja         Deportivo Sanvedra         11     7     0   893
   3  Selmo Quintaleo        Deportivo Sanvedra         11     5     0   965
   4  Anselmo Torvelo        Unión Pinarel              11     4     1   842
   5  Dorvan Fuentalba       Club Brisanda              11     4     1   861
   6  Daveo Lacerna          Unión Embalse Nuevo        11     4     1   902
   7  Tavio Duelmo           Club Granvela              11     4     0   882
   8  Kervo Norvalde         Club Brisanda              11     3     2   871
   9  Fabrel Pardovan        Unión Embalse Nuevo        11     3     2   943
  10  Valdor Roblemar        Club Brisanda              11     3     2   990
  11  Hemiro Jurado          Unión Embalse Nuevo        11     3     1   903
  12  Isidoro Pardovan       Club Granvela              11     3     1   984
  13  Hemiro Escalona        Deportivo Sanvedra         11     3     0   873
  14  Fabrel Sandoral        Academia Isleta Roja       11     3     0   896
* 15  Rodvan Olmedrar        Unión Dorada               11     3     0   896
  16  Dorvan Escalona        Club Brisanda              11     2     3   780
  17  Zenel Olmedrar         Deportivo Sanvedra         11     2     3   840
  18  Joldan Dunavar         Juventud Peñaviva          11     2     3   854
  19  Florian Duelmo         Deportivo Sanvedra         11     2     3   864
  20  Eskardo Norvalde       Academia Isleta Roja       11     2     2   794
  21  Selmo Pedrosel         Club Brisanda              11     2     2   888
  22  Orlandi Norvalde       Academia Isleta Roja       11     2     2   933
  23  Nervio Gallardel       Unión Pinarel              11     2     1   821
  24  Orlandi Bocanegro      Academia Nuevaluna         11     2     1   843
* 25  Xandro Herbasco        Unión Dorada               11     2     1   864

↑/↓ mover (1-25 de 50) · enter ver ficha · * tu club · esc volver
```

## Clasificación: ASISTENTES

```
ASISTENTES · Temporada 4 · Jornada 11 / 18

   #  Jugador                Equipo                     PJ Asist Goles   Min
>  1  Dorvan Escalona        Club Brisanda              11     3     2   780
   2  Zenel Olmedrar         Deportivo Sanvedra         11     3     2   840
   3  Joldan Dunavar         Juventud Peñaviva          11     3     2   854
   4  Florian Duelmo         Deportivo Sanvedra         11     3     2   864
   5  Daveo Duelmo           Unión Embalse Nuevo        11     3     1   875
   6  Jaldin Duelmo          Academia Dorada            11     3     0   891
   7  Rodvan Cordival        Unión Pinarel              11     2     7   857
   8  Kervo Norvalde         Club Brisanda              11     2     3   871
   9  Fabrel Pardovan        Unión Embalse Nuevo        11     2     3   943
  10  Valdor Roblemar        Club Brisanda              11     2     3   990
  11  Eskardo Norvalde       Academia Isleta Roja       11     2     2   794
  12  Selmo Pedrosel         Club Brisanda              11     2     2   888
  13  Orlandi Norvalde       Academia Isleta Roja       11     2     2   933
* 14  Hemiro Calderan        Unión Dorada               11     2     1   815
  15  Marvel Irueta          Unión Pinarel              11     2     1   827
  16  Eskardo Abrelda        Deportivo Sanvedra         11     2     1   881
  17  Hilaro Dunavar         Juventud Peñaviva          11     2     1   911
  18  Nervio Maldovar        Academia Nuevaluna         11     2     1   927
  19  Xandro Yrigaldo        Deportivo Sanvedra         11     2     0   742
  20  Yeraldo Gallardel      Academia Isleta Roja       11     2     0   745
  21  Corvino Olmedrar       Academia Isleta Roja       11     2     0   792
  22  Galdo Aguirel          Club Granvela              11     2     0   807
* 23  Hemiro Jarvale         Unión Dorada               11     2     0   826
  24  Zenel Korvena          Club Granvela              11     2     0   892
  25  Gervasio Valmeron      Deportivo Sanvedra         11     2     0   977

↑/↓ mover (1-25 de 50) · enter ver ficha · * tu club · esc volver
```

## Clasificación: TARJETAS

```
TARJETAS · Temporada 4 · Jornada 11 / 18

   #  Jugador                Equipo                     PJ  Amar  Roja   Pts
>  1  Marvel Irueta          Unión Pinarel              11     1     2     7
   2  Corvino Olmedrar       Academia Isleta Roja       11     4     1     7
   3  Dorvan Roblemar        Club Granvela              11     3     1     6
*  4  Arlen Orbaneja         Unión Dorada               11     2     1     5
   5  Eskardo Sandoral       Club Granvela              11     2     1     5
   6  Galdo Irueta           Academia Nuevaluna         11     2     1     5
   7  Nolven Herbasco        Academia Nuevaluna         11     2     1     5
   8  Hemiro Pedrosel        Juventud Peñaviva          11     2     1     5
   9  Yeraldo Gallardel      Academia Isleta Roja       11     1     1     4
  10  Ivaro Yrigaldo         Unión Embalse Nuevo        11     1     1     4
  11  Kervo Duelmo           Academia Dorada            11     4     0     4
  12  Selmo Pedrosel         Club Brisanda              11     4     0     4
  13  Nervio Maldovar        Academia Nuevaluna         11     4     0     4
  14  Fabrel Pardovan        Unión Embalse Nuevo        11     4     0     4
  15  Orlandi Illanes        Unión Embalse Nuevo        11     4     0     4
  16  Fabrel Torvelo         Club Granvela              11     4     0     4
  17  Marvel Aguirel         Academia Dorada             6     0     1     3
  18  Arlen Zaldumar         Deportivo Sanvedra          6     0     1     3
  19  Joldan Illanes         Unión Embalse Nuevo        11     3     0     3
  20  Wenceo Bastaro         Academia Nuevaluna         11     3     0     3
  21  Joldan Dunavar         Juventud Peñaviva          11     3     0     3
  22  Florian Duelmo         Deportivo Sanvedra         11     3     0     3
  23  Marvel Zaldumar        Club Granvela              11     3     0     3
  24  Jaldin Duelmo          Academia Dorada            11     3     0     3
  25  Baltor Lacerna         Unión Embalse Nuevo        11     3     0     3

↑/↓ mover (1-25 de 50) · enter ver ficha · * tu club · esc volver
```

## Clasificación: PORTEROS · PORTERÍAS IMBATIDAS

```
PORTEROS · PORTERÍAS IMBATIDAS · Temporada 4 · Jornada 11 / 18

   #  Jugador                Equipo                     PJ   Imb   Enc
>  1  Eskardo Herbasco       Club Granvela              11     7    10
   2  Kaelo Irueta           Unión Pinarel              11     6     9
*  3  Nervio Illanes         Unión Dorada               11     5    10
   4  Daveo Montecal         Deportivo Sanvedra         11     5    11
   5  Hemiro Pedrosel        Juventud Peñaviva          11     3     8
   6  Orlandi Pardovan       Academia Nuevaluna         11     3    15
   7  Belron Escalona        Academia Isleta Roja       11     2    13
   8  Selmo Zaldumar         Club Brisanda              11     1    15
   9  Zenel Korvena          Academia Dorada            11     1    17

↑/↓ mover · enter ver ficha · * tu club · esc volver
```

## Clasificación: PORTEROS · MENOS GOLES ENCAJADOS

```
PORTEROS · MENOS GOLES ENCAJADOS · Temporada 4 · Jornada 11 / 18

   #  Jugador                Equipo                     PJ   Enc Enc/PJ
>  1  Hemiro Pedrosel        Juventud Peñaviva          11     8   0.73
   2  Kaelo Irueta           Unión Pinarel              11     9   0.82
   3  Eskardo Herbasco       Club Granvela              11    10   0.91
*  4  Nervio Illanes         Unión Dorada               11    10   0.91
   5  Daveo Montecal         Deportivo Sanvedra         11    11   1.00
   6  Belron Escalona        Academia Isleta Roja       11    13   1.18
   7  Orlandi Pardovan       Academia Nuevaluna         11    15   1.36
   8  Selmo Zaldumar         Club Brisanda              11    15   1.36
   9  Zenel Korvena          Academia Dorada            11    17   1.55
  10  Mendo Zaldumar         Unión Embalse Nuevo        11    23   2.09

↑/↓ mover · enter ver ficha · * tu club · esc volver
```

## Clasificación: MEJORES VALORACIONES

```
MEJORES VALORACIONES · Temporada 4 · Jornada 11 / 18

   #  Jugador                Equipo                     PJ   Val Goles Asist
>  1  Rodvan Cordival        Unión Pinarel              11   6.9     7     2
   2  Mendo Orbaneja         Deportivo Sanvedra         11   6.6     7     0
   3  Anselmo Torvelo        Unión Pinarel              11   6.5     4     1
   4  Selmo Quintaleo        Deportivo Sanvedra         11   6.5     5     0
   5  Tavio Duelmo           Club Granvela              11   6.4     4     0
   6  Dorvan Fuentalba       Club Brisanda              11   6.4     4     1
   7  Isidoro Pardovan       Club Granvela              11   6.4     3     1
   8  Kervo Norvalde         Club Brisanda              11   6.4     3     2
   9  Joldan Dunavar         Juventud Peñaviva          11   6.3     2     3
  10  Nervio Gallardel       Unión Pinarel              11   6.3     2     1
  11  Dorvan Escalona        Club Brisanda              11   6.3     2     3
  12  Valdor Roblemar        Club Brisanda              11   6.3     3     2
  13  Elmiro Dunavar         Unión Pinarel              11   6.3     2     1
  14  Daveo Lacerna          Unión Embalse Nuevo        11   6.3     4     1
  15  Eskardo Sandoral       Club Granvela              11   6.3     1     0
  16  Florian Duelmo         Deportivo Sanvedra         11   6.3     2     3
  17  Xandro Aguirel         Unión Pinarel              11   6.3     1     0
  18  Nolven Abrelda         Club Granvela              11   6.3     2     0
  19  Zenel Olmedrar         Deportivo Sanvedra         11   6.3     2     3
  20  Hilaro Dunavar         Juventud Peñaviva          11   6.3     1     2
  21  Xandro Cordival        Unión Pinarel              11   6.3     0     1
* 22  Anselmo Sandoral       Unión Dorada               11   6.3     2     1
  23  Kaelo Garmendo         Unión Pinarel              11   6.3     2     0
* 24  Rodvan Olmedrar        Unión Dorada               11   6.3     3     0
  25  Eskardo Abrelda        Deportivo Sanvedra         11   6.3     1     2

↑/↓ mover (1-25 de 50) · enter ver ficha · * tu club · esc volver
```

## Comparativa de equipos

```
EQUIPOS · Temporada 4 · Jornada 11 / 18

   #  Equipo                    PJ  GF  GC Imb SM  Am Ro  Goleador
>  1  Unión Pinarel             11  20   9   6  1  14  2  Rodvan Cordival (7)
   2  Club Granvela             11  11  10   7  2  18  2  Tavio Duelmo (4)
   3  Juventud Peñaviva         11  10   8   4  5  20  1  Elmiro Valmeron (2)
*  4  Unión Dorada              11  10  10   5  3  14  1  Rodvan Olmedrar (3)
   5  Deportivo Sanvedra        11  21  11   5  2  26  1  Mendo Orbaneja (7)
   6  Club Brisanda             11  16  15   1  3  19  0  Dorvan Fuentalba (4)
   7  Academia Isleta Roja      11  13  13   2  4  20  2  Fabrel Sandoral (3)
   8  Academia Nuevaluna        11   9  15   3  4  21  2  Nolven Herbasco (2)
   9  Unión Embalse Nuevo       11  15  23   0  3  24  1  Daveo Lacerna (4)
  10  Academia Dorada           11   6  17   1  7  17  1  Belron Garmendo (1)

↑/↓ mover · enter plantilla · tab goleador/asistente · * tu club · esc volver
```

## Comparativa de equipos (tab)

```
EQUIPOS · Temporada 4 · Jornada 11 / 18

   #  Equipo                    PJ  GF  GC Imb SM  Am Ro  Asistente
>  1  Unión Pinarel             11  20   9   6  1  14  2  Marvel Irueta (2)
   2  Club Granvela             11  11  10   7  2  18  2  Galdo Aguirel (2)
   3  Juventud Peñaviva         11  10   8   4  5  20  1  Joldan Dunavar (3)
*  4  Unión Dorada              11  10  10   5  3  14  1  Hemiro Calderan (2)
   5  Deportivo Sanvedra        11  21  11   5  2  26  1  Florian Duelmo (3)
   6  Club Brisanda             11  16  15   1  3  19  0  Dorvan Escalona (3)
   7  Academia Isleta Roja      11  13  13   2  4  20  2  Corvino Olmedrar (2)
   8  Academia Nuevaluna        11   9  15   3  4  21  2  Nervio Maldovar (2)
   9  Unión Embalse Nuevo       11  15  23   0  3  24  1  Daveo Duelmo (3)
  10  Academia Dorada           11   6  17   1  7  17  1  Jaldin Duelmo (3)

↑/↓ mover · enter plantilla · tab goleador/asistente · * tu club · esc volver
```

## Plantilla de un equipo con estadísticas

```
CLUB GRANVELA · estadísticas

  Posición      Jugador                 PJ Tit   Min   G   A  Am  Ro Imb  Val
> Portero       Eskardo Herbasco        11  11   990   0   0   0   0   7  6.2
  Portero       Galdo Herbasco           0   0     0   0   0   0   0   0    -
  Defensa       Marvel Zaldumar         11  11   879   0   1   3   0   6  6.2
  Defensa       Eskardo Sandoral        11  11   836   1   0   2   1   6  6.3
  Defensa       Fabrel Torvelo          11  11   955   0   1   4   0   7  6.2
  Defensa       Zenel Korvena           11  11   892   0   2   0   0   5  6.2
  Defensa       Kervo Urdangal           5   0    69   0   0   0   0   0  6.0
  Defensa       Hemiro Orbaneja          6   0   142   0   0   0   0   0  6.0
  Defensa       Belron Irueta            5   0   112   0   1   0   0   0  6.2
  Mediocampista Galdo Aguirel           11  11   807   0   2   1   0   0  6.2
  Mediocampista Isidoro Pardovan        11  11   984   3   1   2   0   0  6.4
  Mediocampista Dorvan Roblemar         11  11   844   0   0   3   1   0  6.0
  Mediocampista Wenceo Pedrosel          4   0    76   0   0   0   0   0  6.0
  Mediocampista Jaldin Orbaneja          3   0    35   0   0   0   0   0  6.0
  Mediocampista Lisandor Estoral         5   0   169   0   0   0   0   0  6.0
  Mediocampista Wenceo Nebreda           2   0    48   0   0   0   0   0  6.0
  Delantero     Nolven Abrelda          11  11   845   2   0   0   0   0  6.3
  Delantero     Corvino Cordival        11  11   838   1   1   1   0   0  6.2
  Delantero     Tavio Duelmo            11  11   882   4   0   2   0   0  6.4
  Delantero     Orlandi Illanes          3   0    83   0   0   0   0   0  6.1
  Delantero     Pelayo Montecal          8   0   176   0   0   0   0   0  6.0
  Delantero     Tavio Hontoria           5   0   101   0   0   0   0   0  6.0

↑/↓ mover · enter ver ficha · esc volver
```

## Ficha de un portero

```
FICHA · Nervio Illanes
Unión Dorada · Portero · 29 años · Valoración 85
RIT 75  TIR 63  PAS 81  REG 63  DEF 80  FIS 82  REF 89

Temporada 4 (en curso)
  PJ 11 · Tit 11 · Min 990 · G 0 · A 0 · Am 0 · Ro 0 · Imb 5 · Enc 10 · Val 6.0

Trayectoria
  Temp  Club                       PJ   Min   G   A  Am  Ro Imb  Val
     1  Unión Dorada               18  1620   0   0   1   0   5  5.8
     2  Unión Dorada               18  1620   0   0   1   0   8  5.9
     3  Unión Dorada               18  1620   0   0   3   0   9  6.0

Carrera: 65 partidos · 0 goles · 0 asistencias · 5 amarillas · 0 rojas

esc volver
```

## Ficha de un delantero

```
FICHA · Baltor Orbaneja
Unión Dorada · Delantero · 18 años · Valoración 55
RIT 56  TIR 67  PAS 43  REG 53  DEF 17  FIS 43  REF 1

Temporada 4 (en curso)
  PJ 5 · Tit 0 · Min 95 · G 0 · A 0 · Am 0 · Ro 0 · Imb 0 · Val 6.0

Trayectoria
  Temp  Club                       PJ   Min   G   A  Am  Ro Imb  Val
     3  Unión Dorada                7   184   0   0   0   0   0  6.0

Carrera: 12 partidos · 0 goles · 0 asistencias · 0 amarillas · 0 rojas

esc volver
```

## Fin de temporada

```
TEMPORADA 1 TERMINADA

Campeón: Unión Marisal
Tu equipo: Club Marisal, puesto 6 de 10, 22 puntos

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
TEMPORADA 2 · Club Marisal

Valoración del equipo: 80 → 78

Nadie se retira este año.

No llega nadie de la cantera.

En toda la liga: 5 retiros y 5 juveniles.

enter continuar
```

## Confirmar una carrera nueva

```
NUEVA CARRERA

Esto empieza una carrera nueva y pierdes la actual
(Temporada 2 con Club Marisal). Aún no hay guardado automático.
¿Seguro?

> No, volver
  Sí, empezar de cero

↑/↓ mover · enter elegir · esc volver
```
