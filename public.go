package nanomarkup

import (
	"bytes"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// Marshal returns the encoding data for the input value.
//
// It traverses the value recursively.
func Marshal(data any) ([]byte, error) {
	val := reflect.ValueOf(data)
	if val.Kind() == reflect.Pointer {
		val = val.Elem()
	}
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
		lines := strings.Split(val.String(), "\n")
		if len(lines) == 1 {
			return []byte(strings.TrimLeft(val.String(), " \t")), nil
		} else {
			res := "`\n"
			for _, it := range lines {
				res += it + "\n"
			}
			return []byte(res + "`\n"), nil
		}
	case reflect.Bool:
		return []byte(strconv.FormatBool(val.Bool())), nil
	case reflect.Slice, reflect.Array:
		if val.Len() == 0 {
			return []byte("[\n]\n"), nil
		} else {
			return marshalSlice(val)
		}
	case reflect.Map:
		if val.Len() == 0 {
			return []byte("{\n}\n"), nil
		} else {
			return marshalMap(val)
		}
	case reflect.Struct:
		return marshalStruct(data)
	default:
		return []byte(""), nil
	}
}

// MarshalIndent is like Marshal but applies Indent to format the output.
func MarshalIndent(data any, prefix, indent string) ([]byte, error) {
	enc, err := Marshal(data)
	if err != nil {
		return nil, err
	}
	dst := bytes.Buffer{}
	if err = Indent(&dst, enc, prefix, indent); err != nil {
		return nil, err
	} else {
		return dst.Bytes(), nil
	}
}

// Unmarshal parses the encoded data and stores the result in v.
// If v is nil or not a pointer, Unmarshal returns an InvalidArgumentError.
//
// It uses the inverse of the encodings that Marshal uses, allocating
// maps, slices, and pointers as necessary.
func Unmarshal(data []byte, v any) error {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Pointer {
		return &InvalidArgumentError{"Unmarshal", fmt.Errorf("the second argument is not a Pointer")}
	}
	if rv.IsNil() {
		return &InvalidArgumentError{"Unmarshal", fmt.Errorf("the second argument is Nil")}
	}
	elem := rv.Elem()
	if !elem.CanSet() {
		return &InvalidArgumentError{"Unmarshal", fmt.Errorf("the second argument is not settable")}
	}
	d := decoder{}
	d.data = bytes.Split(data, []byte("\n"))
	d.reset()
	return unmarshal(&d, elem, undefined)
}

// Indent function appends to `dst` the nano-encoded source (`src`) in an indented format.
// The data appended to dst does not begin with the prefix nor any indentation,
// to make it easier to embed inside other formatted nano-encoded data.
func Indent(dst *bytes.Buffer, src []byte, prefix, indent string) error {
	// specify a growth factor to reduce probability of allocation memory
	factor := float64(len(prefix)+len(indent)+1)/10 + 1
	dst.Grow(int(float64(len(src)) * factor))
	b := dst.AvailableBuffer()
	b, err := appendIndent(b, src, prefix, indent)
	dst.Write(b)
	return err
}

// Compact appends the nano-encoded src to dst, eliminating insignificant space characters.
func Compact(dst *bytes.Buffer, src []byte) error {
	dst.Grow(len(src))
	b := dst.AvailableBuffer()
	b, err := appendIndent(b, src, "", "")
	dst.Write(b)
	return err
}

// InvalidArgumentError describes an error that occurs when an invalid argument is provided.
type InvalidArgumentError struct {
	Context string
	Err     error
}

// Error returns a string representation of the InvalidArgumentError.
func (e *InvalidArgumentError) Error() string {
	if len(strings.TrimSpace(e.Context)) > 0 {
		return fmt.Sprintf(errorContextFmt, e.Context, e.Err.Error())
	} else {
		return e.Err.Error()
	}
}

// InvalidEntityError describes an error that occurs when an attempt is made with an invalid entity.
type InvalidEntityError struct {
	Context string
	Entity  string
	Err     error
}

// Error returns a string representation of the InvalidEntityError.
func (e *InvalidEntityError) Error() string {
	var s string
	if len(strings.TrimSpace(e.Entity)) > 0 {
		s = fmt.Sprintf("%s: %s", e.Err.Error(), e.Entity)
	} else {
		s = e.Err.Error()
	}

	if len(strings.TrimSpace(e.Context)) > 0 {
		return fmt.Sprintf(errorContextFmt, e.Context, s)
	} else {
		return s
	}
}
