// Package sfft provides a set of method to simplify calculations of 2D and 3D
// FFTs. It is solely based on the 1D FFT method provided by Gonum.
package sfft

import (
	"gonum.org/v1/gonum/fourier"
)

// GonumFT is a type definition of Gonum's Coefficients and Sequence
type GonumFT func(dst []complex128, data []complex128) []complex128

// FFT1 is a data type for 1D FFTs
type FFT1 struct {
	ft *fourier.FFT
	n  int
}

// NewFFT1 creates a new type for FFT1. Size is the length of the array that will be
// Fourier Transformed
func NewFFT1(size int) *FFT1 {
	return &FFT1{
		ft: fourier.NewFFT(size),
		n:  size,
	}
}

// FFT performs forward FFT. The length of the data array has to match
// the size passed when the type was initialized
func (f *FFT1) FFT(data []float64) []complex128 {
	dst := make([]complex128, len(data)/2+1)
	f.ft.Coefficients(dst, data)
	return dst
}

// IFFT performs inverse FFT. The length of the passed slice has to be equal to the
// one returned by FFT (e.g. size/2+1, where size is the value passed on initialization
// to NewFFT1)
func (f *FFT1) IFFT(coeff []complex128) []float64 {
	dst := make([]float64, 2*(len(coeff)-1))
	f.ft.Sequence(dst, coeff)
	return dst
}

// Freq return the frequency corresponding to the coefficient at index i. The spacing
// is assumed to be 1.0
func (f *FFT1) Freq(i int) float64 {
	freq := float64(i) / float64(f.n)
	if i > f.n/2 {
		freq = freq - 1.0
	}
	return freq
}

// FFT2 is a data type for two dimensional Fourier Transforms
type FFT2 struct {
	ftRow *fourier.CmplxFFT
	ftCol *fourier.CmplxFFT
	nr    int
	nc    int
	cols  []int
	rows  []int
}

// NewFFT2 return a new FFT2. nr is the number of rows, and nc is the number of columns
func NewFFT2(nr, nc int) *FFT2 {
	cols := make([]int, nc)
	rows := make([]int, nr)
	for i := 0; i < nc; i++ {
		cols[i] = i
	}
	for i := 0; i < nr; i++ {
		rows[i] = i
	}
	return &FFT2{
		ftRow: fourier.NewCmplxFFT(nc),
		ftCol: fourier.NewCmplxFFT(nr),
		nr:    nr,
		nc:    nc,
		cols:  cols,
		rows:  rows,
	}
}

// FFT performs forward FFT. Data is assumed to be flattened row-major
// (e.g. A(i, j) = data[i*nc + j] where A is the 2D matrix). Therefore,
// the length of the data array has to be nr*nc, where nr and nc is the
// values used on initialization in NewFFT2
func (f *FFT2) FFT(data []complex128) []complex128 {
	if len(data) != f.nr*f.nc {
		panic("FFT: Inconsistent size in 2D FFT")
	}
	f.RowTransform(data, f.ftRow.Coefficients)
	f.ColTransform(data, f.ftCol.Coefficients)
	return data
}

// RowTransform performs inplace transform per row
func (f *FFT2) RowTransform(data []complex128, op GonumFT) []complex128 {
	for _, r := range f.rows {
		row := data[r*f.nc : (r+1)*f.nc]
		op(row, row)
	}
	return data
}

// ColTransform performs in-place transform over columns
func (f *FFT2) ColTransform(data []complex128, op GonumFT) []complex128 {
	for _, c := range f.cols {
		col := extractComplex(data, c, f.nc)
		op(col, col)
		insertComplex(data, col, c, f.nc)
	}
	return data
}

// IFFT performs inverse Fourier Transform. The length of the passed slice has to
// match the one returned by FFT (e.g. nr*nc, where nr and nc are the values passed
// to NewFFT2)
func (f *FFT2) IFFT(coeff []complex128) []complex128 {
	if len(coeff) != f.nr*f.nc {
		panic("FFT: Inconsistent size in 2D FFT")
	}
	f.RowTransform(coeff, f.ftRow.Sequence)
	f.ColTransform(coeff, f.ftCol.Sequence)
	return coeff
}

