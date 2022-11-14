package jsondb

import (
	"reflect"
	"strings"
)

type (
	Model interface {
		PrimaryKey() string
		TableName() string
	}

	column struct {
		Name       string
		JsonName   string
		PrimaryKey bool
		Kind       reflect.Kind
	}

	Schema struct {
		refValue  reflect.Value
		tableName string
		Columns   []*column
	}
)

func parseModelSchema(model Model) (schema *Schema, err error) {
	var (
		pos              int
		hasPrimaryColumn bool
		refType          reflect.Type
		structField      reflect.StructField
	)
	schema = &Schema{
		Columns: make([]*column, 0, 10),
	}
	schema.refValue = reflect.Indirect(reflect.ValueOf(model))
	refType = schema.refValue.Type()
	for i := 0; i < refType.NumField(); i++ {
		structField = refType.Field(i)
		col := &column{
			Kind: structField.Type.Kind(),
		}
		col.Name = structField.Name
		col.JsonName = structField.Tag.Get("json")
		if col.JsonName == "-" {
			continue
		}
		if pos = strings.IndexByte(col.JsonName, ','); pos != -1 {
			col.JsonName = col.JsonName[:i]
		}
		if col.Name == model.PrimaryKey() || col.JsonName == model.PrimaryKey() {
			col.PrimaryKey = true
			hasPrimaryColumn = true
		}
		schema.Columns = append(schema.Columns, col)
	}
	if !hasPrimaryColumn {
		err = ErrPrimaryKeyNotExists
	}
	return
}

func (schema *Schema) PrimaryColumn() *column {
	for _, col := range schema.Columns {
		if col.PrimaryKey {
			return col
		}
	}
	return nil
}

func (schema *Schema) GetColumn(field string) (col *column, err error) {
	for _, c := range schema.Columns {
		if c.Name == field || c.JsonName == field {
			col = c
			break
		}
	}
	if col == nil {
		err = ErrColumnNotExists
	}
	return
}

func (schema *Schema) GetFieldValue(model Model, field string) (v interface{}, col *column, err error) {
	var (
		refValue   reflect.Value
		fieldValue reflect.Value
	)
	for _, c := range schema.Columns {
		if c.Name == field || c.JsonName == field {
			col = c
			break
		}
	}
	if col == nil {
		err = ErrColumnNotExists
		return
	}
	refValue = reflect.Indirect(reflect.ValueOf(model))
	fieldValue = refValue.FieldByName(col.Name)
	v = fieldValue.Interface()
	return
}

func (schema *Schema) GetPrimaryValue(model Model) (v interface{}, err error) {
	primaryColumn := schema.PrimaryColumn()
	if primaryColumn == nil {
		return nil, ErrPrimaryKeyNotExists
	}
	if v, _, err = schema.GetFieldValue(model, primaryColumn.Name); err != nil {
		err = ErrPrimaryKeyNotExists
	}
	return
}
