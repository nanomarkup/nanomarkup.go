package nanometadata

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
