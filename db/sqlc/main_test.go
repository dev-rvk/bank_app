package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

// change to env variable
const (
	dbDriver = "postgres"
	dbSource = "postgres://postgres:postgres@localhost:5432/simple_bank?sslmode=disable"
)

// is required by every function made using sqlc
var testQueries *Queries
var testDB *sql.DB

func TestMain(m *testing.M) {
	var err error
	testDB, err = sql.Open(dbDriver, dbSource) // returns a connection object

	if err != nil {
		log.Fatal("Error while connecting to db: ", err)
	}

	testQueries = New(testDB) //New is defined in the db.go file which creates the Queries object

	os.Exit(m.Run())

}
