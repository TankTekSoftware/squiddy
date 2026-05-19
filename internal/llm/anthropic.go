package llm

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"tankteksoftware.com/squiddy/internal/config"
)

type anthropicStreamer struct {
	cfg config.Config
}

type anthropicMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type anthropicRequest struct {
	Model     string             `json:"model"`
	MaxTokens int                `json:"max_tokens"`
	System    string             `json:"system"`
	Messages  []anthropicMessage `json:"messages"`
	Stream    bool               `json:"stream"`
}

type anthropicDelta struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type anthropicEvent struct {
	Type  string         `json:"type"`
	Delta anthropicDelta `json:"delta"`
}

func (s *anthropicStreamer) Stream(ctx context.Context, system, user string, out io.Writer) error {
	body, err := json.Marshal(anthropicRequest{
		Model:     s.cfg.Model,
		MaxTokens: 1024,
		System:    system,
		Messages:  []anthropicMessage{{Role: "user", Content: user}},
		Stream:    true,
	})
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, s.cfg.BaseURL+"/v1/messages", bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("x-api-key", s.cfg.APIKey)
	req.Header.Set("anthropic-version", "2023-06-01")
	req.Header.Set("content-type", "application/json")
	req.Header.Set("accept", "text/event-stream")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode/100 != 2 {
		return errorFromResponse(resp, "anthropic")
	}

	scanner := bufio.NewScanner(resp.Body)
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)
	for scanner.Scan() {
		line := scanner.Text()
		if !strings.HasPrefix(line, "data: ") {
			continue
		}
		payload := strings.TrimPrefix(line, "data: ")
		var ev anthropicEvent
		if err := json.Unmarshal([]byte(payload), &ev); err != nil {
			continue
		}
		if ev.Type == "content_block_delta" && ev.Delta.Type == "text_delta" {
			if _, err := io.WriteString(out, ev.Delta.Text); err != nil {
				return err
			}
		}
	}
	return scanner.Err()
}

func errorFromResponse(resp *http.Response, provider string) error {
	const max = 2048
	b, _ := io.ReadAll(io.LimitReader(resp.Body, max))
	return fmt.Errorf("%s API returned %s: %s", provider, resp.Status, strings.TrimSpace(string(b)))
}
