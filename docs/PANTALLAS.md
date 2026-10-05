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

Estrella Altamira del Sur    Temporada 4 · Jornada 11 / 18
Valoración del equipo: 76

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

  Unión Fontalva              0 - 0   Club Brisanda
  Deportivo Nuevaluna         1 - 2   Estrella Salinera
> Unión Altamira del Sur      1 - 3   Estrella Altamira del Sur
  Academia Salinera           0 - 1   Juventud Riobravo
  Juventud Altamira del Sur   2 - 3   Deportivo Hondaluz

enter continuar
```

## Tabla de posiciones

```
TABLA · Temporada 4 · Jornada 11 / 18

    #  Equipo                      PJ   G   E   P   GF   GC   DG  Pts
    1  Juventud Altamira del Sur   11   7   3   1   20    8  +12   24
    2  Unión Fontalva              11   7   1   3   15   13   +2   22
    3  Deportivo Hondaluz          11   6   2   3   15    9   +6   20
>   4  Estrella Altamira del Sur   11   5   2   4   13   10   +3   17
    5  Estrella Salinera           11   5   2   4   10   11   -1   17
    6  Unión Altamira del Sur      11   5   1   5   12   13   -1   16
    7  Club Brisanda               11   4   2   5   15   17   -2   14
    8  Academia Salinera           11   4   2   5    5    7   -2   14
    9  Juventud Riobravo           11   2   1   8   11   17   -6    7
   10  Deportivo Nuevaluna         11   2   0   9    7   18  -11    6

esc volver
```

## Plantilla: atributos

```
PLANTILLA · Estrella Altamira del Sur

  Posición       Nombre                  Ed  Val  RIT TIR PAS REG DEF FIS REF
  Portero        Mendo Cordival          30   79   64  52  78  45  72  74  84
  Portero        Fabrel Landrosa         30   69   48  36  64  41  48  68  78
  Defensa        Wenceo Zaldumar         34   79   66  50  72  59  94  74   1
  Defensa        Hilaro Bastaro          34   78   73  61  76  66  85  76   1
  Defensa        Ulmar Hontoria          26   76   74  58  71  55  81  79   1
  Defensa        Corvino Nebreda         31   70   61  42  69  54  77  74   1
  Defensa        Yeraldo Quintaleo       32   64   54  66  56  57  80  50   1
  Defensa        Fabrel Jurado           37   57   43  50  52  65  67  53   1
  Defensa        Cedrio Abrelda          37   53   35  29  38  46  73  45   1
  Mediocampista  Lorvin Dunavar          25   79   68  72  85  87  68  71   1
  Mediocampista  Kaelo Dunavar           30   79   68  74  86  78  76  81   1
  Mediocampista  Dorvan Orbaneja         32   77   70  68  84  80  78  66   1
  Mediocampista  Galdo Abrelda           25   74   71  73  79  74  58  79   1
  Mediocampista  Marvel Dunavar          25   72   65  68  77  75  63  73   1
  Mediocampista  Zenel Rivasel           37   64   55  57  72  67  71  46   1
  Mediocampista  Yeraldo Cordival        20   62   56  49  65  75  52  59   1
  Delantero      Hemiro Jurado           25   73   75  76  73  74  48  68   1
  Delantero      Quirino Zaldumar        33   68   55  85  73  67  42  49   1
  Delantero      Valdor Quintaleo        33   66   53  70  70  76  46  63   1
  Delantero      Dorvan Nebreda          36   62   49  77  52  73  34  37   1
  Delantero      Hemiro Yrigaldo         20   57   55  69  46  56  22  53   1
  Delantero      Elmiro Montecal         18   47   51  54  42  44  10  47   1

