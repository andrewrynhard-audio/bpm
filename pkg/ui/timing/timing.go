package timing

import (
	"fmt"
	"math"

	"github.com/andrewrynhard-audio/bpm/pkg/state"
	"github.com/gdamore/tcell/v2"
)

type Timing struct {
	BPM             float64
	Full            float64
	Half            float64
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

func (t *Timing) Render(sharedState *state.SharedState, screen tcell.Screen, bpm float64) {
	t.BPM = bpm
	t.calculate()
	t.write(sharedState.RoundOutputs, screen)
}

func (t *Timing) Reset(sharedState *state.SharedState, screen tcell.Screen) {
	t.BPM = 0
	t.Full = 0
	t.Half = 0
	t.Quarter = 0
	t.Eighth = 0
	t.Sixteenth = 0
	t.ThirtySecond = 0
	t.SixtyFourth = 0
	t.OneTwentyEighth = 0
	t.TwoFiftySixth = 0
	t.FiveTwelve = 0
	t.TenTwentyFour = 0

	t.write(sharedState.RoundOutputs, screen)
}

func (t *Timing) StateChanged(sharedState *state.SharedState, screen tcell.Screen) {
	// Rerender with the updated rounding state
	t.write(sharedState.RoundOutputs, screen)
}

func (t *Timing) calculate() {
	t.Full = 240000.00 / t.BPM
	t.Half = t.Full * 0.5
	t.Quarter = t.Half * 0.5
	t.Eighth = t.Quarter * 0.5
	t.Sixteenth = t.Quarter * 0.25
	t.ThirtySecond = t.Quarter * 0.125
	t.SixtyFourth = t.Quarter * 0.0625
	t.OneTwentyEighth = t.Quarter * 0.03125
	t.TwoFiftySixth = t.Quarter * 0.015625
	t.FiveTwelve = t.Quarter * 0.0078125
	t.TenTwentyFour = t.Quarter * 0.00390625
}

func (t *Timing) write(roundOutputs bool, screen tcell.Screen) {
	screen.Clear()

	// Table data
	headers := []string{"Note", "Time", "10x", "1/10th", "1/100th", "1/1000th"}
	notes := []string{"1", "1/2", "1/4", "1/8", "1/16", "1/32", "1/64", "1/128", "1/256", "1/512", "1/1024"}
	milliseconds := []float64{
		t.Full, t.Half, t.Quarter, t.Eighth, t.Sixteenth,
		t.ThirtySecond, t.SixtyFourth, t.OneTwentyEighth,
		t.TwoFiftySixth, t.FiveTwelve, t.TenTwentyFour,
	}

	// Column widths (in characters)
	columnWidths := []int{15, 15, 15, 15, 15, 15}
	totalWidth := 0
	for _, w := range columnWidths {
		totalWidth += w
	}

	// Total table height (header + spaced rows)
	rowPadding := 1
	totalHeight := len(notes)*(1+rowPadding) + 1 // +1 for the header row

	// Terminal size
	termWidth, termHeight := screen.Size()

	// Calculate offsets for centering
	startX := (termWidth - totalWidth) / 2
	startY := (termHeight - totalHeight) / 2

	// Render BPM at the top
	renderText(screen, startX, startY-2, fmt.Sprintf("BPM: %d", int(t.BPM)), tcell.StyleDefault.Foreground(tcell.ColorWhite))

	// Render table headers
	currentX := startX
	for i, header := range headers {
		renderText(screen, currentX, startY, header, tcell.StyleDefault.Foreground(tcell.ColorDarkCyan))
		currentX += columnWidths[i]
	}

	// Render each row with padding
	for row := 0; row < len(notes); row++ {
		ms := milliseconds[row]
		ms10x := ms * 10
		ms10 := ms / 10
		ms100 := ms / 100
		ms1000 := ms / 1000

		// Adjust startY for each row with padding
		currentRowY := startY + (row+1)*(1+rowPadding) // Add 1 for header and rowPadding

		currentX := startX
		renderText(screen, currentX, currentRowY, notes[row], tcell.StyleDefault.Foreground(tcell.ColorWhite))
		currentX += columnWidths[0]

		renderText(screen, currentX, currentRowY, formatWithUnit(ms, roundOutputs), tcell.StyleDefault.Foreground(tcell.ColorWhite))
		currentX += columnWidths[1]

		renderText(screen, currentX, currentRowY, formatWithUnit(ms10x, roundOutputs), tcell.StyleDefault.Foreground(tcell.ColorWhite))
		currentX += columnWidths[2]

		renderText(screen, currentX, currentRowY, formatWithUnit(ms10, roundOutputs), tcell.StyleDefault.Foreground(tcell.ColorWhite))
		currentX += columnWidths[3]

		renderText(screen, currentX, currentRowY, formatWithUnit(ms100, roundOutputs), tcell.StyleDefault.Foreground(tcell.ColorWhite))
		currentX += columnWidths[4]

		renderText(screen, currentX, currentRowY, formatWithUnit(ms1000, roundOutputs), tcell.StyleDefault.Foreground(tcell.ColorWhite))
	}

	// Render help message at the bottom
	helpMessage := "Press 'R' to reset, 'ESC' or 'Q' to quit, 'F1' to toggle whole numbers/decimals."
	renderText(screen, (termWidth-len(helpMessage))/2, termHeight-1, helpMessage, tcell.StyleDefault.Foreground(tcell.ColorGreen))

	screen.Show()
}

func renderText(screen tcell.Screen, x, y int, text string, style tcell.Style) {
	for i, ch := range text {
		screen.SetContent(x+i, y, ch, nil, style)
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
