package sfft

import (
	"gonum.org/v1/gonum/floats"
	"math"
	"testing"
)

func TestFFTParForwardBackward(t *testing.T) {
	for i, test := range []struct {
		data     []float64
		nr       int
		nc       int
		nWorkers int
	}{
		{
			data:     []float64{1.0, 2.0, 3.0, 4.0},
			nr:       2,
			nc:       2,
			nWorkers: 2,
		},
		{
			data:     []float64{1.0, 2.0, 3.0, 4.0, 5.0, 6.0, 7.0, 10.0},
			nr:       2,
			nc:       4,
			nWorkers: 2,
		},
		{
			data:     []float64{1.0, 2.0, 3.0, 4.0, 5.0, 6.0, 7.0, 10.0},
			nr:       4,
			nc:       2,
			nWorkers: 2,
		},
	} {
		cdata := make([]complex128, len(test.data))
		for i := range test.data {
			cdata[i] = complex(test.data[i], 0.0)
		}
		ft := NewFFT2Par(test.nr, test.nc, test.nWorkers)
		ft.FFT(cdata)
		ft.IFFT(cdata)
		tol := 1e-10
		for j := range cdata {
			re := real(cdata[j]) / float64(len(cdata))
			im := imag(cdata[j]) / float64(len(cdata))

			if math.Abs(re-test.data[j]) > tol || math.Abs(im) > tol {
				t.Errorf("Test #%d: Inconsistent forward/backward result. Got %f+%f i expected%f+0i", i, re, im, test.data[j])
			}
		}
	}
}

func TestConsistentWithFFT2(t *testing.T) {
	for i, test := range []struct {
		data     []float64
		nr       int
		nc       int
		nWorkers int
	}{
		{
			data:     []float64{1.0, 2.0, 3.0, 4.0},
			nr:       2,
			nc:       2,
			nWorkers: 2,
		},
		{
			data:     []float64{1.0, 2.0, 3.0, 4.0, 5.0, 6.0, 7.0, 10.0},
			nr:       2,
			nc:       4,
			nWorkers: 2,
		},
		{
			data:     []float64{1.0, 2.0, 3.0, 4.0, 5.0, 6.0, 7.0, 10.0},
			nr:       4,
			nc:       2,
			nWorkers: 2,
		},
	} {
		cdata := make([]complex128, len(test.data))
		for i := range test.data {
			cdata[i] = complex(test.data[i], 0.0)
		}
		cdataCpy := make([]complex128, len(cdata))
		copy(cdataCpy, cdata)

		ftPar := NewFFT2Par(test.nr, test.nc, test.nWorkers)
		ft := NewFFT2(test.nr, test.nc)
		ftPar.FFT(cdata)
		ft.FFT(cdataCpy)

		tol := 1e-10
		for j := range cdata {
			diff := cdata[j] - cdataCpy[j]
			if math.Abs(real(diff)) > tol || math.Abs(imag(diff)) > tol {
				t.Errorf("Test #%d: Expected %v got %v\n", i, cdata[j], cdataCpy[j])
			}
		}
	}
}

func TestFFT2ParFreq(t *testing.T) {
	ft := NewFFT2Par(5, 5, 5)
	ftOrig := NewFFT2(5, 5)
	tol := 1e-10
	for i := 0; i < 25; i++ {
		f1 := ft.Freq(i)
		f2 := ftOrig.Freq(i)
		if !floats.EqualApprox(f1, f2, tol) {
			t.Errorf("Expected %v got %v\n", f2, f1)
		}
	}
}

func TestFFT3ParForwardBackward(t *testing.T) {
	for i, test := range []struct {
		data     []float64
		nr       int
		nc       int
		nd       int
		nWorkers int
	}{
		{
			data:     []float64{1.0, 2.0, 3.0, 4.0, 5.0, 6.0, 7.0, 8.0},
			nr:       2,
			nc:       2,
			nd:       2,
			nWorkers: 2,
		},
		{
			data:     []float64{1.0, 2.0, 3.0, 4.0, 5.0, 6.0, 7.0, 8.0},
			nr:       2,
			nc:       4,
			nd:       1,
			nWorkers: 2,
		},
		{
			data:     []float64{1.0, 2.0, 3.0, 4.0, 5.0, 6.0, 7.0, 8.0},
			nr:       4,
			nc:       2,
			nd:       1,
			nWorkers: 2,
		},
	} {
		ft := NewFFT3Par(test.nr, test.nc, test.nd, test.nWorkers)
		cdata := make([]complex128, len(test.data))
		for j := range test.data {
			cdata[j] = complex(test.data[j], 0.0)
		}
		ft.FFT(cdata)
		ft.IFFT(cdata)
		tol := 1e-10
		for j := range test.data {
			re := real(cdata[j]) / float64(len(cdata))
			im := imag(cdata[j]) / float64(len(cdata))
			if math.Abs(re-test.data[j]) > tol || math.Abs(im) > tol {
				t.Errorf("Test #%d: Got (%f,%f) expected %f+0i\n", i, re, im, test.data[j])
			}
		}
	}
}

func TestConsistentWithFFT3(t *testing.T) {
	for i, test := range []struct {
		data     []float64
		nr       int
		nc       int
		nd       int
		nWorkers int
	}{
		{
			data:     []float64{1.0, 2.0, 3.0, 4.0, 5.0, 6.0, 7.0, 8.0},
			nr:       2,
			nc:       2,
			nd:       2,
			nWorkers: 2,
		},
		{
			data:     []float64{1.0, 2.0, 3.0, 4.0, 5.0, 6.0, 7.0, 8.0},
			nr:       2,
			nc:       4,
			nd:       1,
			nWorkers: 2,
		},
		{
			data:     []float64{1.0, 2.0, 3.0, 4.0, 5.0, 6.0, 7.0, 8.0},
			nr:       4,
			nc:       2,
			nd:       1,
			nWorkers: 2,
		},
	} {
		ft := NewFFT3Par(test.nr, test.nc, test.nd, test.nWorkers)
		ftReg := NewFFT3(test.nr, test.nc, test.nd)
		cdata := make([]complex128, len(test.data))
		cdataCpy := make([]complex128, len(cdata))
		for j := range test.data {
			cdata[j] = complex(test.data[j], 0.0)
		}
		copy(cdataCpy, cdata)
		ft.FFT(cdata)
		ftReg.FFT(cdataCpy)
		tol := 1e-10
		for j := range test.data {
			if math.Abs(real(cdata[j]-cdataCpy[j])) > tol || math.Abs(imag(cdata[j]-cdataCpy[j])) > tol {
				t.Errorf("Test #%d: Got %v expected %v\n", i, cdata[j], cdataCpy[j])
			}
		}
	}
}

func TestFFT3ParFreq(t *testing.T) {
	ftPar := NewFFT3Par(12, 12, 12, 4)
	ftReg := NewFFT3(12, 12, 12)
	for i := 0; i < 12*12*12; i++ {
		f1 := ftPar.Freq(i)
		f2 := ftReg.Freq(i)
		if !floats.EqualApprox(f1, f2, 1e-10) {
			t.Errorf("Expected %v got %v\n", f2, f1)
		}
	}
}
