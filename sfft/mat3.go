package sfft

// flattened3 provides functionality of accessing a 3D array represented as a
// slice
type flattened3 struct {
	nr int
	nc int
	nd int
}

// Index returns the index corresponding to (i, j, k) element in the underlying
// data structure
func (f flattened3) Index(i, j, k int) int {
	return k*f.nr*f.nc + i*f.nc + j
}

// Dims returns the overall size in each direction
func (f flattened3) Dims() (int, int, int) {
	return f.nr, f.nc, f.nd
}

// CMat3 is a type that represents a 3D array.
type CMat3 struct {
	Data []complex128
	f    flattened3
}

// Dims returns the dimensions of the 3D array
func (m *CMat3) Dims() (int, int, int) {
	return m.f.Dims()
}

// At return the value at position (i, j, k)
func (m *CMat3) At(i, j, k int) complex128 {
	return m.Data[m.f.Index(i, j, k)]
}

// Set sets the value at position (i, j, k)
func (m *CMat3) Set(i, j, k int, v complex128) {
	m.Data[m.f.Index(i, j, k)] = v
}

// NewCMat3 returns a new CMat3 instance
func NewCMat3(nr, nc, nd int, data []complex128) *CMat3 {
	var v CMat3
	if data == nil {
		v.Data = make([]complex128, nr*nc*nd)
	} else {
		if len(data) != nr*nc*nd {
			panic("sfft: Inconsistent dimension of data passed to NewCMat3")
		}
		v.Data = data
	}
	v.f = flattened3{nr: nr, nc: nc, nd: nd}
	return &v
}
