package test

import "testing"

// TestAdd calls mymath.Add with a variety of inputs, checking
// for correct results.
func TestAdd(t *testing.T) {
	tests := []struct {
		name string
		a, b int
		want int
	}{
		{"two positives", 2, 3, 4},
		{"positive and zero", 1, 0, 1},
		{"two negatives", -1, -2, -3},
		{"positive and negative", 1, -1, 0},
	}

	for _, tt := range tests {
		testname := tt.name
		t.Run(testname, func(t *testing.T) {
			ans := Add(tt.a, tt.b)
			if ans != tt.want {
				t.Errorf("Add(%d, %d) got %d, want %d", tt.a, tt.b, ans, tt.want)
			}
		})
	}
}
