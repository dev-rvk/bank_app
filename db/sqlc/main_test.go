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

func TestMain(m *testing.M){
	var conn *sql.DB
	conn, err := sql.Open(dbDriver, dbSource) // returns a connection object

	if err != nil {
		log.Fatal("Error while connecting to db: ", err)
	}

	testQueries = New(conn) //New is defined in the db.go file which creates the Queries object

	os.Exit(m.Run())

} 