package state

type State struct {
	BPM          float64 // Beats per minute.
	RoundOutputs bool    // Round all calculated outputs to the nearest whole number.
	Locked       bool    // Lock the UI to prevent further input.
}

func (s *State) Reset() {
	s.BPM = 0
	s.RoundOutputs = true
	s.Locked = false
}
