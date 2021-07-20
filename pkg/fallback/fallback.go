package fallback

import (
	"context"

	"github.com/pkg/errors"
)

const (
	// WebhookType fallback execution type
	WebhookType = "webhook"
	// CommandType fallback execution type
	CommandType = "command"
	// ScriptType fallback execution type
	ScriptType = "script"
)

// Config data structure
type Config struct {
	Name         string `yaml:"name"`
	Type         string `yaml:"type"`
	AllowFailure bool   `yaml:"allow_failure"`
	Commands     []Command
}

// Command data structure
type Command struct {
	Command []string `yaml:"command"`
	Args    []string `yaml:"args"`
}

// Execute the fallback actions
func Execute(ctx context.Context, cfg Config) error {
	switch cfg.Type {
	case WebhookType:
		return ExecuteWebHooks(ctx, cfg)
	case CommandType:
		return ExecuteCommands(ctx, cfg)
	case ScriptType:
		return ExecuteScripts(ctx, cfg)
	default:
		return errors.Errorf("invalid proto '%s'", cfg.Type)
	}
}
