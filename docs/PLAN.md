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
  - La estabilidad de la liga a largo plazo se midió primero solo con la media
    global, y eso resultó insuficiente: ver la calibración, más abajo.
- **6b, menús:** opción "Siguiente temporada" al terminar; pantalla de inicio de
  temporada con las bajas y altas del club y la valoración antes y después;
  pantalla de historial (también en el menú principal, con scroll); número de
  temporada en los encabezados; y confirmación al elegir "Nueva carrera", que
  pierde la carrera actual. E2E de varias temporadas seguidas con el teclado.

- **Calibración por posición y por edad (PR aparte):** al revisar la ficha de un
  portero de 20 años con valoración 87 se vio que la media global estable
  ocultaba dos defectos: una inflación (el 40 % de los porteros con valoración
  >= 85) y una dispersión que crecía. Las causas y su corrección:
  - La liga inicial no seguía la curva de la progresión (jugadores de 17 años con
    la calidad de uno de 30, que luego crecían +30): ahora un jugador inicial se
    crea como juvenil de 16 años y envejece con la propia progresión.
  - Los juveniles de 17 a 19 años entraban sin los años de crecimiento previos,
    lo que bajaba cada tramo de edad unos 4 puntos con las temporadas: ahora
    crecen desde los 16 como todos.
  - El crecimiento no tenía techo y los atributos se apilaban en 99: ahora hay
    rendimientos decrecientes cerca del tope (el crecimiento se frena, la caída
    por edad no).
  - La valoración del portero dependía en un 60 % de los reflejos: ahora 50 %, y
    su ventaja en reflejos baja de +25 a +16.
  - Los porteros se generan con edades de 17 a 39 (juegan tres años más).
  - **Retiro por nivel:** el retiro no depende solo de la edad. Desde los 33 años
    (36 en porteros), un jugador con valoración por debajo de 62 suma 8 puntos
    porcentuales de probabilidad de retiro por cada punto que le falta (un
    jugador de 38 años con 47 se retira seguro). Antes, el 31 % de las plantillas
    tenía un veterano de 35 o más con valoración <= 55 y el 25 % de los de 38
    años estaba en 50 o menos; ahora son el 13 % y ninguno. La liga inicial
    descarta a los jugadores que se habrían retirado antes de llegar a su edad.
  Resultado, medido con 30 ligas: la liga de la temporada 20 se parece a la
  inicial por posición (campo 70,5 y porteros 72,7 de media; élite >= 85 de
  ~4 % y ~10 %; dispersión de 9,5 y 11), y cada tramo de edad vale lo mismo en
  la temporada 1 que en la 20. Las pruebas de `aplicacion/calibracion_test.go`
  miden por posición y por edad, no solo la media.

- **Talento y cracks (PR aparte):** cada jugador tiene un **talento** oculto, un
  multiplicador de su crecimiento (100 es normal; se sortea con una
  distribución log-normal, la misma para todas las posiciones). Es lo que separa
  a los cracks. La ficha muestra una **pista** (la "proyección": limitada,
  normal, alta o excepcional) solo en los jugadores de 24 años o menos, nunca el
  número. Se guarda en la base de datos (migración 4; las partidas anteriores
  quedan con talento neutro). Objetivos y resultado, medidos con 40 ligas:
  - De cada cien jugadores de campo, ~3,6 llegan a 90 o más en su carrera y
    ~0,33 a 95 o más, y se reparten por todas las posiciones (defensas 2,6 %,
    medios 3,9 %, delanteros 4,5 %, porteros 4,7 %; antes el 8 % de los porteros
    y 1-2 % del campo). Para lograrlo, la valoración del portero reparte su peso
    entre más atributos (reflejos 35 %), lo que reduce su varianza.
  - **Un crack pesa más en el partido.** La fuerza de un equipo suma un bono por
    cada jugador cuya valoración supera 85, con un coeficiente por posición y
    un tope de saturación por línea. Un crack de 95 suma unos **+8 puntos de
    victoria** en cada posición (entre +7 y +9 según la posición y la muestra;
    antes +1 a +3); tres cracks, ~+21; un once de cracks gana ~75 % y no el
    100 %. La liga sigue abierta: el campeón suma ~36 de 54
    puntos y el equipo de mayor valoración gana la liga el 24-31 % de las veces
    (antes 19 %; el azar puro sería 10 %).
  - Las pruebas garantizan los topes en ambos sentidos, y que un equipo sin
    cracks se comporta igual que antes.

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
- **9b, pantallas:** opción **Estadísticas** en el menú principal y en el de fin
  de temporada. Las seis clasificaciones (goleadores, asistentes, tarjetas,
  porterías imbatidas, menos goles encajados y valoración), la comparativa de
  equipos (con goleador y máximo asistente de cada club), la plantilla de
  cualquier equipo con estadísticas y la **ficha del jugador** (atributos,
  temporada en curso, trayectoria por temporada y totales de carrera). En
  `Plantilla`, `tab` alterna atributos y estadísticas. Cursor con scroll, pila de
  navegación que vuelve a cada pantalla tal como estaba, y tu club marcado con
  `*`. `aplicacion` gana `EstadisticasDeEquipo`, `EstadisticaDeJugador` y
  `Trayectoria`. Una prueba verifica que ninguna pantalla pasa de 80 columnas.
