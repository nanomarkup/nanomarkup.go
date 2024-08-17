package nanomarkup

import (
	"bytes"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"nanomarkup.go/nanocomment"
	"nanomarkup.go/nanodecoder"
	"nanomarkup.go/nanoerror"
	"nanomarkup.go/nanostr"
)

type omitEmpty bool

type unmarshalType int64

const (
	undefined unmarshalType = iota
	entity
	array
)

const (
	tagNameDelim     string = " "
	tagValueDelim    string = ","
	nanoTagName      string = "nano"
	nanoTagIgnore    string = "-"
	nanoTagOmitEmpty string = "omitempty"
)

func marshal(data any, meta *Metadata) ([]byte, error) {
	val := reflect.ValueOf(data)
	if isValueNil(val) {
		return []byte(""), nil
	}
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
		return nanostr.Marshal(val.String()), nil
	case reflect.Bool:
		return []byte(strconv.FormatBool(val.Bool())), nil
	case reflect.Slice, reflect.Array:
		if val.Len() == 0 {
			return []byte("[\n]\n"), nil
		} else {
			return marshalMap(val)
		}
	case reflect.Map:
		if val.Len() == 0 {
			return []byte("{\n}\n"), nil
		} else {
			return marshalSlice(val)
		}
	case reflect.Struct:
		if val.IsZero() {
			return []byte("{\n}\n"), nil
		}
		return marshalStruct(data, meta)
	default:
		return []byte(""), nil
	}
}

func marshalStruct(data any, meta *Metadata) ([]byte, error) {
	val := reflect.ValueOf(data)
	if val.Kind() == reflect.Pointer {
		val = val.Elem()
	}
	typ := val.Type()
	res := []byte("{\n")
	for _, f := range reflect.VisibleFields(typ) {
		if !f.IsExported() {
			continue
		}
		fv := val.Field(f.Index[0])
		if isValueNil(fv) {
			continue
		}
		// handle a nano tag
		data := fv.Interface()
		name := ""
		omitempty := false
		tag, ok := f.Tag.Lookup(nanoTagName)
		if ok {
			if tag == nanoTagIgnore {
				continue
			} else if tag == nanoTagOmitEmpty {
				omitempty = true
			} else {
				items := strings.Split(tag, tagValueDelim)
				l := len(items)
				if l == 1 {
					name = tag
				} else if l > 1 {
					if items[0] == nanoTagOmitEmpty {
						name = items[1]
						omitempty = true
					} else if items[1] == nanoTagOmitEmpty {
						name = items[0]
						omitempty = true
					} else {
						name = items[0]
					}
				}
				if name == nanoTagIgnore {
					continue
				} else if name == nanoTagOmitEmpty {
					name = ""
				}
			}
		}
		if omitempty && isEmpty(data) {
			continue
		}
		if name == "" {
			name = f.Name
		}
		// handle a metadata
		var fmeta *Metadata = nil
		if meta != nil {
			fmeta := meta.GetField(f.Name)
			if fmeta != nil && len(fmeta.Comments) > 0 {
				res = append(res, nanocomment.Marshal(fmeta.Comments)...)
			}
		}
		res = append(res, []byte(name+" ")...)
		v, e := marshal(data, fmeta)
		if e != nil {
			return nil, e
		}
		res = append(res, v...)
		if string(res[len(res)-1]) != "\n" {
			res = append(res, "\n"...)
		}
	}
	res = append(res, []byte("}\n")...)
	return res, nil
}

func marshalSlice(value reflect.Value) ([]byte, error) {
	res := []byte("[\n")
	for i := 0; i < value.Len(); i++ {
		v, e := marshal(value.Index(i).Interface(), nil)
		if e != nil {
			return nil, e
		}
		res = append(res, v...)
		if string(res[len(res)-1]) != "\n" {
			res = append(res, "\n"...)
		}
	}
	res = append(res, []byte("]\n")...)
	return res, nil
}

