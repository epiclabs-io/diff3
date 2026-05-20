package diff3

import (
	"testing"
)

func Test_middleSnake(t *testing.T) {
	for _, tc := range []struct {
		a, b          string
		x, y, u, v, d int
	}{
		{"A", "A", 0, 0, 1, 1, 0},
		{"A", "B", 1, 0, 1, 0, 2},
		{"AB", "AB", 0, 0, 2, 2, 0},
		{"AB", "AC", 2, 1, 2, 1, 2},
		{"AB", "BC", 1, 0, 2, 1, 2},
		{"AB", "CD", 2, 0, 2, 0, 4},
		{"AB", "A", 2, 1, 2, 1, 1},
		{"AB", "B", 1, 0, 2, 1, 1},
		{"AB", "AAB", 1, 2, 2, 3, 1},
		{"AB", "ABB", 2, 3, 2, 3, 1},
		{"", "", 0, 0, 0, 0, 0},
		{"A", "", 1, 0, 1, 0, 1},
		{"", "A", 0, 1, 0, 1, 1},
		{"AB", "", 1, 0, 1, 0, 2},
		{"", "AB", 0, 1, 0, 1, 2},
		{"AAAAAAAABDDDDDDDDDDDDDDD", "AAAAAAAACDDDDDDDDDDDDDDD", 9, 8, 9, 8, 2},
		{"AAAAAAAACCCCCCCCCCCCCCCC", "AAABBBBCCCCCCCCCCCCCCCCC", 5, 0, 8, 3, 10},
	} {
		t.Run(tc.a+"-"+tc.b, func(t *testing.T) {
			x, y, u, v, d := middleSnake([]rune(tc.a), []rune(tc.b))
			if x != tc.x || y != tc.y || u != tc.u || v != tc.v || d != tc.d {
				t.Errorf("expected %d %d %d %d %d, got %d %d %d %d %d", tc.x, tc.y, tc.u, tc.v, tc.d, x, y, u, v, d)
			}
		})
	}
}
