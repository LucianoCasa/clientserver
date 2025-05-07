package db

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var DBConn *sql.DB

func Init(dbPath string) {
	var err error
	DBConn, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatalf("Erro ao conectar ao banco: %v", err)
	}

}
