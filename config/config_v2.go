package config

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/g4s8/gitstrap/ext"
	"github.com/g4s8/gitstrap/github"
	"github.com/g4s8/gitstrap/templates"
	"gopkg.in/yaml.v2"
	"io"
	"log"
	"os"
)

const (
	// V2 - latest version of config
	V2 = "v2"
)

// Config - latest version of the config
type Config struct {
	Version   string              `yaml:"version"`
	Github    *github.Repo        `yaml:"github"`
	Templates templates.Templates `yaml:"templates"`
	Ext       ext.Extensions      `yaml:"extensions,omitempty"`
}

// ParseReader - parse config from reader
func (y *Config) ParseReader(r io.Reader) error {
	dec := yaml.NewDecoder(r)
	if err := dec.Decode(y); err != nil {
		return fmt.Errorf("failed to decode config: %s", err)
	}
	if y.Version != V2 {
		// try v1 fallback
		cv1 := new(ConfigV1)
		if err := dec.Decode(cv1); err != nil {
			return fmt.Errorf("failed to decode v1 config: %s", err)
		}
		if cv1.Gitstrap.Version == V1 {
			log.Printf("v1 configs are deprecated." +
				"Use upgrade command to switch to new version")
			if err := cv1.validate(); err != nil {
				return err
			}
			cv1.upgrade(y)
		}
	}
	if y.Version != V2 {
		return fmt.Errorf("Unsupported version: %s", y.Version)
	}
	y.expand()
	if err := y.validate(); err != nil {
		return err
	}
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

var errNoGithubNode = errors.New("No `github` node")

func (y *Config) validate() error {
	if y.Github == nil {
		return errNoGithubNode
	}
	return nil
}

// expand - expand all strings with environment variable references
func (y *Config) expand() {
	for i, tpl := range y.Templates {
		y.Templates[i].Location = os.ExpandEnv(tpl.Location)
	}
}