tab estadísticas · esc volver
```

## Plantilla: estadísticas (tab)

```
PLANTILLA · Estrella Altamira del Sur · estadísticas

  Posición      Jugador                 PJ Tit   Min   G   A  Am  Ro Imb  Val
  Portero       Mendo Cordival          11  11   990   0   0   0   0   3  5.9
  Portero       Fabrel Landrosa          0   0     0   0   0   0   0   0    -
  Defensa       Wenceo Zaldumar         11  11   886   0   1   1   0   3  6.0
  Defensa       Hilaro Bastaro          11  11   867   0   0   1   0   2  5.9
  Defensa       Ulmar Hontoria          11  11   859   0   1   1   0   2  6.0
  Defensa       Corvino Nebreda         11  11   883   0   0   2   0   3  5.9
  Defensa       Yeraldo Quintaleo        7   0   137   0   0   0   0   0  6.0
  Defensa       Fabrel Jurado            4   0    67   0   0   1   0   0  6.0
  Defensa       Cedrio Abrelda           7   0   196   0   0   1   0   0  5.9
  Mediocampista Lorvin Dunavar          11  11   928   4   2   2   0   0  6.4
  Mediocampista Kaelo Dunavar           11  11   843   1   0   3   0   0  6.0
  Mediocampista Dorvan Orbaneja         11  11   892   3   3   3   0   0  6.4
  Mediocampista Galdo Abrelda            3   0    64   0   0   0   0   0  6.1
  Mediocampista Marvel Dunavar           2   0    64   1   0   0   0   0  6.5
  Mediocampista Zenel Rivasel            3   0   107   0   0   1   0   0  5.9
  Mediocampista Yeraldo Cordival         3   0    66   0   0   1   0   0  5.9
  Delantero     Hemiro Jurado           11  11   850   1   0   0   0   0  6.1
  Delantero     Quirino Zaldumar        11  11   794   0   0   0   0   0  6.0
  Delantero     Valdor Quintaleo        11  11   814   1   0   1   0   0  6.1
  Delantero     Dorvan Nebreda           4   0   110   0   0   0   0   0  6.0
  Delantero     Hemiro Yrigaldo          9   0   269   1   0   0   0   0  6.1
  Delantero     Elmiro Montecal          7   0   204   1   0   0   0   0  6.1

tab atributos · esc volver
```

## Historial de temporadas

```
HISTORIAL · Estrella Altamira del Sur

  Temp  Campeón                     Tu puesto   Pts
     1  Deportivo Hondaluz            9 de 10    20
     2  Juventud Altamira del Sur     5 de 10    24
     3  Academia Salinera             4 de 10    26

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
>  1  Dorvan Landrosa        Deportivo Hondaluz         11     5     0   857
   2  Lisandor Illanes       Juventud Altamira del Sur  11     4     2   897
*  3  Lorvin Dunavar         Estrella Altamira del Sur  11     4     2   928
   4  Selmo Bastaro          Deportivo Hondaluz         11     4     1   957
   5  Tavio Orbaneja         Unión Fontalva             11     4     0   830
*  6  Dorvan Orbaneja        Estrella Altamira del Sur  11     3     3   892
   7  Kaelo Zaldumar         Unión Altamira del Sur     11     3     3   928
   8  Zenel Jarvale          Juventud Altamira del Sur  11     3     2   769
   9  Florian Bocanegro      Unión Fontalva             11     3     2   898
  10  Xandro Bastaro         Unión Fontalva             11     3     2   915
  11  Lorvin Duelmo          Club Brisanda              11     3     2   969
  12  Ivaro Urdangal         Unión Fontalva             11     3     1   865
  13  Orlandi Quintaleo      Club Brisanda              11     3     1   874
  14  Lorvin Lacerna         Estrella Salinera          11     3     0   779
  15  Marvel Norvalde        Juventud Altamira del Sur  11     3     0   906
  16  Lisandor Garmendo      Juventud Altamira del Sur  11     3     0   931
  17  Dorvan Herbasco        Club Brisanda              11     2     3   734
  18  Nervio Bocanegro       Deportivo Nuevaluna        11     2     2   920
  19  Baltor Pardovan        Unión Altamira del Sur     11     2     1   707
  20  Hemiro Pedrosel        Unión Altamira del Sur     11     2     1   811
  21  Wenceo Abrelda         Club Brisanda              11     2     1   905
  22  Galdo Montecal         Estrella Salinera          11     2     1   959
  23  Lorvin Abrelda         Juventud Riobravo           5     2     0   148
  24  Florian Olmedrar       Deportivo Hondaluz          6     2     0   174
  25  Eskardo Olmedrar       Juventud Altamira del Sur  11     2     0   783

