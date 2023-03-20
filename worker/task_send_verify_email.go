package worker

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	db "github.com/gost-codes/sweet_dreams/db/sqlc"
	"github.com/gost-codes/sweet_dreams/util"
	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
)

const TaskSendVerifyEmail = "task:send_verify_email"

type PayloadSendVerifyEmail struct {
	UserID int64 `json:"username"`
}

func (distributor *RedisTaskDistributor) DistributeTaskSendVerifyEmail(ctx context.Context, payload *PayloadSendVerifyEmail, opts ...asynq.Option) error {
	jsonPayload, err := json.Marshal(payload)

	if err != nil {
		return fmt.Errorf("could not marshal payload: %w", err)
	}

	task := asynq.NewTask(TaskSendVerifyEmail, jsonPayload, opts...)
	taskInfo, err := distributor.client.EnqueueContext(ctx, task)

	if err != nil {
		return fmt.Errorf("failed to enque task: %w", err)
	}

	log.Info().Str("type", task.Type()).Bytes("payload", task.Payload()).Str("queue", taskInfo.Queue).Int("max_retries", taskInfo.MaxRetry).Msg("enqueued task")

	return nil
}

func (processor *RedisTaskProcessor) ProcessTaskSendVerifyEmail(ctx context.Context, task *asynq.Task) error {
	var payload PayloadSendVerifyEmail

	err := json.Unmarshal(task.Payload(), &payload)
	if err != nil {
		return fmt.Errorf("failed to unmarshal this payload: %w", asynq.SkipRetry)
	}

	user, err := processor.store.GetUserByID(ctx, payload.UserID)

	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("user doesn't exist: %w", asynq.SkipRetry)
		}

		return fmt.Errorf("failed to get user: %w", err)
	}

	//TODO: send email to user
	log.Info().Str("type", task.Type()).Bytes("payload", task.Payload()).Str("email", user.Email).Msg("Proccess task")
	args := db.CreateVerifyEmailParams{
		Username:  &user.Username,
		Email:     user.Email,
		SecretKey: util.RandomString(64),
	}
	verifyEmail, err := processor.store.CreateVerifyEmail(ctx, args)

	if err != nil {
		return fmt.Errorf("could not create verify email object:%w", err)
	}

	strUrl := fmt.Sprintf("http://localhost:8080/verify_email?id=%v&code=%s", verifyEmail.ID, verifyEmail.SecretKey)

	subject := "Verify Email"
	to := []string{user.Email}
	content := fmt.Sprintf(`
	<h1>Thank you for joining the family</h1>
	<p> Please verify email inorder to have access to your account with this link below</a></p>
	<br>%s
	`, strUrl)

	return processor.mailer.SendEmail(subject, content, to, nil, nil, nil)

}
