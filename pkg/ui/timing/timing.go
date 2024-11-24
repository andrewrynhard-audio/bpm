package timing

import (
	"fmt"
	"math"

	"github.com/andrewrynhard-audio/bpm/pkg/state"
	"github.com/gdamore/tcell/v2"
)

type Timing struct {
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

func New() *Timing {
	return &Timing{}
}

func (t *Timing) Render(sharedState *state.State, screen tcell.Screen) {
	t.calculate(sharedState.BPM)
	t.write(screen, sharedState.BPM, sharedState.RoundOutputs, sharedState.Locked)
}

func (t *Timing) Reset(sharedState *state.State, screen tcell.Screen) {
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

	t.write(screen, sharedState.BPM, sharedState.RoundOutputs, sharedState.Locked)
}

func (t *Timing) StateChanged(sharedState *state.State, screen tcell.Screen) {
	// Rerender with the updated rounding state
	t.write(screen, sharedState.BPM, sharedState.RoundOutputs, sharedState.Locked)
}

func (t *Timing) calculate(bpm float64) {
	t.Full = 240000.00 / bpm
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

func (t *Timing) write(screen tcell.Screen, bpm float64, roundOutputs bool, locked bool) {
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

	if bpm == 0 {
		// Display a centered message when BPM is zero
		message := "Tap to get started"
		startX := (termWidth - len(message)) / 2
		startY := termHeight / 2

		renderText(screen, startX, startY, message, tcell.StyleDefault.Foreground(tcell.ColorGreen).Bold(true))

		return
	}

	// Calculate offsets for centering
	startX := (termWidth - totalWidth) / 2
	startY := (termHeight - totalHeight) / 2

	// Render BPM at the top
	renderText(screen, startX, startY-2, fmt.Sprintf("BPM: %d", int(bpm)), tcell.StyleDefault.Foreground(tcell.ColorGreen).Bold(true))

	// Render table headers
	currentX := startX
	for i, header := range headers {
		renderText(screen, currentX, startY, header, tcell.StyleDefault.Foreground(tcell.ColorDarkCyan).Bold(true))
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
		renderText(screen, currentX, currentRowY, notes[row], tcell.StyleDefault.Foreground(tcell.ColorDarkCyan))
		currentX += columnWidths[0]

		// Render each column with its respective style
		for _, value := range []float64{ms, ms10x, ms10, ms100, ms1000} {
			formatted, style := formatWithUnit(value, roundOutputs)
			renderText(screen, currentX, currentRowY, formatted, style)
			currentX += 15
		}
	}

	// Render help message at the bottom
	helpMessageTpl := "Press 'R' to reset, 'L' to %s, 'ESC' or 'Q' to quit, 'F1' to toggle whole numbers/decimals."
	helpMessage := fmt.Sprintf(helpMessageTpl, "lock")

	if locked {
		helpMessage = fmt.Sprintf(helpMessageTpl, "unlock")
	}

	renderText(screen, (termWidth-len(helpMessage))/2, termHeight-1, helpMessage, tcell.StyleDefault.Foreground(tcell.ColorGreen))
}

func renderText(screen tcell.Screen, x, y int, text string, style tcell.Style) {
	for i, ch := range text {
		screen.SetContent(x+i, y, ch, nil, style)
	}
}

// formatWithUnit formats a number dynamically based on its magnitude and rounding preference.
// - If `roundToWhole` is true, it rounds before determining the unit.
// - For values ≥ 1, it uses milliseconds (ms).
// - For values < 1, it converts to microseconds (us).
func formatWithUnit(value float64, roundToWhole bool) (string, tcell.Style) {
	// Define styles for different units with color-blind-friendly shades
	secondsStyle := tcell.StyleDefault.Foreground(tcell.ColorYellow)                     // Yellow
	millisecondsStyle := tcell.StyleDefault.Foreground(tcell.NewRGBColor(100, 150, 255)) // Light blue
	microsecondsStyle := tcell.StyleDefault.Foreground(tcell.NewRGBColor(255, 165, 0))   // Orange

	switch {
	case roundToWhole:
		rounded := roundHumanCascading(value)

		switch {
		case rounded >= 1000:
			// Handle seconds for values >= 1000 ms

			seconds := rounded / 1000

			// Round to the nearest tenth if greater than 1 second
			if seconds > 1 {
				seconds = math.Round(seconds*10) / 10
			}

			return fmt.Sprintf("%.1f s", seconds), secondsStyle
		case rounded >= 1:
			// Handle milliseconds for values >= 1 ms
			return fmt.Sprintf("%.0f ms", rounded), millisecondsStyle
		default:
			// Handle microseconds for values < 1 ms
			roundedInMicroseconds := roundHumanCascading(value * 1000)
			if roundedInMicroseconds < 1 {
				return "<1 us", microsecondsStyle
			}
			return fmt.Sprintf("%.0f us", roundedInMicroseconds), microsecondsStyle
		}

	default: // Not rounding
		switch {
		case value >= 1000:
			// Handle seconds for values >= 1000 ms
			return fmt.Sprintf("%.3f s", value/1000), secondsStyle
		case value >= 1:
			// Handle milliseconds for values >= 1 ms
			return fmt.Sprintf("%.3f ms", value), millisecondsStyle
		default:
			// Handle microseconds for values < 1 ms
			return fmt.Sprintf("%.3f us", value*1000), microsecondsStyle
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