↑/↓ mover (1-25 de 50) · enter ver ficha · * tu club · esc volver
```

## Clasificación: ASISTENTES

```
ASISTENTES · Temporada 4 · Jornada 11 / 18

   #  Jugador                Equipo                     PJ Asist Goles   Min
>  1  Dorvan Orbaneja        Estrella Altamira del Sur  11     3     3   892
   2  Kaelo Zaldumar         Unión Altamira del Sur     11     3     3   928
   3  Dorvan Herbasco        Club Brisanda              11     3     2   734
   4  Lisandor Bastaro       Deportivo Hondaluz         11     3     1   797
   5  Wenceo Korvena         Juventud Altamira del Sur  11     3     1   846
   6  Mendo Fuentalba        Estrella Salinera          11     3     1   861
   7  Ivaro Cordival         Juventud Riobravo          11     3     0   789
   8  Anselmo Rivasel        Deportivo Nuevaluna        11     3     0   824
   9  Galdo Lacerna          Unión Fontalva             11     3     0   923
  10  Lisandor Illanes       Juventud Altamira del Sur  11     2     4   897
* 11  Lorvin Dunavar         Estrella Altamira del Sur  11     2     4   928
  12  Zenel Jarvale          Juventud Altamira del Sur  11     2     3   769
  13  Florian Bocanegro      Unión Fontalva             11     2     3   898
  14  Xandro Bastaro         Unión Fontalva             11     2     3   915
  15  Lorvin Duelmo          Club Brisanda              11     2     3   969
  16  Nervio Bocanegro       Deportivo Nuevaluna        11     2     2   920
  17  Belron Bocanegro       Club Brisanda              11     2     1   894
  18  Belron Duelmo          Deportivo Hondaluz         11     2     0   789
  19  Isidoro Roblemar       Juventud Riobravo          11     2     0   852
  20  Selmo Bastaro          Deportivo Hondaluz         11     1     4   957
  21  Ivaro Urdangal         Unión Fontalva             11     1     3   865
  22  Orlandi Quintaleo      Club Brisanda              11     1     3   874
  23  Baltor Pardovan        Unión Altamira del Sur     11     1     2   707
  24  Hemiro Pedrosel        Unión Altamira del Sur     11     1     2   811
  25  Wenceo Abrelda         Club Brisanda              11     1     2   905

↑/↓ mover (1-25 de 50) · enter ver ficha · * tu club · esc volver
```

## Clasificación: TARJETAS

```
TARJETAS · Temporada 4 · Jornada 11 / 18

   #  Jugador                Equipo                     PJ  Amar  Roja   Pts
>  1  Galdo Montecal         Estrella Salinera          11     8     1    11
   2  Eskardo Norvalde       Academia Salinera          11     6     1     9
   3  Xandro Gallardel       Unión Fontalva             11     4     1     7
   4  Arlen Sandoral         Deportivo Nuevaluna        11     3     1     6
   5  Eskardo Nebreda        Academia Salinera          11     3     1     6
   6  Belron Duelmo          Deportivo Hondaluz         11     2     1     5
   7  Tavio Orbaneja         Unión Fontalva             11     2     1     5
   8  Valdor Rivasel         Juventud Riobravo          11     2     1     5
   9  Ulmar Hontoria         Deportivo Nuevaluna        11     5     0     5
  10  Joldan Landrosa        Juventud Riobravo          11     1     1     4
  11  Ivaro Garmendo         Estrella Salinera          11     1     1     4
  12  Tavio Torvelo          Deportivo Nuevaluna        11     1     1     4
  13  Joldan Herbasco        Unión Fontalva             11     4     0     4
  14  Tavio Zaldumar         Juventud Altamira del Sur  11     4     0     4
  15  Corvino Bastaro        Academia Salinera          11     4     0     4
  16  Zenel Estoral          Unión Fontalva             11     4     0     4
  17  Ulmar Duelmo           Deportivo Hondaluz         11     4     0     4
  18  Joldan Frondosa        Estrella Salinera          11     4     0     4
  19  Nervio Rivasel         Deportivo Nuevaluna         3     0     1     3
  20  Eskardo Frondosa       Academia Salinera          11     3     0     3