func marshalMap(value reflect.Value) ([]byte, error) {
	res := []byte("{\n")
	iter := value.MapRange()
	for iter.Next() {
		v, e := marshal(iter.Key().Interface(), nil)
		if e != nil {
			return nil, e
		}
		res = append(res, v...)
		res = append(res, 32) // add a space
		v, e = marshal(iter.Value().Interface(), nil)
		if e != nil {
			return nil, e
		}
		res = append(res, v...)
		if string(res[len(res)-1]) != "\n" {
			res = append(res, "\n"...)
		}
	}
	res = append(res, []byte("}\n")...)
	return res, nil
}

func unmarshal(d *nanodecoder.Decoder, v reflect.Value, curr unmarshalType, meta *Metadata) error {
	if meta != nil && meta.fields == nil {
		meta.fields = make(map[string]*Metadata)
	}
	ind := -1
	item, ok := d.Next()
	comments := nanocomment.Comments{}
	for ; ok; item, ok = d.Next() {
		item = bytes.TrimLeft(item, " \t")
		if len(item) == 0 {
			continue
		}
		switch item[0] {
		case 91: // [
			val := item[1:]
			if len(val) > 0 && len(strings.TrimSpace(string(val))) > 0 {
				return &nanoerror.InvalidEntityError{Context: "Unmarshal", Entity: string(item), Err: fmt.Errorf("the data of an array must be started from a new line")}
			}
			if curr == undefined {
				// set the current type and continue the parsing the rest of data
				curr = array
				if meta != nil && len(comments) > 0 {
					meta.Comments.Adds(comments)
					comments = nanocomment.Comments{}
				}
			} else {
				// it is an internal array, parse it using other thread/loop
				e := unmarshal(d, v, array, meta)
				if e != nil {
					return e
				}
			}
		case 93: // ]
			if curr == array {
				return nil
			} else {
				return &nanoerror.InvalidEntityError{Context: "Unmarshal", Entity: "", Err: fmt.Errorf("'[' is missing")}
			}
		case 123: // {
			val := item[1:]
			if len(val) > 0 && len(strings.TrimSpace(string(val))) > 0 {
				return &nanoerror.InvalidEntityError{Context: "Unmarshal", Entity: string(item), Err: fmt.Errorf("the data of an entity must be started from a new line")}
			}
			if curr == undefined {
				// set the current type and continue the parsing the rest of data
				curr = entity
				if meta != nil && len(comments) > 0 {
					meta.Comments.Adds(comments)
					comments = nanocomment.Comments{}
				}
			} else {
				// it is an internal entity, parse it using other thread/loop
				name := v.Elem().String()
				var e error
				if meta == nil {
					e = unmarshal(d, v, entity, nil)
				} else {
					meta.AddField(name, &Metadata{})
					e = unmarshal(d, v, entity, meta.GetField(name))
				}
				if e != nil {
					return e
				}
			}
		case 125: // }
			if curr == entity {
				return nil
			} else {
				return &nanoerror.InvalidEntityError{Context: "Unmarshal", Entity: "", Err: fmt.Errorf("'}' is missing")}
			}
		default:
			comms, err := nanocomment.Unmarshal(d, item)
			if err != nil {
				return err
			} else if len(comms) > 0 {
				if meta != nil {
					comments = comms
				}
				continue
			}
			item, err := nanostr.Unmarshal(d, item)
			if err != nil {
				return err
			}

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
				s := bytes.TrimLeft(item, " \t")
				var ks []byte
				var vs []byte
				space := bytes.Index(s, []byte(" "))
				if space > 0 {
					ks = s[:space]
					vs = bytes.TrimLeft(s[space+1:], " \t")
				} else {
					ks = s
					vs = []byte{}
				}
				vs, err := nanostr.Unmarshal(d, vs)
				if err != nil {
					return err
				}

				if v.Kind() == reflect.Map {
					kv := reflect.New(v.Type().Key()).Elem()
					vv := reflect.New(v.Type().Elem()).Elem()
					if e := unmarshalValue(kv, string(ks)); e != nil {
						return e
					}
					if e := unmarshalValue(vv, string(vs)); e != nil {
						return e
					}
					v.SetMapIndex(kv, vv)
				} else {
					if field, name, omitempty := getField(v, string(ks)); field.IsValid() {
						var vv reflect.Value
						if field.Type().Kind() == reflect.Pointer {
							if field.Type().Elem().Kind() == reflect.Map {
								vv = reflect.MakeMap(field.Type().Elem())
							} else {
								vv = reflect.New(field.Type().Elem()).Elem()
							}
						} else if field.Type().Kind() == reflect.Map {
							vv = reflect.MakeMap(field.Type())
						} else {
							vv = reflect.New(field.Type()).Elem()
						}
						if len(vs) > 0 {
							// check for an inline entity/array
							switch vs[0] {
							case 91: // [
								val := vs[1:]
								if len(val) > 0 && len(strings.TrimSpace(string(val))) > 0 {
									return &nanoerror.InvalidEntityError{Context: "Unmarshal", Entity: string(item), Err: fmt.Errorf("the data of an array must be started from a new line")}
								} else {
									// it is an internal array, parse it using other thread/loop
									e := unmarshal(d, vv, array, meta)
									if e != nil {
										return e
									}
								}
							case 123: // {
								val := vs[1:]
								if len(val) > 0 && len(strings.TrimSpace(string(val))) > 0 {
									return &nanoerror.InvalidEntityError{Context: "Unmarshal", Entity: string(item), Err: fmt.Errorf("the data of an entity must be started from a new line")}
								} else {
									// it is an internal entity, parse it using other thread/loop
									var e error
									if meta == nil {
										e = unmarshal(d, vv, entity, nil)
									} else {
										name := v.Elem().String()
										meta.AddField(name, &Metadata{})
										e = unmarshal(d, vv, entity, meta.GetField(name))
									}
									if e != nil {
										return e
									}
								}
							default:
								if e := unmarshalValue(vv, string(vs)); e != nil {
									return e
								}
							}
						}
						if bool(omitempty) && isEmpty(vv.Interface()) {
							continue
						}
						if field.Type().Kind() == reflect.Pointer {
							field.Set(vv.Addr())
						} else {
							field.Set(vv)
						}
						if meta != nil && len(comments) > 0 {
							m := Metadata{}
							m.Comments.Adds(comments)
							meta.AddField(name, &m)
							comments = nanocomment.Comments{}
						}
					}
				}
			default:
				if meta != nil && len(comments) > 0 {
					meta.Comments.Adds(comments)
					comments = nanocomment.Comments{}
				}
				if e := unmarshalValue(v, string(item)); e != nil {
					return e
				}
			}
		}
	}
	// check the end/close operator
	switch curr {
	case array:
		return &nanoerror.InvalidEntityError{Context: "Unmarshal", Entity: "", Err: fmt.Errorf("']' is missing")}
	case entity:
		return &nanoerror.InvalidEntityError{Context: "Unmarshal", Entity: "", Err: fmt.Errorf("'}' is missing")}
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
		return &nanoerror.InvalidArgumentError{Context: "Unmarshal", Err: fmt.Errorf("uintptr type of the second argument is not supported")}
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
		return &nanoerror.InvalidEntityError{Context: "Unmarshal", Entity: s, Err: fmt.Errorf("cannot decode the entity")}
	}
}

