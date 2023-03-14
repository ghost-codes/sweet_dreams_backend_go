package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/gost-codes/sweet_dreams/api"
	db "github.com/gost-codes/sweet_dreams/db/sqlc"
	"github.com/gost-codes/sweet_dreams/util"
	"github.com/gost-codes/sweet_dreams/worker"
	"github.com/hibiken/asynq"
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
	redisOpts := asynq.RedisClientOpt{
		Addr: config.RedisServerAddress,
	}
	taskDistributor := worker.NewRedisTaskDistributor(redisOpts)
	go runTaskProcessor(redisOpts, *store)
	server, err := api.NewServer(*store, config, taskDistributor)

	server.Start(nil)
	return
}

func runTaskProcessor(redisOpt asynq.RedisClientOpt, store db.Store) {
	taskProcessor := worker.NewRedisTaskProcessor(redisOpt, store)
	fmt.Println("Processor Started")
	err := taskProcessor.Start()

	if err != nil {
		log.Fatal("failed to start task processor: %w", err)
	}
}
