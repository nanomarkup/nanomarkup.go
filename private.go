package nanomarkup

import (
	"reflect"
	"strconv"
	"strings"
)

func marshal(data any) ([]byte, error) {
	val := reflect.ValueOf(data)
	switch val.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return []byte(strconv.FormatInt(val.Int(), 10)), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return []byte(strconv.FormatUint(val.Uint(), 10)), nil
	case reflect.Float32, reflect.Float64:
		return []byte(strconv.FormatFloat(val.Float(), 'g', -1, 64)), nil
	case reflect.Complex64, reflect.Complex128:
		return []byte(strconv.FormatComplex(val.Complex(), 'g', -1, 128)), nil
	case reflect.String:
		return []byte(strings.TrimSpace(val.String())), nil
	case reflect.Bool:
		return []byte(strconv.FormatBool(val.Bool())), nil
	case reflect.Slice, reflect.Array:
		if val.Len() == 0 {
			return []byte(""), nil
		} else {
			return marshalMap(val)
		}
	case reflect.Map:
		if val.Len() == 0 {
			return []byte(""), nil
		} else {
			return marshalSlice(val)
		}
	case reflect.Struct:
		if val.IsZero() {
			return []byte(""), nil
		}
		return marshalStruct(data)
	default:
		return []byte(""), nil
	}
}

func marshalStruct(data any) ([]byte, error) {
	typ := reflect.TypeOf(data)
	val := reflect.ValueOf(data)
	res := []byte("{\n")
	for _, f := range reflect.VisibleFields(typ) {
		if !f.IsExported() {
			continue
		}
		res = append(res, []byte(f.Name+" ")...)
		v, e := marshal(val.Field(f.Index[0]).Interface())
		if e != nil {
			return nil, e
		}
		res = append(res, v...)
		if string(res[len(res)-1]) != "\n" {
			res = append(res, "\n"...)
		}
	}
	if len(res) == 2 {
		return []byte(""), nil
	}
	res = append(res, []byte("}\n")...)
	return res, nil
}

func marshalSlice(value reflect.Value) ([]byte, error) {
	res := []byte("[\n")
	for i := 0; i < value.Len(); i++ {
		v, e := marshal(value.Index(i).Interface())
		if e != nil {
			return nil, e
		}
		res = append(res, v...)
		if string(res[len(res)-1]) != "\n" {
			res = append(res, "\n"...)
		}
	}
	if len(res) == 2 {
		return []byte(""), nil
	}
	res = append(res, []byte("]\n")...)
	return res, nil
}

func marshalMap(value reflect.Value) ([]byte, error) {
	res := []byte("{\n")
	iter := value.MapRange()
	for iter.Next() {
		v, e := marshal(iter.Key().Interface())
		if e != nil {
			return nil, e
		}
		res = append(res, v...)
		res = append(res, 32) // add a space
		v, e = marshal(iter.Value().Interface())
		if e != nil {
			return nil, e
		}
		res = append(res, v...)
		if string(res[len(res)-1]) != "\n" {
			res = append(res, "\n"...)
		}
	}
	if len(res) == 2 {
		return []byte(""), nil
	}
	res = append(res, []byte("}\n")...)
	return res, nil
}
