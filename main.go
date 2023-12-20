package main

import (
	"database/sql"
	"log"

	"github.com/jithinlal/simplebank/api"
	db "github.com/jithinlal/simplebank/db/sqlc"
	_ "github.com/lib/pq"
)

const (
	dbDriver   = "postgres"
	dbSource   = "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable"
	serverAddr = "127.0.0.1:8080"
)

func main() {
	conn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("cannot connect to db: ", err)
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)

	err = server.Start(serverAddr)
	if err != nil {
		log.Fatal("server cannot start")
	}
}
