package ping

import "devcircus.com/cerberus/pkg/target"

// Config data structure
type Config struct {
	target.Target `yaml:"-,inline"`
}
