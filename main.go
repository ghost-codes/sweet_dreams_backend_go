package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/gost-codes/sweet_dreams/api"
	db "github.com/gost-codes/sweet_dreams/db/sqlc"
	"github.com/gost-codes/sweet_dreams/mail"
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
	go runTaskProcessor(redisOpts, *store, config)
	server, err := api.NewServer(*store, config, taskDistributor)

	server.Start(nil)
	return
}

func runTaskProcessor(redisOpt asynq.RedisClientOpt, store db.Store, config util.Config) {
	mailer := mail.NewGmailSender(config.EmailSenderName, config.EmailSenderAddress, config.EmailSenderPassword)
	taskProcessor := worker.NewRedisTaskProcessor(redisOpt, store, mailer)
	fmt.Println("Processor Started")
	err := taskProcessor.Start()

	if err != nil {
		log.Fatal("failed to start task processor: %w", err)
	}
}
