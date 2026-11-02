package email

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

const sendEmailEndpoint = "https://api.brevo.com/v3/smtp/email"

var ErrNotConfigured = errors.New("brevo is not configured")

type Message struct {
	ToEmail string
	ToName  string
	Subject string
	HTML    string
}

type BrevoClient struct {
	apiKey      string
	senderEmail string
	senderName  string
	httpClient  *http.Client
}

func NewBrevoClient(apiKey, senderEmail, senderName string) *BrevoClient {
	return &BrevoClient{
		apiKey:      apiKey,
		senderEmail: senderEmail,
		senderName:  senderName,
		httpClient:  &http.Client{Timeout: 10 * time.Second},
	}
}

func (c *BrevoClient) Configured() bool {
	return c.apiKey != "" && c.senderEmail != ""
}

type brevoContact struct {
	Email string `json:"email"`
	Name  string `json:"name,omitempty"`
}

type brevoPayload struct {
	Sender      brevoContact   `json:"sender"`
	To          []brevoContact `json:"to"`
	Subject     string         `json:"subject"`
	HTMLContent string         `json:"htmlContent"`
}

func (c *BrevoClient) Send(ctx context.Context, msg *Message) error {
	if !c.Configured() {
		return ErrNotConfigured
	}

	body, err := json.Marshal(brevoPayload{
		Sender:      brevoContact{Email: c.senderEmail, Name: c.senderName},
		To:          []brevoContact{{Email: msg.ToEmail, Name: msg.ToName}},
		Subject:     msg.Subject,
		HTMLContent: msg.HTML,
	})
	if err != nil {
		return fmt.Errorf("marshal brevo payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, sendEmailEndpoint, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("create brevo request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("api-key", c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("send brevo request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= http.StatusBadRequest {
		detail, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
		return fmt.Errorf("brevo api status %d: %s", resp.StatusCode, detail)
	}
	return nil
}
