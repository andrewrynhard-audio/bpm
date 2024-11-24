package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/andrewrynhard-audio/bpm/pkg/tap"
	"github.com/nsf/termbox-go"
)

type info struct {
	// bpm is the beats per minute
	bpm             float64
	quarter         float64
	eighth          float64
	sixteenth       float64
	thirtySecond    float64
	sixtyFourth     float64
	oneTwentyEighth float64
	twoFiftySixth   float64
	fiveTwelve      float64
	tenTwentyFour   float64
}

func New(bpm float64) *info {
	i := &info{
		bpm: bpm,
	}
	i.calculate()

	return i
}

func (i *info) calculate() {
	i.quarter = 60000.00 / i.bpm
	i.eighth = i.quarter * 0.5
	i.sixteenth = i.quarter * 0.25
	i.thirtySecond = i.quarter * 0.125
	i.sixtyFourth = i.quarter * 0.0625
	i.oneTwentyEighth = i.quarter * 0.03125
	i.twoFiftySixth = i.quarter * 0.015625
	i.fiveTwelve = i.quarter * 0.0078125
	i.tenTwentyFour = i.quarter * 0.00390625
}

func main() {
	var bpm float64

	flag := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	flag.Float64Var(&bpm, "bpm", 0.0, "beats per minute")
	flag.Parse(os.Args[1:])

	if bpm == 0.0 {
		bpmChan := tap.Loop()

		for bpm := range bpmChan {
			info := New(bpm)
			write(info)
		}
	} else {
		info := New(bpm)
		write(info)
	}
}

func gen(s string, f float64) []any {
	return []any{
		s,
		f,
		f / 10,
		f / 100,
	}
}

func write(info *info) {
	// Clear the screen
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

	// Define table headers and data
	headers := []string{"Note", "Milliseconds", "1/10th", "1/100th"}
	notes := []string{
		"1/4",
		"1/8",
		"1/16",
		"1/32",
		"1/64",
		"1/128",
		"1/256",
		"1/512",
		"1/1024",
	}
	milliseconds := []float64{
		info.quarter,
		info.eighth,
		info.sixteenth,
		info.thirtySecond,
		info.sixtyFourth,
		info.oneTwentyEighth,
		info.twoFiftySixth,
		info.fiveTwelve,
		info.tenTwentyFour,
	}

	// Calculate table dimensions
	tableWidth := 60              // Fixed width for table
	tableHeight := len(notes) + 2 // Number of rows (1 header row + data rows)
	termWidth, termHeight := termbox.Size()
	startX := (termWidth - tableWidth) / 2
	startY := (termHeight - tableHeight) / 2

	// Render BPM at the top of the table
	renderText(startX, startY-2, fmt.Sprintf("BPM: %d", int(info.bpm)), termbox.ColorWhite)

	// Render table headers
	for i, header := range headers {
		renderText(startX+i*15, startY, header, termbox.ColorCyan)
	}

	// Render table rows
	for row := 0; row < len(notes); row++ {
		renderText(startX, startY+row+1, notes[row], termbox.ColorWhite)
		renderText(startX+15, startY+row+1, fmt.Sprintf("%.3f", milliseconds[row]), termbox.ColorWhite)
		renderText(startX+30, startY+row+1, fmt.Sprintf("%.3f", milliseconds[row]/10), termbox.ColorWhite)
		renderText(startX+45, startY+row+1, fmt.Sprintf("%.3f", milliseconds[row]/100), termbox.ColorWhite)
	}

	// Render help message at the bottom of the terminal
	helpMessage := "Press 'R' to reset, 'ESC' or 'Q' to quit."
	renderText((termWidth-len(helpMessage))/2, termHeight-1, helpMessage, termbox.ColorGreen)

	// Flush to display
	termbox.Flush()
}

// renderText renders text at a specific location with a specified color
func renderText(x, y int, text string, color termbox.Attribute) {
	for i, ch := range text {
		termbox.SetCell(x+i, y, ch, color, termbox.ColorDefault)
	}
}
