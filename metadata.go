package nanomarkup

import "fmt"

func (m *Metadata) AddField(name string, data *Metadata) {
	if m.fields == nil {
		m.fields = map[string]*Metadata{}
	}
	m.fields[name] = data
}

func (m *Metadata) GetField(name string) *Metadata {
	if m.fields == nil {
		return nil
	} else {
		return m.fields[name]
	}
}

func (m *Metadata) GetComments() string {
	out := ""
	first := true
	for _, v := range m.comments {
		if first {
			out = v.value
			first = false
		} else {
			out += fmt.Sprintf("\n%s", v.value)
		}
	}
	return out
}

func (m *Metadata) AddComment(value string, multiline bool) {
	m.comments = append(m.comments, &comment{value, multiline})
}

func (m *Metadata) AddComments(comments []*comment) {
	m.comments = append(m.comments, comments...)
}
