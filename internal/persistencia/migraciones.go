package persistencia

import (
	"context"
	"database/sql"
	"fmt"
)

// migraciones son los cambios de esquema en orden; la posición i lleva la base
// de la versión i a la i+1. Nunca se edita una migración ya publicada: se
// agrega otra al final.
var migraciones = []string{
	// 1: esquema inicial.
	`CREATE TABLE partidas (
		ranura           TEXT PRIMARY KEY,
		semilla          INTEGER NOT NULL,
		usuario          INTEGER NOT NULL,
		equipo_usuario   TEXT    NOT NULL,
		jornadas_jugadas INTEGER NOT NULL,
		total_jornadas   INTEGER NOT NULL,
		actualizada      INTEGER NOT NULL
	);
	CREATE TABLE equipos (
		ranura TEXT    NOT NULL REFERENCES partidas(ranura) ON DELETE CASCADE,
		indice INTEGER NOT NULL,
		nombre TEXT    NOT NULL,
		PRIMARY KEY (ranura, indice)
	);
	CREATE TABLE jugadores (
		ranura   TEXT    NOT NULL,
		equipo   INTEGER NOT NULL,
		orden    INTEGER NOT NULL,
		id       INTEGER NOT NULL,
		nombre   TEXT    NOT NULL,
		edad     INTEGER NOT NULL,
		posicion INTEGER NOT NULL,
		ritmo    INTEGER NOT NULL,
		tiro     INTEGER NOT NULL,
		pase     INTEGER NOT NULL,
		regate   INTEGER NOT NULL,
		defensa  INTEGER NOT NULL,
		fisico   INTEGER NOT NULL,
		reflejos INTEGER NOT NULL,
		PRIMARY KEY (ranura, equipo, orden),
		FOREIGN KEY (ranura, equipo) REFERENCES equipos(ranura, indice) ON DELETE CASCADE
	);
	CREATE TABLE resultados (
		ranura          TEXT    NOT NULL REFERENCES partidas(ranura) ON DELETE CASCADE,
		jornada         INTEGER NOT NULL,
		orden           INTEGER NOT NULL,
		local           INTEGER NOT NULL,
		visitante       INTEGER NOT NULL,
		goles_local     INTEGER NOT NULL,
		goles_visitante INTEGER NOT NULL,
		PRIMARY KEY (ranura, jornada, orden)
	);`,
}

// migrar lleva la base al esquema más reciente. Cada migración corre en su
// propia transacción junto con la actualización de user_version, de modo que
// una falla no deja la base a medias.
func migrar(ctx context.Context, db *sql.DB) error {
	var version int
	if err := db.QueryRowContext(ctx, "PRAGMA user_version").Scan(&version); err != nil {
		return fmt.Errorf("leer la version del esquema: %w", err)
	}
	if version > len(migraciones) {
		return fmt.Errorf("la base de datos es de una version mas nueva (%d) que la soportada (%d)",
			version, len(migraciones))
	}
	for i := version; i < len(migraciones); i++ {
		if err := aplicarMigracion(ctx, db, i); err != nil {
			return fmt.Errorf("migracion %d: %w", i+1, err)
		}
	}
	return nil
}

func aplicarMigracion(ctx context.Context, db *sql.DB, i int) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback() //nolint:errcheck // tras Commit es un no-op
	if _, err := tx.ExecContext(ctx, migraciones[i]); err != nil {
		return err
	}
	// PRAGMA no admite parámetros; el valor es un entero propio, no del usuario.
	if _, err := tx.ExecContext(ctx, fmt.Sprintf("PRAGMA user_version = %d", i+1)); err != nil {
		return err
	}
	return tx.Commit()
}