* 21  Kaelo Dunavar          Estrella Altamira del Sur  11     3     0     3
  22  Wenceo Korvena         Juventud Altamira del Sur  11     3     0     3
  23  Ivaro Pedrosel         Deportivo Hondaluz         11     3     0     3
  24  Isidoro Dunavar        Juventud Riobravo          11     3     0     3
  25  Anselmo Hontoria       Club Brisanda              11     3     0     3

↑/↓ mover (1-25 de 50) · enter ver ficha · * tu club · esc volver
```

## Clasificación: PORTEROS · PORTERÍAS IMBATIDAS

```
PORTEROS · PORTERÍAS IMBATIDAS · Temporada 4 · Jornada 11 / 18

   #  Jugador                Equipo                     PJ   Imb   Enc
>  1  Nolven Jarvale         Academia Salinera          11     6     7
   2  Joldan Dunavar         Juventud Altamira del Sur  11     5     8
   3  Isidoro Montecal       Deportivo Hondaluz         11     5     9
   4  Eskardo Cordival       Estrella Salinera          11     5    11
   5  Kervo Illanes          Unión Fontalva             11     4    13
   6  Lisandor Montecal      Unión Altamira del Sur     11     4    13
*  7  Mendo Cordival         Estrella Altamira del Sur  11     3    10
   8  Anselmo Urdangal       Club Brisanda              11     1    17
   9  Orlandi Jarvale        Juventud Riobravo          11     1    17
  10  Gervasio Roblemar      Deportivo Nuevaluna        11     1    18

↑/↓ mover · enter ver ficha · * tu club · esc volver
```

## Clasificación: PORTEROS · MENOS GOLES ENCAJADOS

```
PORTEROS · MENOS GOLES ENCAJADOS · Temporada 4 · Jornada 11 / 18

   #  Jugador                Equipo                     PJ   Enc Enc/PJ
>  1  Nolven Jarvale         Academia Salinera          11     7   0.64
   2  Joldan Dunavar         Juventud Altamira del Sur  11     8   0.73
   3  Isidoro Montecal       Deportivo Hondaluz         11     9   0.82
*  4  Mendo Cordival         Estrella Altamira del Sur  11    10   0.91
   5  Eskardo Cordival       Estrella Salinera          11    11   1.00
   6  Kervo Illanes          Unión Fontalva             11    13   1.18
   7  Lisandor Montecal      Unión Altamira del Sur     11    13   1.18
   8  Anselmo Urdangal       Club Brisanda              11    17   1.55
   9  Orlandi Jarvale        Juventud Riobravo          11    17   1.55
  10  Gervasio Roblemar      Deportivo Nuevaluna        11    18   1.64

↑/↓ mover · enter ver ficha · * tu club · esc volver
```

## Clasificación: MEJORES VALORACIONES

```
MEJORES VALORACIONES · Temporada 4 · Jornada 11 / 18

   #  Jugador                Equipo                     PJ   Val Goles Asist
>  1  Lisandor Illanes       Juventud Altamira del Sur  11   6.6     4     2
   2  Dorvan Landrosa        Deportivo Hondaluz         11   6.5     5     0
   3  Selmo Bastaro          Deportivo Hondaluz         11   6.5     4     1
   4  Zenel Jarvale          Juventud Altamira del Sur  11   6.5     3     2
   5  Florian Bocanegro      Unión Fontalva             11   6.4     3     2
*  6  Lorvin Dunavar         Estrella Altamira del Sur  11   6.4     4     2
   7  Xandro Bastaro         Unión Fontalva             11   6.4     3     2
   8  Ivaro Urdangal         Unión Fontalva             11   6.4     3     1
   9  Marvel Norvalde        Juventud Altamira del Sur  11   6.4     3     0
  10  Lisandor Garmendo      Juventud Altamira del Sur  11   6.4     3     0
  11  Kaelo Zaldumar         Unión Altamira del Sur     11   6.4     3     3
  12  Wenceo Korvena         Juventud Altamira del Sur  11   6.4     1     3
