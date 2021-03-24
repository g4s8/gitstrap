package spec

import (
	"github.com/g4s8/gitstrap/internal/view"
	"gopkg.in/yaml.v3"
)

// ModelFormat - format of model to print
type ModelFormat interface {
	ToView(*Model) (view.Printable, error)
}

type mfYaml struct{}

var MfYaml ModelFormat = &mfYaml{}

func (mf *mfYaml) ToView(m *Model) (view.Printable, error) {
	bin, err := yaml.Marshal(m)
	if err != nil {
		return nil, err
	}
	return &mfYamlView{bin}, nil
}

type mfYamlView struct {
	bytes []byte
}

func (v *mfYamlView) PrintOn(p view.Printer) {
	p.Print(string(v.bytes))
}
