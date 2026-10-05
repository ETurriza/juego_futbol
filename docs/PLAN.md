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
técnico antes del servidor SSH, y para que la carrera de varias temporadas esté
lista antes de conectar la interfaz con el guardado.

| Fase | Contenido                                   | Estado    | Antes     |
|------|---------------------------------------------|-----------|-----------|
| 1    | Modelo y generador                          | hecha     | 1         |
| 2    | Simulación de partidos                      | hecha     | 2         |
| 3    | Liga                                        | hecha     | 3         |
| 4    | Menús locales con Bubble Tea                | hecha     | 4         |
| 5    | Guardado: backend con SQLite                | hecha     | 5a        |
| 6    | Temporadas continuas                        | pendiente | (nueva)   |
| 7    | Guardado en la interfaz                     | pendiente | 5b        |
| 8    | Partido rápido                              | pendiente | (nueva)   |
| 9    | Estadísticas individuales y alineación      | pendiente | (nueva)   |
| 10   | Mercado de fichajes                         | pendiente | 7         |
| 11   | Partido jugable en tiempo real              | pendiente | 8         |
| 12   | Servidor SSH con Wish                       | pendiente | 6         |
| 13   | Modo jugador                                | pendiente | (nueva)   |
| 14   | Estilo visual                               | pendiente | (nueva)   |

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
nuevo calendario y tabla a cero, en el mismo club.
- `Carrera` pasa a llevar el número de temporada y un **historial** (campeón,
  puesto y puntos del usuario de cada año).
- **Edad y evolución:** cada jugador cumple un año; los jóvenes mejoran, los
  veteranos bajan según su edad; hay retiros y entran juveniles nuevos para
  mantener las plantillas, también en los clubes rivales.
- La semilla de cada jornada se deriva también del número de temporada, de modo
  que la carrera sigue siendo reproducible.
- `Guardado` v2 y segunda migración de la base de datos.
*Terminado cuando*: se juegan varias temporadas seguidas, el historial es
correcto y la carrera guardada se reanuda de forma idéntica.

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
Base común de los dos modos.
- La simulación atribuye cada gol y asistencia a jugadores concretos y registra
  minutos y valoración por partido.
- Alineación elegible: formación y once titulares, que la simulación usa (hoy
  elige sola un 4-3-3).
- Estadísticas acumuladas por jugador y por temporada, y tabla de goleadores.
- `Guardado` v3 y nueva migración.
*Terminado cuando*: se puede cambiar la alineación, cambia el rendimiento del
equipo y las estadísticas individuales se acumulan y se guardan.

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
