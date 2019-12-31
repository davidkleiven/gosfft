package main

import (
	"encoding/binary"
	"math/cmplx"
	"os"

	"github.com/davidkleiven/gosfft/sfft"
)

func main() {
	data := sfft.NewMat3(128, 128, 128, nil)

	start := 50
	end := 70
	for i := 0; i < 128; i++ {
		for j := 0; j < 128; j++ {
			for k := 0; k < 128; k++ {
				if i > start && i < end && j > start && j < end && k > start && k < end {
					data.Set(i, j, k, 1.0)
				}
			}
		}
	}

	// Store image
	uint8Rep := data.AsUint8()
	out, _ := os.Create("image3D.bin")
	binary.Write(out, binary.LittleEndian, uint8Rep)
	out.Close()

	// Create a real representation
	cData := data.AsComplex()

	ft := sfft.NewFFT3(128, 128, 128)
	fdData := ft.FFT(cData.Data)

	for i := range fdData {
		data.Data[i] = cmplx.Abs(fdData[i])
	}
	sfft.Center3(data)

	uint8Rep = data.AsUint8()

	// Save such that it can be plotted with for instance Paraview
	outFt, _ := os.Create("ftImg3D.bin")
	binary.Write(outFt, binary.LittleEndian, uint8Rep)
	outFt.Close()
}
