package fallback

import (
	"context"

	"github.com/rs/zerolog/log"
)

// ExecuteWebHooks execute the webhooks block
func ExecuteWebHooks(ctx context.Context, cfg Config) error {
	for _, cmd := range cfg.Commands {
		err := executeWebHook(ctx, cmd)
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

func executeWebHook(ctx context.Context, webHook Command) error {
	return nil
}
