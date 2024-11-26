package main

import (
	"context"
	"fmt"
	"log"
	"math"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

const (
	MaxClicks   = 24  // Number of clicks to consider
	MinInterval = 0.2 // Minimum interval in seconds
)

type DivisionMode int

const (
	NoDivision DivisionMode = iota
	DivideBy100
	DivideBy1000
)

type Timing struct {
	Full            string
	Half            string
	Quarter         string
	Eighth          string
	Sixteenth       string
	ThirtySecond    string
	SixtyFourth     string
	OneTwentyEighth string
}

type App struct {
	ctx          context.Context
	clickTimes   []time.Time
	BPM          float64
	Timings      Timing
	RoundOutputs bool
	DivisionMode DivisionMode
}

func NewApp() *App {
	return &App{
		clickTimes:   []time.Time{},
		Timings:      Timing{},
		RoundOutputs: true,
	}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) Click() float64 {
	log.Println("Version: ", Version)
	now := time.Now()
	a.clickTimes = append(a.clickTimes, now)

	// Keep only the last MaxClicks timestamps
	if len(a.clickTimes) > MaxClicks {
		a.clickTimes = a.clickTimes[len(a.clickTimes)-MaxClicks:]
	}

	// Ensure sufficient data
	if len(a.clickTimes) < 2 {
		return 0 // Not enough data
	}

	// Calculate intervals, ignoring very short ones
	intervals := []float64{}
	for i := 1; i < len(a.clickTimes); i++ {
		interval := a.clickTimes[i].Sub(a.clickTimes[i-1]).Seconds()
		if interval >= MinInterval {
			intervals = append(intervals, interval)
		}
	}

	// Return 0 if no valid intervals
	if len(intervals) == 0 {
		return 0
	}

	// Average the intervals
	averageInterval := 0.0
	for _, interval := range intervals {
		averageInterval += interval
	}
	averageInterval /= float64(len(intervals))

	// Calculate and return BPM
	bpm := 60.0 / averageInterval

	a.BPM = math.Round(bpm)

	return a.BPM
}

func (a *App) Reset() {
	a.clickTimes = []time.Time{}
}

func (a *App) GetTimings() Timing {
	a.calculateTiming(a.BPM)
	return a.Timings
}

func (a *App) calculateTiming(bpm float64) {
	if bpm <= 0 {
		return
	}

	calculateValue := func(value float64) string {
		switch a.DivisionMode {
		case DivideBy100:
			value /= 100
		case DivideBy1000:
			value /= 1000
		}
		formatted := formatWithUnit(value, a.RoundOutputs)
		return formatted
	}

	a.Timings.Full = calculateValue(240000.0 / bpm)
	a.Timings.Half = calculateValue(240000.0 / bpm * 0.5)
	a.Timings.Quarter = calculateValue(240000.0 / bpm * 0.25)
	a.Timings.Eighth = calculateValue(240000.0 / bpm * 0.125)
	a.Timings.Sixteenth = calculateValue(240000.0 / bpm * 0.0625)
	a.Timings.ThirtySecond = calculateValue(240000.0 / bpm * 0.03125)
	a.Timings.SixtyFourth = calculateValue(240000.0 / bpm * 0.015625)
	a.Timings.OneTwentyEighth = calculateValue(240000.0 / bpm * 0.0078125)
}

func (a *App) CycleDivisionMode() {
	a.DivisionMode = (a.DivisionMode + 1) % 3 // Cycle through 0 (NoDivision), 1 (DivideBy100), 2 (DivideBy1000)
	a.calculateTiming(a.BPM)                  // Recalculate timings with the updated mode
}

func (a *App) GetDivisionMode() DivisionMode {
	return a.DivisionMode
}

func (a *App) ToggleRounding() {
	a.RoundOutputs = !a.RoundOutputs
	a.calculateTiming(a.BPM) // Recalculate timings with the updated rounding preference
}

func (a *App) GetRoundState() bool {
	return a.RoundOutputs
}

func formatWithUnit(value float64, roundToWhole bool) string {
	// Handle negative values by returning a placeholder
	if value < 0 {
		return "N/A"
	}

	switch {
	case roundToWhole:
		switch {
		case value >= 1000:
			seconds := value / 1000
			roundedS := roundHumanCascading(seconds)
			return fmt.Sprintf("%.1f s", roundedS)
		case value >= 1:
			roundedMs := roundHumanCascading(value)
			return fmt.Sprintf("%.1f ms", roundedMs)
		case value < 1:
			microseconds := value * 1000
			roundedUs := roundHumanCascading(microseconds)
			return fmt.Sprintf("%.1f µs", roundedUs)
		default:
			// Should not reach here
			return fmt.Sprintf("%.1f", value)
		}

	default: // Not rounding
		switch {
		case value >= 1000:
			return fmt.Sprintf("%.3f s", value/1000)
		case value >= 1:
			return fmt.Sprintf("%.3f ms", value)
		default:
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

	// If the value is already very small, return it directly after rounding
	if value < 1e-10 {
		return math.Round(value)
	}

	adjustedValue := value // Start with the original value
	multiplier := 1.0

	for i := 0; i < 10; i++ { // Limit to 10 decimal places
		multiplier *= 10

		// Shift the value to analyze the next decimal place
		shifted := adjustedValue * multiplier
		intPart, fracPart := math.Modf(shifted)

		// Round the fractional part in its original scale
		roundedFracPart := math.Round(fracPart*10) / 10

		// Update adjusted value with the rounded fractional part
		adjustedValue = (intPart + roundedFracPart) / multiplier

		// Early stop if no significant fractional part remains
		if fracPart < 1e-10 {
			break
		}
	}

	// Final rounding to ensure consistency
	finalValue := math.Round(adjustedValue*10) / 10

	return finalValue
}

func (a *App) CheckForUpdate() UpdateInfo {
	info, err := CheckForUpdate()
	if err != nil {
		info.Message = fmt.Sprintf("Error checking for updates: %v", err)
	}

	return info
}

// OpenURL opens a URL in the user's default browser
func (a *App) OpenURL(url string) {
	runtime.BrowserOpenURL(a.ctx, url)
}
