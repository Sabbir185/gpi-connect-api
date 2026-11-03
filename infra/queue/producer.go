package queue

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/hibiken/asynq"
)

const (
	maxRetry    = 5
	taskTimeout = 30 * time.Second
	// Gives the app a moment to land on its first signed-in screen before the welcome push arrives.
	welcomePushDelay = 5 * time.Second
)

func RedisOpt(redisAddr string) asynq.RedisConnOpt {
	redisAddr = strings.TrimSpace(redisAddr)
	if strings.HasPrefix(redisAddr, "redis://") || strings.HasPrefix(redisAddr, "rediss://") || strings.HasPrefix(redisAddr, "redis-socket://") || strings.HasPrefix(redisAddr, "redis-sentinel://") {
		if opt, err := asynq.ParseRedisURI(redisAddr); err == nil {
			return opt
		}
	}
	if strings.Contains(redisAddr, "@") {
		if opt, err := asynq.ParseRedisURI("redis://" + redisAddr); err == nil {
			return opt
		}
	}
	return asynq.RedisClientOpt{Addr: redisAddr}
}

type Producer struct {
	client *asynq.Client
}

func NewProducer(redisAddr string) *Producer {
	return &Producer{client: asynq.NewClient(RedisOpt(redisAddr))}
}

func (p *Producer) Close() error {
	return p.client.Close()
}

func (p *Producer) EnqueueWelcomeEmail(payload WelcomeEmailPayload) error {
	task, err := NewWelcomeEmailTask(payload)
	if err != nil {
		return err
	}
	return p.enqueue(task, asynq.TaskID(fmt.Sprintf("%s:%d", TypeWelcomeEmail, payload.UserID)))
}

func (p *Producer) EnqueueWelcomePush(payload WelcomePushPayload) error {
	task, err := NewWelcomePushTask(payload)
	if err != nil {
		return err
	}
	return p.enqueue(
		task,
		asynq.TaskID(fmt.Sprintf("%s:%d", TypeWelcomePush, payload.UserID)),
		asynq.ProcessIn(welcomePushDelay),
	)
}

func (p *Producer) enqueue(task *asynq.Task, opts ...asynq.Option) error {
	opts = append(opts, asynq.MaxRetry(maxRetry), asynq.Timeout(taskTimeout))
	_, err := p.client.Enqueue(task, opts...)
	if err != nil && !errors.Is(err, asynq.ErrTaskIDConflict) {
		return fmt.Errorf("enqueue %s: %w", task.Type(), err)
	}
	return nil
}
