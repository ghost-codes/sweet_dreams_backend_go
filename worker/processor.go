package worker

import (
	"context"

	db "github.com/gost-codes/sweet_dreams/db/sqlc"
	"github.com/gost-codes/sweet_dreams/mail"
	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
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
	mailer mail.EmailSender
}

func NewRedisTaskProcessor(redisOpt asynq.RedisClientOpt, store db.Store, mailer mail.EmailSender) TaskProcessor {
	server := asynq.NewServer(redisOpt, asynq.Config{
		Queues: map[string]int{
			CriticalQueue: 10,
			DefaultQueue:  5,
		},
		ErrorHandler: asynq.ErrorHandlerFunc(func(ctx context.Context, task *asynq.Task, err error) {
			log.Info().Err(err).Str("type", task.Type()).Bytes("payload", task.Payload()).Msg("process task failed")
		}),
		Logger: NewLogger(),
	})

	return &RedisTaskProcessor{
		server: server,
		store:  store,
		mailer: mailer,
	}

}

func (processor *RedisTaskProcessor) Start() error {
	serverMux := asynq.NewServeMux()

	serverMux.HandleFunc(TaskSendVerifyEmail, processor.ProcessTaskSendVerifyEmail)
	return processor.server.Start(serverMux)
}
