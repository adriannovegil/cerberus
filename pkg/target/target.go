package target

import "time"

// Target data structure
type Target struct {
	ID         string        `yaml:"id"`
	CheckEvery time.Duration `yaml:"checkEvery"`
	Fallbacks  []string      `yaml:"fallbacks"`
}
