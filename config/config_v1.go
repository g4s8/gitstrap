package config

import (
	"errors"
	"github.com/g4s8/gitstrap/github"
	"github.com/g4s8/gitstrap/templates"
)

const (
	// V1 - first version of config
	V1 = "v1"
)

// ConfigV1 - first version of config (deprecated)
type ConfigV1 struct {
	Gitstrap *struct {
		Version string `yaml:"version"`
		Github  *struct {
			Repo *github.Repo `yaml:"repo"`
		} `yaml:"github"`
		Templates templates.TemplatesV1 `yaml:"templates"`
		Params    map[string]string     `yaml:"params"`
	} `yaml:"gitstrap"`
}

func (c *ConfigV1) upgrade(res *Config) {
	res.Version = V2
	res.Github = c.Gitstrap.Github.Repo
	res.Templates = c.Gitstrap.Templates.Upgrade(c.Gitstrap.Params)
}

var errNoGitstrapNodeV1 = errors.New("No `gitstrap` node")
var errNoGithubNodeV1 = errors.New("No `gitstrap.github` node")
var errNoRepoNodeV1 = errors.New("No `gitstrap.github.repo` node")

func (y *ConfigV1) validate() error {
	if y.Gitstrap == nil {
		return errNoGitstrapNodeV1
	}
	if y.Gitstrap.Github == nil {
		return errNoGithubNodeV1
	}
	if y.Gitstrap.Github.Repo == nil {
		return errNoRepoNodeV1
	}
	return nil
}
