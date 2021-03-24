package spec

import (
	"testing"

	m "github.com/g4s8/go-matchers"
	"gopkg.in/yaml.v3"
)

func Test_unmarshall(t *testing.T) {
	assert := m.Assert(t)
	src := []byte(`
version: "v1"
kind: "Repository"
metadata:
  name: gitstrap
  owner: g4s8
  annotations:
    foo: bar
    baz: 4
spec:
  id: 1
  owner: "testing"
  description: "for test"
`)
	model := new(Model)
	assert.That("Unmarshal model without errors", yaml.Unmarshal(src, &model), m.Nil())
	assert.That("Model kind is OK", model.Kind, m.Eq(KindRepo))
	assert.That("Model version is OK", model.Version, m.Eq(Version))
	t.Run("Model metadata is OK", func(t *testing.T) {
		assert := m.Assert(t)
		assert.That("Name is OK", model.Metadata.Name, m.Eq("gitstrap"))
		assert.That("Owner is OK", model.Metadata.Owner, m.Eq("g4s8"))
		assert.That("Annotations[0] is OK", model.Metadata.Annotations["foo"], m.Eq("bar"))
		assert.That("Annotations[1] is OK", model.Metadata.Annotations["baz"], m.Eq("4"))
	})
	t.Run("Model spec is correct", func(t *testing.T) {
		assert := m.Assert(t)
		repo, typeok := model.Spec.(*Repo)
		assert.That("Spec type is OK", typeok, m.Is(true))
		assert.That("Repo ID is OK", repo.ID, m.Eq(int64(1)))
		assert.That("Repo description is OK", *repo.Description, m.Eq("for test"))
	})
}

func Test_marshall(t *testing.T) {
	assert := m.Assert(t)
	model := new(Model)
	model.Version = Version
	model.Kind = KindRepo
	repo := new(Repo)
	repo.ID = 1
	model.Spec = repo
	_, err := yaml.Marshal(model)
	assert.That("Marshal without errors", err, m.Nil())
}
