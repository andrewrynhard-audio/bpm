package timing

import (
	"fmt"
	"math"

	"github.com/andrewrynhard-audio/bpm/pkg/state"
	"github.com/nsf/termbox-go"
)

type Timing struct {
	BPM             float64
	Quarter         float64
	Eighth          float64
	Sixteenth       float64
	ThirtySecond    float64
	SixtyFourth     float64
	OneTwentyEighth float64
	TwoFiftySixth   float64
	FiveTwelve      float64
	TenTwentyFour   float64
}

func (t *Timing) Render(bpm float64, state *state.SharedState) {
	t.BPM = bpm

	t.calculate()

	t.write(state.RoundOutputs)
}

func (t *Timing) Reset(state *state.SharedState) {
	t.BPM = 0
	t.Quarter = 0
	t.Eighth = 0
	t.Sixteenth = 0
	t.ThirtySecond = 0
	t.SixtyFourth = 0
	t.OneTwentyEighth = 0
	t.TwoFiftySixth = 0
	t.FiveTwelve = 0
	t.TenTwentyFour = 0

	t.write(state.RoundOutputs)
}

func (t *Timing) StateChanged(sharedState *state.SharedState) {
	// Rerender with the updated rounding state
	t.write(sharedState.RoundOutputs)
}

func (i *Timing) calculate() {
	i.Quarter = 60000.00 / i.BPM
	i.Eighth = i.Quarter * 0.5
	i.Sixteenth = i.Quarter * 0.25
	i.ThirtySecond = i.Quarter * 0.125
	i.SixtyFourth = i.Quarter * 0.0625
	i.OneTwentyEighth = i.Quarter * 0.03125
	i.TwoFiftySixth = i.Quarter * 0.015625
	i.FiveTwelve = i.Quarter * 0.0078125
	i.TenTwentyFour = i.Quarter * 0.00390625
}

func (t *Timing) write(roundOutputs bool) {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

	headers := []string{"Note", "Time", "10x", "1/10th", "1/100th", "1/1000th"}
	notes := []string{"1/4", "1/8", "1/16", "1/32", "1/64", "1/128", "1/256", "1/512", "1/1024"}
	milliseconds := []float64{
		t.Quarter, t.Eighth, t.Sixteenth,
		t.ThirtySecond, t.SixtyFourth, t.OneTwentyEighth,
		t.TwoFiftySixth, t.FiveTwelve, t.TenTwentyFour,
	}

	// Get terminal size
	termWidth, termHeight := termbox.Size()

	// Calculate column positions dynamically
	columnWidths := []int{10, 15, 15, 15, 15, 15} // Column widths for headers and data
	totalWidth := 0
	for _, w := range columnWidths {
		totalWidth += w
	}
	startX := (termWidth - totalWidth) / 2
	startY := (termHeight - len(notes) - 2) / 2

	// Display BPM at the top
	renderText(startX, startY-2, fmt.Sprintf("BPM: %d", int(t.BPM)), termbox.ColorWhite)

	// Render table headers
	currentX := startX
	for i, header := range headers {
		renderText(currentX, startY, header, termbox.ColorCyan)
		currentX += columnWidths[i]
	}

	// Render each row
	for row := 0; row < len(notes); row++ {
		ms := milliseconds[row]
		ms10x := ms * 10
		ms10 := ms / 10
		ms100 := ms / 100
		ms1000 := ms / 1000

		currentX := startX
		renderText(currentX, startY+row+1, notes[row], termbox.ColorWhite)
		currentX += columnWidths[0]

		renderText(currentX, startY+row+1, formatWithUnit(ms, roundOutputs), termbox.ColorWhite)
		currentX += columnWidths[1]

		renderText(currentX, startY+row+1, formatWithUnit(ms10x, roundOutputs), termbox.ColorWhite)
		currentX += columnWidths[2]

		renderText(currentX, startY+row+1, formatWithUnit(ms10, roundOutputs), termbox.ColorWhite)
		currentX += columnWidths[3]

		renderText(currentX, startY+row+1, formatWithUnit(ms100, roundOutputs), termbox.ColorWhite)
		currentX += columnWidths[4]

		renderText(currentX, startY+row+1, formatWithUnit(ms1000, roundOutputs), termbox.ColorWhite)
	}

	// Display help message at the bottom
	helpMessage := "Press 'R' to reset, 'ESC' or 'Q' to quit, 'F1' to toggle whole numbers/decimals."
	renderText((termWidth-len(helpMessage))/2, termHeight-1, helpMessage, termbox.ColorGreen)

	termbox.Flush()
}

func renderText(x, y int, text string, color termbox.Attribute) {
	for i, ch := range text {
		termbox.SetCell(x+i, y, ch, color, termbox.ColorDefault)
	}
}

// formatWithUnit formats a number dynamically based on its magnitude and rounding preference.
// - If `roundToWhole` is true, it rounds before determining the unit.
// - For values ≥ 1, it uses milliseconds (ms).
// - For values < 1, it converts to microseconds (μs).
func formatWithUnit(value float64, roundToWhole bool) string {
	switch {
	case roundToWhole:
		// Apply rounding first
		rounded := roundHumanCascading(value)

		switch {
		case rounded >= 1000:
			// Handle seconds for values >= 1000 ms
			return fmt.Sprintf("%.0f s", rounded/1000)
		case rounded >= 1:
			// Handle milliseconds for values >= 1 ms
			return fmt.Sprintf("%.0f ms", rounded)
		default:
			// Handle microseconds for values < 1 ms
			roundedInMicroseconds := roundHumanCascading(value * 1000)
			if roundedInMicroseconds < 1 {
				return "<1 μs"
			}
			return fmt.Sprintf("%.0f μs", roundedInMicroseconds)
		}

	default: // Not rounding
		switch {
		case value >= 1000:
			// Handle seconds for values >= 1000 ms
			return fmt.Sprintf("%.3f s", value/1000)
		case value >= 1:
			// Handle milliseconds for values >= 1 ms
			return fmt.Sprintf("%.3f ms", value)
		default:
			// Handle microseconds for values < 1 ms
			return fmt.Sprintf("%.3f μs", value*1000)
		}
	}
}

// roundHumanCascading rounds a float64 to the nearest whole number
// by applying cascading rounding from the most significant decimal place,
// with a final rounding step to ensure correctness.
func roundHumanCascading(value float64) float64 {
	// Handle negative values by flipping the sign
	if value < 0 {
		return -roundHumanCascading(-value)
	}

	// Start with the original value
	current := value

	// Iterate over decreasing decimal places (tenths, hundredths, etc.)
	multiplier := 1.0
	for i := 0; i < 10; i++ { // Limit to 10 digits
		multiplier *= 10

		// Shift the value to analyze the next decimal place
		shifted := current * multiplier
		_, frac := math.Modf(shifted)

		// Apply rounding if the fractional part >= 0.5
		if frac >= 0.5 {
			current = math.Floor(current*multiplier+0.5) / multiplier
		} else {
			// Stop if no rounding is required
			break
		}
	}

	// Perform one last round to handle any residual fractional part
	return math.Round(current)
}
