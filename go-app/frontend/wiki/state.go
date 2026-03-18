package wiki

// state represents the current position of the Wiki finite state machine
type state int

const (
	ready    state = iota // Setup has been complete; start & goal subjects are known
	playing               // Game is in progress
	paused                // Game has been paused
	finished              // The goal has been reached
)
