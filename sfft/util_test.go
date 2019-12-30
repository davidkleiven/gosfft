package sfft

import (
	"testing"

	"gonum.org/v1/gonum/mat"
)

func TestProd(t *testing.T) {
	for i, test := range []struct {
		slice  []int
		expect int
	}{
		{
			slice:  []int{2},
			expect: 2,
		},
		{
			slice:  []int{2, 4, 5},
			expect: 40,
		},
		{
			slice:  []int{2, 5, 6, 0},
			expect: 0,
		},
	} {
		product := prod(test.slice)

		if product != test.expect {
			t.Errorf("Test #%d: Expected %d got %d", i, test.expect, product)
		}
	}
}

func TestExtractComplex(t *testing.T) {
	for i, test := range []struct {
		seq    []complex128
		expect []complex128
		step   int
	}{
		{
			seq:    []complex128{complex(1, 0), complex(2, 0), complex(3, 0), complex(4, 0)},
			expect: []complex128{complex(1, 0), complex(3, 0)},
			step:   2,
		},
		{
			seq:    []complex128{complex(1, 0), complex(2, 0), complex(3, 0), complex(4, 0)},
			expect: []complex128{complex(1, 0)},
			step:   3,
		},
		{
			seq:    []complex128{complex(1, 0), complex(2, 0), complex(3, 0)},
			expect: []complex128{complex(1, 0), complex(3, 0)},
			step:   2,
		},
	} {
		res := extractComplex(test.seq, 0, test.step)
		tol := 1e-10
		for j := range res {
			if !CmplxEqualApprox(res[j], test.expect[j], tol) {
				t.Errorf("Test #%d: Expected %v got %v\n", i, test.expect, res)
				break
			}
		}
	}
}

func TestCenter2(t *testing.T) {
	matrix := mat.NewCDense(4, 4, nil)
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			matrix.Set(i, j, complex(float64(i*4+j), 0))
		}
	}

	Center2(matrix)

	expect := mat.NewCDense(4, 4, []complex128{
		complex(10, 0), complex(11, 0), complex(8, 0), complex(9, 0),
		complex(14, 0), complex(15, 0), complex(12, 0), complex(13, 0),
		complex(2, 0), complex(3, 0), complex(0, 0), complex(1, 0),
		complex(6, 0), complex(7, 0), complex(4, 0), complex(5, 0),
	})

	if !mat.CEqualApprox(matrix, expect, 1e-10) {
		rMat := realPart(matrix)
		rExp := realPart(expect)
		t.Errorf("Unexpected shift. Expected\n%v\nGot\n%v\n", mat.Formatted(rExp), mat.Formatted(rMat))
	}
}

func realPart(m *mat.CDense) *mat.Dense {
	nr, nc := m.Dims()
	rMat := mat.NewDense(nr, nc, nil)
	for i := 0; i < nr; i++ {
		for j := 0; j < nc; j++ {
			rMat.Set(i, j, real(m.At(i, j)))
		}
	}
	return rMat
}
