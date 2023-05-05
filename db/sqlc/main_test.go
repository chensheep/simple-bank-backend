package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/chensheep/simple-bank-backend/util"
	_ "github.com/lib/pq"
)

const (
	dbDriver = "postgres"
	dbSource = "postgres://postgres:password@localhost:5435/simple-bank-2?sslmode=disable"
)

var testQueries *Queries
var testDB *sql.DB

func TestMain(m *testing.M) {
	var err error

	config, err := util.LoadConfig("../..")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	testDB, err = sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	testQueries = New(testDB)

	os.Exit(m.Run())
}
