package sfft

import "testing"

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