- **9c, alineación, disponibilidad y cansancio** (se parte en tres PR):
  - **9c-1, alineación:** seis formaciones (4-4-2, 4-3-3, 3-5-2, 4-5-1, 5-3-2,
    3-4-3) con perfiles de ataque y defensa de efecto neto parecido (ninguna
    domina a las demás: lo que decide es la plantilla); once y banquillo
    elegidos por el usuario; **jugar fuera de posición** (el rendimiento en el
    puesto sale de los atributos del jugador más una penalización de
    familiaridad según la distancia entre líneas; un portero solo juega de
    portero y un jugador de campo solo ahí en emergencia); alineación
    automática para los rivales (y para el usuario si la pide); y el efecto de
    las rojas en el marcador (diez hombres marcan menos y reciben más). La
    selección automática es un servicio de dominio con parámetros del
    entrenador, para que el modo jugador pueda reutilizarla (ahí decide la IA).
  - **9c-2, disponibilidad y cansancio:** **sanciones** (una roja cuesta de 1 a
    3 partidos; cada 5 amarillas, 1), **lesiones** (con probabilidad que crece
    con la edad, el mal físico y el cansancio; de 1 a 8 jornadas) y **condición
    física**: cada jugador tiene una condición que baja al jugar (más si su
    físico es bajo o es veterano; esa es su "resistencia", derivada del
    atributo físico) y se recupera descansando, y escala su rendimiento. Quien
    no está disponible no puede alinearse; el juego lo sustituye y avisa.
    Guardado v5.
  - **9c-3, pantallas:** alineación (formación, puestos, banquillo y avisos de
    sanción, lesión y cansancio), con pruebas E2E.
  Limitación conocida: con un partido por jornada el cansancio pesa poco por sí
  solo; gana importancia con más partidos (copas) y con el cansancio dentro del
  partido en la fase 11.
*Terminado cuando*: se puede cambiar la alineación, cambia el rendimiento del
equipo y las estadísticas individuales y de equipos se acumulan, se ven y se
guardan.

### 10. Mercado de fichajes
Paquete `mercado`. Requisitos pedidos por el usuario (el diseño detallado se
decide antes de empezar):
- Lista de **agentes libres**, con filtros, para ficharlos.
- **Búsqueda** de jugadores concretos de otros clubes para intentar comprarlos.
- Posibilidad de **pagar la cláusula de rescisión** de un jugador.
- Un **criterio claro** para que un club no quiera vender a un jugador, y que el
  usuario vea el motivo del rechazo.
- Presupuesto, valor de mercado y contratos (salario, duración, cláusula).
- **Rescindir el contrato** de un jugador (liberarlo), por ejemplo un veterano que
  ya no da el nivel.
- **Que un club no pueda acumular estrellas:** precios que suben de forma
  superlineal con el nivel, tope de masa salarial, cláusulas de rescisión altas
  para los cracks y que un crack solo acepte clubes con prestigio o proyecto.
- Pendiente de decidir: si un modelo de lenguaje local (Ollama) tendría algún
  papel; la recomendación es que las reglas decidan siempre y, como mucho, el
  modelo redacte los mensajes de forma opcional.
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
Reutiliza la alineación automática, la disponibilidad y el cansancio de la fase
9c: el entrenador (la IA) decide si juegas según tu nivel, tu condición y su
tendencia a rotar.
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
