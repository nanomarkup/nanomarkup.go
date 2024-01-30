package nanomarkup

import (
	"fmt"
	"reflect"
	"strconv"
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
