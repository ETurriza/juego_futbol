# Plan del proyecto

## Visión

Juego de fútbol para terminal con **modo carrera** (gestión de club: temporadas,
plantilla, fichajes) y, más adelante, **partidos jugables en tiempo real**.

La meta final es servirlo por SSH con Wish para que cualquiera juegue con solo
`ssh`, y tener una **liga en línea entre amigos**.

- Lenguaje: Go
- Interfaz: Bubble Tea
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
- `persistencia` depende de `modelo`; implementa las interfaces de `aplicacion`
  sin importarla (interfaces estructurales de Go).
- `menus` y `red` dependen de `aplicacion`.

No se crean Unit of Work, repositorios ni `aplicacion` antes de que su fase los
necesite.

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

Cada paquete se crea solo cuando llega su fase.

```
cmd/juego/              punto de entrada (modo local / modo servidor)
internal/
  modelo/               fase 1: dominio puro
  generador/            fase 1: generador de jugadores y equipos inventados
  simulacion/           fase 2
  liga/                 fase 3
  aplicacion/           fase 4: casos de uso + puertos
  menus/                fase 4
  persistencia/         fase 5
  red/                  fase 6
  mercado/              fase 7
  partido/              fase 8
tests/e2e/              fase 4 en adelante
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
- **E2E**: en `tests/e2e/`. Desde la fase 4 con `teatest` para la TUI; desde la
  fase 6 levantando el servidor SSH y conectando con un cliente.

`go test ./...` ejecuta solo las pruebas unitarias.

## Fases

### 1. Modelo y generador
Jugador, atributos y equipo en `modelo`. Generador de jugadores con nombres
inventados en `generador`, con aleatoriedad recibida como `*rand.Rand`.
*Terminado cuando*: se genera un equipo completo de forma reproducible con una
semilla fija, con pruebas unitarias.

### 2. Simulación de partidos
Simulación estadística de partidos a partir de atributos (`simulacion`).
Se añade una prueba que verifica los imports de cada paquete contra las reglas
de dependencias.
*Terminado cuando*: el mismo par de equipos y semilla da siempre el mismo
resultado, y mejores atributos producen mejores resultados en promedio.

### 3. Liga
Calendario todos contra todos, jornadas y tabla de posiciones (`liga`).
*Terminado cuando*: se juega una temporada completa en memoria y la tabla es
correcta.

### 4. Menús locales con Bubble Tea
Casos de uso en `aplicacion` y menús en `menus`: plantilla, tabla, avanzar
jornada. Primeras pruebas E2E con `teatest`.
*Terminado cuando*: se puede jugar una temporada desde la terminal.

### 5. Persistencia con SQLite
Adaptador en `persistencia` que implementa los puertos de `aplicacion`.
Pruebas de integración.
*Terminado cuando*: una carrera se guarda y se reanuda.

### 6. Servidor SSH con Wish
Paquete `red`: la misma TUI servida por SSH. Pruebas E2E con cliente SSH.
*Terminado cuando*: `ssh` al servidor abre el juego y varias sesiones conviven.

### 7. Mercado de fichajes
Paquete `mercado`: compra, venta y valoración de jugadores.
*Terminado cuando*: se puede fichar y vender dentro de la carrera.

### 8. Partido jugable en tiempo real
Paquete `partido`: partido en tiempo real que usa los atributos de los
jugadores.
*Terminado cuando*: un partido se juega de principio a fin en la terminal y su
resultado alimenta la liga.

## Convenciones transversales

- Nombres de jugadores y equipos siempre inventados; nada de marcas ni personas
  reales.
- Toda aleatoriedad entra como `*rand.Rand`, nunca global.
- Todo cambio en lógica incluye pruebas.
