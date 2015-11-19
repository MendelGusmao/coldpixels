package coldpixels

type Uint16Slice []uint16
type IntSlice []int

func (f Uint16Slice) Make(len uint16) Uint16Slice {
	f = make([]uint16, len)

	return f
}

func (f Uint16Slice) Set(value uint16) Uint16Slice {
	for i := range f {
		f[i] = value
	}

	return f
}

func (f Uint16Slice) Map(fn func(uint16) uint16) Uint16Slice {
	if fn != nil {
		for i := range f {
			f[i] = fn(f[i])
		}
	}

	return f
}

func (f Uint16Slice) Insert(value uint16, index int) Uint16Slice {
	f = append(f, 0)
	copy(f[index+1:], f[index:])
	f[index] = value

	return f
}

func (f Uint16Slice) Sum() uint16 {
	sum := uint16(0)

	for _, v := range f {
		sum += v
	}

	return sum
}

func (f IntSlice) Make(len int) IntSlice {
	f = make([]int, len)

	return f
}

func (f IntSlice) Set(value int) IntSlice {
	for i := range f {
		f[i] = value
	}

	return f
}

func (f IntSlice) Map(fn func(int) int) IntSlice {
	if fn != nil {
		for i := range f {
			f[i] = fn(f[i])
		}
	}

	return f
}

func (f IntSlice) Insert(value int, index int) IntSlice {
	f = append(f, 0)
	copy(f[index+1:], f[index:])
	f[index] = value

	return f
}
