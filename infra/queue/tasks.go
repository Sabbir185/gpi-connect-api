package queue

import (
	"encoding/json"
	"fmt"

	"github.com/hibiken/asynq"
)

const (
	TypeWelcomeEmail = "email:welcome"
	TypeWelcomePush  = "push:welcome"
)

type WelcomeEmailPayload struct {
	UserID int64  `json:"user_id"`
	Email  string `json:"email"`
	Name   string `json:"name"`
	Locale string `json:"locale"`
}

type WelcomePushPayload struct {
	UserID int64 `json:"user_id"`
}

func NewWelcomeEmailTask(payload WelcomeEmailPayload) (*asynq.Task, error) {
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("marshal welcome email payload: %w", err)
	}
	return asynq.NewTask(TypeWelcomeEmail, data), nil
}

func NewWelcomePushTask(payload WelcomePushPayload) (*asynq.Task, error) {
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("marshal welcome push payload: %w", err)
	}
	return asynq.NewTask(TypeWelcomePush, data), nil
}
