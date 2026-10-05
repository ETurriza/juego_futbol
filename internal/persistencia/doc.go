// Package persistencia implementa aplicacion.RepositorioPartidas sobre SQLite
// (driver modernc.org/sqlite, escrito en Go puro y sin CGO).
//
// Cada partida se guarda por completo en una transacción: la fila de la
// partida, sus equipos, jugadores y resultados. El esquema evoluciona con
// migraciones numeradas registradas en PRAGMA user_version. Depende de modelo
// y de aplicacion (que define el puerto y el tipo Guardado).
package persistencia
