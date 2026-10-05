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

Deportivo Valmora    Temporada 4 · Jornada 11 / 18
Valoración del equipo: 83

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

> Deportivo Valmora           0 - 0   Club Brisanda
  Unión Olivarejo             1 - 2   Estrella Olivarejo
  Deportivo Embalse Nuevo     1 - 4   Academia Castelmar
  Academia Lomaverde          0 - 2   Academia Riobravo
  Academia Altamira del Sur   3 - 1   Estrella Salinera

enter continuar
```

## Tabla de posiciones

```
TABLA · Temporada 4 · Jornada 11 / 18

    #  Equipo                      PJ   G   E   P   GF   GC   DG  Pts
>   1  Deportivo Valmora           11   6   4   1   17    9   +8   22
    2  Academia Castelmar          11   6   3   2   18    7  +11   21
    3  Academia Altamira del Sur   11   5   5   1   20    8  +12   20
    4  Estrella Olivarejo          11   6   2   3   15   13   +2   20
    5  Estrella Salinera           11   5   3   3   11   11   +0   18
    6  Deportivo Embalse Nuevo     11   5   0   6    7   13   -6   15
    7  Club Brisanda               11   4   2   5   14   17   -3   14
    8  Academia Lomaverde          11   3   1   7    9   15   -6   10
    9  Academia Riobravo           11   2   2   7   10   18   -8    8
   10  Unión Olivarejo             11   2   0   9    9   19  -10    6

esc volver
```

## Plantilla: atributos

```
PLANTILLA · Deportivo Valmora

  Posición       Nombre                  Ed  Val  RIT TIR PAS REG DEF FIS REF
  Portero        Gervasio Quintaleo      20   87   58  37  69  45  64  76  99
  Portero        Zenel Garmendo          29   72   44  33  50  37  47  62  84
  Defensa        Anselmo Orbaneja        29   81   73  50  68  54  98  76  15
  Defensa        Marvel Zaldumar         24   81   68  49  66  51  94  85  19
  Defensa        Nolven Aguirel          30   77   69  50  58  55  93  75  17
  Defensa        Kaelo Calderan          22   76   68  45  59  57  89  76   7
  Defensa        Daveo Escalona          30   75   60  51  58  60  91  74  16
  Defensa        Baltor Pedrosel         32   66   57  37  46  36  80  69   7
  Defensa        Tavio Dunavar           35   51   39  25  47  29  66  46  10
  Mediocampista  Mendo Olmedrar          21   93   84  85  99  99  84  87  16
  Mediocampista  Florian Abrelda         26   89   83  83  97  88  85  79  17
  Mediocampista  Corvino Jarvale         28   78   73  64  88  83  73  68   8
  Mediocampista  Elmiro Cordival         24   77   71  72  83  82  67  65   9
  Mediocampista  Quirino Norvalde        33   72   59  64  80  77  63  65  20
  Mediocampista  Orlandi Urdangal        27   67   62  63  76  65  55  61   4
  Mediocampista  Lorvin Norvalde         24   64   55  56  76  67  49  52   1
  Delantero      Anselmo Fuentalba       21   92   91  99  83  97  54  91  23
  Delantero      Elmiro Nebreda          28   81   83  95  71  78  49  63  21
  Delantero      Quirino Valmeron        26   72   70  82  62  79  34  58  12
  Delantero      Galdo Garmendo          29   72   70  83  59  74  35  67  13
  Delantero      Lisandor Yrigaldo       33   71   75  80  69  69  39  58  13
  Delantero      Valdor Cordival         19   63   61  81  47  56  25  49   1

