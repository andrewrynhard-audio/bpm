package main

import (
	"github.com/andrewrynhard-audio/bpm/pkg/state"
	"github.com/andrewrynhard-audio/bpm/pkg/ui"
	"github.com/andrewrynhard-audio/bpm/pkg/ui/timing"
	"github.com/andrewrynhard-audio/bpm/pkg/ui/update"
)

func main() {
	sharedState := &state.State{
		RoundOutputs: true,
	}

	t := timing.New()
	u := update.New()

	ui.Tap(sharedState, t, u)
}
