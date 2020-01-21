package main

import (
	"fmt"
	"time"

	"github.com/davidkleiven/gosfft/sfft"
)

func main() {
	nr := 128
	nc := 128
	nd := 128
	numLoops := 10
	data := make([]complex128, nr*nc*nd)
	dataCpy := make([]complex128, nr*nc*nd)
	for i := range data {
		data[i] = complex(float64(i), 0.0)
	}

	// Time serial version
	ft := sfft.NewFFT3(nr, nc, nd)
	start := time.Now()
	for i := 0; i < numLoops; i++ {
		copy(dataCpy, data)
		ft.FFT(dataCpy)
	}
	ellapsed := time.Since(start)
	fmt.Printf("Time FFT serial: %s\n", ellapsed)

	for _, nWork := range []int{2, 4, 8} {
		ftPar := sfft.NewFFT3Par(nr, nc, nd, nWork)
		start = time.Now()
		for i := 0; i < numLoops; i++ {
			copy(dataCpy, data)
			ftPar.FFT(dataCpy)
		}
		ellapsed = time.Since(start)
		fmt.Printf("Time %d workers: %s\n", nWork, ellapsed)
	}
}
