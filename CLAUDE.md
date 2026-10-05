# CLAUDE.md

Juego de fútbol para terminal (Go + Bubble Tea) con modo carrera y, más
adelante, partidos en tiempo real servidos por SSH con Wish.
La visión, la arquitectura y las fases están en [docs/PLAN.md](docs/PLAN.md).

Módulo: `github.com/ETurriza/juego_futbol`.

## Paquetes

Todo el código vive en `internal/`, con nombres en español sin acentos:
`modelo`, `generador`, `simulacion`, `liga`, `progresion`, `aplicacion`,
`menus`, `persistencia`, `mercado`, `partido`, `red`. Además, `arquitectura`
solo contiene la prueba que verifica las reglas de dependencias de abajo.

Al crear un paquete nuevo hay que agregarlo al mapa `permitidos` de
`internal/arquitectura/arquitectura_test.go`; si no, `go test ./...` falla.
Cada paquete se crea solo cuando llega su fase (ver `docs/PLAN.md`). No crear
carpetas ni código por adelantado. El punto de entrada está en `cmd/juego/`.

## Dependencias entre paquetes

- `internal/modelo` no importa ningún otro paquete del proyecto.
- `progresion` solo depende de `modelo`; `generador`, de `modelo` y `progresion`.
- `simulacion`, `liga` y `mercado` solo dependen de `modelo` y entre sí; nunca
  de `aplicacion`, `menus`, `persistencia` o `red`.
- `aplicacion` depende de `modelo` y de los servicios de dominio. Define los
  puertos (interfaces) que consume.
- `persistencia` depende de `modelo` y de `aplicacion`: implementa sus puertos
  y usa su tipo `Guardado`. Usa SQLite con `modernc.org/sqlite` (sin CGO).
- `menus` y `red` dependen de `aplicacion`.
- Las capas externas dependen de las internas, nunca al revés.

## Ejecutar

`go run ./cmd/juego [--semilla N] [--equipos N]` (necesita una terminal real).
La misma semilla da siempre la misma liga.

## Interfaz (`menus`)

- Cada pantalla cabe en 80 columnas; una prueba lo verifica con datos reales.
- Teclas: `↑/↓` o `j/k`, `enter`, `esc` o `backspace` para volver, `tab` para
  alternar vistas, `q` solo en el menú principal y el de fin de temporada,
  `ctrl+c` en cualquier pantalla.
- Las opciones de un menú se indexan con constantes con nombre (`opAvanzar`,
  `opFinNueva`...), no con números. Las pruebas navegan con esas constantes.
- Las pantallas con cursor guardan su estado en una pila para que `esc` vuelva a
  cada pantalla tal como estaba. El club del usuario va marcado con `*`.
- Si se agrega o se mueve una opción de menú, hay que actualizar las secuencias
  de teclas de las pruebas unitarias y de `tests/e2e/`.
- Las pantallas clave tienen archivos de referencia en
  `internal/menus/testdata/pantallas/` y una página generada, `docs/PANTALLAS.md`
  (semilla fija, terminal de 80x30). Tras un cambio visual intencionado:
  `go test ./internal/menus -update` y revisar el diff. Una pantalla nueva
  importante se agrega a `galeria` en `galeria_test.go`.

## Calibración

Los modelos con números (progresión por edad, sucesos de partido, valoración)
tienen sus constantes con nombre al inicio del archivo. Sus pruebas imprimen con
`go test -v` las tablas y porcentajes medidos, para juzgar la calibración a
simple vista, y comprueban rangos realistas. Cambiar una constante exige
revisar esas pruebas.

Las pruebas de calibración también fijan cuántos cracks hay y cuánto pesan:
~3-4 % de los jugadores de campo con pico >= 90 repartidos por todas las
posiciones, y un crack de 95 suma ~+8 puntos de victoria en cada una. Esos
números salen de ligas generadas de verdad; para simular muchos partidos sin
detalle se usa `simulacion.Marcador`, que es mucho más rápido que `Simular`.

Una media global estable no basta: la liga se calibra por posición y por edad
(`aplicacion/calibracion_test.go`): dispersión, porcentaje de estrellas,
atributos en el tope, y que la liga inicial se parezca a la de muchas
temporadas después. Los jugadores iniciales y los juveniles se crean con la
propia progresión (`generador` usa `progresion`), nunca con una distribución
aparte.

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
- E2E: en `tests/e2e/`, con un `tea.Program` real (ver `docs/PLAN.md`).
- Antes de terminar una tarea se corren `go build ./...` y `go test ./...`. Si
  se toca `persistencia` o un puerto de `aplicacion`, también
  `go test -tags=integration ./...`.
- Cada implementación de un puerto pasa las pruebas de contrato de
  `internal/aplicacion/contrato`.
- Una prueba nueva debe poder fallar: tras escribirla, romper a propósito el
  código que cubre y comprobar que falla (y restaurar el archivo).

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
