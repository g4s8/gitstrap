package config

import (
	"bytes"
	"gopkg.in/yaml.v2"
)

func (c *Config) ToYaml() (string, error) {
	var buf bytes.Buffer
	enc := yaml.NewEncoder(&buf)
	defer enc.Close()
	if err := enc.Encode(c); err != nil {
		return "", err
	}
	return buf.String(), nil
}
