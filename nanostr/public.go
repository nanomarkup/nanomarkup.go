package nanostr

import (
	"fmt"
	"strings"

	"nanomarkup.go/nanodecoder"
	"nanomarkup.go/nanoerror"
)

func Marshal(value string) []byte {
	lines := strings.Split(value, "\n")
	if len(lines) == 1 {
		return []byte(strings.TrimLeft(value, " \t"))
	} else {
		res := "`\n"
		for _, it := range lines {
			res += it + "\n"
		}
		return []byte(res + "`\n")
	}
}

func Unmarshal(d *nanodecoder.Decoder, item []byte) ([]byte, error) {
	if len(item) > 0 && item[0] == 96 { // `
		// update the item variable by a multi-line value
		val := item[1:]
		if len(val) > 0 && len(strings.TrimSpace(string(val))) > 0 {
			return item, &nanoerror.InvalidEntityError{"Parse", string(item), fmt.Errorf("the data of a multi-line value must be started from a new line")}
		}
		mval := []byte{}
		first := true
		completed := false
		item, ok := d.Next()
		for ; ok; item, ok = d.Next() {
			if len(item) == 0 {
				mval = append(mval, "\n"...)
			} else if len(item) == 1 && item[0] == 96 { // `
				completed = true
				break
			} else {
				if !first {
					mval = append(mval, "\n"...)
				}
				mval = append(mval, item...)
			}
			if first {
				first = false
			}
		}
		if completed {
			return mval, nil
		} else {
			return item, &nanoerror.InvalidEntityError{"Parse", "", fmt.Errorf("'`' is missing")}
		}
	} else {
		return item, nil
	}
}
