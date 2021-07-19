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
	Name    string   `yaml:"name"`
	Type    string   `yaml:"type"`
	Actions []string `yaml:"actions"`
}

// Execute the fallback actions
func Execute(ctx context.Context, cfg Config) error {
	switch cfg.Type {
	case WebhookType:
		return ExecuteWebHook(ctx, cfg)
	case CommandType:
		return ExecuteCommand(ctx, cfg)
	case ScriptType:
		return ExecuteScript(ctx, cfg)
	default:
		return errors.Errorf("invalid proto '%s'", cfg.Type)
	}
}
