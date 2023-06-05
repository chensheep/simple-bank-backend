package worker

import (
	"context"

	db "github.com/chensheep/simple-bank-backend/db/sqlc"
	"github.com/chensheep/simple-bank-backend/email"
	"github.com/hibiken/asynq"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

const (
	QueueCritical = "critical"
	QueueDefault  = "default"
	QueueLow      = "low"
)

type TaskProcessor interface {
	Start() error
	ProcessTaskSendVerifyEmail(context.Context, *asynq.Task) error
}

type RedisTaskProcessor struct {
	server      *asynq.Server
	store       db.Store
	emailSender email.EmailSender
}

func NewRedisTaskProcessor(r asynq.RedisConnOpt, store db.Store, emailSender email.EmailSender) *RedisTaskProcessor {
	logger := NewLogger()
	redis.SetLogger(logger)

	server := asynq.NewServer(r, asynq.Config{
		Queues: map[string]int{
			QueueCritical: 10,
			QueueDefault:  5,
			QueueLow:      1,
		},
		ErrorHandler: asynq.ErrorHandlerFunc(func(ctx context.Context, task *asynq.Task, err error) {
			log.Error().Err(err).Str("type", task.Type()).
				Bytes("payload", task.Payload()).Msg("failed to process task")
		}),
		Logger: logger,
	})
	return &RedisTaskProcessor{
		server:      server,
		store:       store,
		emailSender: emailSender,
	}
}

func (processor *RedisTaskProcessor) Start() error {
	mux := asynq.NewServeMux()
	mux.HandleFunc(TaskSendVerifyEmail, processor.ProcessTaskSendVerifyEmail)
	// ...register other handlers...

	if err := processor.server.Start(mux); err != nil {
		return err
	}

	return nil
}
