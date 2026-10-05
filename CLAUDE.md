# CLAUDE.md

Juego de fútbol para terminal (Go + Bubble Tea) con modo carrera y, más
adelante, partidos en tiempo real servidos por SSH con Wish.
La visión, la arquitectura y las fases están en [docs/PLAN.md](docs/PLAN.md).

Módulo: `github.com/ETurriza/juego_futbol`.

## Paquetes

Todo el código vive en `internal/`, con nombres en español sin acentos:
`modelo`, `generador`, `simulacion`, `liga`, `aplicacion`, `menus`,
`persistencia`, `mercado`, `partido`, `red`. Además, `arquitectura` solo contiene
la prueba que verifica las reglas de dependencias de abajo.

Al crear un paquete nuevo hay que agregarlo al mapa `permitidos` de
`internal/arquitectura/arquitectura_test.go`; si no, `go test ./...` falla.
Cada paquete se crea solo cuando llega su fase (ver `docs/PLAN.md`). No crear
carpetas ni código por adelantado. El punto de entrada está en `cmd/juego/`.

## Dependencias entre paquetes

- `internal/modelo` no importa ningún otro paquete del proyecto.
- `simulacion`, `liga` y `mercado` solo dependen de `modelo` y entre sí; nunca
  de `aplicacion`, `menus`, `persistencia` o `red`.
- `aplicacion` depende de `modelo` y de los servicios de dominio. Define los
  puertos (interfaces) que consume.
- `persistencia` depende de `modelo` e implementa los puertos de `aplicacion`
  sin importarla.
- `menus` y `red` dependen de `aplicacion`.
- Las capas externas dependen de las internas, nunca al revés.

## Aleatoriedad

Toda aleatoriedad se recibe como parámetro (`*rand.Rand`), nunca global (nada de
`rand.Intn` a nivel de paquete). Así las pruebas son reproducibles con una
semilla fija.

## Datos

Nombres de jugadores y equipos siempre inventados. Nada de marcas ni personas
reales.

## Pruebas

- Todo cambio en lógica incluye pruebas.
- Unit: `*_test.go` junto al código, con semilla fija, sin disco ni red.
- Integration: `*_integration_test.go` con `//go:build integration`; se corren
  con `go test -tags=integration ./...`.
- E2E: en `tests/e2e/`, desde la fase 4.
- Antes de terminar una tarea se corren `go build ./...` y `go test ./...`.

## CLAUDE.md por paquete

Los paquetes con reglas propias pueden tener su `CLAUDE.md` local, creado junto
con el paquete en su fase. No repetir ahí las reglas de este archivo.

## Git

- Flujo: `feat/*` -> PR a `dev` -> PR de `dev` a `main` (prod). Las ramas de
  trabajo salen de `dev`. Ambas ramas largas tienen un ruleset en GitHub (PR
  obligatorio, sin push directo).
- Nunca trabajar en `main` ni en `dev`.
- Se pueden hacer commits y crear ramas, pero antes de cada commit se entrega
  al usuario un resumen detallado de lo que se va a hacer (rama, archivos
  incluidos, mensaje).
- Se pueden abrir PR con `gh`, también previo resumen detallado al usuario
  (rama origen y destino, título, descripción). Para eso se permite hacer push
  de la rama de trabajo, y solo de ella.
- Nunca hacer push a `main` ni a `dev`, y nunca hacer merge. De eso se encarga
  el usuario.
