package main

import (
	"github.com/andrewrynhard-audio/bpm/pkg/state"
	"github.com/andrewrynhard-audio/bpm/pkg/ui"
	"github.com/andrewrynhard-audio/bpm/pkg/ui/timing"
)

func main() {
	sharedState := &state.State{
		RoundOutputs: true,
	}

	t := &timing.Timing{}

	ui.Tap(sharedState, t)
}
