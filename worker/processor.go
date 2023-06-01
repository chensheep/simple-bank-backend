package worker

import (
	"context"

	db "github.com/chensheep/simple-bank-backend/db/sqlc"
	"github.com/hibiken/asynq"
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
	server *asynq.Server
	store  db.Store
}

func NewRedisTaskProcessor(r asynq.RedisConnOpt, store db.Store) *RedisTaskProcessor {
	server := asynq.NewServer(r, asynq.Config{
		Queues: map[string]int{
			QueueCritical: 10,
			QueueDefault:  5,
			QueueLow:      1,
		},
	})
	return &RedisTaskProcessor{
		server: server,
		store:  store,
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
