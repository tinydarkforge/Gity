package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/tinydarkforge/intake/internal/jsonx"
)

type OllamaClient struct {
	Host    string
	Model   string
	Timeout time.Duration
	HTTP    *http.Client
}

func NewOllama(host, model string, timeout time.Duration) *OllamaClient {
	return &OllamaClient{
		Host:    host,
		Model:   model,
		Timeout: timeout,
		HTTP:    &http.Client{Timeout: timeout},
	}
}

type generateRequest struct {
	Model   string                 `json:"model"`
	Prompt  string                 `json:"prompt"`
	Format  string                 `json:"format,omitempty"`
	Stream  bool                   `json:"stream"`
	Options map[string]interface{} `json:"options,omitempty"`
}

type generateResponse struct {
	Response string `json:"response"`
	Done     bool   `json:"done"`
}

// GenerateJSON calls Ollama with format:json and stream:false, returns the
// model's text response (which is expected to be a JSON blob).
func (c *OllamaClient) GenerateJSON(ctx context.Context, prompt string) (string, error) {
	body, err := json.Marshal(generateRequest{
		Model:  c.Model,
		Prompt: prompt,
		Format: "json",
		Stream: false,
		Options: map[string]interface{}{
			"temperature": 0.2,
			"num_ctx":     8192,
		},
	})
	if err != nil {
		return "", err
	}
	req, err := http.NewRequestWithContext(ctx, "POST", c.Host+"/api/generate", bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.HTTP.Do(req)
	if err != nil {
		if isConnRefused(err) {
			return "", fmt.Errorf("ollama is not running — start it with: ollama serve")
		}
		return "", fmt.Errorf("ollama call: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		data, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("ollama http %d: %s", resp.StatusCode, strings.TrimSpace(string(data)))
	}
	var env generateResponse
	if err := json.NewDecoder(resp.Body).Decode(&env); err != nil {
		return "", fmt.Errorf("decode ollama envelope: %w", err)
	}
	return env.Response, nil
}

// ParseInto runs Extract on model text into v.
func ParseInto(raw string, v any) error {
	return jsonx.Extract(raw, v)
}

func isConnRefused(err error) bool {
	return err != nil && strings.Contains(err.Error(), "connection refused")
}

// Ping checks the /api/tags endpoint is reachable.
func (c *OllamaClient) Ping(ctx context.Context) error {
	req, err := http.NewRequestWithContext(ctx, "GET", c.Host+"/api/tags", nil)
	if err != nil {
		return err
	}
	resp, err := c.HTTP.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return fmt.Errorf("ollama ping %d", resp.StatusCode)
	}
	return nil
}
