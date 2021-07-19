package fallback

import (
	"context"

	"github.com/rs/zerolog/log"
)

// ExecuteCommand execute the commands block
func ExecuteCommand(ctx context.Context, cfg Config) error {
	log.Warn().Msgf("Executing fallback action: %s", cfg.Name)
	return nil
}
