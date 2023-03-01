package core

import (
	"fmt"
	"reflect"
	"time"
)

type Fields map[string]interface{}

func NewFields() *Fields {
	return &Fields{}
}

func (f *Fields) AddPair(key string, value interface{}) *Fields {
	(*f)[key] = value
	return f
}

func (f *Fields) AddArray(key string, values ...interface{}) *Fields {
	(*f)[key] = values
	return f
}

func (f *Fields) AddFields(key string, value *Fields) *Fields {
	(*f)[key] = *value
	return f
}

func (f *Fields) Export() Fields {
	return *f
}

func FlatMap(source map[string]interface{}) map[string]interface{} {
	return flatMap("", source)
}

func flatMap(root string, source map[string]interface{}) map[string]interface{} {
	var (
		keyPath string
		target  map[string]interface{}
	)
	target = make(map[string]interface{})
	for key, value := range source {
		if len(root) == 0 {
			keyPath = key
		} else {
			keyPath = fmt.Sprintf("%s.%s", root, key)
		}
		//filter struct and struct ptr
		rv := reflect.ValueOf(value)
		if rv.Kind() == reflect.Ptr {
			if rv.IsNil() {
				continue
			}
			value = rv.Elem().Interface()
			rv = reflect.ValueOf(value)
		}
		if rv.Kind() == reflect.Struct && rv.Type().String() != "time.Time" {
			continue
		}
		switch rv.Kind() {
		case reflect.Slice, reflect.Array:
			{
				switch value.(type) {
				case []interface{}:
				default:
					{
						temp := make([]interface{}, rv.Len())
						for i := 0; i < rv.Len(); i++ {
							temp[i] = rv.Index(i).Interface()
						}
						value = temp
					}
				}
			}
		case reflect.Map:
			{
				switch value.(type) {
				case map[string]interface{}:
				default:
					{
						temp := make(map[string]interface{}, rv.Len())
						for _, k := range rv.MapKeys() {
							temp[fmt.Sprintf("%v", k.Interface())] = rv.MapIndex(k).Interface()
						}
						value = temp
					}
				}
			}
		}

		switch value.(type) {
		case map[string]interface{}:
			{
				value = flatMap(keyPath, value.(map[string]interface{}))
				for rKey, rValue := range value.(map[string]interface{}) {
					target[rKey] = rValue
				}
			}
		case Fields:
			{
				value = flatMap(keyPath, value.(Fields))
				for rKey, rValue := range value.(map[string]interface{}) {
					target[rKey] = rValue
				}
			}
		case []interface{}:
			{
				value = flatArray(keyPath, value.([]interface{}))
				for rKey, rValue := range value.(map[string]interface{}) {
					target[rKey] = rValue
				}
			}
		case time.Time:
			{
				t := value.(time.Time)
				target[keyPath] = t.Format(time.RFC3339Nano)
			}
		default:
			{
				target[keyPath] = value
			}
		}
	}
	return target
}

func flatArray(root string, source []interface{}) map[string]interface{} {
	var (
		keyPath string
		target  map[string]interface{}
	)
	target = make(map[string]interface{})
	for index, value := range source {
		if len(root) == 0 {
			keyPath = fmt.Sprintf("%d", index)
		} else {
			keyPath = fmt.Sprintf("%s.%d", root, index)
		}
		//filter struct and struct ptr
		rv := reflect.ValueOf(value)
		if rv.Kind() == reflect.Ptr {
			value = rv.Elem().Interface()
			rv = reflect.ValueOf(value)
		}
		if rv.Kind() == reflect.Struct && rv.Type().String() != "time.Time" {
			continue
		}
		switch rv.Kind() {
		case reflect.Slice, reflect.Array:
			{
				switch value.(type) {
				case []interface{}:
				default:
					{
						temp := make([]interface{}, rv.Len())
						for i := 0; i < rv.Len(); i++ {
							temp[i] = rv.Index(i).Interface()
						}
						value = temp
					}
				}
			}
		case reflect.Map:
			{
				switch value.(type) {
				case map[string]interface{}:
				default:
					{
						temp := make(map[string]interface{}, rv.Len())
						for _, k := range rv.MapKeys() {
							temp[fmt.Sprintf("%v", k.Interface())] = rv.MapIndex(k).Interface()
						}
						value = temp
					}
				}
			}
		}

		switch value.(type) {
		case map[string]interface{}:
			{
				value = flatMap(keyPath, value.(map[string]interface{}))
				for rKey, rValue := range value.(map[string]interface{}) {
					target[rKey] = rValue
				}
			}
		case Fields:
			{
				value = flatMap(keyPath, value.(Fields))
				for rKey, rValue := range value.(map[string]interface{}) {
					target[rKey] = rValue
				}
			}
		case []interface{}:
			{
				value = flatArray(keyPath, value.([]interface{}))
				for rKey, rValue := range value.(map[string]interface{}) {
					target[rKey] = rValue
				}
			}
		case time.Time:
			{
				t := value.(time.Time)
				target[keyPath] = t.Format(time.RFC3339Nano)
			}
		default:
			{
				target[keyPath] = value
			}
		}
	}
	return target
}
