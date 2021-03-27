package spec

import (
	"bufio"
	"errors"
	"gopkg.in/yaml.v3"
	"io"
	"os"
	"path/filepath"
)

var errReadNoDocument = errors.New("Not YAML document")

// FromDecoder creates model from yaml decoder
func (m *Model) FromDecoder(d *yaml.Decoder) error {
	type Doc struct {
		Model `yaml:"inline"`
	}
	doc := new(Doc)
	if err := d.Decode(&doc); err != nil {
		return err
	}
	if d == nil {
		return errReadNoDocument
	}
	*m = doc.Model
	if m.Metadata == nil {
		m.Metadata = new(Metadata)
	}
	return nil
}

// FromReader creates model from io reader
func (m *Model) FromReader(r io.Reader) error {
	return m.FromDecoder(yaml.NewDecoder(r))
}

// ReadFile models from file
func ReadFile(name string) ([]*Model, error) {
	fn, err := filepath.Abs(name)
	if err != nil {
		return nil, err
	}
	f, err := os.Open(fn)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return ReadStream(bufio.NewReader(f))
}

// ReadStream of models from reader
func ReadStream(r io.Reader) ([]*Model, error) {
	dec := yaml.NewDecoder(r)
	res := make([]*Model, 0)
	for {
		model := new(Model)
		err := model.FromDecoder(dec)
		if err == nil {
			res = append(res, model)
		} else if errors.Is(err, errReadNoDocument) {
			continue
		} else if errors.Is(err, io.EOF) {
			break
		} else {
			return nil, err
		}
	}
	return res, nil
}
