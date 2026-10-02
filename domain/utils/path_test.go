package utils

import "testing"

func TestNormalizePath(t *testing.T) {
	tests := map[string]string{
		`  "C:\Games\Example"  `: `C:\Games\Example`,
		`  D:\Game Files\Demo  `: `D:\Game Files\Demo`,
		`E:\Games\Plain`:         `E:\Games\Plain`,
	}

	for input, want := range tests {
		if got := NormalizePath(input); got != want {
			t.Errorf("NormalizePath(%q) = %q, want %q", input, got, want)
		}
	}
}
