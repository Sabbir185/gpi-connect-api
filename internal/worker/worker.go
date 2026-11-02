package worker

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"firebase.google.com/go/v4/messaging"
	"github.com/Sabbir185/gpi/infra/email"
	"github.com/Sabbir185/gpi/infra/queue"
	"github.com/Sabbir185/gpi/internal/device"
	"github.com/hibiken/asynq"
)

type Worker struct {
	server     *asynq.Server
	mailer     *email.BrevoClient
	messaging  *messaging.Client
	deviceRepo device.DeviceRepository
}

func New(redisAddr string, concurrency int, mailer *email.BrevoClient, messagingClient *messaging.Client, deviceRepo device.DeviceRepository) *Worker {
	server := asynq.NewServer(
		queue.RedisOpt(redisAddr),
		asynq.Config{
			Concurrency: concurrency,
			ErrorHandler: asynq.ErrorHandlerFunc(func(ctx context.Context, task *asynq.Task, err error) {
				log.Printf("task %s failed: %v", task.Type(), err)
			}),
		},
	)
	return &Worker{server: server, mailer: mailer, messaging: messagingClient, deviceRepo: deviceRepo}
}

func (w *Worker) Run() error {
	mux := asynq.NewServeMux()
	mux.HandleFunc(queue.TypeWelcomeEmail, w.handleWelcomeEmail)
	mux.HandleFunc(queue.TypeWelcomePush, w.handleWelcomePush)
	return w.server.Run(mux)
}

func (w *Worker) handleWelcomeEmail(ctx context.Context, task *asynq.Task) error {
	var payload queue.WelcomeEmailPayload
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return fmt.Errorf("unmarshal welcome email payload: %v: %w", err, asynq.SkipRetry)
	}
	if !w.mailer.Configured() {
		log.Printf("welcome email for user %d skipped: BREVO_API_KEY or BREVO_SENDER_EMAIL is not set", payload.UserID)
		return nil
	}

	text := messagesFor(payload.Locale)
	err := w.mailer.Send(ctx, &email.Message{
		ToEmail: payload.Email,
		ToName:  payload.Name,
		Subject: text.EmailSubject,
		HTML:    welcomeEmailHTML(text, payload.Name),
	})
	if err != nil {
		return fmt.Errorf("send welcome email: %w", err)
	}

	log.Printf("welcome email sent to user %d", payload.UserID)
	return nil
}

func (w *Worker) handleWelcomePush(ctx context.Context, task *asynq.Task) error {
	var payload queue.WelcomePushPayload
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return fmt.Errorf("unmarshal welcome push payload: %v: %w", err, asynq.SkipRetry)
	}

	targets, err := w.deviceRepo.ListPushTargets(ctx, payload.UserID)
	if err != nil {
		return fmt.Errorf("load push targets: %w", err)
	}
	if len(targets) == 0 {
		log.Printf("welcome push for user %d skipped: no device with a push token", payload.UserID)
		return nil
	}

	messages := make([]*messaging.Message, 0, len(targets))
	for _, target := range targets {
		text := messagesFor(target.Locale)
		messages = append(messages, &messaging.Message{
			Token:        target.Token,
			Notification: &messaging.Notification{Title: text.PushTitle, Body: text.PushBody},
			Data:         map[string]string{"type": "welcome"},
			Android:      &messaging.AndroidConfig{Priority: "high"},
			APNS: &messaging.APNSConfig{
				Payload: &messaging.APNSPayload{Aps: &messaging.Aps{Sound: "default"}},
			},
		})
	}

	response, err := w.messaging.SendEach(ctx, messages)
	if err != nil {
		return fmt.Errorf("send welcome push: %w", err)
	}

	var failed error
	for i, result := range response.Responses {
		if result.Success {
			continue
		}
		if messaging.IsUnregistered(result.Error) || messaging.IsSenderIDMismatch(result.Error) {
			if err := w.deviceRepo.DeletePushToken(ctx, targets[i].Token); err != nil {
				log.Printf("remove stale push token for user %d: %v", payload.UserID, err)
			}
			continue
		}
		failed = errors.Join(failed, result.Error)
	}

	log.Printf("welcome push for user %d: %d sent, %d failed", payload.UserID, response.SuccessCount, response.FailureCount)
	if failed != nil && response.SuccessCount == 0 {
		return fmt.Errorf("send welcome push: %w", failed)
	}
	return nil
}
