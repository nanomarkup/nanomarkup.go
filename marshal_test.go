package nanomarkup

import (
	"encoding"
	"fmt"
	"testing"
	"time"

	"github.com/nanomarkup/nanomarkup.go/nanometadata"
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
		out, err := Marshal(item.v, nil)
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
		out, err := Marshal(item.v, nil)
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
		out, err := Marshal(item.v, nil)
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
		out, err := Marshal(item.v, nil)
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

	out, err := Marshal(sin, nil)
	if s := checkMarshal(sin, out, swant, err); s != "" {
		t.Error(s)
	}

	// test an array/slice
	ain := []string{sin}
	awant := "[\n`\n" + sin + "\n`\n]\n"

	out, err = Marshal(ain, nil)
	if s := checkMarshal(ain, out, awant, err); s != "" {
		t.Error(s)
	}

	// test a struct
	type st struct {
		MultiValue string
	}
	tin := st{sin}
	twant := "{\nMultiValue `\n" + sin + "\n`\n}\n"

	out, err = Marshal(tin, nil)
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
		out, err := Marshal(item.v, nil)
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
		out, err := Marshal(item.v, nil)
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
		out, err := Marshal(item.v, nil)
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
		out, err := Marshal(item.v, nil)
		if s := checkMarshal(item.v, out, item.want, err); s != "" {
			t.Error(s)
		}
	}
}

func TestStructMapMarshal(t *testing.T) {
	type test struct {
		Log map[string]string
	}
	in := test{map[string]string{"Username": "John"}}
	want := "{\nLog {\nUsername John\n}\n}\n"
	out, err := Marshal(in, nil)
	if s := checkMarshal(in, out, want, err); s != "" {
		t.Error(s)
	}
}

func TestMapMapMarshal(t *testing.T) {
	in := map[string]map[string]string{"Log": map[string]string{"Username": "John"}}
	want := "{\nLog {\nUsername John\n}\n}\n"
	out, err := Marshal(in, nil)
	if s := checkMarshal(in, out, want, err); s != "" {
		t.Error(s)
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
		out, err := Marshal(item.v, nil)
		if s := checkMarshal(item.v, out, item.want, err); s != "" {
			t.Error(s)
		}
	}
}

func TestMetaIntMarshal(t *testing.T) {
	in := 1983
	comment := " A Birthday date"
	want := fmt.Sprintf("//%s\n%d", comment, in)
	meta := nanometadata.CreateMetadata(comment, false)
	out, err := Marshal(in, meta)
	if s := checkMarshal(in, out, want, err); s != "" {
		t.Error(s)
	}
}

func TestMetaStringMarshal(t *testing.T) {
	in := "Hello World!"
	comment := " Hello World comment"
	want := fmt.Sprintf("//%s\n%s", comment, in)
	meta := nanometadata.CreateMetadata(comment, false)
	out, err := Marshal(in, meta)
	if s := checkMarshal(in, out, want, err); s != "" {
		t.Error(s)
	}
}

func TestMetaMultiLineMarshal(t *testing.T) {
	in := `testing
a multi
line
value`
	comment := " A multi-line value"
	want := fmt.Sprintf("//%s\n`\n%s\n`\n", comment, in)
	meta := nanometadata.CreateMetadata(comment, false)
	out, err := Marshal(in, meta)
	if s := checkMarshal(in, out, want, err); s != "" {
		t.Error(s)
	}
}

func TestMetaBooleanMarshal(t *testing.T) {
	in := true
	comment := " Check type of boolean"
	want := fmt.Sprintf("//%s\n%t", comment, in)
	meta := nanometadata.CreateMetadata(comment, false)
	out, err := Marshal(in, meta)
	if s := checkMarshal(in, out, want, err); s != "" {
		t.Error(s)
	}
}

func TestMetaArrayMarshal(t *testing.T) {
	in := [3]int{1, 2, 3}
	comment := " Check type of array"
	want := fmt.Sprintf("//%s\n[\n1\n2\n3\n]\n", comment)
	meta := nanometadata.CreateMetadata(comment, false)
	out, err := Marshal(in, meta)
	if s := checkMarshal(in, out, want, err); s != "" {
		t.Error(s)
	}
}

