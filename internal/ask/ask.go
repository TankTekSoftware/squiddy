package ask

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"tankteksoftware.com/squiddy/internal/config"
	"tankteksoftware.com/squiddy/internal/llm"
)

func Run(question string) error {
	cfg, err := config.Load()
	if err != nil {
		if errors.Is(err, config.ErrNoConfig) {
			return fmt.Errorf("squiddy isn't set up yet. Run:\n  squiddy provider <anthropic|openai>\n  squiddy api_key <your-api-key>")
		}
		return err
	}

	streamer, err := llm.New(cfg)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	if err := streamer.Stream(ctx, llm.SystemPrompt, question, os.Stdout); err != nil {
		return err
	}
	fmt.Println()
	return nil
}
