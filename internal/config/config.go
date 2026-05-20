package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Provider string

const (
	ProviderAnthropic Provider = "anthropic"
	ProviderOpenAI    Provider = "openai"
)

type Config struct {
	Provider Provider `json:"provider"`
	APIKey   string   `json:"api_key"`
	Model    string   `json:"model"`
	BaseURL  string   `json:"base_url"`
}

// ErrNoConfig is returned when no saved config file exists yet.
var ErrNoConfig = errors.New("no squiddy configuration found")

// Path returns the absolute path to the squiddy config file.
// On Windows this is %AppData%\squiddy\config.json; elsewhere ~/.config/squiddy/config.json.
func Path() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "squiddy", "config.json"), nil
}

// Load reads config from the saved file.
// Returns ErrNoConfig if no config has been saved yet.
func Load() (Config, error) {
	cfg, err := loadFromFile()
	if errors.Is(err, os.ErrNotExist) {
		return Config{}, ErrNoConfig
	}
	if err != nil {
		return Config{}, err
	}
	if err := validate(cfg); err != nil {
		return Config{}, err
	}
	applyDefaults(&cfg)
	return cfg, nil
}

// SetAPIKey merges the given key into the saved config and writes it to disk.
func SetAPIKey(key string) error {
	key = strings.TrimSpace(key)
	if key == "" {
		return fmt.Errorf("api key cannot be empty")
	}
	cfg, err := loadFromFileLenient()
	if err != nil {
		return err
	}
	cfg.APIKey = key
	applyDefaults(&cfg)
	return Save(cfg)
}

// SetProvider validates and saves the provider into the config file.
// Switching providers resets model + base_url so the new provider's defaults
// get filled in (otherwise you'd be left pointing at the old provider's host).
func SetProvider(name string) error {
	p, err := parseProvider(name)
	if err != nil {
		return err
	}
	cfg, err := loadFromFileLenient()
	if err != nil {
		return err
	}
	if cfg.Provider != p {
		cfg.Model = ""
		cfg.BaseURL = ""
	}
	cfg.Provider = p
	applyDefaults(&cfg)
	return Save(cfg)
}

// Save writes the config to disk with restrictive permissions.
func Save(cfg Config) error {
	path, err := Path()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return err
	}
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o600)
}

func loadFromFile() (Config, error) {
	path, err := Path()
	if err != nil {
		return Config{}, err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return Config{}, err
	}
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return Config{}, fmt.Errorf("parse %s: %w", path, err)
	}
	cfg.Provider = Provider(strings.ToLower(strings.TrimSpace(string(cfg.Provider))))
	cfg.APIKey = strings.TrimSpace(cfg.APIKey)
	cfg.Model = strings.TrimSpace(cfg.Model)
	cfg.BaseURL = strings.TrimRight(strings.TrimSpace(cfg.BaseURL), "/")
	return cfg, nil
}

// loadFromFileLenient returns an empty Config (no error) if the file doesn't exist yet.
// Used by Set* helpers so the first save creates the file.
func loadFromFileLenient() (Config, error) {
	cfg, err := loadFromFile()
	if errors.Is(err, os.ErrNotExist) {
		return Config{}, nil
	}
	return cfg, err
}

func validate(cfg Config) error {
	if string(cfg.Provider) == "" {
		return fmt.Errorf("no provider set. Run: squiddy provider <anthropic|openai>")
	}
	if _, err := parseProvider(string(cfg.Provider)); err != nil {
		return err
	}
	if cfg.APIKey == "" {
		return fmt.Errorf("no api key set. Run: squiddy api_key <your-api-key>")
	}
	return nil
}

func parseProvider(s string) (Provider, error) {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "anthropic":
		return ProviderAnthropic, nil
	case "openai":
		return ProviderOpenAI, nil
	default:
		return "", fmt.Errorf("provider %q is not supported (expected 'anthropic' or 'openai')", s)
	}
}

func applyDefaults(cfg *Config) {
	if cfg.Model == "" {
		switch cfg.Provider {
		case ProviderAnthropic:
			cfg.Model = "claude-haiku-4-5"
		case ProviderOpenAI:
			cfg.Model = "gpt-4o-mini"
		}
	}
	if cfg.BaseURL == "" {
		switch cfg.Provider {
		case ProviderAnthropic:
			cfg.BaseURL = "https://api.anthropic.com"
		case ProviderOpenAI:
			cfg.BaseURL = "https://api.openai.com"
		}
	}
}
