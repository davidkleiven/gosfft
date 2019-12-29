package sfft

import (
	"gonum.org/v1/gonum/fourier"
)

// FFT1 is a data type for 1D FFTs
type FFT1 struct {
	ft *fourier.FFT
	n  int
}

// NewFFT1 creates a new type for FFT1. Size is the length of the size
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

// Freq return the frequency corresponding to the coefficient at index i
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
}

// NewFFT2 return a new FFT2. nr is the number of rows, and nc is the number of columns
func NewFFT2(nr, nc int) *FFT2 {
	return &FFT2{
		ftRow: fourier.NewCmplxFFT(nc),
		ftCol: fourier.NewCmplxFFT(nr),
		nr:    nr,
		nc:    nc,
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

	// Perform first axis transform
	for r := 0; r < f.nr; r++ {
		row := data[r*f.nc : (r+1)*f.nc]
		f.ftRow.Coefficients(row, row)
	}

	// Perform second axis transform
	for c := 0; c < f.nc; c++ {
		col := extractComplex(data, c, f.nc)
		f.ftCol.Coefficients(col, col)
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

	// Perform first axis transform
	for r := 0; r < f.nr; r++ {
		row := coeff[r*f.nc : (r+1)*f.nc]
		f.ftRow.Sequence(row, row)
	}

	// Perform second axis transform
	for c := 0; c < f.nc; c++ {
		col := extractComplex(coeff, c, f.nc)
		f.ftCol.Sequence(col, col)
		insertComplex(coeff, col, c, f.nc)
	}
	return coeff
}

// Freq return the 2D frequency corresponding to index i in the array returned by
// FFT
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
