package main

import "testing"

func Test_FormatURL(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"https://example.com", "example.com:443"},
		{"https://www.example.com/", "www.example.com:443"},
		{"https://sub.example.com/", "sub.example.com:443"},
		{"https://example.com/my-test/", "example.com:443"},            // Only get base path
		{"http://example.com", "http://example.com"},                   // No change for non-https URL
		{"http://example.com:80", "http://example.com:80"},             // No change for non-https URL with port
		{"http://sub.example.com/", "http://sub.example.com/"},         // No change for subdomain
		{"http://example.com/my-test/", "http://example.com/my-test/"}, // No change for non-https URL with path
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			result := FormatURL(test.input)
			if result != test.expected {
				t.Errorf("Input: %s, Expected: %s, Got: %s", test.input, test.expected, result)
			}
		})
	}
}