tab estadísticas · esc volver
```

## Plantilla: estadísticas (tab)

```
PLANTILLA · Deportivo Valmora · estadísticas

  Posición      Jugador                 PJ Tit   Min   G   A  Am  Ro Imb  Val
  Portero       Gervasio Quintaleo      11  11   990   0   0   0   0   5  6.1
  Portero       Zenel Garmendo           0   0     0   0   0   0   0   0    -
  Defensa       Anselmo Orbaneja        11  11   904   0   0   3   0   3  6.0
  Defensa       Marvel Zaldumar         11  11   940   0   0   3   1   4  5.9
  Defensa       Nolven Aguirel          11  11   887   0   2   2   0   4  6.3
  Defensa       Kaelo Calderan          11  11   933   0   0   1   0   5  6.2
  Defensa       Daveo Escalona           8   0   172   0   0   0   0   0  6.0
  Defensa       Baltor Pedrosel          5   0   147   0   0   1   0   0  5.9
  Defensa       Tavio Dunavar            5   0    73   0   0   0   0   0  6.0
  Mediocampista Mendo Olmedrar          11  11   910   3   3   2   0   0  6.5
  Mediocampista Florian Abrelda         11  11   779   3   2   1   0   0  6.5
  Mediocampista Corvino Jarvale         11  11   797   1   1   1   0   0  6.2
  Mediocampista Elmiro Cordival          4   0    88   0   0   1   0   0  6.0
  Mediocampista Quirino Norvalde         5   0   139   0   1   1   0   0  6.1
  Mediocampista Orlandi Urdangal         6   0   176   1   0   0   0   0  6.1
  Mediocampista Lorvin Norvalde          3   0    58   0   1   0   0   0  6.1
  Delantero     Anselmo Fuentalba       11  11   953   3   3   2   0   0  6.5
  Delantero     Elmiro Nebreda          11  11   779   3   1   0   0   0  6.5
  Delantero     Quirino Valmeron        11  11   857   3   1   1   0   0  6.4
  Delantero     Galdo Garmendo           5   0   115   0   0   0   0   0  6.1
  Delantero     Lisandor Yrigaldo        3   0    81   0   0   0   0   0  6.0
  Delantero     Valdor Cordival          3   0   105   0   0   0   0   0  6.1

tab atributos · esc volver
```

## Historial de temporadas

```
HISTORIAL · Deportivo Valmora

  Temp  Campeón                     Tu puesto   Pts
     1  Academia Altamira del Sur     9 de 10    19
     2  Academia Altamira del Sur     6 de 10    23
     3  Academia Lomaverde            2 de 10    32

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
>  1  Pelayo Abrelda         Estrella Salinera          11     6     0   801
   2  Hilaro Valmeron        Academia Lomaverde         11     5     0   810
   3  Nolven Quintaleo       Academia Castelmar         11     4     3   951
   4  Rodvan Hontoria        Academia Altamira del Sur  11     4     0   858
   5  Rodvan Calderan        Academia Altamira del Sur  11     4     0   883
   6  Zenel Rivasel          Academia Altamira del Sur  11     3     3   768
*  7  Mendo Olmedrar         Deportivo Valmora          11     3     3   910
   8  Ivaro Aguirel          Academia Castelmar         11     3     3   912
*  9  Anselmo Fuentalba      Deportivo Valmora          11     3     3   953
* 10  Florian Abrelda        Deportivo Valmora          11     3     2   779
  11  Pelayo Hontoria        Estrella Olivarejo         11     3     2   878
  12  Galdo Dunavar          Club Brisanda              11     3     2   948
* 13  Elmiro Nebreda         Deportivo Valmora          11     3     1   779
  14  Quirino Olmedrar       Estrella Olivarejo         11     3     1   784
  15  Cedrio Illanes         Academia Castelmar         11     3     1   829
  16  Hemiro Bocanegro       Academia Castelmar         11     3     1   851
* 17  Quirino Valmeron       Deportivo Valmora          11     3     1   857
  18  Wenceo Landrosa        Academia Riobravo          11     3     0   790
  19  Xandro Valmeron        Unión Olivarejo            11     3     0   801
  20  Gervasio Jarvale       Club Brisanda              11     3     0   916
  21  Yeraldo Landrosa       Academia Riobravo          11     2     4   791
  22  Joldan Yrigaldo        Estrella Olivarejo         11     2     3   911
  23  Xandro Fuentalba       Club Brisanda              11     2     2   836
  24  Yeraldo Quintaleo      Academia Altamira del Sur  11     2     2   943
  25  Quirino Cordival       Academia Lomaverde         11     2     2   973

