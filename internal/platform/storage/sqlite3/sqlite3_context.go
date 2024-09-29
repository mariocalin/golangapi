package sqlite3

import (
	"database/sql"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

type Sqlite3Context struct {
	Db *sql.DB
}

func CreateSqlite3Context() *Sqlite3Context {
	dbpath := os.Getenv("SQLITE3_DATABASE_FILE")
	if dbpath == "" {
		panic("SQlite3 db path not provided")
	}

	db, err := sql.Open("sqlite3", dbpath)
	if err != nil {
		panic(err)
	}

	return &Sqlite3Context{Db: db}
}
