package bootstrap

import (
	"database/sql"
	"errors"
	"library-api/kit/env"
	"library-api/kit/log"
	"os"
	"path"

	_ "github.com/mattn/go-sqlite3"
)

func getDatabase() *sql.DB {
	dbPath := env.GetStringOrDefault("SQLITE3_DATABASE_FILE", "database.db")

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.FatalErr("Error getting database", err)
	}

	dropAndCreateTables := env.GetBoolOrDefault("SQLITE3_DB_DROP_AND_CREATE_TABLES", false)
	if dropAndCreateTables {
		schemaFolder, found := env.LookupString("SCHEMA_FOLDER")
		if !found {
			log.Fatal("Error getting schema folder")
		}

		runScript(db, path.Join(schemaFolder, "ddl/drop.sql"))
		runScript(db, path.Join(schemaFolder, "ddl/tables.sql"))
	}

	return db
}

func runScript(db *sql.DB, scriptPath string) {
	content, err := os.ReadFile(scriptPath)
	if err != nil {
		log.FatalErr("Error al leer el fichero", err)
	}

	script := string(content)

	transaction, err := db.Begin()
	if err != nil {
		log.FatalErr("Error running script", err)
	}

	_, err = transaction.Exec(script)
	if err != nil {
		err = errors.Join(err, transaction.Rollback())
		log.FatalErr("Error running script", err)
	}

	if err := transaction.Commit(); err != nil {
		log.FatalErr("Error running script", err)
	}
}
