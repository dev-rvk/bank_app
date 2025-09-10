package main

import (
	"database/sql"
	"log"

	"github.com/devrvk/simplebank/api"
	db "github.com/devrvk/simplebank/db/sqlc"
	"github.com/devrvk/simplebank/util"
	_ "github.com/lib/pq"
)

// connect db to store pass the store to server struct and then start the server

// entrypoint for server
func main(){
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("Error reading config file:", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource) // returns a connection object

	if err != nil {
		log.Fatal("Error while connecting to db: ", err)
	}
	store := db.NewStore(conn)

	server, err := api.NewServer(config, store)
	if err != nil{
		log.Fatal("Failed to create new server: ", err)
	}

	err = server.Start(config.ServerAddress)

	if err != nil{
		log.Fatal("Cannot Start Server: ", err)
	}

}