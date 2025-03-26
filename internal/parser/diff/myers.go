package diff

// source: https://github.com/MFAshby/myers
// Myer's diff algorithm in golang
// Ported from https://blog.robertelder.org/diff-algorithm/

type OpType int

const (
	OpInsert OpType = iota
	OpDelete
)

type Op struct {
	OpType OpType // Insert or delete, as above
	OldPos int    // Position in the old list of item to be inserted or deleted
	NewPos int    // Position in the _new_ list of item to be inserted
	Elem   any    // Actual value to be inserted or deleted
}

// Returns a minimal list of differences between 2 lists e and f
// requiring O(min(len(e),len(f))) space and O(min(len(e),len(f)) * D)
// worst-case execution time where D is the number of differences.
func myersDiff(e, f []any, equals func(any, any) bool) []Op {
	return diffInternal(e, f, equals, 0, 0)
}

func diffInternal(e, f []any, equals func(any, any) bool, i, j int) []Op {
	N := len(e)
	M := len(f)
	L := N + M
	Z := 2*min(N, M) + 2
	switch {
	case N > 0 && M > 0:
		w := N - M
		g := make([]int, Z)
		p := make([]int, Z)

		hMax := L/2 + L%2 + 1
		for h := range hMax {
			for r := range 2 {
				var c, d []int
				var o, m int
				if r == 0 {
					c = g
					d = p
					o = 1
					m = 1
				} else {
					c = p
					d = g
					o = 0
					m = -1
				}
				kMin := -(h - 2*max(0, h-M))
				kMax := h - 2*max(0, h-N) + 1
				for k := kMin; k < kMax; k += 2 {
					var a int
					if k == -h || k != h && c[pyMod((k-1), Z)] < c[pyMod((k+1), Z)] {
						a = c[pyMod((k+1), Z)]
					} else {
						a = c[pyMod((k-1), Z)] + 1
					}
					b := a - k
					s, t := a, b

					for a < N && b < M && equals(e[(1-o)*N+m*a+(o-1)], f[(1-o)*M+m*b+(o-1)]) {
						a, b = a+1, b+1
					}
					c[pyMod(k, Z)] = a
					z := -(k - w)
					if pyMod(L, 2) == o && z >= -(h-o) && z <= h-o && c[pyMod(k, Z)]+d[pyMod(z, Z)] >= N {
						var D, x, y, u, v int
						if o == 1 {
							D = 2*h - 1
							x = s
							y = t
							u = a
							v = b
						} else {
							D = 2 * h
							x = N - a
							y = M - b
							u = N - s
							v = M - t
						}
						switch {
						case D > 1 || (x != u && y != v):
							return append(diffInternal(e[0:x], f[0:y], equals, i, j), diffInternal(e[u:N], f[v:M], equals, i+u, j+v)...)
						case M > N:
							return diffInternal(make([]any, 0), f[N:M], equals, i+N, j+N)
						case M < N:
							return diffInternal(e[M:N], make([]any, 0), equals, i+M, j+M)
						default:
							return make([]Op, 0)
						}
					}
				}
			}
		}
	case N > 0:
		res := make([]Op, N)
		for n := range N {
			res[n] = Op{OpDelete, i + n, -1, e[n]}
		}
		return res
	default:
		res := make([]Op, M)
		for n := range M {
			res[n] = Op{OpInsert, i, j + n, f[n]}
		}
		return res
	}
	panic("Should never hit this!")
}

/**
 * The remainder op in python always matches the sign of the _denominator_
 * e.g -1%3 = 2.
 * In golang it matches the sign of the numerator.
 * See https://en.wikipedia.org/wiki/Modulo_operation#Variants_of_the_definition
 * Since we always have a positive denominator here, we can emulate the
 * pyMod x%y as (x+y) % y
 */
func pyMod(x, y int) int {
	return (x + y) % y
}

// Let us map element in same way as in

// Convenient wrapper for string lists
func myersDiffStr(e, f []string, compareSpace bool) []Op {
	e1, f1 := make([]any, len(e)), make([]any, len(f))
	for i, ee := range e {
		e1[i] = ee
	}
	for i, fe := range f {
		f1[i] = fe
	}
	return myersDiff(e1, f1, func(s1, s2 any) bool {
		return compareStrings(s1.(string), s2.(string), compareSpace)
	})
}
