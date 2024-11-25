package main

import (
	"fmt"

	"github.com/andrewrynhard-audio/bpm/pkg/state"
	"github.com/andrewrynhard-audio/bpm/pkg/ui"
	"github.com/andrewrynhard-audio/bpm/pkg/ui/timing"
	"github.com/andrewrynhard-audio/bpm/pkg/update"
)

func main() {
	err := update.Check()
	if err != nil {
		fmt.Printf("Uh-oh, looks like we couldn't check for updates: %v\n", err)
	}

	sharedState := &state.State{
		RoundOutputs: true,
	}

	t := &timing.Timing{}

	ui.Tap(sharedState, t)
}
