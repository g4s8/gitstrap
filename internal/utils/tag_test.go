package utils

import (
	"reflect"
	"testing"

	m "github.com/g4s8/go-matchers"
)

type original struct {
	A int    `yaml:"a,omitempty,test"`
	B string `yaml:"b,omitempty"`
	C bool   `yaml:"c,test,omitempty"`
	D int    `json:"d,omitempty"`
	E string
}

type target struct {
	A int    `yaml:"a,test"`
	B string `yaml:"b"`
	C bool   `yaml:"c,test"`
	D int    `json:"d,omitempty"`
	E string
}

func TestRemoveTagsOmitempty(t *testing.T) {
	assert := m.Assert(t)
	o := original{1, "a", true, 2, "b"}
	target := target{1, "a", true, 2, "b"}
	modified, err := RemoveTagsOmitempty(o, "yaml")
	if err != nil {
		t.Error(err)
	}
	tVal := reflect.ValueOf(target)
	mVal := reflect.ValueOf(modified)
	tType := tVal.Type()
	mType := mVal.Type()
	for i := 0; i < tVal.NumField(); i++ {
		wantField := tVal.Field(i).Interface()
		gotField := mVal.Field(i).Interface()
		assert.That("Fields are equal", gotField, m.Eq(wantField))
		wantT := tType.Field(i)
		gotT := mType.Field(i)
		assert.That("StructFields are equal", gotT, m.Eq(wantT))
	}
}