↑/↓ mover (1-25 de 50) · enter ver ficha · * tu club · esc volver
```

## Clasificación: ASISTENTES

```
ASISTENTES · Temporada 4 · Jornada 11 / 18

   #  Jugador                Equipo                     PJ Asist Goles   Min
>  1  Cedrio Hontoria        Club Brisanda              11     6     1   914
   2  Fabrel Korvena         Unión Olivarejo            11     5     0   896
   3  Yeraldo Landrosa       Academia Riobravo          11     4     2   791
   4  Nolven Quintaleo       Academia Castelmar         11     3     4   951
   5  Zenel Rivasel          Academia Altamira del Sur  11     3     3   768
*  6  Mendo Olmedrar         Deportivo Valmora          11     3     3   910
   7  Ivaro Aguirel          Academia Castelmar         11     3     3   912
*  8  Anselmo Fuentalba      Deportivo Valmora          11     3     3   953
   9  Joldan Yrigaldo        Estrella Olivarejo         11     3     2   911
  10  Ulmar Frondosa         Estrella Salinera          11     3     0   750
* 11  Florian Abrelda        Deportivo Valmora          11     2     3   779
  12  Pelayo Hontoria        Estrella Olivarejo         11     2     3   878
  13  Galdo Dunavar          Club Brisanda              11     2     3   948
  14  Xandro Fuentalba       Club Brisanda              11     2     2   836
  15  Yeraldo Quintaleo      Academia Altamira del Sur  11     2     2   943
  16  Quirino Cordival       Academia Lomaverde         11     2     2   973
  17  Lisandor Hontoria      Estrella Salinera          11     2     1   823
  18  Ivaro Norvalde         Academia Altamira del Sur  11     2     1   826
  19  Ulmar Jarvale          Estrella Salinera          11     2     1   836
  20  Wenceo Torvelo         Estrella Olivarejo         11     2     1   893
  21  Zenel Duelmo           Academia Altamira del Sur  11     2     1   956
* 22  Nolven Aguirel         Deportivo Valmora          11     2     0   887
  23  Quirino Escalona       Academia Lomaverde         11     2     0   962
  24  Tavio Duelmo           Academia Castelmar         11     2     0   967
* 25  Elmiro Nebreda         Deportivo Valmora          11     1     3   779

↑/↓ mover (1-25 de 50) · enter ver ficha · * tu club · esc volver
```

## Clasificación: TARJETAS

```
TARJETAS · Temporada 4 · Jornada 11 / 18

   #  Jugador                Equipo                     PJ  Amar  Roja   Pts
>  1  Fabrel Korvena         Unión Olivarejo            11     4     1     7
   2  Quirino Korvena        Academia Lomaverde         11     4     1     7
   3  Hemiro Bastaro         Unión Olivarejo            11     4     1     7
   4  Lorvin Korvena         Academia Riobravo          11     4     1     7
   5  Mendo Maldovar         Estrella Olivarejo         11     3     1     6
*  6  Marvel Zaldumar        Deportivo Valmora          11     3     1     6
   7  Baltor Montecal        Club Brisanda              11     2     1     5
   8  Mendo Estoral          Estrella Salinera          11     2     1     5
   9  Rodvan Calderan        Deportivo Embalse Nuevo    11     2     1     5
  10  Xandro Zaldumar        Estrella Salinera          11     5     0     5
  11  Quirino Escalona       Academia Lomaverde         11     5     0     5
  12  Kervo Pardovan         Estrella Salinera           3     1     1     4
  13  Yeraldo Sandoral       Unión Olivarejo            11     4     0     4
  14  Ivaro Aguirel          Academia Castelmar         11     4     0     4
  15  Zenel Duelmo           Academia Altamira del Sur  11     4     0     4
  16  Mendo Roblemar         Academia Altamira del Sur  11     4     0     4
  17  Isidoro Quintaleo      Deportivo Embalse Nuevo     6     0     1     3
  18  Dorvan Fuentalba       Academia Riobravo          11     0     1     3
  19  Xandro Valmeron        Unión Olivarejo            11     3     0     3
  20  Ivaro Korvena          Academia Castelmar         11     3     0     3
  21  Ivaro Hontoria         Unión Olivarejo            11     3     0     3
  22  Anselmo Bocanegro      Club Brisanda              11     3     0     3
  23  Xandro Fuentalba       Club Brisanda              11     3     0     3
  24  Xandro Jarvale         Academia Altamira del Sur  11     3     0     3
  25  Rodvan Hontoria        Academia Altamira del Sur  11     3     0     3

