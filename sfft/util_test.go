package sfft

import "testing"

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
