package ui

import (
	"fmt"
	"math"
	"time"

	"github.com/andrewrynhard-audio/bpm/pkg/state"
	"github.com/gdamore/tcell/v2"
)

type Element interface {
	Render(*state.SharedState, tcell.Screen, float64)
	Reset(*state.SharedState, tcell.Screen)
	StateChanged(*state.SharedState, tcell.Screen)
}

func Tap(sharedState *state.SharedState, elements ...Element) {
	screen, err := tcell.NewScreen()
	if err != nil {
		fmt.Println("Error creating tcell screen:", err)
		return
	}
	defer screen.Fini()

	err = screen.Init()
	if err != nil {
		fmt.Println("Error initializing tcell screen:", err)
		return
	}

	screen.EnableMouse()
	screen.Clear()

	var lastClick time.Time
	var intervals []time.Duration

	for {
		// Poll for events
		event := screen.PollEvent()

		switch ev := event.(type) {
		case *tcell.EventMouse:
			if sharedState.Locked {
				continue
			}

			if ev.Buttons() == tcell.Button1 { // Left mouse button
				now := time.Now()

				if !lastClick.IsZero() {
					intervals = append(intervals, now.Sub(lastClick))
				}
				lastClick = now

				if len(intervals) > 1 {
					var totalInterval time.Duration
					for _, interval := range intervals {
						totalInterval += interval
					}
					averageInterval := totalInterval / time.Duration(len(intervals))
					bpm := 60.0 / averageInterval.Seconds()

					roundedBpm := math.Round(bpm)

					for _, element := range elements {
						element.Render(sharedState, screen, roundedBpm)
					}
				}
			}

		case *tcell.EventKey:
			switch ev.Key() {
			case tcell.KeyEscape, tcell.KeyCtrlC:
				return
			case tcell.KeyF1:
				sharedState.RoundOutputs = !sharedState.RoundOutputs
				for _, element := range elements {
					element.StateChanged(sharedState, screen)
				}
			case tcell.KeyRune:
				switch ev.Rune() {
				case 'q', 'Q':
					return
				case 'r', 'R':
					intervals = nil
					lastClick = time.Time{}
					sharedState.Reset()
					for _, element := range elements {
						element.Reset(sharedState, screen)
					}
				case 'F':
					sharedState.RoundOutputs = !sharedState.RoundOutputs
					for _, element := range elements {
						element.StateChanged(sharedState, screen)
					}
				case 'l', 'L':
					sharedState.Locked = !sharedState.Locked
					for _, element := range elements {
						element.StateChanged(sharedState, screen)
					}
				}
			}

		case *tcell.EventResize:
			screen.Clear()
			for _, element := range elements {
				element.Render(sharedState, screen, 0)
			}

		default:
			// Ignore other events
		}
	}
}
