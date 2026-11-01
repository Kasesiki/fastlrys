package fastlrys

import (
	"encoding/json"
	"reflect"
)

type resq struct {
	Err     bool
	Message string
	Body    map[string]any
}

func BuildResq() *resq {
	return &resq{
		Body: make(map[string]any),
	}
}

func (re *resq) Error(message string) *resq {
	re.Err = true
	re.Message = message
	return re
}

func (re *resq) CoverBody(data any) *resq {
	v := reflect.ValueOf(data)
	t := reflect.TypeOf(data)
	StructDrop(t, v)
	return re
}

func StructDrop(t reflect.Type, v reflect.Value) any {
	if v.Kind() == reflect.Struct {
		data := make(map[string]any)
		for i := 0; i < v.NumField(); i++ {
			data[t.Field(i).Name] = StructDrop(t.Field(i).Type, v.Field(i))
		}
		return data
	} else {
		return v.Interface()
	}
}

func (re *resq) AddBody(key string, value any) *resq {
	re.Body[key] = value
	return re
}

func (re *resq) DelBody(key string) *resq {
	delete(re.Body, key)
	return re
}

func (re resq) BuildJson() string {
	js, _ := json.Marshal(re)
	return string(js)
}
