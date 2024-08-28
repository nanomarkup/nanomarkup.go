package nanometadata

import "github.com/nanomarkup/nanomarkup.go/nanocomment"

type Metadata struct {
	fields   map[string]*Metadata
	Comments nanocomment.Comments
}

func CreateMetadata(comment string, multiline bool) *Metadata {
	m := Metadata{}
	m.Comments.Add(comment, multiline)
	return &m
}
