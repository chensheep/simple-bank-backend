package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

const (
	dbDriver = "postgres"
	dbSource = "postgres://postgres:password@localhost:5435/simple-bank-2?sslmode=disable"
)

var testQueries *Queries

func TestMain(m *testing.M) {
	db, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	testQueries = New(db)

	os.Exit(m.Run())
}
