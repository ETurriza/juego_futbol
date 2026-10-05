# Plan del proyecto

## Visión

Juego de fútbol para terminal con **modo carrera de varias temporadas**: la
carrera no termina al acabar una temporada, sigues en el mismo club año tras
año, con plantilla que envejece y se renueva, historial y premios. Hay dos
modos de juego:

- **Técnico** (primero): gestionas el club — alineaciones, fichajes, tácticas —
  y, más adelante, juegas los partidos en tiempo real.
- **Jugador** (después): eres un futbolista; sigues tus estadísticas, aspiras a
  premios individuales y colectivos y negocias tu futuro. El técnico (la IA)
  decide la alineación.

Además habrá un **partido rápido** independiente de la carrera.

La meta final es servirlo por SSH con Wish para que cualquiera juegue con solo
`ssh`, y tener una **liga en línea entre amigos**.

- Lenguaje: Go
- Interfaz: Bubble Tea v2
- Módulo: `github.com/ETurriza/juego_futbol`

## Arquitectura por capas

Inspirada en *Architecture Patterns with Python*: dominio puro en el centro,
casos de uso encima, y adaptadores/entrypoints en el borde. Las capas externas
dependen de las internas, nunca al revés.

```
  menus (TUI)      red (SSH/Wish)        <- entrypoints
        \              /
         aplicacion                      <- casos de uso + puertos (interfaces)
        /     |      \
  simulacion  liga  mercado  partido     <- servicios de dominio
        \     |      /
          modelo                         <- entidades, sin dependencias

  persistencia (SQLite)                  <- adaptador: implementa los puertos
        \
         modelo
```

| Cosmic Python        | Aquí                                                   |
|----------------------|--------------------------------------------------------|
| Modelo de dominio    | `modelo`, más `simulacion`, `liga`, `mercado`, `partido` |
| Capa de servicios    | `aplicacion`                                           |
| Repositorio / puertos| Interfaces definidas en `aplicacion` (donde se consumen) |
| Adaptadores          | `persistencia`                                         |
| Entrypoints          | `menus`, `red`                                         |

Reglas de dependencias:

- `modelo` no importa ningún otro paquete del proyecto.
- `simulacion`, `liga` y `mercado` solo dependen de `modelo` y entre sí.
- `aplicacion` depende de `modelo` y de los servicios de dominio.
- `persistencia` depende de `modelo` y de `aplicacion`: implementa sus puertos y
  usa su tipo `Guardado` (el adaptador, capa externa, depende de la capa interna
  que define el puerto).
- `menus` y `red` dependen de `aplicacion`.

No se crean Unit of Work, repositorios ni paquetes antes de que su fase los
necesite. Si una fase necesita un paquete nuevo (por ejemplo, uno para la
evolución de los jugadores), se decide en su diseño y se registra en el mapa
`permitidos` de la prueba de arquitectura.

## Ramas y flujo de trabajo

```
feat/*  --PR-->  dev  --PR-->  main (prod)
```

- `main` es prod: código estable, solo recibe PR desde `dev`. Es lo que se
  despliega para el servidor SSH.
- `dev` es la integración: recibe PR de las ramas `feat/*` (una por fase o
  tarea), que salen de `dev`.
- Ambas tienen un ruleset: PR obligatorio, sin push directo, sin force-push ni
  borrado, y comprobación de `go build`, `go vet` y `go test`.

## Estructura de carpetas

Cada paquete se crea solo cuando llega su fase. Entre paréntesis, la fase en que
se creó o se creará (ver la tabla de fases).

```
cmd/juego/              punto de entrada (modo local / modo servidor)
internal/
  modelo/               (1) dominio puro
  generador/            (1) jugadores y equipos inventados
  simulacion/           (2)
  arquitectura/         (2) prueba de reglas de dependencias entre paquetes
  liga/                 (3)
  progresion/           (6) edad, evolución y retiro de los jugadores
  aplicacion/           (4) casos de uso + puertos
    contrato/           (5) pruebas de contrato de los puertos
  menus/                (4)
  persistencia/         (5)
  mercado/              (10)
  partido/              (11)
  red/                  (12)
tests/e2e/              (4) en adelante
docs/                   documentación (este plan)
```

