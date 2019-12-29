package sfft

// prod returns the product of all the elements in the passed slice
func prod(x []int) int {
	res := 1
	for _, v := range x {
		res *= v
	}
	return res
}
