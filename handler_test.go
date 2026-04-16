package main

import "testing"

func TestCleanChirp(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"hello world", "hello world"},
		{"kerfuffle is silly", "**** is silly"},
		{"KERFUFFLE is loud", "**** is loud"},
		{"sharbert!", "sharbert!"}, // punctuation should NOT match
		{"fornax is here", "**** is here"},
	}

	for _, tc := range tests {
		output := cleanChirp(tc.input)
		if output != tc.expected {
			t.Errorf("cleanChirp(%q) = %q; expected %q", tc.input, output, tc.expected)
		}
	}
}
