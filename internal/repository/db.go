package repository

import (
	"database/sql"
	"log"
)

func NewDB(dbDriver string, dbSource string) *sql.DB {
	db, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("Cannot connect to db: ", err)
	}
	return db
}
