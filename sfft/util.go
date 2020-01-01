package sfft

import "math"

// prod returns the product of all the elements in the passed slice
func prod(x []int) int {
	res := 1
	for _, v := range x {
		res *= v
	}
	return res
}

// ToComplex returns a complex representation of a real slice. The real part of the returned array
// is equal to the passed array, and the imaginary part is set to zero
func ToComplex(data []float64) []complex128 {
	res := make([]complex128, len(data))
	for i := range data {
		res[i] = complex(data[i], 0.0)
	}
	return res
}

// extract returns every n-th element of a sequence of complex number. The step size is given
// step
func extractComplex(data []complex128, start int, step int) []complex128 {
	if step == 0 {
		panic("extractComplex: Stepsize has to be larger than zero")
	}
	n := len(data) / step
	res := make([]complex128, n)
	for i := 0; i < n; i++ {
		res[i] = data[start+i*step]
	}
	return res
}

// insertComplex inserts elements into dst. It does the opposite of
// extractComplex
func insertComplex(dst []complex128, data []complex128, start int, step int) {
	for i := 0; i < len(data); i++ {
		dst[start+i*step] = data[i]
	}
}

// CmplxEqualApprox checks if a and b is equal. a and b is considered equal if
// their real and imaginary parts are equal within tol
func CmplxEqualApprox(a complex128, b complex128, tol float64) bool {
	return math.Abs(real(a)-real(b)) < tol && math.Abs(imag(a)-imag(b)) < tol
}

// Mutable2 is a generic interface for a type where elements can be accessed via the
// At method and altered via the Set method
type Mutable2 interface {
	// At returns the (i, j) element
	At(i, j int) complex128

	// Set sets a new value for the element at (i, j)
	Set(i, j int, v complex128)

	// Dims return the size of the data
	Dims() (int, int)
}

// Center2 puts the origin a(i.e. zero frequency) t the center of the image of a 2D transform
func Center2(data Mutable2) {
	nr, nc := data.Dims()
	for i := 0; i < nr; i++ {
		for j := 0; j < nc/2; j++ {
			v1 := data.At(i, j)
			v2 := data.At(i, nc/2+j)
			data.Set(i, j, v2)
			data.Set(i, nc/2+j, v1)
		}
	}

	for i := 0; i < nr/2; i++ {
		for j := 0; j < nc; j++ {
			v1 := data.At(i, j)
			v2 := data.At(nr/2+i, j)
			data.Set(i, j, v2)
			data.Set(nr/2+i, j, v1)
		}
	}
}

// Center3 brings the center (i.e. zero frequency) of a 3D transform to the center of the 3D array
func Center3(data *Mat3) {
	nr, nc, nd := data.Dims()
	for i := 0; i < nr; i++ {
		for j := 0; j < nc; j++ {
			for k := 0; k < nd/2; k++ {
				v1 := data.At(i, j, k)
				v2 := data.At(i, j, nd/2+k)
				data.Set(i, j, k, v2)
				data.Set(i, j, nd/2+k, v1)
			}
		}
	}

	for i := 0; i < nr; i++ {
		for j := 0; j < nc/2; j++ {
			for k := 0; k < nd; k++ {
				v1 := data.At(i, j, k)
				v2 := data.At(i, j+nc/2, k)
				data.Set(i, j, k, v2)
				data.Set(i, j+nc/2, k, v1)
			}
		}
	}

	for i := 0; i < nr/2; i++ {
		for j := 0; j < nc; j++ {
			for k := 0; k < nd; k++ {
				v1 := data.At(i, j, k)
				v2 := data.At(i+nr/2, j, k)
				data.Set(i, j, k, v2)
				data.Set(i+nr/2, j, k, v1)
			}
		}
	}
}
