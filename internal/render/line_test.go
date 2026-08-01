package render

import "testing"

func TestSourceLineCount(t *testing.T) {
	tests := []struct {
		in   string
		want int
	}{
		{"", 0},
		{"a", 1},
		{"a\n", 2},
		{"a\nb\nc", 3},
	}
	for _, tt := range tests {
		if got := SourceLineCount(tt.in); got != tt.want {
			t.Errorf("SourceLineCount(%q) = %d, want %d", tt.in, got, tt.want)
		}
	}
}

func TestApproximateSourceLine(t *testing.T) {
	tests := []struct {
		name   string
		pct    float64
		lines  int
		want   int
	}{
		{"empty", 0.5, 0, 1},
		{"top", 0, 100, 1},
		{"bottom", 1, 100, 100},
		{"mid", 0.5, 100, 50},
		{"clamp low", -1, 10, 1},
		{"clamp high", 2, 10, 10},
		{"single", 0, 1, 1},
		{"single end", 1, 1, 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ApproximateSourceLine(tt.pct, tt.lines)
			if got != tt.want {
				t.Fatalf("got %d, want %d", got, tt.want)
			}
		})
	}
}
