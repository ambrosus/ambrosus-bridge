package helpers

type Pair[T, U any] struct {
	First  T
	Second U
}

func minmax(x, y int) (min int, max int) {
	if x < y {
		return x, y
	}
	return y, x
}

// ZipDiff is like Python's `zip_longest`
func ZipDiff[T, U any](ts []T, us []U) []Pair[T, U] {
	// identify the minimum and maximum lengths
	lmin, lmax := minmax(len(ts), len(us))

	pairs := make([]Pair[T, U], lmax)
	// build tuples up to the minimum length
	for i := 0; i < lmin; i++ {
		pairs[i] = Pair[T, U]{ts[i], us[i]}
	}
	if lmin == lmax {
		return pairs
	}

	// build tuples with one zero value for [lmin,lmax) range
	for i := lmin; i < lmax; i++ {
		p := Pair[T, U]{}
		if len(ts) == lmax {
			p.First = ts[i]
		} else {
			p.Second = us[i]
		}
		pairs[i] = p
	}
	return pairs
}