func isEmpty(v interface{}) bool {
	if v == nil {
		return true
	}
	val := reflect.ValueOf(v)
	switch val.Kind() {
	case reflect.Array, reflect.Slice:
		return val.Len() == 0
	default:
		return val.IsZero()
	}
}

func isValueNil(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Chan, reflect.Func, reflect.Map, reflect.Pointer, reflect.UnsafePointer, reflect.Interface, reflect.Slice:
		return v.IsNil()
	default:
		return false
	}
}

func getField(src reflect.Value, name string) (reflect.Value, string, omitEmpty) {
	rValue := reflect.Value{}
	var rEmpty omitEmpty = true
	if src.Kind() != reflect.Struct {
		return rValue, name, rEmpty
	}
	// nano tag has more priority than a field of struct
	for _, f := range reflect.VisibleFields(src.Type()) {
		if !f.IsExported() {
			continue
		}
		fv := src.Field(f.Index[0])
		if isValueNil(fv) {
			continue
		}
		tag, ok := f.Tag.Lookup(nanoTagName)
		if !ok || tag == nanoTagIgnore || tag == nanoTagOmitEmpty {
			continue
		}
		items := strings.Split(tag, tagValueDelim)
		fn := ""
		l := len(items)
		if l == 1 {
			fn = tag
		} else if l > 1 {
			if items[0] == nanoTagOmitEmpty {
				fn = items[1]
				rEmpty = true
			} else if items[1] == nanoTagOmitEmpty {
				fn = items[0]
				rEmpty = true
			} else {
				fn = items[0]
			}
		}
		if fn != nanoTagIgnore && fn != nanoTagOmitEmpty && fn == name {
			return fv, f.Name, rEmpty
		}
	}
	// check field
	sf, ok := src.Type().FieldByName(name)
	if !ok {
		return rValue, sf.Name, rEmpty
	} else {
		rValue = src.FieldByName(sf.Name)
		tag, ok := sf.Tag.Lookup(nanoTagName)
		if !ok {
			return rValue, sf.Name, false
		} else if tag == nanoTagIgnore {
			return reflect.Value{}, sf.Name, false
		} else if tag == nanoTagOmitEmpty {
			return rValue, sf.Name, true
		}
		items := strings.Split(tag, tagValueDelim)
		l := len(items)
		if l == 1 {
			if tag == nanoTagIgnore {
				return reflect.Value{}, sf.Name, false
			} else if tag == nanoTagOmitEmpty {
				return rValue, sf.Name, true
			}
		} else if l > 1 {
			if items[0] == nanoTagIgnore || items[1] == nanoTagIgnore {
				return reflect.Value{}, sf.Name, false
			} else if items[0] == nanoTagOmitEmpty || items[1] == nanoTagOmitEmpty {
				return rValue, sf.Name, true
			}
		}
		return rValue, sf.Name, false
	}
}

