package worker

import (
	"context"

	db "github.com/gost-codes/sweet_dreams/db/sqlc"
	"github.com/gost-codes/sweet_dreams/util"
	"github.com/hibiken/asynq"
)

const (
	CriticalQueue = "critical"
	DefaultQueue  = "default"
)

type TaskProcessor interface {
	Start() error
	ProcessTaskSendVerifyEmail(ctx context.Context, task *asynq.Task) error
}

type RedisTaskProcessor struct {
	server *asynq.Server
	store  db.Store
	config util.Config
}

func NewRedisTaskProcessor(redisOpt asynq.RedisClientOpt, store db.Store, config util.Config) TaskProcessor {
	server := asynq.NewServer(redisOpt, asynq.Config{
		Queues: map[string]int{
			CriticalQueue: 10,
			DefaultQueue:  5,
		},
	})

	return &RedisTaskProcessor{
		server: server,
		store:  store,
		config: config,
	}

}

func (processor *RedisTaskProcessor) Start() error {
	serverMux := asynq.NewServeMux()

	serverMux.HandleFunc(TaskSendVerifyEmail, processor.ProcessTaskSendVerifyEmail)
	return processor.server.Start(serverMux)
}
