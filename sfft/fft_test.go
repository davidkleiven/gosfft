package sfft

import (
	"math"
	"math/cmplx"
	"testing"

	"gonum.org/v1/gonum/floats"
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

	freqs := make([]float64, 5)
	for i := 0; i < 5; i++ {
		freqs[i] = ft.Freq(i)
	}
	expectFreq := []float64{0.0, 1. / 8., 1. / 4., 3. / 8., 1. / 2.}
	if !floats.EqualApprox(expectFreq, freqs, 1e-10) {
		t.Errorf("Unexpected freq. Expected\n%v\ngot\n%v\n", expectFreq, freqs)
	}
}

func TestFFT2(t *testing.T) {
	data := []float64{0.0, 0.0, 0.0, 0.0, 0.0,
		0.0, 1.0, 1.0, 0.0, 0.0,
		0.0, 0.0, 0.0, 0.0, 0.0,
		0.0, 0.0, 0.0, 0.0, 0.0}
	expectFFT := make([]complex128, 20)
	for j := 0; j < 4; j++ {
		for i := 0; i < 5; i++ {
			exp1 := complex(0, -0.5*math.Pi*float64(j))
			exp2 := complex(0, -2.0*math.Pi*float64(i)/5.0)
			exp3 := complex(0, -4.0*math.Pi*float64(i)/5.0)
			expectFFT[j*5+i] = cmplx.Exp(exp1) * (cmplx.Exp(exp2) + cmplx.Exp(exp3))
		}
	}
	cmplxData := ToComplex(data)
	ft := NewFFT2(4, 5)
	res := ft.FFT(cmplxData)
	tol := 1e-10

	for i := range res {
		if !CmplxEqualApprox(res[i], expectFFT[i], tol) {
			t.Errorf("Expected %v got%v", expectFFT[i], res[i])
		}
	}
	ift := ft.IFFT(res)

	for i := range res {
		if math.Abs(real(ift[i])-20*data[i]) > tol || math.Abs(imag(ift[i])) > tol {
			t.Errorf("Expected\n%v\ngot\n%v\n", data, cmplxData)
			break
		}
	}
}

func TestFFT2Freq(t *testing.T) {
	data := []float64{1.0, -1.0, 2.0, 3.0,
		2.0, 3.0, -5.0, 1.0,
		4.0, 5.0, 6.0, 7.0,
		2.0, 1.0, 3.0, 4.0}
	cData := ToComplex(data)
	ft := NewFFT2(4, 4)
	coeff := ft.FFT(cData)
	tol := 1e-10

	// Check complex conjugate property
	num := 0
	for i := 1; i < len(coeff); i++ {
		f1 := ft.Freq(i)
		f1[0] = f1[0]
		f1[1] = f1[1]
		for j := 1; j < len(coeff); j++ {
			f2 := ft.Freq(j)
			f2[0] = -f2[0]
			f2[1] = -f2[1]
			if floats.EqualApprox(f1, f2, tol) {
				num++
				if math.Abs(cmplx.Abs(coeff[i])-cmplx.Abs(coeff[j])) > tol {
					t.Errorf("Conjugation is not valid. (%v): (%v, %v)\n", f1, coeff[i], coeff[j])
				}
			}
		}
	}

	if num != 8 {
		t.Errorf("Unexpected number of conjugate pairs. Expected %d got %d\n", 8, num)
	}

}
