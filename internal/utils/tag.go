package utils

import (
	"fmt"
	"reflect"

	"github.com/fatih/structtag"
)

func RemoveTagsOmitempty(s interface{}, key string) (interface{}, error) {
	var value reflect.Value
	if reflect.TypeOf(s).Kind() == reflect.Ptr {
		value = reflect.Indirect(reflect.ValueOf(s))
	} else {
		value = reflect.ValueOf(s)
	}
	t := value.Type()
	sf := make([]reflect.StructField, 0)
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tag, err := removeOmitempty(field.Tag, "yaml")
		if err != nil {
			return nil, err
		}
		field.Tag = tag
		sf = append(sf, field)
	}
	newType := reflect.StructOf(sf)
	newValue := value.Convert(newType)
	return newValue.Interface(), nil
}

func removeOmitempty(tag reflect.StructTag, key string) (reflect.StructTag, error) {
	newTag := *new(reflect.StructTag)
	tags, err := structtag.Parse(string(tag))
	if err != nil {
		return newTag, err
	}
	yamlTag, err := tags.Get(key)
	if err != nil {
		return newTag, err
	}
	for i, v := range yamlTag.Options {
		if v == "omitempty" {
			yamlTag.Options = append(yamlTag.Options[:i], yamlTag.Options[i+1:]...)
		}
	}
	stringTags := fmt.Sprintf(`%v`, tags)
	newTag = reflect.StructTag(stringTags)
	return newTag, nil
}
