package sfft

import (
	"math"
	"testing"
)

func TestFFT1(t *testing.T) {
	ft := NewFFT1(8)
	data := make([]float64, 8)
	data[0] = 1.0
	data[2] = 1.0

	expect := make([]complex128, 5)
	expect[0] = complex(2.0, 0.0)
	expect[1] = complex(1.0, -1.0)
	expect[2] = complex(0.0, 0.0)
	expect[3] = complex(1.0, 1.0)
	expect[4] = complex(2.0, 0.0)

	coeff := ft.FFT(data)
	tol := 1e-10

	for i := range expect {
		if math.Abs(real(expect[i])-real(coeff[i])) > tol || math.Abs(imag(expect[i])-imag(coeff[i])) > tol {
			t.Errorf("Expected %v got %v", expect[i], coeff[i])
		}
	}

	inv := ft.IFFT(coeff)

	for i := range inv {
		if math.Abs(inv[i]-float64(len(data))*data[i]) > tol {
			t.Errorf("Expected %v got %v\n", inv, data)
			break
		}
	}
}
