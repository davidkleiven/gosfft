package sfft

import "gonum.org/v1/gonum/fourier"

// FFT1 is a data type for 1D FFTs
type FFT1 struct {
	ft *fourier.FFT
}

// NewFFT1 creates a new type for FFT1
func NewFFT1(size int) *FFT1 {
	return &FFT1{
		ft: fourier.NewFFT(size),
	}
}

// FFT performs forward FFT
func (f *FFT1) FFT(data []float64) []complex128 {
	dst := make([]complex128, len(data)/2+1)
	f.ft.Coefficients(dst, data)
	return dst
}

// IFFT performs inverse FFT
func (f *FFT1) IFFT(coeff []complex128) []float64 {
	dst := make([]float64, 2*(len(coeff)-1))
	f.ft.Sequence(dst, coeff)
	return dst
}
