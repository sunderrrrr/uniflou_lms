package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/uniflou/uniflou-lms/backend/internal/config"
)

type Client struct {
	url    string
	key    string
	model  string
	client *http.Client
}

func NewClient(cfg *config.Config) *Client {
	return &Client{
		url:   cfg.LLMAPIURL,
		key:   cfg.LLMAPIKey,
		model: cfg.LLMModel,
		client: &http.Client{
			Timeout: 30 * time.Second,
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				if len(via) >= 3 || !allowedHost(req.URL.String()) {
					return http.ErrUseLastResponse
				}
				return nil
			},
		},
	}
}

func (c *Client) Complete(ctx context.Context, prompt string) (string, error) {
	if c.key == "" || c.url == "" {
		return "Локальный режим: внешний ИИ не настроен. Сверьте ответ с пояснением к заданию, отметьте правило, которое вызвало затруднение, и повторите похожее задание позже.", nil
	}
	if !allowedHost(c.url) {
		return "", fmt.Errorf("llm host not allowed")
	}
	body := map[string]any{
		"modelUri":          fmt.Sprintf("gpt://%s", c.model),
		"completionOptions": map[string]any{"stream": false, "temperature": 0.3, "maxTokens": 800},
		"messages":          []map[string]string{{"role": "user", "text": prompt}},
	}
	raw, _ := json.Marshal(body)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.url, bytes.NewReader(raw))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Api-Key "+c.key)
	resp, err := c.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return "", fmt.Errorf("LLM request failed with status %d", resp.StatusCode)
	}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	var parsed struct {
		Result struct {
			Alternatives []struct {
				Message struct {
					Text string `json:"text"`
				} `json:"message"`
			} `json:"alternatives"`
		} `json:"result"`
	}
	if err := json.Unmarshal(b, &parsed); err != nil {
		return "", err
	}
	if len(parsed.Result.Alternatives) == 0 || parsed.Result.Alternatives[0].Message.Text == "" {
		return "", fmt.Errorf("LLM returned an empty response")
	}
	return parsed.Result.Alternatives[0].Message.Text, nil
}

func allowedHost(rawURL string) bool {
	parsed, err := url.Parse(rawURL)
	if err != nil || parsed.Scheme != "https" || parsed.Hostname() == "" {
		return false
	}
	if port := parsed.Port(); port != "" && port != "443" {
		return false
	}
	switch strings.ToLower(parsed.Hostname()) {
	case "api.yandexcloud.net", "gigachat.devices.sberbank.ru", "llm.api.cloud.yandex.net":
		return true
	default:
		return false
	}
}
