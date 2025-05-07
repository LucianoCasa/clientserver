package db

import (
	"database/sql"
	"fmt"
)

func MigrateUp(db *sql.DB) error {
	sql := `
	CREATE TABLE IF NOT EXISTS cotacoes (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		valor REAL NOT NULL,
		created_at TEXT NOT NULL
	);`
	_, err := db.Exec(sql)
	if err != nil {
		return fmt.Errorf("erro ao executar MigrateUp: %w", err)
	}
	return nil
}

func MigrateDown(db *sql.DB) error {
	sql := `DROP TABLE IF EXISTS cotacoes;`
	_, err := db.Exec(sql)
	if err != nil {
		return fmt.Errorf("erro ao executar MigrateDown: %w", err)
	}
	return nil
}
