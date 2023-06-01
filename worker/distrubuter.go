package worker

import (
	"context"

	"github.com/hibiken/asynq"
)

type TaskDistrubutor interface {
	DistrubuteTaskSendVerifyEmailTask(
		context.Context,
		*SendVerifyEmailPayload,
		...asynq.Option) error
}

type RedisDistrubutor struct {
	client *asynq.Client
}

func NewRedisDistrubuter(option asynq.RedisConnOpt) TaskDistrubutor {
	return &RedisDistrubutor{
		client: asynq.NewClient(option),
	}
}
