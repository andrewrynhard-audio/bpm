package tap

import (
	"fmt"
	"math"
	"time"

	"github.com/nsf/termbox-go"
)

// Loop starts the tap-to-tempo loop and returns a channel to send calculated BPMs.
// The BPM calculation resets explicitly when the spacebar is pressed.
func Loop() chan float64 {
	bpmChan := make(chan float64)

	go func() {
		err := termbox.Init()
		if err != nil {
			fmt.Println("Error initializing termbox:", err)
			close(bpmChan)
			return
		}
		defer termbox.Close()

		var lastClick time.Time
		var intervals []time.Duration

		termbox.SetInputMode(termbox.InputMouse)

		for {
			// Poll for an event
			ev := termbox.PollEvent()

			switch ev.Type {
			case termbox.EventMouse:
				if ev.Key == termbox.MouseLeft { // Detect left mouse button clicks
					now := time.Now()

					if !lastClick.IsZero() {
						intervals = append(intervals, now.Sub(lastClick))
					}
					lastClick = now

					// Calculate BPM
					if len(intervals) > 1 {
						var totalInterval time.Duration
						for _, interval := range intervals {
							totalInterval += interval
						}
						averageInterval := totalInterval / time.Duration(len(intervals))
						bpm := 60.0 / averageInterval.Seconds()

						// Round BPM to the nearest whole number
						roundedBpm := math.Round(bpm)

						// Send the rounded BPM to the channel
						bpmChan <- roundedBpm
					}
				}

			case termbox.EventKey:
				if ev.Ch == 'r' || ev.Ch == 'R' { // Reset on spacebar press
					intervals = nil
					lastClick = time.Time{}
					bpmChan <- 0.0
				} else if ev.Key == termbox.KeyEsc || ev.Ch == 'q' { // Exit on ESC or 'q'
					close(bpmChan)
					return
				}

			case termbox.EventError:
				fmt.Println("Termbox error:", ev.Err)
				close(bpmChan)
				return
			}
		}
	}()

	return bpmChan
}
