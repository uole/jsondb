package jsondb

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

func compare(x, y interface{}) (int, error) {
	var (
		intVal   int64
		uintVal  uint64
		floatVal float64
		s        string
		err      error
	)
	a := reflect.ValueOf(x)
	b := reflect.ValueOf(y)
	switch a.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		switch b.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			if a.Int() == b.Int() {
				return 0, nil
			}
			if a.Int() < b.Int() {
				return -1, nil
			}
			return +1, nil
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			if a.Int() == int64(b.Uint()) {
				return 0, nil
			}
			if a.Int() < int64(b.Uint()) {
				return -1, nil
			}
			return +1, nil
		case reflect.Float32, reflect.Float64:
			if a.Int() == int64(b.Float()) {
				return 0, nil
			}
			if a.Int() < int64(b.Float()) {
				return -1, nil
			}
			return +1, nil
		case reflect.String:
			if intVal, err = strconv.ParseInt(b.String(), 10, 64); err == nil {
				if a.Int() == intVal {
					return 0, nil
				}
				if a.Int() < intVal {
					return -1, nil
				}
				return +1, nil
			}
		default:
			err = fmt.Errorf("uncomparable kind %s", b.Kind().String())
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		switch b.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			if a.Uint() == uint64(b.Int()) {
				return 0, nil
			}
			if a.Uint() < uint64(b.Int()) {
				return -1, nil
			}
			return +1, nil
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			if a.Uint() == b.Uint() {
				return 0, nil
			}
			if a.Uint() < b.Uint() {
				return -1, nil
			}
			return +1, nil
		case reflect.Float32, reflect.Float64:
			if a.Uint() == uint64(b.Float()) {
				return 0, nil
			}
			if a.Uint() < uint64(b.Float()) {
				return -1, nil
			}
			return +1, nil
		case reflect.String:
			if uintVal, err = strconv.ParseUint(b.String(), 10, 64); err == nil {
				if a.Uint() == uintVal {
					return 0, nil
				}
				if a.Uint() < uintVal {
					return -1, nil
				}
				return +1, nil
			}
		default:
			err = fmt.Errorf("uncomparable kind %s", b.Kind().String())
		}
	case reflect.Float32, reflect.Float64:
		switch b.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			if a.Float() == float64(b.Int()) {
				return 0, nil
			}
			if a.Float() < float64(b.Int()) {
				return -1, nil
			}
			return +1, nil
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			if a.Float() == float64(b.Uint()) {
				return 0, nil
			}
			if a.Float() < float64(b.Uint()) {
				return -1, nil
			}
			return +1, nil
		case reflect.Float32, reflect.Float64:
			if a.Float() == b.Float() {
				return 0, nil
			}
			if a.Float() < b.Float() {
				return -1, nil
			}
			return +1, nil
		case reflect.String:
			if floatVal, err = strconv.ParseFloat(b.String(), 64); err == nil {
				if a.Float() == floatVal {
					return 0, nil
				}
				if a.Float() < floatVal {
					return -1, nil
				}
				return +1, nil
			}
		default:
			err = fmt.Errorf("uncomparable kind %s", b.Kind().String())
		}
	case reflect.String:
		switch b.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			s = strconv.FormatInt(b.Int(), 10)
			return strings.Compare(a.String(), s), nil
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			s = strconv.FormatUint(b.Uint(), 10)
			return strings.Compare(a.String(), s), nil
		case reflect.Float32, reflect.Float64:
			s = strconv.FormatFloat(b.Float(), 'f', -1, 64)
			return strings.Compare(a.String(), s), nil
		case reflect.String:
			return strings.Compare(a.String(), b.String()), nil
		default:
			err = fmt.Errorf("uncomparable kind %s", b.Kind().String())
		}
	default:
		err = fmt.Errorf("uncomparable kind %s", a.Kind().String())
	}
	return 0, fmt.Errorf("uncomparable kind")
}
