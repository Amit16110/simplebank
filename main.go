package main

import (
	"database/sql"
	"log"

	"github.com/amit16110/simplebank/api"
	db "github.com/amit16110/simplebank/db/sqlc"
	"github.com/amit16110/simplebank/util"
	_ "github.com/lib/pq" // with this driver code talk to database.
)

const (
	//Similar as a main.test file
	dbDriver      = "postgres"
	dbSource      = "postgresql://postgres:secret@localhost:5432/simplebank?sslmode=disable"
	serverAddress = "0.0.0.0:8080"
)

func main() {
	// config file
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}
	conn, err := sql.Open(config.DBDriver, config.DBSource)

	if err != nil {
		log.Fatal("Cannot connect to db:", err)
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)

	err = server.Start(config.ServerAddress)

	if err != nil {
		log.Fatal("Cannot start the server:", err)
	}

}
