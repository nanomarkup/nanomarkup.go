package nanocomment

import "fmt"

func (c *Comments) Add(value string, multiline bool) {
	*c = append(*c, &Comment{value, multiline})
}

func (c *Comments) Adds(comments Comments) {
	*c = append(*c, comments...)
}

func (c Comments) String() string {
	out := ""
	first := true
	for _, v := range c {
		if first {
			out = v.value
			first = false
		} else {
			out += fmt.Sprintf("\n%s", v.value)
		}
	}
	return out
}
