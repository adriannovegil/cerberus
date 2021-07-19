package fallback

import (
	"context"

	"github.com/rs/zerolog/log"
)

// ExecuteScript execute the scripts block
func ExecuteScript(ctx context.Context, cfg Config) error {
	log.Warn().Msgf("Executing fallback action: %s", cfg.Name)
	return nil
}
