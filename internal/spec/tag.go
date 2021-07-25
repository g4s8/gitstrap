package spec

import (
	"fmt"
	"reflect"

	"github.com/fatih/structtag"
)

func RemoveTagsOmitempty(s interface{}) (interface{}, error) {
	var value reflect.Value
	if reflect.TypeOf(s).Kind() == reflect.Ptr {
		value = reflect.Indirect(reflect.ValueOf(s))
	} else {
		value = reflect.ValueOf(s)
	}
	t := value.Type()
	sf := make([]reflect.StructField, 0)
	for i := 0; i < t.NumField(); i++ {
		sf = append(sf, t.Field(i))
		tag := t.Field(i).Tag
		tag, err := editTag(tag)
		if err != nil {
			return nil, err
		}
		sf[i].Tag = tag
	}
	newType := reflect.StructOf(sf)
	newValue := value.Convert(newType)
	return newValue.Interface(), nil
}

func editTag(tag reflect.StructTag) (reflect.StructTag, error) {
	newTag := *new(reflect.StructTag)
	tags, err := structtag.Parse(string(tag))
	if err != nil {
		return newTag, err
	}
	yamlTag, err := tags.Get("yaml")
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