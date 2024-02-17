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
		{v: func() {}, want: ""},
		{v: interface{}(nil), want: ""},
	}

	for _, item := range testCases {
		out, err := Marshal(item.v)
		if s := checkMarshal(item.v, out, item.want, err); s != "" {
			t.Error(s)
		}
	}
}

func TestEmptyMarshal(t *testing.T) {
	testCases := []struct {
		v    any
		want string
	}{
		{v: new(float64), want: "0"},
		{v: []any(nil), want: "[\n]\n"},
		{v: []string(nil), want: "[\n]\n"},
		{v: map[string]string(nil), want: "{\n}\n"},
		{v: []byte(nil), want: "[\n]\n"},
		{v: struct{}{}, want: "{\n}\n"},
		{v: struct{ M string }{"gopher"}, want: "{\nM gopher\n}\n"},
		{v: struct{ M testing.B }{}, want: "{\nM {\n}\n}\n"},
		{v: struct{ M encoding.TextMarshaler }{}, want: "{\n}\n"},
		{v: struct{ M any }{(nil)}, want: "{\n}\n"},
	}

	for _, item := range testCases {
		out, err := Marshal(item.v)
		if s := checkMarshal(item.v, out, item.want, err); s != "" {
			t.Error(s)
		}
	}
}

func TestNumberMarshal(t *testing.T) {
	testCases := []struct {
		v    any
		want string
	}{
		{v: 40, want: "40"},
		{v: -40, want: "-40"},
		{v: 40.4, want: "40.4"},
		{v: -40.4, want: "-40.4"},
		{v: 40.0, want: "40"},
		{v: -40.0, want: "-40"},
	}

	for _, item := range testCases {
		out, err := Marshal(item.v)
		if s := checkMarshal(item.v, out, item.want, err); s != "" {
			t.Error(s)
		}
	}
}

func TestStringMarshal(t *testing.T) {
	testCases := []struct {
		v    any
		want string
	}{
		{v: "go", want: "go"},
		{v: " go", want: "go"},
		{v: "go ", want: "go "},
		{v: " go ", want: "go "},
		{v: "hello world", want: "hello world"},
	}

	for _, item := range testCases {
		out, err := Marshal(item.v)
		if s := checkMarshal(item.v, out, item.want, err); s != "" {
			t.Error(s)
		}
	}
}

func TestMultiLineMarshal(t *testing.T) {
	// test a string
	sin := `testing
a multi
line
value`
	swant := "`\n" + sin + "\n`\n"

	out, err := Marshal(sin)
	if s := checkMarshal(sin, out, swant, err); s != "" {
		t.Error(s)
	}

	// test an array/slice
	ain := []string{sin}
	awant := "[\n`\n" + sin + "\n`\n]\n"

	out, err = Marshal(ain)
	if s := checkMarshal(ain, out, awant, err); s != "" {
		t.Error(s)
	}

	// test a struct
	type st struct {
		MultiValue string
	}
	tin := st{sin}
	twant := "{\nMultiValue `\n" + sin + "\n`\n}\n"

	out, err = Marshal(tin)
	if s := checkMarshal(tin, out, twant, err); s != "" {
		t.Error(s)
	}
}

func TestBooleanMarshal(t *testing.T) {
	testCases := []struct {
		v    any
		want string
	}{
		{v: true, want: "true"},
		{v: false, want: "false"},
	}

	for _, item := range testCases {
		out, err := Marshal(item.v)
		if s := checkMarshal(item.v, out, item.want, err); s != "" {
			t.Error(s)
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
		if s := checkMarshal(item.v, out, item.want, err); s != "" {
			t.Error(s)
		}
	}
}

func TestSliceMarshal(t *testing.T) {
	testCases := []struct {
		v    any
		want string
	}{
		// test arrays
		{v: [0]int{}, want: "[\n]\n"},
		{v: [3]int{}, want: "[\n0\n0\n0\n]\n"},
		{v: [3]string{"apple", "banana", "cherry"}, want: "[\napple\nbanana\ncherry\n]\n"},
		// test slices
		{v: []int{}, want: "[\n]\n"},
		{v: []int{1, 2, 3}, want: "[\n1\n2\n3\n]\n"},
		{v: []string{"apple", "banana", "cherry"}, want: "[\napple\nbanana\ncherry\n]\n"},
	}

	for _, item := range testCases {
		out, err := Marshal(item.v)
		if s := checkMarshal(item.v, out, item.want, err); s != "" {
			t.Error(s)
		}
	}
}

func TestMapMarshal(t *testing.T) {
	testCases := []struct {
		v    any
		want string
	}{
		{v: map[int]interface{}{}, want: "{\n}\n"},
		{v: map[int]interface{}{1: "Hi!"}, want: "{\n1 Hi!\n}\n"},
		{v: map[string]interface{}{"key": "value"}, want: "{\nkey value\n}\n"},
	}

	for _, item := range testCases {
		out, err := Marshal(item.v)
		if s := checkMarshal(item.v, out, item.want, err); s != "" {
			t.Error(s)
		}
	}
}

func TestNanoTagMarshal(t *testing.T) {
	testCases := []struct {
		v    any
		want string
	}{
		{v: struct {
			Field1 int    `nano:"omitempty"`
			Field2 string `nano:"omitempty"`
		}{0, ""}, want: "{\n}\n"},
		{v: struct {
			Field1 int    `nano:"-,omitempty"`
			Field2 string `nano:"omitempty,-"`
		}{1, "1"}, want: "{\n}\n"},
		{v: struct {
			Field1 int    `nano:"test1"`
			Field2 string `nano:"test2"`
		}{0, ""}, want: "{\ntest1 0\ntest2 \n}\n"},
		{v: struct {
			Field1 int    `nano:"omitempty,test1"`
			Field2 string `nano:"test2,omitempty"`
			Field3 string `nano:"omitempty,omitempty"`
		}{0, "", ""}, want: "{\n}\n"},
		{v: struct {
			Field1 int `nano:"omitempty"`
			Field2 int
			Field3 int `nano:"-"`
		}{1, 2, 3}, want: "{\nField1 1\nField2 2\n}\n"},
		{v: struct {
			Field1 int    `nano:"omitempty"`
			Field2 string `nano:"omitempty"`
			Field3 int
		}{0, "", 3}, want: "{\nField3 3\n}\n"},
	}

	for _, item := range testCases {
		out, err := Marshal(item.v)
		if s := checkMarshal(item.v, out, item.want, err); s != "" {
			t.Error(s)
		}
	}
}
