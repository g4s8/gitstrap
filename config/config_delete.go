package config

import (
	"errors"
	"github.com/g4s8/gitstrap/context"
)

var errNoRepo = errors.New("no repo node in config")
var errWrongVersion = errors.New("expected V2 version")

// Delete repository
func (c *Config) Delete(ctx *context.Context) error {
	if c.Version != V2 {
		return errWrongVersion
	}
	if c.Github == nil {
		return errNoRepo
	}
	return c.Github.Delete(ctx)
}
