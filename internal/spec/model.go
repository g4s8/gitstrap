package spec

import (
	"errors"
	"fmt"

	"strings"

	"gopkg.in/yaml.v3"
)

// Model of spec
type Model struct {
	Version  string      `yaml:"version"`
	Kind     Kind        `yaml:"kind"`
	Metadata *Metadata   `yaml:"metadata,omitempty"`
	Spec     interface{} `yaml:"-"`
}

const (
	// Version of spec
	Version = "v2.0-alpha"
)

// Kind of specification
type Kind string

func (k Kind) validate() error {
	for _, v := range [...]Kind{KindRepo, KindReadme, KindOrg} {
		if k == v {
			return nil
		}
	}
	return &errUnknownKind{k}
}

const (
	// KindRepo - repository model kind
	KindRepo = Kind("Repository")
	// KindReadme - repository readme model kind
	KindReadme = Kind("Readme")
	// KindOrg - organization model kind
	KindOrg = Kind("Organization")
)

// NewModel with kind
func NewModel(kind Kind) (*Model, error) {
	if err := kind.validate(); err != nil {
		return nil, err
	}
	m := new(Model)
	m.Version = Version
	m.Kind = kind
	m.Metadata = new(Metadata)
	return m, nil
}

type errUnknownKind struct {
	kind Kind
}

func (e *errUnknownKind) Error() string {
	return fmt.Sprintf("unknown spec kind: `%s`", e.kind)
}

type errInvalidKind struct {
	actual Kind
	expect Kind
}

func (e *errInvalidKind) Error() string {
	return fmt.Sprintf("Invalid model kind: expects `%s` but was `%s`", e.expect, e.actual)
}

type errInvalidSpecType struct {
	expect interface{}
	spec   interface{}
}

func (e *errInvalidSpecType) Error() string {
	return fmt.Sprintf("Invalid spec type: expects `%T` but was `%T`", e.expect, e.spec)
}

var errSpecIsNil = errors.New("Model spec is nil")

// RepoSpec extracted from model
func (m *Model) RepoSpec(r *Repo) error {
	if m.Kind != KindRepo {
		return &errInvalidKind{m.Kind, KindRepo}
	}
	if m.Spec == nil {
		return errSpecIsNil
	}
	switch spec := m.Spec.(type) {
	case *Repo:
		*r = *spec
		return nil
	default:
		return &errInvalidSpecType{r, spec}
	}
}

// ReadmeSpec extracted from model
func (m *Model) ReadmeSpec(r *Readme) error {
	if m.Kind != KindReadme {
		return &errInvalidKind{m.Kind, KindReadme}
	}
	if m.Spec == nil {
		return errSpecIsNil
	}
	switch spec := m.Spec.(type) {
	case *Readme:
		*r = *spec
		return nil
	default:
		return &errInvalidSpecType{r, spec}
	}
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

func (m *Model) Info() string {
	sb := new(strings.Builder)
	sb.WriteString(string(m.Kind))
	sb.WriteString(": ")
	sb.WriteString(m.Metadata.Info())
	return sb.String()
}
