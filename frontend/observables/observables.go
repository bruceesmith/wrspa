/*
Package observables contains string constants representing named States
that are registered by ObserveState() and updated by SetState
*/
package observables

const (
	ElapsedTime  = "elapsedTime"  // ElapsedTime is updated by a timer
	GameSelected = "gameSelected" // GameSelected is updated when either Custom or Random is chosen
	WikiState    = "wikiState"    // WikiState is updated according to button or anchor clicks
)
