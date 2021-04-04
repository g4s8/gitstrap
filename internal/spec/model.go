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
	Version = "v2"
)

// Kind of specification
type Kind string

func (k Kind) validate() error {
	for _, v := range [...]Kind{KindRepo, KindReadme, KindOrg, KindHook, KindTeam} {
		if k == v {
			return nil
		}
	}
	return &errUnknownKind{k}
}

// Require this kind to be another kind
// panics with ErrInvalidKind error if doesn't.
// Could be handlerd with ErrInvalidKind.RecoverHandler
func (k Kind) Require(req Kind) {
	if k != req {
		panic(&ErrInvalidKind{Expected: req, Actual: k})
	}
}

const (
	// KindRepo - repository model kind
	KindRepo = Kind("Repository")
	// KindReadme - repository readme model kind
	KindReadme = Kind("Readme")
	// KindOrg - organization model kind
	KindOrg = Kind("Organization")
	// KindHook - repository webhook
	KindHook = Kind("WebHook")
	KindTeam = Kind("Team")
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

// ErrInvalidKind - error that kind is not the value as expected
type ErrInvalidKind struct {
	// Expected and Actual values of kind
	Expected, Actual Kind
}

func (e *ErrInvalidKind) Error() string {
	return fmt.Sprintf("Invalid model kind: expects `%s` but was `%s`", e.Expected, e.Actual)
}

// RecoverHandler could be used to catch this error on panic with defer
func (e *ErrInvalidKind) RecoverHandler(out *error) {
	if rec := recover(); rec != nil {
		if err, ok := rec.(error); ok && errors.Is(err, e) {
			*out = err
		} else {
			panic(rec)
		}
	}
}

type errInvalidSpecType struct {
	spec interface{}
}

func (e *errInvalidSpecType) Error() string {
	return fmt.Sprintf("Invalid spec type `%T`", e.spec)
}

var errSpecIsNil = errors.New("Model spec is nil")

// GetSpec extracts a spec from model
func (m *Model) GetSpec(out interface{}) (re error) {
	if m.Spec == nil {
		return errSpecIsNil
	}
	errh := new(ErrInvalidKind)
	defer errh.RecoverHandler(&re)
	var ok bool
	switch s := out.(type) {
	case *Repo:
		m.Kind.Require(KindRepo)
		var t *Repo
		t, ok = m.Spec.(*Repo)
		*s = *t
	case *Readme:
		var t *Readme
		m.Kind.Require(KindReadme)
		t, ok = m.Spec.(*Readme)
		*s = *t
	case *Org:
		var t *Org
		m.Kind.Require(KindOrg)
		t, ok = m.Spec.(*Org)
		*s = *t
	case *Hook:
		var t *Hook
		m.Kind.Require(KindHook)
		t, ok = m.Spec.(*Hook)
		*s = *t
	case *Team:
		var t *Team
		m.Kind.Require(KindTeam)
		t, ok = m.Spec.(*Team)
		*s = *t
	default:
		return &errInvalidSpecType{s}
	}
	if !ok {
		return &errInvalidSpecType{m.Spec}
	}
	return nil
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
	case KindOrg:
		m.Spec = new(Org)
	case KindHook:
		m.Spec = new(Hook)
	case KindTeam:
		m.Spec = new(Team)
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
