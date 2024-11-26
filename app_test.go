package main

import (
	"testing"
)

// TestRoundHumanCascading tests the roundHumanCascading function for correct rounding behavior.
func TestRoundHumanCascading(t *testing.T) {
	tests := []struct {
		input    float64
		expected float64
		desc     string // Description of the test case for debugging
	}{
		// Standard cases
		{5.445, 5.5, "Rounding up near halfway point"},
		{5.444, 5.4, "Rounding down near halfway point"},
		{-5.445, -5.5, "Negative rounding up near halfway point"},
		{-5.444, -5.4, "Negative rounding down near halfway point"},

		// Edge cases
		{0.0, 0.0, "Zero input"},
		{-0.0, 0.0, "Negative zero input"},
		{1.499, 1.5, "Rounding up small positive number"},
		{1.001, 1.0, "Rounding down small positive number"},
		{-1.499, -1.5, "Rounding up small negative number"},
		{-1.001, -1.0, "Rounding down small negative number"},

		// Large values
		{12345.6789, 12345.7, "Rounding large positive number"},
		{-12345.6789, -12345.7, "Rounding large negative number"},

		// Small values (close to zero)
		{0.0005, 0.0, "Rounding down very small positive number"},
		{-0.0005, -0.0, "Rounding down very small negative number"},
		{0.00051, 0.0, "Rounding up very small positive number"},
		{-0.00051, -0.0, "Rounding up very small negative number"},

		// Extreme cases
		{1e10, 1e10, "Large number, no rounding needed"},
		{-1e10, -1e10, "Large negative number, no rounding needed"},
		{1e-10, 0.0, "Tiny positive number rounded to zero"},
		{-1e-10, 0.0, "Tiny negative number rounded to zero"},
	}

	for _, test := range tests {
		output := roundHumanCascading(test.input)
		if output != test.expected {
			t.Errorf("Test failed: %s\nroundHumanCascading(%f) = %f; want %f", test.desc, test.input, output, test.expected)
		}
	}
}

// TestFormatWithUnit tests the formatWithUnit function for correct unit formatting and rounding behavior.
func TestFormatWithUnit(t *testing.T) {
	tests := []struct {
		value        float64
		roundToWhole bool
		expected     string
		desc         string // Description of the test case for debugging
	}{
		// Rounding cases
		{1500.0, true, "1.5 s", "Rounding: Convert milliseconds to seconds"},
		{1234.567, true, "1.2 s", "Rounding: Round milliseconds to seconds"},
		{999.9, true, "999.9 ms", "Rounding: Stay in milliseconds just below 1000"},
		{0.999, true, "999.0 µs", "Rounding: Convert milliseconds to microseconds"},
		{0.036408, true, "36.4 µs", "Rounding: Microseconds with one decimal place"},
		{0.123, true, "123.0 µs", "Rounding: Larger microseconds with one decimal place"},

		// Non-rounding cases
		{1500.0, false, "1.500 s", "Non-rounding: Convert milliseconds to seconds"},
		{1234.567, false, "1.235 s", "Non-rounding: Truncate to 3 decimal places in seconds"},
		{999.9, false, "999.900 ms", "Non-rounding: Keep milliseconds below 1000"},
		{0.999, false, "999.000 μs", "Non-rounding: Convert milliseconds to microseconds"},
		{0.036408, false, "36.408 μs", "Non-rounding: Microseconds with three decimal places"},
		{0.123, false, "123.000 μs", "Non-rounding: Larger microseconds with three decimal places"},

		// Negative values
		{-1500.0, true, "N/A", "Invalid: Negative milliseconds"},
		{-0.036408, true, "N/A", "Invalid: Negative microseconds"},

		// Edge cases
		{0.0, true, "0.0 µs", "Rounding: Zero value"},
		{0.0, false, "0.000 μs", "Non-rounding: Zero value"},
		{1000.0, true, "1.0 s", "Rounding: Exact transition to seconds"},
		{1000.0, false, "1.000 s", "Non-rounding: Exact transition to seconds"},
	}

	for _, test := range tests {
		output := formatWithUnit(test.value, test.roundToWhole)
		if output != test.expected {
			t.Errorf("Test failed: %s\nformatWithUnit(%f, %t) = %s; want %s",
				test.desc, test.value, test.roundToWhole, output, test.expected)
		}
	}
}
