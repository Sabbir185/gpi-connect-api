package firebase

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"firebase.google.com/go/v4/messaging"
	"google.golang.org/api/option"
)

type Client struct {
	Auth      *auth.Client
	Messaging *messaging.Client
}

func NewClient(ctx context.Context, credentialsBase64 string) (*Client, error) {
	credentials, err := base64.StdEncoding.DecodeString(credentialsBase64)
	if err != nil {
		return nil, fmt.Errorf("decode firebase credentials: %w", err)
	}

	var serviceAccount struct {
		ProjectID string `json:"project_id"`
	}
	if err := json.Unmarshal(credentials, &serviceAccount); err != nil || serviceAccount.ProjectID == "" {
		return nil, fmt.Errorf("firebase credentials are not a service account json")
	}

	app, err := firebase.NewApp(
		ctx,
		&firebase.Config{ProjectID: serviceAccount.ProjectID},
		option.WithAuthCredentialsJSON(option.ServiceAccount, credentials),
	)
	if err != nil {
		return nil, fmt.Errorf("init firebase app: %w", err)
	}

	authClient, err := app.Auth(ctx)
	if err != nil {
		return nil, fmt.Errorf("init firebase auth: %w", err)
	}

	messagingClient, err := app.Messaging(ctx)
	if err != nil {
		return nil, fmt.Errorf("init firebase messaging: %w", err)
	}

	return &Client{Auth: authClient, Messaging: messagingClient}, nil
}
