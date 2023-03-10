package main

import (
	"database/sql"
	"log"

	"github.com/gost-codes/sweet_dreams/api"
	db "github.com/gost-codes/sweet_dreams/db/sqlc"
	"github.com/gost-codes/sweet_dreams/util"
	_ "github.com/lib/pq"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("unable to load config:", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource())
	if err != nil {
		log.Fatal("unable to connect database:", err)
	}

	store := db.NewStore(conn)
	server, err := api.NewServer(*store, config)

	server.Start(nil)
	return
}
