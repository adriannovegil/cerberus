package fallback

import (
	"context"

	"github.com/rs/zerolog/log"
)

// ExecuteScripts execute the scripts block
func ExecuteScripts(ctx context.Context, cfg Config) error {
	for _, cmd := range cfg.Commands {
		err := executeScript(ctx, cmd)
		if err != nil {
			if cfg.AllowFailure {
				log.Warn().Msgf("Executing fallback action: %s", cfg.Name)
				continue
			}
			return err
		}
	}
	return nil
}

func executeScript(ctx context.Context, script Command) error {
	return nil
}
