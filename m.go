package m

import "math/cmplx"

func M(c complex128, l int) int {
	var i int
	var n complex128 = 0+0i
	for i = 0; i < l && cmplx.Abs(n) < 2; i++ {
		n = n * n + c
	}
	return i
}