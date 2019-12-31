package sfft

import "gonum.org/v1/gonum/floats"

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

// Mat3 is a structure represen {ting a 3D array with floats
type Mat3 struct {
	Data []float64
	f    flattened3
}

// Dims return the size of the matrix
func (m *Mat3) Dims() (int, int, int) {
	return m.f.Dims()
}

// At returns the value at position i
func (m *Mat3) At(i, j, k int) float64 {
	return m.Data[m.f.Index(i, j, k)]
}

// Set sets the a value at position i, j, k
func (m *Mat3) Set(i, j, k int, v float64) {
	m.Data[m.f.Index(i, j, k)] = v
}

// NewMat3 returns a new instance of the Mat3 struct
func NewMat3(nr, nc, nd int, data []float64) *Mat3 {
	var v Mat3
	if data == nil {
		v.Data = make([]float64, nr*nc*nd)
	} else {
		if len(data) != nr*nc*nd {
			panic("Inconsistent length of array passed to NewMat3")
		}
		v.Data = data
	}
	v.f = flattened3{nr: nr, nc: nc, nd: nd}
	return &v
}

// AsUint8 return the underlying slice as uint8. The data is scaled such that the largest value
// equals 255 and the smallest value equals 0
func (m *Mat3) AsUint8() []uint8 {
	maxval := floats.Max(m.Data)
	minval := floats.Min(m.Data)
	res := make([]uint8, len(m.Data))
	for i := range m.Data {
		res[i] = uint8(255 * (m.Data[i] - minval) / (maxval - minval))
	}
	return res
}

// AsComplex converts the matrix into a complex matrix with imaginary part
// set to zero
func (m *Mat3) AsComplex() *CMat3 {
	nr, nc, nd := m.Dims()
	cmat := NewCMat3(nr, nc, nd, nil)
	for i := range m.Data {
		cmat.Data[i] = complex(m.Data[i], 0.0)
	}
	return cmat
}