* 13  Dorvan Orbaneja        Estrella Altamira del Sur  11   6.4     3     3
  14  Tavio Orbaneja         Unión Fontalva             11   6.3     4     0
  15  Lisandor Bastaro       Deportivo Hondaluz         11   6.3     1     3
  16  Dorvan Herbasco        Club Brisanda              11   6.3     2     3
  17  Lorvin Duelmo          Club Brisanda              11   6.3     3     2
  18  Eskardo Olmedrar       Juventud Altamira del Sur  11   6.3     2     0
  19  Mendo Fuentalba        Estrella Salinera          11   6.3     1     3
  20  Lorvin Lacerna         Estrella Salinera          11   6.3     3     0
  21  Lorvin Abrelda         Juventud Riobravo           5   6.3     2     0
  22  Florian Olmedrar       Deportivo Hondaluz          6   6.2     2     0
  23  Orlandi Quintaleo      Club Brisanda              11   6.2     3     1
  24  Galdo Lacerna          Unión Fontalva             11   6.2     0     3
  25  Baltor Pardovan        Unión Altamira del Sur     11   6.2     2     1

↑/↓ mover (1-25 de 50) · enter ver ficha · * tu club · esc volver
```

## Comparativa de equipos

```
EQUIPOS · Temporada 4 · Jornada 11 / 18

   #  Equipo                    PJ  GF  GC Imb SM  Am Ro  Goleador
>  1  Juventud Altamira del Sur 11  20   8   5  2  22  0  Lisandor Illanes (4)
   2  Unión Fontalva            11  15  13   4  1  27  2  Tavio Orbaneja (4)
   3  Deportivo Hondaluz        11  15   9   5  3  19  1  Dorvan Landrosa (5)
*  4  Estrella Altamira del Sur 11  13  10   3  4  18  0  Lorvin Dunavar (4)
   5  Estrella Salinera         11  10  11   5  3  21  2  Lorvin Lacerna (3)
   6  Unión Altamira del Sur    11  12  13   4  3  12  0  Kaelo Zaldumar (3)
   7  Club Brisanda             11  15  17   1  3  19  0  Lorvin Duelmo (3)
   8  Academia Salinera         11   5   7   6  6  22  2  Galdo Montecal (2)
   9  Juventud Riobravo         11  11  17   1  4  18  2  Joldan Landrosa (2)
  10  Deportivo Nuevaluna       11   7  18   1  6  24  3  Nervio Bocanegro (2)

↑/↓ mover · enter plantilla · tab goleador/asistente · * tu club · esc volver
```

## Comparativa de equipos (tab)

```
EQUIPOS · Temporada 4 · Jornada 11 / 18

   #  Equipo                    PJ  GF  GC Imb SM  Am Ro  Asistente
>  1  Juventud Altamira del Sur 11  20   8   5  2  22  0  Wenceo Korvena (3)
   2  Unión Fontalva            11  15  13   4  1  27  2  Galdo Lacerna (3)
   3  Deportivo Hondaluz        11  15   9   5  3  19  1  Lisandor Bastaro (3)
*  4  Estrella Altamira del Sur 11  13  10   3  4  18  0  Dorvan Orbaneja (3)
   5  Estrella Salinera         11  10  11   5  3  21  2  Mendo Fuentalba (3)
   6  Unión Altamira del Sur    11  12  13   4  3  12  0  Kaelo Zaldumar (3)
   7  Club Brisanda             11  15  17   1  3  19  0  Dorvan Herbasco (3)
   8  Academia Salinera         11   5   7   6  6  22  2  Fabrel Garmendo (1)
   9  Juventud Riobravo         11  11  17   1  4  18  2  Ivaro Cordival (3)
  10  Deportivo Nuevaluna       11   7  18   1  6  24  3  Anselmo Rivasel (3)

↑/↓ mover · enter plantilla · tab goleador/asistente · * tu club · esc volver
```

## Plantilla de un equipo con estadísticas

```
UNIÓN FONTALVA · estadísticas

  Posición      Jugador                 PJ Tit   Min   G   A  Am  Ro Imb  Val
