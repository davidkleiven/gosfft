package sfft

import (
	"math"
	"testing"

	"gonum.org/v1/gonum/floats"
)

func TestCMat3(t *testing.T) {

	// Test initialize with nil
	mat1 := NewCMat3(1, 2, 3, nil)
	mat1.Set(0, 0, 0, complex(1.0, 0))
	mat1.Set(0, 1, 2, complex(2.0, 0))
	expectData := make([]complex128, 6)
	expectData[0] = complex(1.0, 0)
	expectData[5] = complex(2.0, 0)
	for i := range expectData {
		if !CmplxEqualApprox(expectData[i], mat1.Data[i], 1e-10) {
			t.Errorf("Underlying flattened array is not correct")
		}
	}

	nr, nc, nd := mat1.Dims()

	if nr != 1 || nc != 2 || nd != 3 {
		t.Errorf("Unexpected dimensions")
	}

	// Test initialize with an array
	mat2 := NewCMat3(1, 2, 3, expectData)
	if !CmplxEqualApprox(mat2.At(0, 0, 0), complex(1.0, 0.0), 1e-10) {
		t.Errorf("Element at (0, 0, 0) does not match")
	}

	if !CmplxEqualApprox(mat2.At(0, 1, 2), complex(2.0, 0.0), 1e-10) {
		t.Errorf("Element at (0, 1, 2) does not match")
	}

	// Test initialize with an array of wrong size
	init := make([]complex128, 40)

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("NewCMat did not panic")
		}
	}()

	NewCMat3(3, 4, 5, init)
}

func TestMat3(t *testing.T) {
	// Test initializing with nil
	mat1 := NewMat3(3, 2, 1, nil)
	mat1.Set(0, 0, 0, 1.0)
	mat1.Set(2, 1, 0, 2.0)
	expectArray := make([]float64, 6)
	expectArray[0] = 1.0
	expectArray[5] = 2.0
	if !floats.EqualApprox(mat1.Data, expectArray, 1e-10) {
		t.Errorf("Unexpected array")
	}

	// Test initialize with a sequence
	mat2 := NewMat3(3, 2, 1, expectArray)
	if math.Abs(mat2.At(0, 0, 0)-1.0) > 1e-10 {
		t.Errorf("Unexpected value at (0, 0, 0)")
	}
	if math.Abs(mat2.At(2, 1, 0)-2.0) > 1e-10 {
		t.Errorf("Unexpected valud at (2, 1, 0)")
	}

	// Initialize with wrong size
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Did not panic")
		}
	}()
	NewMat3(4, 4, 4, expectArray)
}

func TestAsUint8(t *testing.T) {
	data := []float64{1.0, 0.0, 0.0,
		-1.0, 0.0, 0.0,

		2.0, 0.0, 0.0,
		-2.0, 0.0, 0.0}
	mat3 := NewMat3(2, 3, 2, data)
	uint8Rep := mat3.AsUint8()

	expectUint8 := []uint8{
		191, 127, 127,
		63, 127, 127,
		255, 127, 127,
		0, 127, 127,
	}
	for i := range expectUint8 {
		if expectUint8[i] != uint8Rep[i] {
			t.Errorf("Unexpecte uint8 rep. Expected %d got %d\n", expectUint8[i], uint8Rep[i])
		}
	}
}

func TestAsComplex(t *testing.T) {
	data := []float64{1.0, 2.0, 3.0, 4.0, 5.0, 6.0, 7.0, 8.0}
	mat3 := NewMat3(2, 2, 2, data)
	cmat3 := mat3.AsComplex()
	for i := range mat3.Data {
		if math.Abs(real(cmat3.Data[i])-data[i]) > 1e-10 || math.Abs(imag(cmat3.Data[i])) > 1e-10 {
			t.Errorf("Real part of complex matrix does not match the data it was created from")
		}
	}
}
