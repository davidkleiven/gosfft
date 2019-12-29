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

// ToComplex returns a complex representation of a real slice
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

// CmplxEqualApprox checks if a and b is equal
func CmplxEqualApprox(a complex128, b complex128, tol float64) bool {
	return math.Abs(real(a)-real(b)) < tol && math.Abs(imag(a)-imag(b)) < tol
}
