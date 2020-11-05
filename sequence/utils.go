package sequence

import "math"

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}
func maxFloat3(f1, f2, f3 float64) float64 {
	return math.Max(f1, math.Max(f2, f3))
}

func maxFloat3Dir(up, left, upLeft float64) (float64, allgAction) {
	m := up
	act := actionUp
	if left >= m {
		m = left
		act = actionLeft
	}
	if upLeft >= m {
		m = upLeft
		act = actionUpLeft
	}

	return m, act
}

func maxFloat3DirAlt(
	f1 float64, a1 allgAction,
	f2 float64, a2 allgAction,
	f3 float64, a3 allgAction,
) (float64, allgAction) {
	f := f1
	a := a1
	if f2 >= f {
		f = f2
		a = a2
	}
	if f3 >= f {
		f = f3
		a = a3
	}

	return f, a
}