↑/↓ mover (1-25 de 50) · enter ver ficha · * tu club · esc volver
```

## Clasificación: PORTEROS · PORTERÍAS IMBATIDAS

```
PORTEROS · PORTERÍAS IMBATIDAS · Temporada 4 · Jornada 11 / 18

   #  Jugador                Equipo                     PJ   Imb   Enc
>  1  Yeraldo Yrigaldo       Academia Castelmar         11     6     7
   2  Nolven Gallardel       Estrella Salinera          11     6    11
   3  Hilaro Irueta          Academia Altamira del Sur  11     5     8
*  4  Gervasio Quintaleo     Deportivo Valmora          11     5     9
   5  Rodvan Calderan        Deportivo Embalse Nuevo    11     3    12
   6  Dorvan Illanes         Estrella Olivarejo         11     3    13
   7  Nervio Orbaneja        Academia Lomaverde         11     3    15
   8  Joldan Fuentalba       Club Brisanda              11     2    17
   9  Hilaro Torvelo         Unión Olivarejo            11     2    19
  10  Lisandor Torvelo       Academia Riobravo          11     1    18

↑/↓ mover · enter ver ficha · * tu club · esc volver
```

## Clasificación: PORTEROS · MENOS GOLES ENCAJADOS

```
PORTEROS · MENOS GOLES ENCAJADOS · Temporada 4 · Jornada 11 / 18

   #  Jugador                Equipo                     PJ   Enc Enc/PJ
>  1  Yeraldo Yrigaldo       Academia Castelmar         11     7   0.64
   2  Hilaro Irueta          Academia Altamira del Sur  11     8   0.73
*  3  Gervasio Quintaleo     Deportivo Valmora          11     9   0.82
   4  Nolven Gallardel       Estrella Salinera          11    11   1.00
   5  Rodvan Calderan        Deportivo Embalse Nuevo    11    12   1.09
   6  Dorvan Illanes         Estrella Olivarejo         11    13   1.18
   7  Nervio Orbaneja        Academia Lomaverde         11    15   1.36
   8  Joldan Fuentalba       Club Brisanda              11    17   1.55
   9  Lisandor Torvelo       Academia Riobravo          11    18   1.64
  10  Hilaro Torvelo         Unión Olivarejo            11    19   1.73

↑/↓ mover · enter ver ficha · * tu club · esc volver
```

## Clasificación: MEJORES VALORACIONES

```
MEJORES VALORACIONES · Temporada 4 · Jornada 11 / 18

   #  Jugador                Equipo                     PJ   Val Goles Asist
>  1  Nolven Quintaleo       Academia Castelmar         11   6.6     4     3
   2  Pelayo Abrelda         Estrella Salinera          11   6.6     6     0
*  3  Anselmo Fuentalba      Deportivo Valmora          11   6.5     3     3
*  4  Mendo Olmedrar         Deportivo Valmora          11   6.5     3     3
   5  Isidoro Jurado         Academia Castelmar         11   6.5     2     1
*  6  Florian Abrelda        Deportivo Valmora          11   6.5     3     2
   7  Zenel Rivasel          Academia Altamira del Sur  11   6.5     3     3
*  8  Elmiro Nebreda         Deportivo Valmora          11   6.5     3     1
   9  Hemiro Bocanegro       Academia Castelmar         11   6.4     3     1
  10  Rodvan Calderan        Academia Altamira del Sur  11   6.4     4     0
  11  Ivaro Aguirel          Academia Castelmar         11   6.4     3     3
  12  Cedrio Illanes         Academia Castelmar         11   6.4     3     1