Los paquetes con reglas propias pueden tener su propio `CLAUDE.md`, creado
junto con el paquete en su fase.

## Estrategia de pruebas

- **Unit**: archivos `*_test.go` junto al código. Rápidos, con semilla fija, sin
  disco ni red.
- **Integration**: junto al código, en `*_integration_test.go` con
  `//go:build integration`. Se corren con `go test -tags=integration ./...`.
  Aplican sobre todo a `persistencia` (SQLite real en directorio temporal) y a
  `aplicacion` con repositorio real.
- **E2E**: en `tests/e2e/`. Desde la fase 4, con un `tea.Program` real
  alimentado con bytes de teclado (no se usa `teatest`: no está etiquetado para
  Bubble Tea v2); desde la fase 12, levantando el servidor SSH y conectando con
  un cliente.

`go test ./...` ejecuta solo las pruebas unitarias.

## Fases

El orden cambió respecto al plan original para poner la jugabilidad del modo
técnico antes del servidor SSH, y para que la carrera de varias temporadas y las
estadísticas estén listas antes de conectar la interfaz con el guardado.

| Fase | Contenido                                   | Estado    | Antes     |
|------|---------------------------------------------|-----------|-----------|
| 1    | Modelo y generador                          | hecha     | 1         |
| 2    | Simulación de partidos                      | hecha     | 2         |
| 3    | Liga                                        | hecha     | 3         |
| 4    | Menús locales con Bubble Tea                | hecha     | 4         |
| 5    | Guardado: backend con SQLite                | hecha     | 5a        |
| 6    | Temporadas continuas                        | hecha     | (nueva)   |
| 9    | Estadísticas individuales y alineación      | en curso  | (nueva)   |
| 7    | Guardado en la interfaz                     | pendiente | 5b        |
| 8    | Partido rápido                              | pendiente | (nueva)   |
| 10   | Mercado de fichajes                         | pendiente | 7         |
| 11   | Partido jugable en tiempo real              | pendiente | 8         |
| 12   | Servidor SSH con Wish                       | pendiente | 6         |
| 13   | Modo jugador                                | pendiente | (nueva)   |
| 14   | Estilo visual                               | pendiente | (nueva)   |

La tabla va en orden de ejecución: la fase 9 se hace antes que la 7 y la 8,
porque cambia el `Guardado` y la base de datos y conviene migrarlos antes de que
existan partidas reales guardadas.

### 1. Modelo y generador (hecha)
Jugador, atributos y equipo en `modelo`. Generador de jugadores con nombres
inventados en `generador`, con aleatoriedad recibida como `*rand.Rand`.

### 2. Simulación de partidos (hecha)
Simulación estadística de partidos a partir de atributos (`simulacion`), y una
prueba que verifica los imports de cada paquete contra las reglas de
dependencias.

### 3. Liga (hecha)
Calendario todos contra todos, ida y vuelta (la segunda vuelta invierte las
localías); con un número impar de equipos se agrega un descanso por jornada.
Tabla ordenada por puntos, diferencia de goles, goles a favor y nombre.

### 4. Menús locales con Bubble Tea (hecha)
- **4a, `aplicacion`:** `Carrera` (liga de 10 equipos inventados y un equipo del
  usuario), `AvanzarJornada`, `Tabla`, `Plantilla`, `Campeon`. Cada jornada usa
  un `*rand.Rand` derivado de la semilla y del número de jornada, de modo que
  una carrera queda determinada por su semilla y su avance.
- **4b, `menus` y `cmd/juego`:** menú principal, plantilla, tabla, avanzar
  jornada y fin de temporada; banderas `--semilla` y `--equipos`; pruebas de
  `Update`/`View` y primeras pruebas E2E.

