package utils

import (
	"fmt"
	"reflect"

	"github.com/fatih/structtag"
)

//RemoveTagsOmitempty returns struct of new type with same fields and values as s,
//but removes option "omitempty" from tags with specified key.
func RemoveTagsOmitempty(s interface{}, key string) (interface{}, error) {
	var value reflect.Value
	if reflect.TypeOf(s).Kind() == reflect.Ptr {
		value = reflect.Indirect(reflect.ValueOf(s))
	} else {
		value = reflect.ValueOf(s)
	}
	t := value.Type()
	nf := t.NumField()
	sf := make([]reflect.StructField, nf)
	for i := 0; i < nf; i++ {
		field := t.Field(i)
		tag, err := removeOmitempty(field.Tag, key)
		if err != nil {
			return nil, err
		}
		field.Tag = *tag
		sf[i] = field
	}
	newType := reflect.StructOf(sf)
	newValue := value.Convert(newType)
	return newValue.Interface(), nil
}

func removeOmitempty(tag reflect.StructTag, key string) (*reflect.StructTag, error) {
	tags, err := structtag.Parse(string(tag))
	if err != nil {
		return nil, err
	}
	yamlTag, err := tags.Get(key)
	if err != nil {
		return &tag, nil
	}
	for i, v := range yamlTag.Options {
		if v == "omitempty" {
			yamlTag.Options = append(yamlTag.Options[:i], yamlTag.Options[i+1:]...)
			break
		}
	}
	stringTags := fmt.Sprintf(`%v`, tags)
	newTag := reflect.StructTag(stringTags)
	return &newTag, nil
}
