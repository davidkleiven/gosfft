package main

import (
	"math/cmplx"

	"github.com/davidkleiven/gosfft/sfft"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

func main() {
	data := make([]float64, 256)
	for i := 50; i < 100; i++ {
		data[i] = 1.0
	}

	fft := sfft.NewFFT1(len(data))
	ftData := fft.FFT(data)

	ftAmp := make([]float64, len(data)/2+1)
	for i := range ftData {
		ftAmp[i] = cmplx.Abs(ftData[i])
	}

	// Plot the signal and the result
	pltSignal := plot.New()
	pltSignal.X.Label.Text = "Time (s)"
	pltSignal.Y.Label.Text = "Amplitude"
	pts := make(plotter.XYs, len(data))
	for i := range data {
		pts[i] = plotter.XY{X: float64(i), Y: data[i]}
	}

	l, err := plotter.NewLine(pts)
	if err != nil {
		panic(err)
	}
	pltSignal.Add(l)
	// Plot the frequency spectrum
	pltFreq := plot.New()
	if err != nil {
		panic(err)
	}

	pltFreq.X.Label.Text = "Frequency (Hz)"
	pltFreq.Y.Label.Text = "Amplitude"
	freqData := make(plotter.XYs, len(ftAmp))
	for i := range ftAmp {
		freqData[i] = plotter.XY{X: fft.Freq(i), Y: ftAmp[i]}
	}

	l, err = plotter.NewLine(freqData)
	if err != nil {
		panic(err)
	}
	pltFreq.Add(l)

	// Save result
	if err := pltSignal.Save(4*vg.Inch, 4*vg.Inch, "signal1D.png"); err != nil {
		panic(err)
	}

	if err := pltFreq.Save(4*vg.Inch, 4*vg.Inch, "fft1D.png"); err != nil {
		panic(err)
	}
}
