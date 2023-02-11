package main

import (
	"database/sql"
	"log"

	"github.com/beabear/simplebank/api"
	db "github.com/beabear/simplebank/db/sqlc"
	_ "github.com/go-sql-driver/mysql"
)

const (
	dbDriver      = "mysql"
	dbSource      = "root:root@tcp(127.0.0.1:3307)/simple_bank?parseTime=true"
	serverAddress = "0.0.0.0:8080"
)

func main() {
	conn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)

	err = server.Start(serverAddress)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}
}
