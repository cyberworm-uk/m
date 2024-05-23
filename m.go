package m

import "math/cmplx"

func z(n, c complex128) complex128 {
	return n*n + c
}

func iter(n, c complex128, l int) int {
	var i int
	for i = 0; i < l; i++ {
		n = z(n, c)
		if cmplx.Abs(n) > 2 {
			return i
		}
	}
	return i
}

const n0 = complex(0, 0)

func M(c complex128, l int) int {
	return iter(n0, c, l)
}
