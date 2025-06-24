package main

import "testing"

func TestWordWrap(t *testing.T) {
	cases := []struct {
		desc     string
		input    string
		width    int
		expected string
	}{
		{
			desc:     "wrap short text",
			input:    "Hello world",
			width:    20,
			expected: "Hello world",
		},
		{
			desc:     "wrap long text",
			input:    "This is a long text that should be wrapped correctly.",
			width:    20,
			expected: "This is a long text\nthat should be\nwrapped correctly.",
		},
		{
			desc:     "wrap with punctuation",
			input:    "Hello, world! This is a test.",
			width:    20,
			expected: "Hello, world! This\nis a test.",
		},
	}
	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			got := wordWrap(tc.input, tc.width)
			if got != tc.expected {
				t.Errorf("expected %q, got %q", tc.expected, got)
			}
		})
	}
}
