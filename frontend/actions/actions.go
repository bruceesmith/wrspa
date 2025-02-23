/*
Package actions contains functions that return strings representing named Actions
that are processed by handlers defined in Handle() and invoked by NewAction or NewActionWithValue
*/
package actions

// PageLoaded is invoked when the API server has returned an HTML page
func PageLoaded() string {
	return "pageLoaded"
}
