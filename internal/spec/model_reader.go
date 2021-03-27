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
	if err := d.Decode(m); err != nil {
		return err
	}
	if m.Metadata == nil {
		m.Metadata = new(Metadata)
	}
	return nil
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
