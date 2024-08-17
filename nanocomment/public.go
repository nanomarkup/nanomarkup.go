package nanocomment

import (
	"bytes"

	"nanomarkup.go/nanodecoder"
)

type Comment struct {
	value     string
	multiline bool
}

type Comments []*Comment

const (
	CommentOpCode string = "//"
)

func Marshal(comments Comments) []byte {
	out := []byte{}
	for _, v := range comments {
		if v.multiline {

		} else if v.value == "" {
			out = append(out, []byte("\n")...)
		} else {
			out = append(out, []byte(CommentOpCode+v.value+"\n")...)
		}
	}
	return out
}

func Unmarshal(d *nanodecoder.Decoder, item []byte) (Comments, error) {
	if len(item) > 1 && item[0] == 47 && item[1] == 47 {
		// it is a comment
		out := Comments{}
		out.Add(string(item[2:]), false)
		// check next data
		item, ok := d.Next()
		for ; ok; item, ok = d.Next() {
			item = bytes.TrimLeft(item, " \t")
			if len(item) == 0 {
				out.Add("", false)
			} else if len(item) > 1 && item[0] == 47 && item[1] == 47 {
				out.Add(string(item[2:]), false)
			} else {
				break
			}
		}
		if ok {
			d.Prev()
		}
		return out, nil
	} else {
		return nil, nil
	}
}
