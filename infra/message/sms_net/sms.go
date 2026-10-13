package smsnet

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const sendSMSEndPoint = "https://api.sms.net.bd/sendsms"

type SendSMSPayload struct {
	To  string
	Msg string
}

type SendSMSResponse struct {
	Error int    `json:"error"`
	Msg   string `json:"msg"`
	Data  struct {
		ResquestID int `json:"request_id"`
	} `json:"data"`
}

type SMSClient struct {
	apiKey     string
	httpClient *http.Client
}

func NewSMSClient(apiKey string) *SMSClient {
	return &SMSClient{
		apiKey: apiKey,
		httpClient: &http.Client{
			Timeout: 5 * time.Second, // Network timeout
		},
	}
}

func (c *SMSClient) SendOTP(ctx context.Context, payload *SendSMSPayload) (*SendSMSResponse, error) {
	// Preapre form data, url encoded
	formData := url.Values{}
	formData.Set("api_key", c.apiKey)
	formData.Set("msg", payload.Msg)
	formData.Set("to", payload.To)

	// create http POST request with context
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		sendSMSEndPoint,
		strings.NewReader(formData.Encode()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create sms http req: %w", err)
	}

	// set headers
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// send the request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// read and decode response
	var response SendSMSResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode sms response: %w", err)
	}

	// check api error field
	if response.Error != 0 {
		return nil, fmt.Errorf("sms api error: %s (code: %d)", response.Msg, response.Error)
	}

	return &response, nil
}
