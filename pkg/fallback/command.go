package fallback

import (
	"context"
	"strings"

	"devcircus.com/cerberus/pkg/util/shell"
	"github.com/rs/zerolog/log"
)

// ExecuteCommands execute the commands block
func ExecuteCommands(ctx context.Context, cfg Config) error {
	for _, cmd := range cfg.Commands {
		err := executeCommand(ctx, cmd)
		if err != nil {
			if cfg.AllowFailure {
				log.Warn().Msgf("Error executing the command: %s %s ")
				continue
			}
			return err
		}
	}
	return nil
}

func executeCommand(ctx context.Context, cmd Command) error {
	err := shell.RunCommand(
		strings.Join(cmd.Command, " "),
		cmd.Args...)
	if err != nil {
		return err
	}
	return nil
}
