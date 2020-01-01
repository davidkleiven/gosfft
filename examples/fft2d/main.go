package main

import (
	"math/cmplx"

	"github.com/davidkleiven/gosfft/sfft"
	"gonum.org/v1/gonum/mat"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/palette"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

// GridMatrix is a type used to be able to plot a 2D matrix with Gonum's heatmap
type GridMatrix struct {
	mat mat.Matrix
}

// Dims return the size of the underlying matrix
func (g GridMatrix) Dims() (int, int) { r, c := g.mat.Dims(); return r, c }

// X returns the coordinate value corresponding to column c
func (g GridMatrix) X(c int) float64 { return float64(c) }

// Y returns the coordinate value corresponiding to the row r
func (g GridMatrix) Y(r int) float64 { return float64(r) }

// Z returns the function value corresponding to the (r, c) element in the underlying matrix
func (g GridMatrix) Z(c, r int) float64 { return g.mat.At(r, c) }

func main() {
	cmat := mat.NewCDense(64, 64, nil)
	for i := 0; i < 64; i++ {
		for j := 0; j < 64; j++ {
			if i > 25 && i < 35 && j > 25 && j < 35 {
				cmat.Set(i, j, complex(1.0, 0.0))
			}
		}
	}

	// Store the real-part which will be used for plotting later
	img := mat.NewDense(64, 64, nil)
	for i := 0; i < 64; i++ {
		for j := 0; j < 64; j++ {
			img.Set(i, j, real(cmat.At(i, j)))
		}
	}

	ft := sfft.NewFFT2(64, 64)
	ftData := ft.FFT(cmat.RawCMatrix().Data)
	ftMat := mat.NewCDense(64, 64, ftData)
	sfft.Center2(ftMat)

	// Extract amplitude
	amp := mat.NewDense(64, 64, nil)
	for i := 0; i < 64; i++ {
		for j := 0; j < 64; j++ {
			amp.Set(i, j, cmplx.Abs(ftMat.At(i, j)))
		}
	}

	// Plot the results

	gridImg := GridMatrix{img}
	pltImg, err := plot.New()
	if err != nil {
		panic(err)
	}

	im := plotter.NewHeatMap(gridImg, palette.Heat(10, 1))
	pltImg.Add(im)

	pltFT, err := plot.New()
	if err != nil {
		panic(err)
	}

	gridFT := GridMatrix{mat: amp}
	h := plotter.NewHeatMap(gridFT, palette.Heat(10, 1))
	pltFT.Add(h)

	// Save result
	if err := pltImg.Save(4*vg.Inch, 4*vg.Inch, "img.png"); err != nil {
		panic(err)
	}
	if err := pltFT.Save(4*vg.Inch, 4*vg.Inch, "fft2D.png"); err != nil {
		panic(err)
	}

}
