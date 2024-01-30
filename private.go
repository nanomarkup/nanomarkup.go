package nanomarkup

import (
	"bytes"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type decoder struct {
	data  [][]byte
	index int
}

type unmarshalType int64

const (
	undefined unmarshalType = iota
	entity
	array
)

const (
	errorContextFmt string = "[%s] %s"
)

func marshal(data any) ([]byte, error) {
	val := reflect.ValueOf(data)
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
		return []byte(strings.TrimSpace(val.String())), nil
	case reflect.Bool:
		return []byte(strconv.FormatBool(val.Bool())), nil
	case reflect.Slice, reflect.Array:
		if val.Len() == 0 {
			return []byte(""), nil
		} else {
			return marshalMap(val)
		}
	case reflect.Map:
		if val.Len() == 0 {
			return []byte(""), nil
		} else {
			return marshalSlice(val)
		}
	case reflect.Struct:
		if val.IsZero() {
			return []byte(""), nil
		}
		return marshalStruct(data)
	default:
		return []byte(""), nil
	}
}

func marshalStruct(data any) ([]byte, error) {
	typ := reflect.TypeOf(data)
	val := reflect.ValueOf(data)
	res := []byte("{\n")
	for _, f := range reflect.VisibleFields(typ) {
		if !f.IsExported() {
			continue
		}
		res = append(res, []byte(f.Name+" ")...)
		v, e := marshal(val.Field(f.Index[0]).Interface())
		if e != nil {
			return nil, e
		}
		res = append(res, v...)
		if string(res[len(res)-1]) != "\n" {
			res = append(res, "\n"...)
		}
	}
	if len(res) == 2 {
		return []byte(""), nil
	}
	res = append(res, []byte("}\n")...)
	return res, nil
}

func marshalSlice(value reflect.Value) ([]byte, error) {
	res := []byte("[\n")
	for i := 0; i < value.Len(); i++ {
		v, e := marshal(value.Index(i).Interface())
		if e != nil {
			return nil, e
		}
		res = append(res, v...)
		if string(res[len(res)-1]) != "\n" {
			res = append(res, "\n"...)
		}
	}
	if len(res) == 2 {
		return []byte(""), nil
	}
	res = append(res, []byte("]\n")...)
	return res, nil
}

func marshalMap(value reflect.Value) ([]byte, error) {
	res := []byte("{\n")
	iter := value.MapRange()
	for iter.Next() {
		v, e := marshal(iter.Key().Interface())
		if e != nil {
			return nil, e
		}
		res = append(res, v...)
		res = append(res, 32) // add a space
		v, e = marshal(iter.Value().Interface())
		if e != nil {
			return nil, e
		}
		res = append(res, v...)
		if string(res[len(res)-1]) != "\n" {
			res = append(res, "\n"...)
		}
	}
	if len(res) == 2 {
		return []byte(""), nil
	}
	res = append(res, []byte("}\n")...)
	return res, nil
}

func unmarshal(d *decoder, v reflect.Value, curr unmarshalType) error {
	ind := -1
	item, ok := d.next()
	for ; ok; item, ok = d.next() {
		item = bytes.TrimLeft(item, " ")
		if len(item) == 0 {
			continue
		}
		switch item[0] {
		case 91: // [
			val := item[1:]
			if len(val) > 0 && len(strings.TrimSpace(string(val))) > 0 {
				return &InvalidEntityError{"Unmarshal", string(item), fmt.Errorf("the data of an array must be started from a new line")}
			}
			if curr == undefined {
				// set the current type and continue the parsing the rest of data
				curr = array
			} else {
				// it is an internal array, parse it using other thread/loop
				e := unmarshal(d, v, array)
				if e != nil {
					return e
				}
			}
		case 93: // ]
			if curr == array {
				return nil
			} else {
				return &InvalidEntityError{"Unmarshal", "", fmt.Errorf("'[' is missing")}
			}
		case 123: // {
			val := item[1:]
			if len(val) > 0 && len(strings.TrimSpace(string(val))) > 0 {
				return &InvalidEntityError{"Unmarshal", string(item), fmt.Errorf("the data of an entity must be started from a new line")}
			}
			if curr == undefined {
				// set the current type and continue the parsing the rest of data
				curr = entity
			} else {
				// it is an internal entity, parse it using other thread/loop
				e := unmarshal(d, v, entity)
				if e != nil {
					return e
				}
			}
		case 125: // }
			if curr == entity {
				return nil
			} else {
				return &InvalidEntityError{"Unmarshal", "", fmt.Errorf("'}' is missing")}
			}
		default:
			switch v.Kind() {
			case reflect.Array:
				ind++
				elem := v.Index(ind)
				if e := unmarshalValue(elem, string(item)); e != nil {
					return e
				}
			case reflect.Slice:
				val := reflect.New(v.Type().Elem()).Elem()
				if e := unmarshalValue(val, string(item)); e != nil {
					return e
				}
				v.Set(reflect.Append(v, val))
			case reflect.Map, reflect.Struct:
				s := strings.TrimLeft(string(item), " ")
				ks := ""
				vs := ""
				space := strings.Index(s, " ")
				if space > 0 {
					ks = s[:space]
					vs = strings.TrimLeft(s[space+1:], " ")
				} else {
					ks = s
				}
				if v.Kind() == reflect.Map {
					kv := reflect.New(v.Type().Key()).Elem()
					vv := reflect.New(v.Type().Elem()).Elem()
					if e := unmarshalValue(kv, ks); e != nil {
						return e
					}
					if e := unmarshalValue(vv, vs); e != nil {
						return e
					}
					v.SetMapIndex(kv, vv)
				} else {
					if field := v.FieldByName(ks); field.IsValid() {
						vv := reflect.New(field.Type()).Elem()
						// check for an inline entity/array
						switch vs[0] {
						case 91: // [
							val := vs[1:]
							if len(val) > 0 && len(strings.TrimSpace(string(val))) > 0 {
								return &InvalidEntityError{"Unmarshal", string(item), fmt.Errorf("the data of an array must be started from a new line")}
							} else {
								// it is an internal array, parse it using other thread/loop
								e := unmarshal(d, vv, array)
								if e != nil {
									return e
								}
							}
						case 123: // {
							val := vs[1:]
							if len(val) > 0 && len(strings.TrimSpace(string(val))) > 0 {
								return &InvalidEntityError{"Unmarshal", string(item), fmt.Errorf("the data of an entity must be started from a new line")}
							} else {
								// it is an internal entity, parse it using other thread/loop
								e := unmarshal(d, vv, entity)
								if e != nil {
									return e
								}
							}
						default:
							if e := unmarshalValue(vv, vs); e != nil {
								return e
							}
						}
						field.Set(vv)
					}
				}
			default:
				if e := unmarshalValue(v, string(item)); e != nil {
					return e
				}
			}
		}
	}
	// check the end/close operator
	switch curr {
	case array:
		return &InvalidEntityError{"Unmarshal", "", fmt.Errorf("']' is missing")}
	case entity:
		return &InvalidEntityError{"Unmarshal", "", fmt.Errorf("'}' is missing")}
	}
	return nil
}

