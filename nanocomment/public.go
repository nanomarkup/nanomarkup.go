package nanocomment

import (
	"bytes"
	"fmt"

	"nanomarkup.go/nanodecoder"
)

type Comment struct {
	value     string
	multiline bool
}

type Comments []*Comment

const (
	SingleCommentOpCode       string = "//"
	MultilineCommentBegOpCode string = "/*"
	MultilineCommentEndOpCode string = "*/"
)

func Marshal(comments Comments) []byte {
	out := []byte{}
	for _, v := range comments {
		if v.multiline {
			out = append(out, []byte(MultilineCommentBegOpCode+v.value+MultilineCommentEndOpCode+"\n")...)
		} else if v.value == "" {
			out = append(out, []byte("\n")...)
		} else {
			out = append(out, []byte(SingleCommentOpCode+v.value+"\n")...)
		}
	}
	return out
}

func Unmarshal(d *nanodecoder.Decoder, item []byte) (Comments, error) {
	l := 0
	ok := true
	out := Comments{}
	mcomm := ""
	multi := false
	handled := false
	var b []byte
	for ; ok; item, ok = d.Next() {
		if multi {
			// check MultilineCommentEndOpCode
			b = bytes.TrimRight(item, " ")
			l = len(b)
			if l > 1 && b[l-1] == 47 && b[l-2] == 42 {
				// add a multiline comment
				multi = false
				mcomm += string(b[0 : l-2])
				out.Add(mcomm, true)
				continue
			}
			// add item to multiline comment
			mcomm = fmt.Sprintf("%s%s\n", mcomm, string(item))
		} else {
			// check SingleCommentOpCode and MultilineCommentBegOpCode, where / = 47, * = 42
			item = bytes.TrimLeft(item, " \t")
			l = len(item)
			if l == 0 {
				out.Add("", false)
			} else if l > 1 && item[0] == 47 && (item[1] == 42 || item[1] == 47) {
				if item[1] == 47 {
					// add a single comment
					out.Add(string(item[2:]), false)
				} else {
					// check MultilineCommentEndOpCode
					b = bytes.TrimRight(item[2:], " ")
					l = len(b)
					if l > 1 && b[l-1] == 47 && b[l-2] == 42 {
						// add a multiline comment
						out.Add(string(b[0:l-2]), true)
					} else {
						// add item to multiline comment
						multi = true
						mcomm = fmt.Sprintf("%s\n", string(item[2:]))
					}
				}
			} else {
				// it is not a single comment or multiline comment
				break
			}
		}
		handled = true
	}
	if ok && handled {
		d.Prev()
	}
	return out, nil
}
