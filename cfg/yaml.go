package cfg

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io"
)

const (
	V1 = "v1"
)

type YamlConfig struct {
	Gitstrap Gitstrap `yaml:"gitstrap"`
}

type Gitstrap struct {
	Version   string            `yaml:"version"`
	Github    *Github           `yaml:"github"`
	Templates []Template        `yaml:"templates"`
	Params    map[string]string `yaml:"params"`
}

type Github struct {
	Repo *GithubRepo `yaml:"repo"`
}

type GithubRepo struct {
	Name          *string    `yaml:"name"`
	Description   *string    `yaml:"description"`
	Private       *bool      `yaml:"private"`
	AutoInit      *bool      `yaml:"autoInit"`
	Hooks         []RepoHook `yaml:"hooks"`
	Collaborators []string   `yaml:"collaborators"`
}

type Template struct {
	Name     string `yaml:"name"`
	Location string `yaml:"location"`
}

type RepoHook struct {
	Url    string   `yaml:"url"`
	Type   string   `yaml:"type"`
	Events []string `yaml:"events"`
	Active *bool    `yaml:"active"`
}

func (y *YamlConfig) Parse(r io.Reader) error {
	err := yaml.NewDecoder(r).Decode(y)
	if y.Gitstrap.Version != V1 {
		return fmt.Errorf("Unsupported version: %s", y.Gitstrap.Version)
	}
	return err
}
