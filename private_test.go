package nanomarkup

import (
	"fmt"
	"reflect"
	"strconv"
	"testing"
)

func anyToStr(v any) string {
	val := reflect.ValueOf(v)
	switch val.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(val.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return strconv.FormatUint(val.Uint(), 10)
	case reflect.Float32, reflect.Float64:
		return strconv.FormatFloat(val.Float(), 'g', -1, 64)
	case reflect.Complex64, reflect.Complex128:
		return strconv.FormatComplex(val.Complex(), 'g', -1, 128)
	case reflect.String:
		return v.(string)
	case reflect.Bool:
		return strconv.FormatBool(val.Bool())
	case reflect.Slice, reflect.Array:
		if val.Len() == 0 {
			return ""
		} else {
			switch reflect.TypeOf(v).Elem().Kind() {
			case reflect.Uint8:
				return string(v.([]byte))
			default:
				return fmt.Sprintf("%#v", v)
			}
		}
	case reflect.Map:
		if val.Len() == 0 {
			return ""
		} else {
			return fmt.Sprintf("%#v", v)
		}
	case reflect.Struct:
		return fmt.Sprintf("%#v", v)
	default:
		return fmt.Sprintf("%#v", v)
	}
}

func checkMarshal(in any, out any, want any, e error) string {
	sout := anyToStr(out)
	swant := anyToStr(want)
	if e == nil && sout == swant {
		return ""
	}
	res := fmt.Sprintf("[Marshal] in: %#v; out: %s; want: %s", in, sout, swant)
	if e != nil {
		res += "; error: " + e.Error()
	}
	return res
}

func checkUnmarshalInt(in string, out int64, want int64, e error) string {
	if e == nil && out == want {
		return ""
	}
	res := fmt.Sprintf("[Marshal] in: %s; out: %d; want: %d", in, out, want)
	if e != nil {
		res += "; error: " + e.Error()
	}
	return res
}

func checkUnmarshalUint(in string, out uint64, want uint64, e error) string {
	if e == nil && out == want {
		return ""
	}
	res := fmt.Sprintf("[Marshal] in: %s; out: %d; want: %d", in, out, want)
	if e != nil {
		res += "; error: " + e.Error()
	}
	return res
}

func checkUnmarshalFloat(in string, out float64, want float64, e error) string {
	if e == nil && out == want {
		return ""
	}
	res := fmt.Sprintf("[Marshal] in: %s; out: %f; want: %f", in, out, want)
	if e != nil {
		res += "; error: " + e.Error()
	}
	return res
}

func checkUnmarshalComplex(in string, out complex128, want complex128, e error) string {
	if e == nil && out == want {
		return ""
	}
	res := fmt.Sprintf("[Marshal] in: %s; out: %v; want: %v", in, out, want)
	if e != nil {
		res += "; error: " + e.Error()
	}
	return res
}

func checkUnmarshalString(in string, out string, want string, e error) string {
	if e == nil && out == want {
		return ""
	}
	res := fmt.Sprintf("[Marshal] in: %s; out: %s; want: %s", in, out, want)
	if e != nil {
		res += "; error: " + e.Error()
	}
	return res
}

func checkUnmarshalBool(in string, out bool, want bool, e error) string {
	if e == nil && out == want {
		return ""
	}
	res := fmt.Sprintf("[Marshal] in: %s; out: %t; want: %t", in, out, want)
	if e != nil {
		res += "; error: " + e.Error()
	}
	return res
}

func testStructs(t *testing.T, st1 any, st2 any) {
	typ := reflect.TypeOf(st1).Elem()
	isNilF1 := false
	isNilF2 := false
	var f1 reflect.Value
	for i := 0; i < typ.NumField(); i++ {
		f := typ.Field(i)
		// check if the field is public
		if f.PkgPath == "" {
			f1 = reflect.ValueOf(st1).Elem().FieldByName(f.Name)
			f2 := reflect.ValueOf(st2).Elem().FieldByName(f.Name)
			isNilF1 = isValueNil(f1)
			isNilF2 = isValueNil(f2)
			if isNilF1 && isNilF2 {
				continue
			} else if (isNilF1 && !isNilF2) || (!isNilF1 && isNilF2) {
				if isNilF1 {
					t.Errorf("%s field of the request data is Nil", f.Name)
				} else if isNilF2 {
					t.Errorf("%s field of the decoded data is Nil", f.Name)
				}
				continue
			}
			if f1.Kind() == reflect.Pointer {
				testStructs(t, f1.Interface(), f2.Interface())
			} else {
				v1 := f1.Interface()
				v2 := f2.Interface()
				if !reflect.DeepEqual(v1, v2) {
					t.Errorf("A value of '%s' field is different:\n'%v'\n'%v'", f.Name, v1, v2)
				}
			}
		}
	}
}
