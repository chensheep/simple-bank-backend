package worker

import (
	"context"
	"encoding/json"
	"fmt"

	db "github.com/chensheep/simple-bank-backend/db/sqlc"
	"github.com/chensheep/simple-bank-backend/util"
	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
)

const (
	TaskSendVerifyEmail = "task:send_verify_email"
)

type SendVerifyEmailPayload struct {
	Username string `json:"username"`
}

func (d *RedisDistrubutor) DistrubuteTaskSendVerifyEmailTask(ctx context.Context,
	payload *SendVerifyEmailPayload,
	opts ...asynq.Option) error {
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marhal task payload: %w", err)
	}

	// task options can be passed to NewTask, which can be overridden at enqueue time.
	task := asynq.NewTask(TaskSendVerifyEmail, jsonPayload, opts...)
	taskInfo, err := d.client.EnqueueContext(ctx, task)
	if err != nil {
		return fmt.Errorf("could not enqueue task: %w", err)
	}

	log.Info().Str("type", task.Type()).Bytes("payload", task.Payload()).Str("queue", taskInfo.Queue).
		Int("max_retry", taskInfo.MaxRetry).Msg("enqueued task")

	return nil
}

func (processor *RedisTaskProcessor) ProcessTaskSendVerifyEmail(ctx context.Context, t *asynq.Task) error {
	var p SendVerifyEmailPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}

	user, err := processor.store.GetUser(ctx, p.Username)
	if err != nil {
		// Retry if user not found because create user transaction may not be committed yet
		// if err == sql.ErrNoRows {
		// 	return fmt.Errorf("user %s not found: %w", p.Username, asynq.SkipRetry)
		// }
		return fmt.Errorf("failed to get user: %w", err)
	}

	arg := db.CreateVerifyEmailParams{
		Username:   user.Username,
		Email:      user.Email,
		SecretCode: util.RandomString(32),
	}
	verifyEmail, err := processor.store.CreateVerifyEmail(ctx, arg)
	if err != nil {
		return fmt.Errorf("failed to create verify email: %w", err)
	}

	verfifyUrl := fmt.Sprintf("http://localhost:8080/v1/verify_email?email_id=%d&secret_code=%s", verifyEmail.ID, verifyEmail.SecretCode)
	to := []string{user.Email}
	subject := "Welcome to Simple Bank! Please verify your email"
	content := fmt.Sprintf(`
		Welcome to Simple Bank, %s! <br/>
		Please <a href="%s"> click here </a> to verify your email.<br/>
	`, user.FullName, verfifyUrl)
	err = processor.emailSender.SendEmail(to, []string{}, []string{}, subject, content, []string{})
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	log.Info().Str("type", t.Type()).Bytes("payload", t.Payload()).
		Str("email", user.Email).Msg("processed task")

	return nil
}