func appendIndent(dst, src []byte, prefix, indent string) ([]byte, error) {
	first := true
	level := 0
	origLen := len(dst)
	currIndent := prefix
	var err error
	d := nanodecoder.Decoder{}
	d.Init(bytes.Split(src, []byte("\n")))
	item, ok := d.Next()
	for ; ok; item, ok = d.Next() {
		item = bytes.TrimLeft(item, " \t")
		if len(item) == 0 {
			continue
		}
		switch item[0] {
		case 91: // [
			val := item[1:]
			if len(val) > 0 && len(strings.TrimSpace(string(val))) > 0 {
				err = &nanoerror.InvalidEntityError{Context: "Indent", Entity: string(item), Err: fmt.Errorf("the data of an array must be started from a new line")}
				break
			}
			if first {
				first = false
				dst = append(dst, []byte(string(item)+"\n")...)
			} else {
				dst = append(dst, []byte(currIndent+string(item)+"\n")...)
			}
			level++
			currIndent += indent
		case 93, 125: // ], }
			if level == 0 {
				err = &nanoerror.InvalidEntityError{Context: "Indent", Entity: string(item), Err: fmt.Errorf("invalid data")}
				break
			}
			level--
			if level == 0 {
				currIndent = prefix
			} else {
				currIndent = prefix + strings.Repeat(indent, level)
			}
			dst = append(dst, []byte(currIndent+string(item)+"\n")...)
		case 123: // {
			val := item[1:]
			if len(val) > 0 && len(strings.TrimSpace(string(val))) > 0 {
				err = &nanoerror.InvalidEntityError{Context: "Indent", Entity: string(item), Err: fmt.Errorf("the data of an entity must be started from a new line")}
				break
			}
			if first {
				first = false
				dst = append(dst, []byte(string(item)+"\n")...)
			} else {
				dst = append(dst, []byte(currIndent+string(item)+"\n")...)
			}
			level++
			currIndent += indent
		case 96: // `
			val, err := appendIndentMultiline(&d, prefix, currIndent, first)
			if err != nil {
				break
			}
			if val != nil {
				dst = append(dst, val...)
				first = false
			}
		default:
			if len(item) > 2 && item[0] == 47 && item[1] == 47 {
				// it is a comment
				if first {
					first = false
					dst = append(dst, []byte(string(item)+"\n")...)
				} else {
					dst = append(dst, []byte(currIndent+string(item)+"\n")...)
				}
				continue
			}
			up := false
			str := string(item)
			space := strings.Index(str, " ")
			if space > 0 {
				str = strings.TrimLeft(str[space+1:], " \t")
				if len(str) > 0 {
					if str[0] == 91 || str[0] == 123 { // [, {
						up = true
					} else if str[0] == 96 { // `
						val, err := appendIndentMultiline(&d, prefix, currIndent, first)
						if err != nil {
							break
						} else if val != nil {
							dst = append(dst, val...)
							first = false
							continue
						}
					}
				}
			}
			if first {
				first = false
				dst = append(dst, []byte(string(item)+"\n")...)
			} else {
				dst = append(dst, []byte(currIndent+string(item)+"\n")...)
			}
			if up {
				level++
				currIndent += indent
			}
		}
	}
	if err != nil {
		return dst[:origLen], err
	}
	return dst, nil
}

