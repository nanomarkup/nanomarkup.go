package nanomarkup

import (
	"fmt"
	"testing"
)

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

func TestMultiLineUnmarshal(t *testing.T) {
	// test a string
	swant := `testing
a multi
line
value`
	sin := "`\n" + swant + "\n`\n"
	sout := ""
	err := Unmarshal([]byte(sin), &sout)
	if s := checkUnmarshalString(sin, sout, swant, err); s != "" {
		t.Error(s)
	}

	// test an array/slice
	ain := "[\n`\n" + swant + "\n`\n]\n"
	awant := []string{"testing\na multi\nline\nvalue"}
	aout := []string{}
	err = Unmarshal([]byte(ain), &aout)
	pass := len(aout) == len(awant)
	if pass {
		for i := range aout {
			if aout[i] != awant[i] {
				pass = false
				break
			}
		}
	}
	mes := ""
	if !pass {
		mes = fmt.Sprintf("[Unmarshal] in: %s; out: %v; want: %v", ain, aout, awant)
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

	// test a struct
	type st struct {
		MultiValue string
	}
	tin := "{\nMultiValue `\n" + swant + "\n`\n}\n"
	twant := st{swant}
	tout := st{}
	err = Unmarshal([]byte(tin), &tout)
	if err != nil {
		t.Error(err)
	} else {
		testStructs(t, &tout, &twant)
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

func TestNanoTagUnmarshal(t *testing.T) {
	type t1 struct {
		Field1  int
		Field2  string
		Field3  int    `nano:"-"`
		Field4  int    `nano:"omitempty"`
		Field5  int    `nano:"-,omitempty"`
		Field6  string `nano:"omitempty,-"`
		Field7  int    `nano:"omitempty,test7"`
		Field8  string `nano:"test8,omitempty"`
		Field9  string `nano:"omitempty,omitempty"`
		Field10 int    `nano:"test10"`
		Field11 string `nano:"test11"`
	}
	in := `{
Field1 1
Field2 2
Field3 3
Field4 0
Field5 5
Field6 6
test7 0
test8 
Field9 
test10 10
test11 11
}
`
	out := t1{
		0,
		"",
		33,
		44,
		55,
		"66",
		77,
		"88",
		"99",
		1010,
		"1111",
	}
	want := t1{1, "2", 33, 44, 55, "66", 77, "88", "99", 10, "11"}

	err := Unmarshal([]byte(in), &out)
	mes := ""
	if out.Field1 != want.Field1 ||
		out.Field2 != want.Field2 ||
		out.Field3 != want.Field3 ||
		out.Field4 != want.Field4 ||
		out.Field5 != want.Field5 ||
		out.Field6 != want.Field6 ||
		out.Field7 != want.Field7 ||
		out.Field8 != want.Field8 ||
		out.Field9 != want.Field9 ||
		out.Field10 != want.Field10 ||
		out.Field11 != want.Field11 {
		mes = fmt.Sprintf("[Unmarshal] in: %s; out: %v; want: %v", in, out, want)
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
