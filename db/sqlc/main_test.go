package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/devrvk/simplebank/util"
	_ "github.com/lib/pq"
)

// is required by every function made using sqlc
var testQueries *Queries
var testDB *sql.DB

func TestMain(m *testing.M) {
	var err error
	config, err := util.LoadConfig("../..")
	if err != nil {
		log.Fatal("Error reading config file:", err)
	}

	
	testDB, err = sql.Open(config.DBDriver, config.DBSource) // returns a connection object

	if err != nil {
		log.Fatal("Error while connecting to db: ", err)
	}

	testQueries = New(testDB) //New is defined in the db.go file which creates the Queries object

	os.Exit(m.Run())

}
