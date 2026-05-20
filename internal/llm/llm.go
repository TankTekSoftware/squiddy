package llm

import (
	"context"
	"fmt"
	"io"

	"tankteksoftware.com/squiddy/internal/config"
)

type Streamer interface {
	Stream(ctx context.Context, system, user string, out io.Writer) error
}

func New(cfg config.Config) (Streamer, error) {
	switch cfg.Provider {
	case config.ProviderAnthropic:
		return &anthropicStreamer{cfg: cfg}, nil
	case config.ProviderOpenAI:
		return &openAIStreamer{cfg: cfg}, nil
	default:
		return nil, fmt.Errorf("unknown provider %q", cfg.Provider)
	}
}
