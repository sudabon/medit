package render

import "testing"

func TestContentWidth(t *testing.T) {
	tests := []struct {
		term int
		want int
	}{
		{20, 20},
		{44, 40}, // 44-4=40
		{100, PreferredContentWidth},
		{200, PreferredContentWidth},
		{90, 86}, // 90-4=86 < preferred
	}
	for _, tt := range tests {
		if got := ContentWidth(tt.term); got != tt.want {
			t.Errorf("ContentWidth(%d) = %d, want %d", tt.term, got, tt.want)
		}
	}
}

func TestLeftPad(t *testing.T) {
	if got := LeftPad(100, 88); got != 6 {
		t.Fatalf("LeftPad = %d, want 6", got)
	}
	if got := LeftPad(80, 88); got != 0 {
		t.Fatalf("LeftPad narrow = %d", got)
	}
}
