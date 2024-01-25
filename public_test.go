package nanomarkup

import (
	"encoding"
	"testing"
)

func TestNilMarshal(t *testing.T) {
	testCases := []struct {
		v    any
		want string
	}{
		{v: nil, want: ""},
		{v: new(float64), want: ""},
		{v: []any(nil), want: ""},
		{v: []string(nil), want: ``},
		{v: map[string]string(nil), want: ""},
		{v: []byte(nil), want: ""},
		{v: func() {}, want: ""},
		{v: struct{}{}, want: ""},
		{v: interface{}(nil), want: ""},
		{v: struct{ M string }{"gopher"}, want: "{\nM gopher\n}\n"},
		{v: struct{ M testing.B }{}, want: "{\nM \n}\n"},
		{v: struct{ M encoding.TextMarshaler }{}, want: "{\nM \n}\n"},
		{v: struct{ M any }{(nil)}, want: "{\nM \n}\n"},
	}

	for _, item := range testCases {
		out, err := Marshal(item.v)
		if err != nil || string(out) != item.want {
			t.Errorf("Marshal(%#v) = %#q, %#v, want %#q, nil", item.v, out, err, item.want)
			continue
		}
	}
}

func TestNumberMarshal(t *testing.T) {
	testCases := []struct {
		v    any
		want string
	}{
		{v: 40, want: "int 40"},
		{v: -40, want: "int -40"},
		{v: 40.4, want: "float64 40.4"},
		{v: -40.4, want: "float64 -40.4"},
		{v: 40.0, want: "float64 40"},
		{v: -40.0, want: "float64 -40"},
	}

	for _, item := range testCases {
		out, err := Marshal(item.v)
		if err != nil || string(out) != item.want {
			t.Errorf("Marshal(%#v) = %#q, %#v, want %#q, nil", item.v, out, err, item.want)
			continue
		}
	}
}

func TestStringMarshal(t *testing.T) {
	testCases := []struct {
		v    any
		want string
	}{
		{v: "go", want: "string go"},
		{v: " go", want: "string go"},
		{v: "go ", want: "string go"},
		{v: " go ", want: "string go"},
		{v: "hello world", want: "string hello world"},
	}

	for _, item := range testCases {
		out, err := Marshal(item.v)
		if err != nil || string(out) != item.want {
			t.Errorf("Marshal(%#v) = %#q, %#v, want %#q, nil", item.v, out, err, item.want)
			continue
		}
	}
}

func TestBooleanMarshal(t *testing.T) {
	testCases := []struct {
		v    any
		want string
	}{
		{v: true, want: "bool true"},
		{v: false, want: "bool false"},
	}

	for _, item := range testCases {
		out, err := Marshal(item.v)
		if err != nil || string(out) != item.want {
			t.Errorf("Marshal(%#v) = %#q, %#v, want %#q, nil", item.v, out, err, item.want)
			continue
		}
	}
}

func TestStructMarshal(t *testing.T) {
	type t1 struct {
		TestInt int
		TestStr string
	}
	type t2 struct {
		T1 t1
	}
	v1 := t1{7, "test 7"}
	v2 := t2{v1}
	testCases := []struct {
		v    any
		want string
	}{
		{v: struct {
			Field1 int
			field2 int
			Field3 int
		}{1, 2, 3}, want: "{\nField1 1\nField3 3\n}\n"},
		{v: struct {
			String string
		}{"Hello World!"}, want: "{\nString Hello World!\n}\n"},
		{v: v1, want: "{\nTestInt 7\nTestStr test 7\n}\n"},
		{v: v2, want: "{\nT1 {\nTestInt 7\nTestStr test 7\n}\n}\n"},
	}

	for _, item := range testCases {
		out, err := Marshal(item.v)
		if err != nil || string(out) != item.want {
			t.Errorf("Marshal(%#v) = %#q, %#v, want %#q, nil", item.v, out, err, item.want)
			continue
		}
	}
}

func TestSliceMarshal(t *testing.T) {
	testCases := []struct {
		v    any
		want string
	}{
		// test arrays
		{v: [0]int{}, want: ""},
		{v: [3]int{}, want: "[\n0\n0\n0\n]\n"},
		{v: [3]string{"apple", "banana", "cherry"}, want: "[\napple\nbanana\ncherry\n]\n"},
		// test slices
		{v: []int{}, want: ""},
		{v: []int{1, 2, 3}, want: "[\n1\n2\n3\n]\n"},
		{v: []string{"apple", "banana", "cherry"}, want: "[\napple\nbanana\ncherry\n]\n"},
	}

	for _, item := range testCases {
		out, err := Marshal(item.v)
		if err != nil || string(out) != item.want {
			t.Errorf("Marshal(%#v) = %#q, %#v, want %#q, nil", item.v, out, err, item.want)
			continue
		}
	}
}

func TestMapMarshal(t *testing.T) {
	testCases := []struct {
		v    any
		want string
	}{
		{v: map[int]interface{}{}, want: ""},
		{v: map[int]interface{}{1: "Hi!"}, want: "{\n1 Hi!\n}\n"},
		{v: map[string]interface{}{"key": "value"}, want: "{\nkey value\n}\n"},
	}

	for _, item := range testCases {
		out, err := Marshal(item.v)
		if err != nil || string(out) != item.want {
			t.Errorf("Marshal(%#v) = %#q, %#v, want %#q, nil", item.v, out, err, item.want)
			continue
		}
	}
}