* 13  Quirino Valmeron       Deportivo Valmora          11   6.4     3     1
  14  Pelayo Hontoria        Estrella Olivarejo         11   6.4     3     2
  15  Quirino Olmedrar       Estrella Olivarejo         11   6.4     3     1
  16  Hilaro Valmeron        Academia Lomaverde         11   6.4     5     0
  17  Joldan Yrigaldo        Estrella Olivarejo         11   6.4     2     3
  18  Yeraldo Quintaleo      Academia Altamira del Sur  11   6.4     2     2
  19  Cedrio Hontoria        Club Brisanda              11   6.4     1     6
  20  Rodvan Hontoria        Academia Altamira del Sur  11   6.4     4     0
  21  Tavio Duelmo           Academia Castelmar         11   6.3     0     2
  22  Ivaro Norvalde         Academia Altamira del Sur  11   6.3     1     2
  23  Kaelo Gallardel        Academia Castelmar         11   6.3     2     0
  24  Galdo Dunavar          Club Brisanda              11   6.3     3     2
  25  Cedrio Quintaleo       Academia Castelmar         11   6.3     0     0

↑/↓ mover (1-25 de 50) · enter ver ficha · * tu club · esc volver
```

## Comparativa de equipos

```
EQUIPOS · Temporada 4 · Jornada 11 / 18

   #  Equipo                    PJ  GF  GC Imb SM  Am Ro  Goleador
>  1  Deportivo Valmora         11  17   9   5  3  19  1  Anselmo Fuentalba (3)
   2  Academia Castelmar        11  18   7   6  4  16  0  Nolven Quintaleo (4)
   3  Academia Altamira del Sur 11  20   8   5  2  27  0  Rodvan Calderan (4)
   4  Estrella Olivarejo        11  15  13   3  1  15  1  Pelayo Hontoria (3)
   5  Estrella Salinera         11  11  11   6  2  22  2  Pelayo Abrelda (6)
   6  Deportivo Embalse Nuevo   11   7  13   3  6  17  2  Baltor Sandoral (2)
   7  Club Brisanda             11  14  17   2  2  18  1  Galdo Dunavar (3)
   8  Academia Lomaverde        11   9  15   3  5  18  1  Hilaro Valmeron (5)
   9  Academia Riobravo         11  10  18   1  5  19  2  Wenceo Landrosa (3)
  10  Unión Olivarejo           11   9  19   2  6  25  2  Xandro Valmeron (3)

↑/↓ mover · enter plantilla · tab goleador/asistente · * tu club · esc volver
```

## Comparativa de equipos (tab)

```
EQUIPOS · Temporada 4 · Jornada 11 / 18

   #  Equipo                    PJ  GF  GC Imb SM  Am Ro  Asistente
>  1  Deportivo Valmora         11  17   9   5  3  19  1  Anselmo Fuentalba (3)
   2  Academia Castelmar        11  18   7   6  4  16  0  Ivaro Aguirel (3)
   3  Academia Altamira del Sur 11  20   8   5  2  27  0  Zenel Rivasel (3)
   4  Estrella Olivarejo        11  15  13   3  1  15  1  Joldan Yrigaldo (3)
   5  Estrella Salinera         11  11  11   6  2  22  2  Ulmar Frondosa (3)
   6  Deportivo Embalse Nuevo   11   7  13   3  6  17  2  Baltor Sandoral (1)
   7  Club Brisanda             11  14  17   2  2  18  1  Cedrio Hontoria (6)
   8  Academia Lomaverde        11   9  15   3  5  18  1  Quirino Cordival (2)
   9  Academia Riobravo         11  10  18   1  5  19  2  Yeraldo Landrosa (4)
  10  Unión Olivarejo           11   9  19   2  6  25  2  Fabrel Korvena (5)

↑/↓ mover · enter plantilla · tab goleador/asistente · * tu club · esc volver
```

## Plantilla de un equipo con estadísticas

```
ACADEMIA CASTELMAR · estadísticas

  Posición      Jugador                 PJ Tit   Min   G   A  Am  Ro Imb  Val
