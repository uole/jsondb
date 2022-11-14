package jsondb

import (
	"reflect"
)

const (
	ExprEqual        = "="
	ExprNotEqual     = "!="
	ExprGreater      = ">"
	ExprGreaterEqual = ">="
	ExprLess         = "<"
	ExprLessEqual    = "<="
)

type (
	expression struct {
		Field string
		Expr  string
		Value interface{}
	}

	Query struct {
		model       Model
		schema      *Schema
		documents   []Document
		offset      int
		limit       int
		expressions []*expression
	}
)

func (query *Query) isComparable(kind reflect.Kind) bool {
	vs := []reflect.Kind{
		reflect.Int,
		reflect.Int8,
		reflect.Int16,
		reflect.Int32,
		reflect.Int64,
		reflect.Uint,
		reflect.Uint8,
		reflect.Uint16,
		reflect.Uint32,
		reflect.Uint64,
		reflect.Float32,
		reflect.Float64,
		reflect.String,
	}
	for _, i := range vs {
		if i == kind {
			return true
		}
	}
	return false
}

func (query *Query) filter(e *expression, m Document) (ok bool) {
	var (
		n   int
		col *column
		v   interface{}
		err error
	)
	if col, err = query.schema.GetColumn(e.Field); err != nil {
		//if epr is not equal return true, otherwise return false
		if e.Expr == ExprNotEqual {
			return true
		}
		return false
	}
	v = m.Getter(col.JsonName, col.Name)
	switch e.Expr {
	case ExprEqual:
		if n, err = compare(v, e.Value); err == nil && n == 0 {
			ok = true
		}
	case ExprNotEqual:
		if n, err = compare(v, e.Value); err != nil || n != 0 {
			ok = true
		}
	case ExprGreater:
		if query.isComparable(col.Kind) {
			if n, err = compare(v, e.Value); err == nil && n > 0 {
				ok = true
			}
		}
	case ExprGreaterEqual:
		if n, err = compare(v, e.Value); err == nil && n >= 0 {
			ok = true
		}
	case ExprLess:
		if n, err = compare(v, e.Value); err == nil && n < 0 {
			ok = true
		}
	case ExprLessEqual:
		if n, err = compare(v, e.Value); err == nil && n <= 0 {
			ok = true
		}
	}
	return
}

//Where add query condition
func (query *Query) Where(field string, expr string, value interface{}) *Query {
	if query.expressions == nil {
		query.expressions = make([]*expression, 0, 10)
	}
	query.expressions = append(query.expressions, &expression{
		Field: field,
		Expr:  expr,
		Value: value,
	})
	return query
}

//Offset set data offset
func (query *Query) Offset(n int) *Query {
	query.offset = n
	return query
}

//Limit set record limit
func (query *Query) Limit(n int) *Query {
	query.limit = n
	return query
}

func (query *Query) Documents() (documents []Document, err error) {
	var (
		index   int
		count   int
		isMatch bool
	)
	if len(query.documents) <= 0 {
		err = ErrRecordNotFound
		return
	}
	documents = make([]Document, 0, len(query.documents))
	for _, document := range query.documents {
		isMatch = true
		for _, expr := range query.expressions {
			if !query.filter(expr, document) {
				isMatch = false
				break
			}
		}
		if isMatch {
			if index >= query.offset {
				documents = append(documents, document)
				count++
				if query.limit > 0 && count >= query.limit {
					goto __end
				}
			}
			index++
		}
	}
__end:
	return
}

func (query *Query) Find(v interface{}) (err error) {
	var (
		index   int
		count   int
		isMatch bool
		mv      reflect.Value
	)
	if len(query.documents) <= 0 {
		return ErrRecordNotFound
	}
	refType := reflect.TypeOf(v)
	refPtr := reflect.ValueOf(v)
	if refType.Kind() != reflect.Ptr {
		return
	}
	refVal := refPtr.Elem()
	if refVal.Kind() != reflect.Slice && refVal.Kind() != reflect.Array {
		return
	}
	result := make([]reflect.Value, 0, len(query.documents))
	for _, document := range query.documents {
		isMatch = true
		for _, expr := range query.expressions {
			if !query.filter(expr, document) {
				isMatch = false
				break
			}
		}
		if isMatch {
			if index >= query.offset {
				if mv, err = document.Unmarshal(refType.Elem().Elem()); err == nil {
					result = append(result, mv)
					count++
					if query.limit > 0 && count >= query.limit {
						goto __end
					}
				}
			}
			index++
		}
	}
__end:
	refVal.Set(reflect.Append(refVal, result...))
	return
}

func newQuery(model Model, schema *Schema, documents []Document) *Query {
	return &Query{
		model:     model,
		schema:    schema,
		documents: documents,
	}
}
