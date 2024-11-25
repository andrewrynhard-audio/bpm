package state

type SharedState struct {
	RoundOutputs bool // Round all calculated outputs to the nearest whole number.
	Locked       bool // Lock the UI to prevent further input.
}

func (s *SharedState) Reset() {
	s.RoundOutputs = true
	s.Locked = false
}
