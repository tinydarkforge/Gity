package services

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	tea "github.com/charmbracelet/bubbletea"
)

// TokenMsg is emitted for each streamed chunk.
type TokenMsg struct {
	Chunk string
}

// StreamDoneMsg is emitted once Ollama finishes sending.
type StreamDoneMsg struct {
	Full string
}

// StreamErrMsg is emitted on any transport or decode error.
type StreamErrMsg struct {
	Err error
}

// StreamCmd launches a goroutine that POSTs to /api/generate with stream:true
// and returns a tea.Cmd that delivers a single StreamDoneMsg / StreamErrMsg.
// Token chunks are pushed to the Program via the provided Send function.
//
// Using a send-channel keeps the API aligned with Bubble Tea's
// (*Program).Send — callers hold a program reference and pass its Send.
func (c *OllamaClient) StreamCmd(ctx context.Context, prompt string, send func(tea.Msg)) tea.Cmd {
	return func() tea.Msg {
		body, err := json.Marshal(generateRequest{
			Model:  c.Model,
			Prompt: prompt,
			Stream: true,
			Options: map[string]interface{}{
				"temperature": 0.2,
				"num_ctx":     8192,
			},
		})
		if err != nil {
			return StreamErrMsg{Err: err}
		}
		req, err := http.NewRequestWithContext(ctx, "POST", c.Host+"/api/generate", bytes.NewReader(body))
		if err != nil {
			return StreamErrMsg{Err: err}
		}
		req.Header.Set("Content-Type", "application/json")

		resp, err := c.HTTP.Do(req)
		if err != nil {
			if isConnRefused(err) {
				return StreamErrMsg{Err: fmt.Errorf("ollama is not running — start it with: ollama serve")}
			}
			return StreamErrMsg{Err: fmt.Errorf("ollama stream: %w", err)}
		}
		defer resp.Body.Close()
		if resp.StatusCode != 200 {
			return StreamErrMsg{Err: fmt.Errorf("ollama http %d", resp.StatusCode)}
		}

		scanner := bufio.NewScanner(resp.Body)
		buf := make([]byte, 0, 64*1024)
		scanner.Buffer(buf, 1024*1024)

		var full bytes.Buffer
		for scanner.Scan() {
			line := scanner.Bytes()
			if len(line) == 0 {
				continue
			}
			var chunk generateResponse
			if err := json.Unmarshal(line, &chunk); err != nil {
				continue
			}
			if chunk.Response != "" {
				full.WriteString(chunk.Response)
				if send != nil {
					send(TokenMsg{Chunk: chunk.Response})
				}
			}
			if chunk.Done {
				break
			}
		}
		if err := scanner.Err(); err != nil {
			return StreamErrMsg{Err: err}
		}
		return StreamDoneMsg{Full: full.String()}
	}
}