> Portero       Yeraldo Yrigaldo        11  11   990   0   0   0   0   6  6.2
  Portero       Tavio Garmendo           0   0     0   0   0   0   0   0    -
  Defensa       Tavio Duelmo            11  11   967   0   2   1   0   6  6.3
  Defensa       Cedrio Quintaleo        11  11   874   0   0   0   0   6  6.3
  Defensa       Ivaro Korvena           11  11   814   0   1   3   0   5  6.2
  Defensa       Isidoro Jurado          11  11   911   2   1   0   0   6  6.5
  Defensa       Hemiro Valmeron          5   0   102   0   0   0   0   0  6.0
  Defensa       Daveo Irueta             7   0   122   0   0   0   0   0  6.0
  Defensa       Eskardo Estoral          7   0   204   0   0   1   0   0  6.0
  Mediocampista Hilaro Gallardel        11  11   779   1   0   2   0   0  6.1
  Mediocampista Ivaro Aguirel           11  11   912   3   3   4   0   0  6.4
  Mediocampista Nolven Quintaleo        11  11   951   4   3   2   0   0  6.6
  Mediocampista Belron Escalona          5   0    87   0   0   0   0   0  6.0
  Mediocampista Kaelo Lacerna            3   0    87   0   0   0   0   0  6.1
  Mediocampista Pelayo Landrosa          6   0   153   0   0   2   0   0  5.9
  Mediocampista Mendo Herbasco           2   0    34   0   0   0   0   0  6.0
  Delantero     Cedrio Illanes          11  11   829   3   1   0   0   0  6.4
  Delantero     Hemiro Bocanegro        11  11   851   3   1   0   0   0  6.4
  Delantero     Kaelo Gallardel         11  11   862   2   0   1   0   0  6.3
  Delantero     Galdo Escalona           4   0   128   0   0   0   0   0  6.0
  Delantero     Galdo Landrosa           7   0   146   0   0   0   0   0  6.0
  Delantero     Nolven Olmedrar          3   0    87   0   0   0   0   0  6.0

↑/↓ mover · enter ver ficha · esc volver
```

## Ficha de un portero

```
FICHA · Gervasio Quintaleo
Deportivo Valmora · Portero · 20 años · Valoración 87
RIT 58  TIR 37  PAS 69  REG 45  DEF 64  FIS 76  REF 99

Temporada 4 (en curso)
  PJ 11 · Tit 11 · Min 990 · G 0 · A 0 · Am 0 · Ro 0 · Imb 5 · Enc 9 · Val 6.1

Trayectoria
  Temp  Club                       PJ   Min   G   A  Am  Ro Imb  Val
     1  Deportivo Valmora          18  1620   0   0   0   0   3  5.6
     2  Deportivo Valmora          18  1620   0   0   0   0   6  5.8
     3  Deportivo Valmora          18  1620   0   0   1   0   5  5.9

Carrera: 65 partidos · 0 goles · 0 asistencias · 1 amarillas · 0 rojas

esc volver
```

## Ficha de un delantero

```
FICHA · Valdor Cordival
Deportivo Valmora · Delantero · 19 años · Valoración 63
RIT 61  TIR 81  PAS 47  REG 56  DEF 25  FIS 49  REF 1

Temporada 4 (en curso)
  PJ 3 · Tit 0 · Min 105 · G 0 · A 0 · Am 0 · Ro 0 · Imb 0 · Val 6.1

Trayectoria
  Temp  Club                       PJ   Min   G   A  Am  Ro Imb  Val
     2  Deportivo Valmora           5   100   0   0   0   0   0  6.0
     3  Deportivo Valmora           9   252   1   1   1   0   0  6.1

Carrera: 17 partidos · 1 goles · 1 asistencias · 1 amarillas · 0 rojas

esc volver
```

## Fin de temporada

```
TEMPORADA 1 TERMINADA

Campeón: Unión Embalse Nuevo
Tu equipo: Club Pinarel, puesto 2 de 10, 34 puntos

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
TEMPORADA 2 · Club Pinarel

Valoración del equipo: 79 → 79

Nadie se retira este año.

No llega nadie de la cantera.

En toda la liga: 5 retiros y 5 juveniles.

enter continuar
```

## Confirmar una carrera nueva

```
NUEVA CARRERA

Esto empieza una carrera nueva y pierdes la actual
(Temporada 2 con Club Pinarel). Aún no hay guardado automático.
¿Seguro?

> No, volver
  Sí, empezar de cero

↑/↓ mover · enter elegir · esc volver
```