### 5. Guardado: backend con SQLite (hecha)
Se guarda el estado real de la carrera (equipos, jugadores, resultados), no solo
la semilla, porque los fichajes y la evolución de los jugadores cambian las
plantillas. En `aplicacion`: `Guardado` (foto de la carrera con tipos simples),
`Exportar`/`Importar` con validación, el puerto `RepositorioPartidas` (por
*ranura*: una en local, una por usuario en SSH), `GuardarCarrera`/`CargarCarrera`
y un repositorio en memoria, con pruebas de contrato reutilizables. En
`persistencia`: adaptador SQLite (`modernc.org/sqlite`, sin CGO), migraciones por
`PRAGMA user_version`, guardado transaccional y pruebas de integración.

### 6. Temporadas continuas
Al terminar una temporada la carrera sigue: acción "Siguiente temporada" con
nuevo calendario y tabla a cero, en el mismo club. Se entrega en dos PR:

- **6a, dominio y guardado:**
  - `progresion` (solo depende de `modelo`): evolución anual y retiro por edad.
    Los atributos técnicos tienen una meseta larga (de los 28 a los 33, con un
    20 % de probabilidad anual de "gran año" entre los 31 y los 34) y el
    declive es rápido desde los 34; el físico baja antes (desde los 28). Los
    porteros envejecen tres años más tarde. No hay retiros antes de los 34;
    después la probabilidad sube hasta el 100 % a los 39 (42 en porteros). Las
    constantes de calibración están al inicio de `progresion.go`.
  - `generador.Juvenil` y `generador.MinimoPorPosicion`: los retiros se reponen
    con juveniles de 16 a 19 años de la misma posición, hasta la composición
    mínima de 2 porteros, 7 defensas, 7 mediocampistas y 6 delanteros (nunca se
    recorta una plantilla que ya la supera).
  - `Carrera` pasa a llevar el número de temporada, un historial (campeón,
    puesto y puntos del usuario de cada año) y `ProximoID` (una ID de jugador
    nunca se reutiliza). `SiguienteTemporada()` es atómica y devuelve las bajas
    y altas del club del usuario.
  - La semilla de la temporada 1 es la de la carrera; las siguientes derivan la
    suya, y la evolución de las plantillas usa otra semilla derivada.
  - `Guardado` v2 y migración 2 de la base de datos (las partidas existentes
    quedan en la temporada 1).
  - Se verificó con pruebas que la liga es estable a largo plazo: promediando
    12 semillas, la valoración media va de 70,8 a 71,1 en 40 temporadas y la
    edad media se queda en 27.
- **6b, menús:** opción "Siguiente temporada" al terminar; pantalla de inicio de
  temporada con las bajas y altas del club y la valoración antes y después;
  pantalla de historial (también en el menú principal, con scroll); número de
  temporada en los encabezados; y confirmación al elegir "Nueva carrera", que
  pierde la carrera actual. E2E de varias temporadas seguidas con el teclado.

*Terminado cuando*: se juegan varias temporadas seguidas desde la terminal, el
historial es correcto y la carrera guardada se reanuda de forma idéntica.

### 7. Guardado en la interfaz
En `menus` y `cmd/juego`:
- Autoguardado tras cada jornada y al crear una carrera.
- Pantalla de inicio con "Continuar" / "Nueva carrera" si hay partida guardada;
  "Nueva carrera" pide confirmación porque reemplaza la guardada.
- Si el guardado falla, aviso en pantalla sin cerrar el juego.
- Banderas `--db` (por defecto en el directorio de datos del usuario) y
  `--ranura` (nombre de la partida, por defecto `principal`).
*Terminado cuando*: se juega, se cierra y se reanuda la carrera, con E2E.

### 8. Partido rápido
Opción en el menú principal, independiente de la carrera y del guardado: elegir
dos equipos, simular el partido y ver el resultado. Cuando exista la fase 11
ofrecerá también "jugar" el partido.
*Terminado cuando*: se simula un partido rápido desde el menú, con pruebas.

### 9. Estadísticas individuales y alineación
Base común de los dos modos. La simulación ya no devuelve solo un marcador:
genera una línea de tiempo con los sucesos de cada partido. Se entrega en tres PR:

