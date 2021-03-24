package spec

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

// Model of spec
type Model struct {
	Version  string      `yaml:"version"`
	Kind     string      `yaml:"kind"`
	Metadata *Metadata   `yaml:"metadata,omitempty"`
	Spec     interface{} `yaml:"-"`
}

const (
	// Version of spec
	Version = "v2.0-alpha"
)

// Metadata for spec
type Metadata struct {
	Name        string            `yaml:"name,omitempty"`
	Owner       string            `yaml:"owner,omitempty"`
	ID          *int64            `yaml:"id,omitempty"`
	Annotations map[string]string `yaml:"annotations,omitempty"`
}

const (
	// KindRepo - repository model kind
	KindRepo = "Repository"
	// KindReadme - repository readme model kind
	KindReadme = "Readme"
)

type errUnknownKind struct {
	kind string
}

func (e *errUnknownKind) Error() string {
	return fmt.Sprintf("unknown spec kind: `%s`", e.kind)
}

func (m *Model) MarshalYAML() (interface{}, error) {
	type M Model
	type temp struct {
		*M   `yaml:",inline"`
		Spec interface{} `yaml:"spec"`
	}
	t := &temp{(*M)(m), m.Spec}
	return t, nil
}

func (m *Model) UnmarshalYAML(value *yaml.Node) error {
	type M Model
	type temp struct {
		*M   `yaml:",inline"`
		Spec yaml.Node `yaml:"spec"`
	}
	obj := &temp{M: (*M)(m)}
	if err := value.Decode(obj); err != nil {
		return err
	}
	switch m.Kind {
	case KindRepo:
		m.Spec = new(Repo)
	case KindReadme:
		m.Spec = new(Readme)
	default:
		return &errUnknownKind{m.Kind}
	}
	return obj.Spec.Decode(m.Spec)
}
