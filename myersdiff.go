package diff3

import (
	"iter"
)

func middleSnake[T comparable](a []T, b []T) (int, int, int, int, int) {
	n := len(a)
	m := len(b)
	maxD := (n + m + 1) / 2
	delta := n - m
	deltaIsOdd := delta%2 != 0

	vSize := 2*min(m, n) + 2
	vf := make([]int, vSize)
	vb := make([]int, vSize)

	offset := maxD + 1
	vf[(1+offset)%vSize] = 0
	vb[(1+offset)%vSize] = 0

	for d := 0; d <= maxD; d++ {
		for k := -(d - 2*max(0, d-m)); k <= d-2*max(0, d-n); k += 2 {
			var x, y int
			idxK := (k + offset) % vSize
			idxKPlus1 := (k + 1 + offset) % vSize
			idxKMinus1 := (k - 1 + offset) % vSize
			if k == -d || (k != d && vf[idxKMinus1] < vf[idxKPlus1]) {
				x = vf[idxKPlus1]
			} else {
				x = vf[idxKMinus1] + 1
			}
			y = x - k
			xStart, yStart := x, y
			for x < n && y < m && a[x] == b[y] {
				x++
				y++
			}
			vf[idxK] = x

			if deltaIsOdd {
				inverseK := -(k - delta)
				if inverseK >= -(d-1) && inverseK <= (d-1) {
					idxInverseK := (inverseK + offset) % vSize
					if vb[idxInverseK] != -1 && vf[idxK]+vb[idxInverseK] >= n {
						return xStart, yStart, x, y, 2*d - 1
					}
				}
			}
		}

		for k := -(d - 2*max(0, d-m)); k <= d-2*max(0, d-n); k += 2 {
			var x, y int
			idxK := (k + offset) % vSize
			idxKPlus1 := (k + 1 + offset) % vSize
			idxKMinus1 := (k - 1 + offset) % vSize
			if k == -d || (k != d && vb[idxKMinus1] < vb[idxKPlus1]) {
				x = vb[idxKPlus1]
			} else {
				x = vb[idxKMinus1] + 1
			}
			y = x - k
			xStartRev, yStartRev := x, y
			for x < n && y < m && a[n-1-x] == b[m-1-y] {
				x++
				y++
			}
			vb[idxK] = x

			if !deltaIsOdd {
				inverseK := -(k - delta)
				if inverseK >= -d && inverseK <= d {
					idxInverseK := (inverseK + offset) % vSize
					if vf[idxInverseK] != -1 && vb[idxK]+vf[idxInverseK] >= n {
						return n - x, m - y, n - xStartRev, m - yStartRev, 2 * d
					}
				}
			}
		}
	}

	panic("unreachable")
}

func shortestEditScript[T comparable](a, b []T, currentX, currentY int) iter.Seq[*diffIndicesResult] {
	n, m := len(a), len(b)
	if n == 0 && m == 0 {
		return func(yield func(*diffIndicesResult) bool) {}
	} else if n == 0 || m == 0 {
		return func(yield func(*diffIndicesResult) bool) {
			yield(&diffIndicesResult{
				file1: [2]int{currentX, n},
				file2: [2]int{currentY, m},
			})
		}
	}
	x, y, u, v, d := middleSnake(a, b)
	if d > 1 || (x != u && y != v) {
		it1 := make(chan iter.Seq[*diffIndicesResult], 1)
		go func() {
			it1 <- shortestEditScript[T](a[:x], b[:y], currentX, currentY)
		}()
		it2 := make(chan iter.Seq[*diffIndicesResult], 1)
		go func() {
			it2 <- shortestEditScript[T](a[u:], b[v:], currentX+u, currentY+v)
		}()
		return func(yield func(*diffIndicesResult) bool) {
			for diff := range <-it1 {
				if !yield(diff) {
					return
				}
			}
			for diff := range <-it2 {
				if !yield(diff) {
					return
				}
			}
		}
	} else if m > n {
		return shortestEditScript[T](a[n:], b[n:], currentX+n, currentY+n)
	} else if n > m {
		return shortestEditScript[T](a[m:], b[m:], currentX+m, currentY+m)
	} else {
		return func(yield func(*diffIndicesResult) bool) {}
	}
}
