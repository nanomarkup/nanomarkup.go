package nanodecoder

func (d *Decoder) Init(data [][]byte) {
	d.data = data
	d.index = -1
}

func (d *Decoder) Curr() ([]byte, bool) {
	if d.index >= len(d.data) {
		return nil, false
	} else {
		return d.data[d.index], true
	}
}

func (d *Decoder) Prev() ([]byte, bool) {
	if d.index < 1 {
		return nil, false
	} else {
		d.index--
		return d.data[d.index], true
	}
}

func (d *Decoder) Next() ([]byte, bool) {
	if d.index+1 >= len(d.data) {
		return nil, false
	} else {
		d.index++
		return d.data[d.index], true
	}
}
