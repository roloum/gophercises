package main

import "testing"

func TestNormalize(t *testing.T) {

	cases := []struct {
		input    string
		expected string
	}{
		{"1234567890", "1234567890"},
		{"123 456 7891", "1234567891"},
		{"(123) 456 7892", "1234567892"},
		{"(123) 456-7893", "1234567893"},
		{"123-456-7894", "1234567894"},
		{"1234567892", "1234567892"},
		{"(123)456-7892", "1234567892"},
	}

	for _, c := range cases {
		t.Run(c.input, func(t *testing.T) {
			output := normalize(c.input)
			if output != c.expected {
				t.Errorf("Received: %s, expected: %s", output, c.expected)
			}
		})
	}
	_ = normalize("")
}
