package gitstrap

import (
	"bufio"
	"errors"
	"fmt"
	"gopkg.in/yaml.v2"
	"io"
	"os"
)

const (
	// V1 - first version of config
	V1 = "v1"
)

// Config - gitstrap config
type Config struct {
	Gitstrap *struct {
		Version string `yaml:"version"`
		Github  *struct {
			Repo *struct {
				Name        *string `yaml:"name"`
				Description *string `yaml:"description"`
				Private     *bool   `yaml:"private"`
				AutoInit    *bool   `yaml:"autoInit"`
				Hooks       []struct {
					URL    string   `yaml:"url"`
					Type   string   `yaml:"type"`
					Events []string `yaml:"events"`
					Active *bool    `yaml:"active"`
				} `yaml:"hooks"`
				Collaborators []string `yaml:"collaborators"`
			} `yaml:"repo"`
		} `yaml:"github"`
		Templates []struct {
			Name     string `yaml:"name"`
			Location string `yaml:"location"`
			URL      string `yaml:"url"`
		} `yaml:"templates"`
		Params map[string]string `yaml:"params"`
	} `yaml:"gitstrap"`
}

// ParseReader - parse config from reader
func (y *Config) ParseReader(r io.Reader) error {
	if err := yaml.NewDecoder(r).Decode(y); err != nil {
		return err
	}
	if y.Gitstrap.Version != V1 {
		return fmt.Errorf("Unsupported version: %s", y.Gitstrap.Version)
	}
	y.expand()
	return nil
}

// ParseFile - parse config from file
func (y *Config) ParseFile(name string) error {
	f, err := os.Open(name)
	if err != nil {
		return fmt.Errorf("Failed to open config file: %s", err)
	}
	if err = y.ParseReader(bufio.NewReader(f)); err != nil {
		return err
	}
	err = f.Close()
	return err
}

var errNoGitstrapNode = errors.New("No `gitstrap` node")
var errNoGithubNode = errors.New("No `gitstrap.github` node")
var errNoRepoNode = errors.New("No `gitstrap.github.repo` node")

// Validate - validate config, return error if invalid
func (y *Config) Validate() error {
	if y.Gitstrap == nil {
		return errNoGitstrapNode
	}
	if y.Gitstrap.Github == nil {
		return errNoGithubNode
	}
	if y.Gitstrap.Github.Repo == nil {
		return errNoRepoNode
	}
	return nil
}

// expand - expand all strings with environment variable references
func (y *Config) expand() {
	for i, tpl := range y.Gitstrap.Templates {
		y.Gitstrap.Templates[i].Location = os.ExpandEnv(tpl.Location)
	}
}