func unmarshalValue(v reflect.Value, s string) error {
	switch v.Kind() {
	case reflect.Int:
		n, e := strconv.ParseInt(s, 0, 0)
		if e != nil {
			return e
		}
		v.SetInt(int64(n))
		return nil
	case reflect.Int8:
		n, e := strconv.ParseInt(s, 0, 8)
		if e != nil {
			return e
		}
		v.SetInt(int64(n))
		return nil
	case reflect.Int16:
		n, e := strconv.ParseInt(s, 0, 16)
		if e != nil {
			return e
		}
		v.SetInt(int64(n))
		return nil
	case reflect.Int32:
		n, e := strconv.ParseInt(s, 0, 32)
		if e != nil {
			return e
		}
		v.SetInt(int64(n))
		return nil
	case reflect.Int64:
		n, e := strconv.ParseInt(s, 0, 64)
		if e != nil {
			return e
		}
		v.SetInt(int64(n))
		return nil
	case reflect.Uint:
		n, e := strconv.ParseUint(s, 0, 0)
		if e != nil {
			return e
		}
		v.SetUint(uint64(n))
		return nil
	case reflect.Uint8:
		n, e := strconv.ParseUint(s, 0, 8)
		if e != nil {
			return e
		}
		v.SetUint(uint64(n))
		return nil
	case reflect.Uint16:
		n, e := strconv.ParseUint(s, 0, 16)
		if e != nil {
			return e
		}
		v.SetUint(uint64(n))
		return nil
	case reflect.Uint32:
		n, e := strconv.ParseUint(s, 0, 32)
		if e != nil {
			return e
		}
		v.SetUint(uint64(n))
		return nil
	case reflect.Uint64:
		n, e := strconv.ParseUint(s, 0, 64)
		if e != nil {
			return e
		}
		v.SetUint(uint64(n))
		return nil
	case reflect.Uintptr:
		return &InvalidArgumentError{"Unmarshal", fmt.Errorf("uintptr type of the second argument is not supported")}
	case reflect.Float32:
		n, e := strconv.ParseFloat(s, 32)
		if e != nil {
			return e
		}
		v.SetFloat(n)
		return nil
	case reflect.Float64:
		n, e := strconv.ParseFloat(s, 64)
		if e != nil {
			return e
		}
		v.SetFloat(n)
		return nil
	case reflect.Complex64:
		n, e := strconv.ParseComplex(s, 64)
		if e != nil {
			return e
		}
		v.SetComplex(n)
		return nil
	case reflect.Complex128:
		n, e := strconv.ParseComplex(s, 128)
		if e != nil {
			return e
		}
		v.SetComplex(n)
		return nil
	case reflect.String:
		v.SetString(s)
		return nil
	case reflect.Bool:
		n, e := strconv.ParseBool(s)
		if e != nil {
			return e
		}
		v.SetBool(n)
		return nil
	default:
		return &InvalidEntityError{"Unmarshal", s, fmt.Errorf("cannot decode the entity")}
	}
}
