package jsondb

import (
	"encoding/json"
	"reflect"
	"strings"
)

type Document map[string]interface{}

func (document Document) Getter(names ...string) interface{} {
	var (
		ok bool
		v  interface{}
	)
	for _, name := range names {
		if v, ok = document[name]; ok {
			return v
		}
	}
	return v
}

func (document Document) Decode(v interface{}) (err error) {
	var (
		buf []byte
	)
	if buf, err = json.Marshal(document); err != nil {
		return
	}
	err = json.Unmarshal(buf, v)
	return
}

func (document Document) Unmarshal(refType reflect.Type) (val reflect.Value, err error) {
	var (
		buf []byte
	)
	val = reflect.New(refType)
	if buf, err = json.Marshal(document); err != nil {
		return
	}
	if err = json.Unmarshal(buf, val.Interface()); err != nil {
		return
	}
	if refType.Kind() == reflect.Ptr {
		return val.Elem(), nil
	} else {
		return val, nil
	}
}

func toDocument(v interface{}) Document {
	var (
		pos         int
		fieldName   string
		refValue    reflect.Value
		refType     reflect.Type
		structField reflect.StructField
		fieldValue  reflect.Value
	)
	doc := make(map[string]interface{})
	refValue = reflect.Indirect(reflect.ValueOf(v))
	refType = refValue.Type()
	for i := 0; i < refType.NumField(); i++ {
		structField = refType.Field(i)
		fieldValue = refValue.Field(i)
		fieldName = structField.Tag.Get("json")
		if fieldName == "-" {
			continue
		}
		if pos = strings.IndexByte(fieldName, ','); pos > -1 {
			fieldName = fieldName[:i]
		}
		if fieldName == "" {
			fieldName = structField.Name
		}
		switch fieldValue.Kind() {
		case reflect.Struct:
			doc[fieldName] = toDocument(fieldValue.Interface())
		default:
			doc[fieldName] = fieldValue.Interface()
		}
	}
	return doc
}