func appendIndentMultiline(d *nanodecoder.Decoder, prefix, currIndent string, first bool) ([]byte, error) {
	item, ok := d.Curr()
	if !ok {
		return nil, nil
	}

	dst := []byte{}
	str := string(bytes.TrimLeft(item, " \t"))
	space := strings.Index(str, " ")
	if space > 0 {
		// process key/value data
		item = []byte(strings.TrimLeft(str[space+1:], " \t"))
		if len(item) == 0 || item[0] != 96 { // `
			return nil, nil
		} else {
			str = str[:space]
			if first {
				dst = append(dst, []byte(str+" ")...)
			} else {
				dst = append(dst, []byte(currIndent+str+" ")...)
				first = true
			}
		}
	} else if item[0] != 96 { // `
		return nil, nil
	}

	val := item[1:]
	if len(val) > 0 && len(strings.TrimSpace(string(val))) > 0 {
		return nil, &nanoerror.InvalidEntityError{Context: "Indent", Entity: string(item), Err: fmt.Errorf("the data of a multi-line value must be started from a new line")}
	}
	if first {
		first = false
		dst = append(dst, []byte(string(item)+"\n")...)
	} else {
		dst = append(dst, []byte(currIndent+string(item)+"\n")...)
	}

	completed := false
	item, ok = d.Next()
	for ; ok; item, ok = d.Next() {
		if len(item) == 0 {
			dst = append(dst, "\n"...)
		} else if len(item) == 1 && item[0] == 96 { // `
			completed = true
			break
		} else {
			dst = append(dst, []byte(prefix+string(item)+"\n")...)
		}
	}
	if completed {
		dst = append(dst, []byte(currIndent+string(item)+"\n")...)
		return dst, nil
	} else {
		return nil, &nanoerror.InvalidEntityError{Context: "Indent", Entity: "", Err: fmt.Errorf("'`' is missing")}
	}
}
