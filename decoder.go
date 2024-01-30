package nanomarkup

func (d *decoder) reset() {
	d.index = -1
}

func (d *decoder) next() ([]byte, bool) {
	if d.index+1 >= len(d.data) {
		return nil, false
	} else {
		d.index++
		return d.data[d.index], true
	}
}
