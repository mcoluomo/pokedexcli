package main

import (
	"testing"
)

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input    string
		expected []string
	}{
		{
			input:    "  Hello  World  ",
			expected: []string{"hello", "world"},
		},
		{
			input:    "  This  is a   Test  case  ",
			expected: []string{"this", "is", "a", "test", "case"},
		},
		{
			input:    "    ",
			expected: []string{""},
		},
	}

	for _, c := range cases {
		actual := cleanInput(c.input)
		if len(actual) != len(c.expected) {
			t.Errorf(`
			Expecting length:  %d
			Actual length:     %d
			Fail
			`, len(c.expected), len(actual))
		}
		// Check the length of the actual slice against the expected slice
		// if they don't match, use t.Errorf to print an error message
		// and fail the test
		for i := range actual {
			word := actual[i]
			expectedWord := c.expected[i]
			if word != expectedWord {
				t.Errorf(`
				Expecting:  %s
				Actual:     %s
				Fail
			`, expectedWord, word)
			}
			// Check each word in the slice
			// if they don't match, use t.Errorf to print an error message
			// and fail the test
		}
	}
}
