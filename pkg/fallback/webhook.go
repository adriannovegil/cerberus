package fallback

import (
	"context"

	"github.com/rs/zerolog/log"
)

// ExecuteWebHook execute the webhooks block
func ExecuteWebHook(ctx context.Context, cfg Config) error {
	log.Warn().Msgf("Executing fallback action: %s", cfg.Name)
	return nil
}
