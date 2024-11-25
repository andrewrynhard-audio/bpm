package timing

import (
	"testing"
)

// TestRoundHalfUp tests the roundHalfUp function for correct rounding behavior.
func TestRoundHalfUp(t *testing.T) {
	tests := []struct {
		input    float64
		expected float64
	}{
		{15.469, 16},   // Cascading rounds up
		{15.496, 16},   // Properly rounds up
		{77.444, 77},   // Cascading rounds down
		{1.484, 2},     // Cascading stops early
		{-15.496, -16}, // Negative cascading rounds up
		{-15.469, -16}, // Negative rounds up
		{1.500, 2},     // Exact boundary rounds up
		{-1.500, -2},   // Exact boundary rounds down
	}

	for _, test := range tests {
		output := roundHumanCascading(test.input)
		if output != test.expected {
			t.Errorf("roundHalfUp(%f) = %f; want %f", test.input, output, test.expected)
		}
	}
}

// TestFormatWithUnit tests the formatWithUnit function for correct formatting and unit behavior.
func TestFormatWithUnit(t *testing.T) {
	tests := []struct {
		value          float64
		roundToWhole   bool
		expectedOutput string
	}{
		// Tests for milliseconds
		{15.496, true, "16 ms"},
		{15.496, false, "15.496 ms"},
		{15.000, true, "15 ms"},
		{15.000, false, "15.000 ms"},
		// Tests for microseconds
		{0.996, true, "1 ms"},
		{0.996, false, "996.000 us"},
		{0.0015, true, "2 us"},
		{0.0015, false, "1.500 us"},
	}

	for _, test := range tests {
		output, _ := formatWithUnit(test.value, test.roundToWhole)
		if output != test.expectedOutput {
			t.Errorf("formatWithUnit(%f, %v) = %q; want %q", test.value, test.roundToWhole, output, test.expectedOutput)
		}
	}
}
