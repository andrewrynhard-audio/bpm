package ui

import (
	"fmt"
	"math"
	"time"

	"github.com/andrewrynhard-audio/bpm/pkg/state"
	"github.com/nsf/termbox-go"
)

type Element interface {
	Render(float64, *state.SharedState)
	Reset(*state.SharedState)
	StateChanged(sharedState *state.SharedState)
}

func Tap(sharedState *state.SharedState, elements ...Element) {
	stop := make(chan struct{})
	defer close(stop)

	go func() {
		err := termbox.Init()
		if err != nil {
			fmt.Println("Error initializing termbox:", err)
			return
		}
		defer termbox.Close()

		var lastClick time.Time
		var intervals []time.Duration

		termbox.SetInputMode(termbox.InputMouse)

		for {
			ev := termbox.PollEvent()

			switch ev.Type {
			case termbox.EventMouse:
				if ev.Key == termbox.MouseLeft {
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
							element.Render(roundedBpm, sharedState)
						}
					}
				}

			case termbox.EventKey:
				if ev.Key == termbox.KeyEsc || ev.Ch == 'q' {
					stop <- struct{}{}
					return
				} else if ev.Ch == 'r' || ev.Ch == 'R' {
					intervals = nil
					lastClick = time.Time{}

					for _, element := range elements {
						element.Reset(sharedState)
					}
				} else if ev.Key == termbox.KeyF1 {
					sharedState.RoundOutputs = !sharedState.RoundOutputs

					for _, element := range elements {
						element.StateChanged(sharedState)
					}

				}

			case termbox.EventError:
				fmt.Println("Termbox error:", ev.Err)
				return
			}
		}
	}()

	<-stop
}