// Freq return the 2D frequency corresponding to index i in the array returned by
// FFT. The spacing is assumed to be 1.0
func (f *FFT2) Freq(i int) []float64 {
	col := i % f.nc
	row := i / f.nc

	freqs := make([]float64, 2)
	freqs[0] = float64(row) / float64(f.nr)
	freqs[1] = float64(col) / float64(f.nc)

	if row > f.nr/2 {
		freqs[0] -= 1.0
	}

	if col > f.nc/2 {
		freqs[1] -= 1.0
	}
	return freqs
}

// FFT3 is a structure for performing 3D FFTs
type FFT3 struct {
	row   *fourier.CmplxFFT
	col   *fourier.CmplxFFT
	depth *fourier.CmplxFFT
}

// NewFFT3 returns a new 3D Fourier transform object. nr is the number of rows,
// nc is the number of columns and nd is the number of nr x nc "sheets"
func NewFFT3(nr, nc, nd int) *FFT3 {
	return &FFT3{
		row:   fourier.NewCmplxFFT(nc),
		col:   fourier.NewCmplxFFT(nr),
		depth: fourier.NewCmplxFFT(nd),
	}
}

// RowTransform performs FFT over rows
func (f *FFT3) RowTransform(data []complex128, op GonumFT) []complex128 {
	nc := f.row.Len()
	nr := f.col.Len()

	// Perform FFT over first axis
	for r := 0; r < nr; r++ {
		for d := 0; d < f.depth.Len(); d++ {
			row := data[d*nr*nc+r*nc : d*nr*nc+(r+1)*nc]
			op(row, row)
		}
	}
	return data
}

// ColTransform performs FFT over columns
func (f *FFT3) ColTransform(data []complex128, op GonumFT) []complex128 {
	nc := f.row.Len()
	nr := f.col.Len()
	for d := 0; d < f.depth.Len(); d++ {
		plane := data[d*nr*nc : (d+1)*nr*nc]
		for c := 0; c < nc; c++ {
			col := extractComplex(plane, c, nc)
			op(col, col)
			insertComplex(plane, col, c, nc)
		}
	}
	return data
}

// DepthTransform performs FFT over the "depth" of a 3D matrix
func (f *FFT3) DepthTransform(data []complex128, op GonumFT) []complex128 {
	nc := f.row.Len()
	nr := f.col.Len()
	for r := 0; r < f.row.Len(); r++ {
		for c := 0; c < f.col.Len(); c++ {
			start := c*nc + r
			seq := extractComplex(data, start, nr*nc)
			op(seq, seq)
			insertComplex(data, seq, start, nr*nc)
		}
	}
	return data
}

// fourierTransform performs forward FFT or backward FFT depending on the functions passed. tRow
// is the function used to perform FT over rows, tCol is the function used to perform FT over columns
// and tDepth is the function used to perform FT in the third direction
func (f *FFT3) fourierTransform(data []complex128, tRow GonumFT, tCol GonumFT, tDepth GonumFT) []complex128 {
	if len(data) != f.row.Len()*f.col.Len()*f.depth.Len() {
		panic("FFT3: Inconsistent length of data")
	}
	f.RowTransform(data, tRow)
	f.ColTransform(data, tCol)
	f.DepthTransform(data, tDepth)
	return data
}

// FFT performs forward fourier transform. The length of the passed array has to be equal to
// nr*nc*nd, where nr, nc and nd are the values passed to NewFFT3
func (f *FFT3) FFT(data []complex128) []complex128 {
	return f.fourierTransform(data, f.row.Coefficients, f.col.Coefficients, f.depth.Coefficients)
}

// IFFT performs the inverse fourier transform. The length of the passed array has to match
// the one returned by FFT (e.g. nr*nc*nd)
func (f *FFT3) IFFT(data []complex128) []complex128 {
	return f.fourierTransform(data, f.row.Sequence, f.col.Sequence, f.depth.Sequence)
}

// Freq returns the frequency correpsondex to index i in the array returned
// by FFT
func (f *FFT3) Freq(i int) []float64 {
	nr := f.col.Len()
	nc := f.row.Len()
	nd := f.depth.Len()

	c := i % nc
	r := (i / nc) % nr
	d := i / (nr * nc)

	freq := make([]float64, 3)
	freq[0] = float64(c) / float64(nc)
	freq[1] = float64(r) / float64(nr)
	freq[2] = float64(d) / float64(nd)

	if c > nc/2 {
		freq[0] -= 1.0
	}

	if r > nr/2 {
		freq[1] -= 1.0
	}

	if d > nd/2 {
		freq[2] -= 1.0
	}
	return freq
}
