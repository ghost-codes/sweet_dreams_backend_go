package worker

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
)

const TaskSendAdminEmail = "task:send_admin_email"

type PayloadSendAdminEmail struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (distributor *RedisTaskDistributor) DistributeTaskSendAdminEmail(ctx context.Context, payload *PayloadSendAdminEmail, opts ...asynq.Option) error {
	jsonPayload, err := json.Marshal(payload)

	if err != nil {
		return fmt.Errorf("could not marshal payload: %w", err)
	}

	task := asynq.NewTask(TaskSendAdminEmail, jsonPayload, opts...)
	taskInfo, err := distributor.client.EnqueueContext(ctx, task)

	if err != nil {
		return fmt.Errorf("failed to enque task: %w", err)
	}

	log.Info().Str("type", task.Type()).Bytes("payload", task.Payload()).Str("queue", taskInfo.Queue).Int("max_retries", taskInfo.MaxRetry).Msg("enqueued task")

	return nil
}

func (processor *RedisTaskProcessor) ProcessTaskSendAdminEmail(ctx context.Context, task *asynq.Task) error {
	var payload PayloadSendAdminEmail

	err := json.Unmarshal(task.Payload(), &payload)
	if err != nil {
		return fmt.Errorf("failed to unmarshal this payload: %w", asynq.SkipRetry)
	}

	admin, err := processor.store.GetAdmin(ctx, payload.Username)

	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("user doesn't exist: %w", asynq.SkipRetry)
		}

		return fmt.Errorf("failed to get user: %w", err)
	}

	//TODO: send email to user
	log.Info().Str("type", task.Type()).Bytes("payload", task.Payload()).Str("email", admin.Email).Msg("Proccess task")

	if err != nil {
		return fmt.Errorf("could not create verify email object:%w", err)
	}

	subject := "Admin account notification"
	to := []string{admin.Email}
	content := fmt.Sprintf(`
	<h1>Admin privillege</h1>
	<p> You have been added as an admin to the sweet dreams project here is your password to access the admin panel</p>
	<h4>Password:</h4>
	<h2>%s</h2>
	`, payload.Password)

	return processor.mailer.SendEmail(subject, content, to, nil, nil, nil)

}
