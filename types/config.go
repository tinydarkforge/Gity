package types

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
)

type Config struct {
	Repo        string `json:"repo"`
	Model       string `json:"model"`
	OllamaHost  string `json:"ollamaHost"`
	TimeoutSec  int    `json:"timeoutSec"`
	MaxTurns    int    `json:"maxTurns"`
	Sound       bool   `json:"sound"`
	Debug       bool   `json:"debug"`
	TemplateDir string `json:"templateDir"`
}

func DefaultConfig() Config {
	return Config{
		Repo:        "",
		Model:       "llama3",
		OllamaHost:  "http://localhost:11434",
		TimeoutSec:  120,
		MaxTurns:    6,
		Sound:       true,
		Debug:       false,
		TemplateDir: ".github/ISSUE_TEMPLATES",
	}
}

func ConfigPath() (string, error) {
	if p := os.Getenv("INTAKE_CONFIG"); p != "" {
		return p, nil
	}
	base := os.Getenv("XDG_CONFIG_HOME")
	if base == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		base = filepath.Join(home, ".config")
	}
	return filepath.Join(base, "intake", "config.json"), nil
}

func LoadConfig() (Config, error) {
	cfg := DefaultConfig()
	path, err := ConfigPath()
	if err != nil {
		return cfg, err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return applyEnv(cfg), nil
		}
		return cfg, err
	}
	if err := json.Unmarshal(data, &cfg); err != nil {
		return cfg, fmt.Errorf("parse config %s: %w", path, err)
	}
	return applyEnv(cfg), nil
}

func SaveConfig(cfg Config) error {
	path, err := ConfigPath()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o600)
}

func applyEnv(cfg Config) Config {
	if v := os.Getenv("INTAKE_REPO"); v != "" {
		cfg.Repo = v
	}
	if v := os.Getenv("INTAKE_MODEL"); v != "" {
		cfg.Model = v
	}
	if v := os.Getenv("OLLAMA_HOST"); v != "" {
		cfg.OllamaHost = v
	}
	if v := os.Getenv("INTAKE_TIMEOUT"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			cfg.TimeoutSec = n
		}
	}
	if v := os.Getenv("INTAKE_MAX_TURNS"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			cfg.MaxTurns = n
		}
	}
	if v := os.Getenv("INTAKE_TEMPLATE_DIR"); v != "" {
		cfg.TemplateDir = v
	}
	if v := os.Getenv("INTAKE_DEBUG"); v == "1" || v == "true" {
		cfg.Debug = true
	}
	return cfg
}