- **9a, sucesos y estadísticas (dominio y guardado):**
  - `modelo`: `Evento`, `DetallePartido` (alineaciones titulares y sucesos, de
    donde se derivan los minutos de cada jugador) y `Estadisticas`.
  - `simulacion`: para cada equipo, un once automático (portero, 4 defensas, 3
    medios y 3 delanteros), de 3 a 5 **sustituciones automáticas** entre los
    minutos 46 y 85 (los porteros no se cambian), tarjetas (unas 3,4 amarillas
    por partido, 0,14 rojas con las segundas amarillas) y los goles con sus
    asistencias (el 70 % tiene asistente), repartidos entre quienes están en el
    campo en cada minuto: seis de cada diez goles son de delanteros.
  - `liga`: guarda el detalle de cada partido y calcula, a partir de él, las
    estadísticas de cada jugador (partidos, titularidades, minutos, goles,
    asistencias, amarillas, rojas, porterías imbatidas, goles encajados y
    valoración del partido de 1 a 10) y de cada equipo (resultados, goles,
    imbatidas, partidos sin marcar y tarjetas).
  - `aplicacion`: consultas de jugadores y equipos (con goleador y máximo
    asistente de cada club), clasificaciones (goleadores, asistentes,
    tarjetas, porterías imbatidas, goles encajados por partido y valoración) y
    estadísticas de carrera. Al terminar cada temporada se **archivan las
    estadísticas de toda la liga** por jugador.
  - `Guardado` v3 y migración 3 de la base de datos. Los partidos guardados
    antes no tienen detalle: cuentan en la tabla pero no en las estadísticas.
  - Limitaciones: los porteros no reciben rojas directas; una roja reduce los
    minutos del expulsado pero no cambia el marcador; no hay lesiones.
  - Como la simulación consume más aleatoriedad, una semilla da una liga
    distinta a la de antes de esta fase.
- **9b, pantallas:** estadísticas de jugadores y de equipos, clasificaciones y
  estadísticas de carrera en `menus`.
- **9c, alineación elegible:** formación y once titulares elegidos por el
  usuario, sustituciones elegidas, y el efecto de jugar con diez hombres en el
  marcador.
*Terminado cuando*: se puede cambiar la alineación, cambia el rendimiento del
equipo y las estadísticas individuales y de equipos se acumulan, se ven y se
guardan.

### 10. Mercado de fichajes
Paquete `mercado`: compra, venta y valoración de jugadores, con presupuesto.
*Terminado cuando*: se puede fichar y vender dentro de la carrera.

### 11. Partido jugable en tiempo real
Paquete `partido`: partido en tiempo real que usa la alineación y los atributos
de los jugadores. El resultado alimenta la liga y las estadísticas.
*Terminado cuando*: un partido se juega de principio a fin en la terminal.
Con esta fase queda completo el **modo técnico**.

### 12. Servidor SSH con Wish
Paquete `red`: la misma TUI servida por SSH, con una ranura de guardado por
usuario. Pruebas E2E con cliente SSH. Se decide aquí el diseño de la liga en
línea entre amigos.
*Terminado cuando*: `ssh` al servidor abre el juego y varias sesiones conviven.

### 13. Modo jugador
Al empezar una carrera se elige el modo. En el modo jugador controlas a un
futbolista: sigues tus estadísticas, aspiras a premios individuales (goleador,
mejor jugador) y colectivos, y negocias contratos; el técnico (la IA) decide la
alineación, así que puedes no ser titular.
*Terminado cuando*: se juega una carrera completa como futbolista, con premios y
estadísticas personales.

### 14. Estilo visual
Menús más llamativos: paleta de colores, banner, barras de valoración y
atributos, filas alternadas, insignias de resultado, trofeo y podio. Se hace una
maqueta en texto antes de escribir código.
*Terminado cuando*: el usuario aprueba el resultado en su terminal.

## Convenciones transversales

- Nombres de jugadores y equipos siempre inventados; nada de marcas ni personas
  reales.
- Toda aleatoriedad entra como `*rand.Rand`, nunca global.
- Todo cambio en lógica incluye pruebas.