func TestMetaSliceMarshal(t *testing.T) {
	in := []int{1, 2, 3}
	comment := " Check type of slice"
	want := fmt.Sprintf("//%s\n[\n1\n2\n3\n]\n", comment)
	meta := nanometadata.CreateMetadata(comment, false)
	out, err := Marshal(in, meta)
	if s := checkMarshal(in, out, want, err); s != "" {
		t.Error(s)
	}
}

func TestMetaMapMarshal(t *testing.T) {
	in := map[int]int{1: 1, 2: 2}
	comment := " Check type of map"
	meta := nanometadata.CreateMetadata(comment, false)
	out, err := Marshal(in, meta)
	// a map does not keep an element order
	want1 := fmt.Sprintf("//%s\n{\n1 1\n2 2\n}\n", comment)
	want2 := fmt.Sprintf("//%s\n{\n2 2\n1 1\n}\n", comment)
	s1 := checkMarshal(in, out, want1, err)
	s2 := checkMarshal(in, out, want2, err)
	if s1 == "" || s2 == "" {
		return
	} else {
		if s1 != "" {
			t.Error(s1)
		} else {
			t.Error(s2)
		}
	}
}

func TestMetaStructMarshal(t *testing.T) {
	in := struct {
		Field1 int
		Field2 string `nano:"omitempty"`
		Field3 string `nano:"-"`
		Field4 int    `nano:"Year"`
	}{0, "Hi!", "Hello!", 2024}
	meta := nanometadata.CreateMetadata(" Object's comment", false)
	meta.AddField("Field1", nanometadata.CreateMetadata(" Testing a comment...", false))
	meta.AddField("Field2", nanometadata.CreateMetadata(" It cannot be empty", false))
	meta.AddField("Field4", nanometadata.CreateMetadata(" Current year is", false))
	want := `// Object's comment
{
// Testing a comment...
Field1 0
// It cannot be empty
Field2 Hi!
// Current year is
Year 2024
}
`
	out, err := Marshal(in, meta)
	if s := checkMarshal(in, out, want, err); s != "" {
		t.Error(s)
	}
}

func TestCommentsMarshal(t *testing.T) {
	in := "Hello World!"
	empty := ""
	multi1 := " First multiline comment "
	multi2 := " Second\nmultiline\ncomment "
	single1 := " First single comment"
	single2 := " Second signle comment"
	want := fmt.Sprintf("//%s\n%s\n/*%s*/\n//%s\n/*%s*/\n/*%s*/\n%s", single1, empty, multi1, single2, multi2, empty, in)
	meta := nanometadata.CreateMetadata(single1, false)
	meta.Comments.Add(empty, false)
	meta.Comments.Add(multi1, true)
	meta.Comments.Add(single2, false)
	meta.Comments.Add(multi2, true)
	meta.Comments.Add(empty, true)
	out, err := Marshal(in, meta)
	if s := checkMarshal(in, out, want, err); s != "" {
		t.Error(s)
	}
}

func TestCustomMarshal(t *testing.T) {
	in := struct {
		Today time.Time
	}{time.Date(1983, 12, 20, 19, 30, 0, 0, time.Local)}
	want := "{\nToday {\n1983-12-20T19:30:00-05:00\n}\n}\n"
	out, err := Marshal(in, nil)
	if s := checkMarshal(in, out, want, err); s != "" {
		t.Error(s)
	}

	in2 := struct {
		Today customTime
	}{customTime{time.Date(1983, 12, 20, 19, 30, 0, 0, time.Local)}}
	want = fmt.Sprintf("{\nToday {\n%s\n}\n}\n", in2.Today.Format(in2.Today.getTimeFormat()))
	out, err = Marshal(in2, nil)
	if s := checkMarshal(in, out, want, err); s != "" {
		t.Error(s)
	}
}
