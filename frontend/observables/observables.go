/*
Package observables contains functions that return strings representing named States
that are registered by ObserveState() and updated by SetState
*/
package observables

// ElapsedTime is updated by a timer
func ElapsedTime() string {
	return "elapsedTime"
}

// GameSelected is updated when either Custom or Random is chosen
func GameSelected() string {
	return "gameSelected"
}

// WikiState is updated according to button or anchor clicks
func WikiState() string {
	return "wikiState"
}
