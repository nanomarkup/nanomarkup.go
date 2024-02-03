package nanomarkup

import (
	"bytes"
	"encoding"
	"fmt"
	"net/http"
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

func TestIntUnmarshal(t *testing.T) {
	testCases := []struct {
		v    string
		want int
	}{
		{v: "40", want: 40},
		{v: "-40", want: -40},
	}

	r := new(int)
	for _, item := range testCases {
		err := Unmarshal([]byte(item.v), r)
		if s := checkUnmarshalInt(item.v, int64(*r), int64(item.want), err); s != "" {
			t.Error(s)
		}
	}
}

func TestInt8Unmarshal(t *testing.T) {
	testCases := []struct {
		v    string
		want int8
	}{
		{v: "40", want: 40},
		{v: "-40", want: -40},
	}

	r := new(int8)
	for _, item := range testCases {
		err := Unmarshal([]byte(item.v), r)
		if s := checkUnmarshalInt(item.v, int64(*r), int64(item.want), err); s != "" {
			t.Error(s)
		}
	}
}

func TestInt16Unmarshal(t *testing.T) {
	testCases := []struct {
		v    string
		want int16
	}{
		{v: "40", want: 40},
		{v: "-40", want: -40},
	}

	r := new(int16)
	for _, item := range testCases {
		err := Unmarshal([]byte(item.v), r)
		if s := checkUnmarshalInt(item.v, int64(*r), int64(item.want), err); s != "" {
			t.Error(s)
		}
	}
}

func TestInt32Unmarshal(t *testing.T) {
	testCases := []struct {
		v    string
		want int32
	}{
		{v: "40", want: 40},
		{v: "-40", want: -40},
	}

	r := new(int32)
	for _, item := range testCases {
		err := Unmarshal([]byte(item.v), r)
		if s := checkUnmarshalInt(item.v, int64(*r), int64(item.want), err); s != "" {
			t.Error(s)
		}
	}
}

func TestInt64Unmarshal(t *testing.T) {
	testCases := []struct {
		v    string
		want int64
	}{
		{v: "40", want: 40},
		{v: "-40", want: -40},
	}

	r := new(int64)
	for _, item := range testCases {
		err := Unmarshal([]byte(item.v), r)
		if s := checkUnmarshalInt(item.v, int64(*r), int64(item.want), err); s != "" {
			t.Error(s)
		}
	}
}

func TestUintUnmarshal(t *testing.T) {
	testCases := []struct {
		v    string
		want uint
	}{
		{v: "40", want: 40},
	}

	r := new(uint)
	for _, item := range testCases {
		err := Unmarshal([]byte(item.v), r)
		if s := checkUnmarshalUint(item.v, uint64(*r), uint64(item.want), err); s != "" {
			t.Error(s)
		}
	}
}

func TestUint8Unmarshal(t *testing.T) {
	testCases := []struct {
		v    string
		want uint8
	}{
		{v: "40", want: 40},
	}

	r := new(uint8)
	for _, item := range testCases {
		err := Unmarshal([]byte(item.v), r)
		if s := checkUnmarshalUint(item.v, uint64(*r), uint64(item.want), err); s != "" {
			t.Error(s)
		}
	}
}

func TestUint16Unmarshal(t *testing.T) {
	testCases := []struct {
		v    string
		want uint16
	}{
		{v: "40", want: 40},
	}

	r := new(uint16)
	for _, item := range testCases {
		err := Unmarshal([]byte(item.v), r)
		if s := checkUnmarshalUint(item.v, uint64(*r), uint64(item.want), err); s != "" {
			t.Error(s)
		}
	}
}

func TestUint32Unmarshal(t *testing.T) {
	testCases := []struct {
		v    string
		want uint32
	}{
		{v: "40", want: 40},
	}

	r := new(uint32)
	for _, item := range testCases {
		err := Unmarshal([]byte(item.v), r)
		if s := checkUnmarshalUint(item.v, uint64(*r), uint64(item.want), err); s != "" {
			t.Error(s)
		}
	}
}

func TestUint64Unmarshal(t *testing.T) {
	testCases := []struct {
		v    string
		want uint64
	}{
		{v: "40", want: 40},
	}

	r := new(uint64)
	for _, item := range testCases {
		err := Unmarshal([]byte(item.v), r)
		if s := checkUnmarshalUint(item.v, uint64(*r), uint64(item.want), err); s != "" {
			t.Error(s)
		}
	}
}

func TestFloat32Unmarshal(t *testing.T) {
	testCases := []struct {
		v    string
		want float32
	}{
		{v: "40.4", want: 40.4},
		{v: "-40.4", want: -40.4},
		{v: "40.0", want: 40},
		{v: "-40.0", want: -40},
	}

	r := new(float32)
	for _, item := range testCases {
		err := Unmarshal([]byte(item.v), r)
		if s := checkUnmarshalFloat(item.v, float64(*r), float64(item.want), err); s != "" {
			t.Error(s)
		}
	}
}

func TestFloat64Unmarshal(t *testing.T) {
	testCases := []struct {
		v    string
		want float64
	}{
		{v: "40.4", want: 40.4},
		{v: "-40.4", want: -40.4},
		{v: "40.0", want: 40},
		{v: "-40.0", want: -40},
	}

	r := new(float64)
	for _, item := range testCases {
		err := Unmarshal([]byte(item.v), r)
		if s := checkUnmarshalFloat(item.v, float64(*r), float64(item.want), err); s != "" {
			t.Error(s)
		}
	}
}

func TestComplex64Unmarshal(t *testing.T) {
	testCases := []struct {
		v    string
		want complex64
	}{
		{v: "40+41i", want: 40 + 41i},
	}

	r := new(complex64)
	for _, item := range testCases {
		err := Unmarshal([]byte(item.v), r)
		if s := checkUnmarshalComplex(item.v, complex128(*r), complex128(item.want), err); s != "" {
			t.Error(s)
		}
	}
}

func TestComplex128Unmarshal(t *testing.T) {
	testCases := []struct {
		v    string
		want complex128
	}{
		{v: "40+41i", want: 40 + 41i},
	}

	r := new(complex128)
	for _, item := range testCases {
		err := Unmarshal([]byte(item.v), r)
		if s := checkUnmarshalComplex(item.v, complex128(*r), complex128(item.want), err); s != "" {
			t.Error(s)
		}
	}
}

func TestStringUnmarshal(t *testing.T) {
	testCases := []struct {
		v    string
		want string
	}{
		{v: "go", want: "go"},
		{v: " go", want: "go"},
		{v: "go ", want: "go "},
		{v: " go ", want: "go "},
		{v: "hello world", want: "hello world"},
	}

	r := new(string)
	for _, item := range testCases {
		err := Unmarshal([]byte(item.v), r)
		if s := checkUnmarshalString(item.v, *r, item.want, err); s != "" {
			t.Error(s)
		}
	}
}

func TestBooleanUnmarshal(t *testing.T) {
	testCases := []struct {
		v    string
		want bool
	}{
		{v: "true", want: true},
		{v: "false", want: false},
	}

	r := new(bool)
	for _, item := range testCases {
		err := Unmarshal([]byte(item.v), r)
		if s := checkUnmarshalBool(item.v, *r, item.want, err); s != "" {
			t.Error(s)
		}
	}
}

func TestArrayUnmarshal(t *testing.T) {
	// test an empty array
	array := [0]int{}
	err := Unmarshal([]byte("[\n]"), &array)
	mes := ""
	if len(array) > 0 {
		mes = fmt.Sprintf("[Unmarshal] in: %s; out: %v; want: %v", "[]", array, [0]int{})
	}
	if err != nil {
		if mes == "" {
			mes = "[Unmarshal]: " + err.Error()
		} else {
			mes += "; error: " + err.Error()
		}
	}
	if mes != "" {
		t.Error(mes)
	}
	// test array of int
	arrayInt := []struct {
		v    string
		want [3]int
	}{
		{v: "[\n1\n2\n3\n]\n", want: [3]int{1, 2, 3}},
		{v: "[\n1\n-2\n3\n]", want: [3]int{1, -2, 3}},
	}
	for _, item := range arrayInt {
		rai := [3]int{}
		err = Unmarshal([]byte(item.v), &rai)
		if err == nil && rai == item.want {
			continue
		}
		mes = fmt.Sprintf("[Unmarshal] in: %s; out: %v; want: %v", item.v, rai, item.want)
		if err != nil {
			if mes == "" {
				mes = "[Unmarshal]: " + err.Error()
			} else {
				mes += "; error: " + err.Error()
			}
		}
		t.Error(mes)
	}
	// test array of string
	arrayString := []struct {
		v    string
		want [3]string
	}{
		{v: "[\napple\nbanana\ncherry\n]", want: [3]string{"apple", "banana", "cherry"}},
	}
	for _, item := range arrayString {
		ras := [3]string{}
		err = Unmarshal([]byte(item.v), &ras)
		if err == nil && ras == item.want {
			continue
		}
		mes = fmt.Sprintf("[Unmarshal] in: %s; out: %v; want: %v", item.v, ras, item.want)
		if err != nil {
			if mes == "" {
				mes = "[Unmarshal]: " + err.Error()
			} else {
				mes += "; error: " + err.Error()
			}
		}
		t.Error(mes)
	}
}

func TestSliceUnmarshal(t *testing.T) {
	// test an empty slice
	slice := []int{}
	err := Unmarshal([]byte("[\n]"), &slice)
	mes := ""
	if len(slice) > 0 {
		mes = fmt.Sprintf("[Unmarshal] in: %s; out: %v; want: %v", "[]", slice, []int{})
	}
	if err != nil {
		if mes == "" {
			mes = "[Unmarshal]: " + err.Error()
		} else {
			mes += "; error: " + err.Error()
		}
	}
	if mes != "" {
		t.Error(mes)
	}
	// test slice of int
	sliceInt := []struct {
		v    string
		want []int
	}{
		{v: "[\n1\n2\n3\n]\n", want: []int{1, 2, 3}},
		{v: "[\n1\n-2\n3\n]", want: []int{1, -2, 3}},
	}
	for _, item := range sliceInt {
		rsi := []int{}
		err = Unmarshal([]byte(item.v), &rsi)
		pass := len(rsi) == len(item.want)
		if pass {
			for i := range rsi {
				if rsi[i] != item.want[i] {
					pass = false
					break
				}
			}
		}
		if err == nil && pass {
			continue
		}
		mes = fmt.Sprintf("[Unmarshal] in: %s; out: %v; want: %v", item.v, rsi, item.want)
		if err != nil {
			if mes == "" {
				mes = "[Unmarshal]: " + err.Error()
			} else {
				mes += "; error: " + err.Error()
			}
		}
		t.Error(mes)
	}
	// test slice of string
	sliceString := []struct {
		v    string
		want []string
	}{
		{v: "[\nHello\nWorld\n!\n]", want: []string{"Hello", "World", "!"}},
		{v: "[\napple\nbanana\ncherry\n]", want: []string{"apple", "banana", "cherry"}},
	}
	for _, item := range sliceString {
		rss := []string{}
		err = Unmarshal([]byte(item.v), &rss)
		pass := len(rss) == len(item.want)
		if pass {
			for i := range rss {
				if rss[i] != item.want[i] {
					pass = false
					break
				}
			}
		}
		if err == nil && pass {
			continue
		}
		mes = fmt.Sprintf("[Unmarshal] in: %s; out: %v; want: %v", item.v, rss, item.want)
		if err != nil {
			if mes == "" {
				mes = "[Unmarshal]: " + err.Error()
			} else {
				mes += "; error: " + err.Error()
			}
		}
		t.Error(mes)
	}
}

func TestMapUnmarshal(t *testing.T) {
	// test an empty map
	m := map[int]interface{}{}
	err := Unmarshal([]byte("{\n}"), &m)
	mes := ""
	if len(m) > 0 {
		mes = fmt.Sprintf("[Unmarshal] in: %s; out: %v; want: %v", "[]", m, []int{})
	}
	if err != nil {
		if mes == "" {
			mes = "[Unmarshal]: " + err.Error()
		} else {
			mes += "; error: " + err.Error()
		}
	}
	if mes != "" {
		t.Error(mes)
	}
	// test map where key and value are int
	mapInt := []struct {
		v    string
		want map[int]int
	}{
		{v: "{\n1 1\n2 2\n3 3\n}\n", want: map[int]int{1: 1, 2: 2, 3: 3}},
		{v: "{\n1 1\n2 -2\n3 3\n}", want: map[int]int{1: 1, 2: -2, 3: 3}},
	}
	for _, item := range mapInt {
		rmi := map[int]int{}
		err = Unmarshal([]byte(item.v), &rmi)
		pass := len(rmi) == len(item.want)
		if pass {
			for i := range item.want {
				if rmi[i] != item.want[i] {
					pass = false
					break
				}
			}
		}
		if err == nil && pass {
			continue
		}
		mes = fmt.Sprintf("[Unmarshal] in: %s; out: %v; want: %v", item.v, rmi, item.want)
		if err != nil {
			if mes == "" {
				mes = "[Unmarshal]: " + err.Error()
			} else {
				mes += "; error: " + err.Error()
			}
		}
		t.Error(mes)
	}
	// test map where key and value are string
	mapString := []struct {
		v    string
		want map[string]string
	}{
		{v: "{\nkey value\nHi!\n}", want: map[string]string{"key": "value", "Hi!": ""}},
	}
	for _, item := range mapString {
		rms := map[string]string{}
		err = Unmarshal([]byte(item.v), &rms)
		pass := len(rms) == len(item.want)
		if pass {
			for i := range item.want {
				if rms[i] != item.want[i] {
					pass = false
					break
				}
			}
		}
		if err == nil && pass {
			continue
		}
		mes = fmt.Sprintf("[Unmarshal] in: %s; out: %v; want: %v", item.v, rms, item.want)
		if err != nil {
			if mes == "" {
				mes = "[Unmarshal]: " + err.Error()
			} else {
				mes += "; error: " + err.Error()
			}
		}
		t.Error(mes)
	}
}

func TestStructUnmarshal(t *testing.T) {
	type t1 struct {
		TestInt int
		TestStr string
	}
	type t2 struct {
		T1 t1
	}
	st1 := "{\nTestInt 7\nTestStr test 7\n}\n"
	rv1 := t1{}
	want1 := t1{7, "test 7"}
	st2 := "{\nT1 {\nTestInt 7\nTestStr test 7\n}\n}\n"
	rv2 := t2{rv1}
	want2 := t2{want1}

	err := Unmarshal([]byte(st1), &rv1)
	mes := ""
	if rv1.TestInt != want1.TestInt &&
		rv1.TestStr != want1.TestStr {
		mes = fmt.Sprintf("[Unmarshal] in: %s; out: %v; want: %v", st1, rv1, want1)
	}
	if err != nil {
		if mes == "" {
			mes = "[Unmarshal]: " + err.Error()
		} else {
			mes += "; error: " + err.Error()
		}
	}
	if mes != "" {
		t.Error(mes)
	}

	err = Unmarshal([]byte(st2), &rv2)
	mes = ""
	if rv2.T1.TestInt != want2.T1.TestInt &&
		rv2.T1.TestStr != want2.T1.TestStr {
		mes = fmt.Sprintf("[Unmarshal] in: %s; out: %v; want: %v", st2, rv2, want2)
	}
	if err != nil {
		if mes == "" {
			mes = "[Unmarshal]: " + err.Error()
		} else {
			mes += "; error: " + err.Error()
		}
	}
	if mes != "" {
		t.Error(mes)
	}
}

func TestTabUnmarshal(t *testing.T) {
	type t1 struct {
		TestInt int
		TestStr string
	}
	type t2 struct {
		T1 t1
	}
	rv1 := t1{}
	want1 := t1{7, "test 7"}
	st2 := "{\n\tT1 {\n\t\tTestInt 7\n\t\tTestStr test 7\n\t}\n}\n"
	rv2 := t2{rv1}
	want2 := t2{want1}

	err := Unmarshal([]byte(st2), &rv2)
	mes := ""
	if rv2.T1.TestInt != want2.T1.TestInt &&
		rv2.T1.TestStr != want2.T1.TestStr {
		mes = fmt.Sprintf("[Unmarshal] in: %s; out: %v; want: %v", st2, rv2, want2)
	}
	if err != nil {
		if mes == "" {
			mes = "[Unmarshal]: " + err.Error()
		} else {
			mes += "; error: " + err.Error()
		}
	}
	if mes != "" {
		t.Error(mes)
	}
}

func TestHTTPRequestStruct(t *testing.T) {
	req, err := http.NewRequest("GET", "https://google.com", nil)
	if err != nil {
		t.Error(err)
	}
	enc, err := Marshal(req)
	if err != nil {
		t.Error(err)
	}
	dec := &http.Request{}
	err = Unmarshal(enc, dec)
	if err != nil {
		t.Error(err)
	}
	testStructs(t, req, dec)
}

func TestIndentIndent(t *testing.T) {
	req, err := http.NewRequest("GET", "https://google.com", nil)
	if err != nil {
		t.Error(err)
		return
	}
	enc, err := Marshal(req)
	if err != nil {
		t.Error(err)
		return
	}
	dst := bytes.Buffer{}
	if err = Indent(&dst, enc, "", "    "); err != nil {
		t.Error(err)
		return
	}
	want := `{
    Method GET
    URL {
        Scheme https
        Opaque 
        Host google.com
        Path 
        RawPath 
        OmitHost false
        ForceQuery false
        RawQuery 
        Fragment 
        RawFragment 
    }
    Proto HTTP/1.1
    ProtoMajor 1
    ProtoMinor 1
    Header {
    }
    ContentLength 0
    Close false
    Host google.com
    RemoteAddr 
    RequestURI 
}
`
	out := dst.String()
	if out != want {
		t.Errorf("[Indent] in: %s; out: %s; want: %s", enc, out, want)
	}
}

func TestIndentPrefix(t *testing.T) {
	req, err := http.NewRequest("GET", "https://google.com", nil)
	if err != nil {
		t.Error(err)
		return
	}
	enc, err := Marshal(req)
	if err != nil {
		t.Error(err)
		return
	}
	dst := bytes.Buffer{}
	if err = Indent(&dst, enc, "##", "  "); err != nil {
		t.Error(err)
		return
	}
	want := `{
##  Method GET
##  URL {
##    Scheme https
##    Opaque 
##    Host google.com
##    Path 
##    RawPath 
##    OmitHost false
##    ForceQuery false
##    RawQuery 
##    Fragment 
##    RawFragment 
##  }
##  Proto HTTP/1.1
##  ProtoMajor 1
##  ProtoMinor 1
##  Header {
##  }
##  ContentLength 0
##  Close false
##  Host google.com
##  RemoteAddr 
##  RequestURI 
##}
`
	out := dst.String()
	if out != want {
		t.Errorf("[Indent] in: %s; out: %s; want: %s", enc, out, want)
	}
}

func TestMarshalIndent(t *testing.T) {
	req, err := http.NewRequest("GET", "https://google.com", nil)
	if err != nil {
		t.Error(err)
		return
	}
	enc, err := MarshalIndent(req, "\t", "\t")
	if err != nil {
		t.Error(err)
		return
	}
	want := "{\n" +
		"\t\tMethod GET\n" +
		"\t\tURL {\n" +
		"\t\t\tScheme https\n" +
		"\t\t\tOpaque \n" +
		"\t\t\tHost google.com\n" +
		"\t\t\tPath \n" +
		"\t\t\tRawPath \n" +
		"\t\t\tOmitHost false\n" +
		"\t\t\tForceQuery false\n" +
		"\t\t\tRawQuery \n" +
		"\t\t\tFragment \n" +
		"\t\t\tRawFragment \n" +
		"\t\t}\n" +
		"\t\tProto HTTP/1.1\n" +
		"\t\tProtoMajor 1\n" +
		"\t\tProtoMinor 1\n" +
		"\t\tHeader {\n" +
		"\t\t}\n" +
		"\t\tContentLength 0\n" +
		"\t\tClose false\n" +
		"\t\tHost google.com\n" +
		"\t\tRemoteAddr \n" +
		"\t\tRequestURI \n" +
		"\t}\n"
	if string(enc) != want {
		t.Errorf("[MarshalIndent] in: %s; out: %s; want: %s", enc, string(enc), want)
	}
}

func TestCompact(t *testing.T) {
	req, err := http.NewRequest("GET", "https://google.com", nil)
	if err != nil {
		t.Error(err)
		return
	}
	enc, err := Marshal(req)
	if err != nil {
		t.Error(err)
		return
	}
	ind := bytes.Buffer{}
	if err = Indent(&ind, enc, "\t", " "); err != nil {
		t.Error(err)
		return
	}
	dst := bytes.Buffer{}
	if err = Compact(&dst, ind.Bytes()); err != nil {
		t.Error(err)
		return
	}
	want := `{
Method GET
URL {
Scheme https
Opaque 
Host google.com
Path 
RawPath 
OmitHost false
ForceQuery false
RawQuery 
Fragment 
RawFragment 
}
Proto HTTP/1.1
ProtoMajor 1
ProtoMinor 1
Header {
}
ContentLength 0
Close false
Host google.com
RemoteAddr 
RequestURI 
}
`
	out := dst.String()
	if out != want {
		t.Errorf("[Compact] in: %s; out: %s; want: %s", ind.String(), out, want)
	}
	if out != string(enc) {
		t.Errorf("[Compact] in: %s; out: %s; want: %s", ind.String(), out, string(enc))
	}
}
