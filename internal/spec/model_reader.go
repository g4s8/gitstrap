package spec

import (
	"bufio"
	"gopkg.in/yaml.v3"
	"io"
	"os"
	"path/filepath"
)

// FromDecoder creates model from yaml decoder
func (m *Model) FromDecoder(d *yaml.Decoder) error {
	return d.Decode(m)
}

// FromReader creates model from io reader
func (m *Model) FromReader(r io.Reader) error {
	return m.FromDecoder(yaml.NewDecoder(r))
}

// FromFile creates model from file
func (m *Model) FromFile(name string) error {
	fn, err := filepath.Abs(name)
	if err != nil {
		return err
	}
	f, err := os.Open(fn)
	if err != nil {
		return err
	}
	defer f.Close()
	return m.FromReader(bufio.NewReader(f))
}
