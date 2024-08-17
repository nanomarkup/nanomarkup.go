package nanomarkup

import (
	"bytes"
	"fmt"
	"reflect"
	"strconv"

	"nanomarkup.go/nanocomment"
	"nanomarkup.go/nanodecoder"
	"nanomarkup.go/nanoerror"
	"nanomarkup.go/nanostr"
)

type Metadata struct {
	fields   map[string]*Metadata
	Comments nanocomment.Comments
}

func CreateMetadata(comment string, multiline bool) *Metadata {
	m := Metadata{}
	m.Comments.Add(comment, multiline)
	return &m
}

// Marshal returns the encoding data for the input value.
//
// It traverses the value recursively.
func Marshal(data any, meta *Metadata) ([]byte, error) {
	val := reflect.ValueOf(data)
	if val.Kind() == reflect.Pointer {
		val = val.Elem()
	}
	out := []byte("")
	var err error = nil
	if meta != nil && len(meta.Comments) > 0 {
		out = append(out, nanocomment.Marshal(meta.Comments)...)
	}
	switch val.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		out = append(out, []byte(strconv.FormatInt(val.Int(), 10))...)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		out = append(out, []byte(strconv.FormatUint(val.Uint(), 10))...)
	case reflect.Float32, reflect.Float64:
		out = append(out, []byte(strconv.FormatFloat(val.Float(), 'g', -1, 64))...)
	case reflect.Complex64, reflect.Complex128:
		out = append(out, []byte(strconv.FormatComplex(val.Complex(), 'g', -1, 128))...)
	case reflect.String:
		out = append(out, nanostr.Marshal(val.String())...)
	case reflect.Bool:
		out = append(out, []byte(strconv.FormatBool(val.Bool()))...)
	case reflect.Slice, reflect.Array:
		if val.Len() == 0 {
			out = append(out, []byte("[\n]\n")...)
		} else {
			if o, err := marshalSlice(val); err == nil {
				out = append(out, o...)
			} else {
				out = []byte("")
			}
		}
	case reflect.Map:
		if val.Len() == 0 {
			out = append(out, []byte("{\n}\n")...)
		} else {
			if o, err := marshalMap(val); err == nil {
				out = append(out, o...)
			} else {
				out = []byte("")
			}
		}
	case reflect.Struct:
		if o, err := marshalStruct(data, meta); err == nil {
			out = append(out, o...)
		} else {
			out = []byte("")
		}
	}
	return out, err
}

// MarshalIndent is like Marshal but applies Indent to format the output.
func MarshalIndent(data any, prefix, indent string) ([]byte, error) {
	enc, err := Marshal(data, nil)
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
func Unmarshal(data []byte, v any, meta *Metadata) error {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Pointer {
		return &nanoerror.InvalidArgumentError{"Unmarshal", fmt.Errorf("the second argument is not a Pointer")}
	}
	if rv.IsNil() {
		return &nanoerror.InvalidArgumentError{"Unmarshal", fmt.Errorf("the second argument is Nil")}
	}
	elem := rv.Elem()
	if !elem.CanSet() {
		return &nanoerror.InvalidArgumentError{"Unmarshal", fmt.Errorf("the second argument is not settable")}
	}
	d := nanodecoder.Decoder{}
	d.Init(bytes.Split(data, []byte("\n")))
	return unmarshal(&d, elem, undefined, meta)
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
