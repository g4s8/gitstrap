package config

import (
	"github.com/g4s8/gitstrap/context"
	"github.com/g4s8/gitstrap/github"
)

// Get config from remote repository and fill fields
// with actual values.
func (c *Config) Get(ctx *context.Context, name string) error {
	c.Version = V2
	c.Github = new(github.Repo)
	if err := c.Github.Get(ctx, name); err != nil {
		return err
	}
	return nil
}
