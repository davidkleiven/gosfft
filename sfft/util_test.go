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