> Portero       Kervo Illanes           11  11   990   0   0   0   0   4  5.9
  Portero       Arlen Fuentalba          0   0     0   0   0   0   0   0    -
  Defensa       Fabrel Frondosa         11  11   901   0   0   2   0   3  6.0
  Defensa       Xandro Gallardel        11  11   865   0   1   4   1   2  5.8
  Defensa       Zenel Estoral           11  11   880   0   1   4   0   4  6.1
  Defensa       Joldan Herbasco         11  11   833   0   0   4   0   3  5.9
  Defensa       Nolven Orbaneja          5   0   144   0   0   0   0   0  5.9
  Defensa       Elmiro Norvalde          8   0   208   0   0   2   0   0  5.9
  Defensa       Kaelo Torvelo            7   0   149   0   0   0   0   0  6.0
  Mediocampista Galdo Lacerna           11  11   923   0   3   2   0   0  6.2
  Mediocampista Xandro Bastaro          11  11   915   3   2   2   0   0  6.4
  Mediocampista Nervio Rivasel          11  11   744   0   1   1   0   0  6.1
  Mediocampista Florian Quintaleo        4   0   101   0   1   1   0   0  6.0
  Mediocampista Kervo Garmendo           3   0    48   0   0   0   0   0  6.0
  Mediocampista Tavio Sandoral           8   0   215   1   0   1   0   0  6.1
  Mediocampista Hemiro Dunavar           1   0    28   0   0   0   0   0  6.1
  Delantero     Florian Bocanegro       11  11   898   3   2   2   0   0  6.4
  Delantero     Ivaro Urdangal          11  11   865   3   1   0   0   0  6.4
  Delantero     Tavio Orbaneja          11  11   830   4   0   2   1   0  6.3
  Delantero     Baltor Lacerna           4   0   106   0   0   0   0   0  6.0
  Delantero     Nervio Garmendo          5   0   106   0   0   0   0   0  6.0
  Delantero     Nervio Calderan          4   0    75   1   0   0   0   0  6.0

↑/↓ mover · enter ver ficha · esc volver
```

## Ficha de un portero

```
FICHA · Mendo Cordival
Estrella Altamira del Sur · Portero · 30 años · Valoración 79
RIT 64  TIR 52  PAS 78  REG 45  DEF 72  FIS 74  REF 84

Temporada 4 (en curso)
  PJ 11 · Tit 11 · Min 990 · G 0 · A 0 · Am 0 · Ro 0 · Imb 3 · Enc 10 · Val 5.9

Trayectoria
  Temp  Club                       PJ   Min   G   A  Am  Ro Imb  Val
     1  Estrella Altamira del Sur  18  1620   0   0   0   0   5  5.6
     2  Estrella Altamira del Sur  18  1620   0   0   0   0   5  5.7
     3  Estrella Altamira del Sur  18  1620   0   0   0   0   5  5.8

Carrera: 65 partidos · 0 goles · 0 asistencias · 0 amarillas · 0 rojas

esc volver
```

## Ficha de un delantero

```
FICHA · Elmiro Montecal
Estrella Altamira del Sur · Delantero · 18 años · Valoración 47
RIT 51  TIR 54  PAS 42  REG 44  DEF 10  FIS 47  REF 1

Temporada 4 (en curso)
  PJ 7 · Tit 0 · Min 204 · G 1 · A 0 · Am 0 · Ro 0 · Imb 0 · Val 6.1

Trayectoria
  Temp  Club                       PJ   Min   G   A  Am  Ro Imb  Val
     3  Estrella Altamira del Sur   8   204   0   0   0   0   0  6.0

Carrera: 15 partidos · 1 goles · 0 asistencias · 0 amarillas · 0 rojas

esc volver
```

## Fin de temporada

```
TEMPORADA 1 TERMINADA

Campeón: Juventud Brisanda
Tu equipo: Juventud Sanvedra, puesto 10 de 10, 17 puntos

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
TEMPORADA 2 · Juventud Sanvedra

Valoración del equipo: 75 → 76

Nadie se retira este año.

No llega nadie de la cantera.

En toda la liga: 9 retiros y 9 juveniles.

enter continuar
```

## Confirmar una carrera nueva

```
NUEVA CARRERA

Esto empieza una carrera nueva y pierdes la actual
(Temporada 2 con Juventud Sanvedra). Aún no hay guardado automático.
¿Seguro?

> No, volver
  Sí, empezar de cero

↑/↓ mover · enter elegir · esc volver
```
