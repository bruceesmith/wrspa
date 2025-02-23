/*
Package observables contains functions that return strings representing named States
that are registered by ObserveState() and updated by SetState
*/
package observables

// ElapsedTime is updated by a timer
func ElapsedTime() string {
	return "elapsedTime"
}

// WikiState is updated according to button or anchor clicks
func WikiState() string {
	return "wikiState"
}
