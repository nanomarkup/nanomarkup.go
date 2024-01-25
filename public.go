package nanomarkup

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// Marshal returns the Nano Markup encoding of data.
//
// It traverses the value data recursively.
func Marshal(data any) ([]byte, error) {
	val := reflect.ValueOf(data)
	switch val.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return []byte(fmt.Sprintf("%s %s", val.Kind().String(), strconv.FormatInt(val.Int(), 10))), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return []byte(fmt.Sprintf("%s %s", val.Kind().String(), strconv.FormatUint(val.Uint(), 10))), nil
	case reflect.Float32, reflect.Float64:
		return []byte(fmt.Sprintf("%s %s", val.Kind().String(), strconv.FormatFloat(val.Float(), 'g', -1, 64))), nil
	case reflect.Complex64, reflect.Complex128:
		return []byte(fmt.Sprintf("%s %s", val.Kind().String(), strconv.FormatComplex(val.Complex(), 'g', -1, 128))), nil
	case reflect.String:
		return []byte(fmt.Sprintf("%s %s", val.Kind().String(), strings.TrimSpace(val.String()))), nil
	case reflect.Bool:
		return []byte(fmt.Sprintf("%s %s", val.Kind().String(), strconv.FormatBool(val.Bool()))), nil
	case reflect.Slice, reflect.Array:
		if val.Len() == 0 {
			return []byte(""), nil
		} else {
			return marshalSlice(val)
		}
	case reflect.Map:
		if val.Len() == 0 {
			return []byte(""), nil
		} else {
			return marshalMap(val)
		}
	case reflect.Struct:
		return marshalStruct(data)
	default:
		return []byte(""), nil
	}
}
