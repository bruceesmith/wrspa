package wrserver

// SettingsResponse is the response for the settings endpoint
// It contains the log level and the trace IDs that are used for tracing
type SettingsResponse struct {
	LogLevel string   `json:"loglevel"`
	TraceIDs []string `json:"traceids"`
}

// SpecialRandomResponse is the response for the specialrandom endpoint
// It contains the random Wikipedia start and goal subjects
type SpecialRandomResponse struct {
	Start string `json:"start"`
	Goal  string `json:"goal"`
}

// EndPoint is the type for the endpoint names
type EndPoint string

const (
	Settings      EndPoint = "settings"      // Settings endpoint
	SpecialRandom EndPoint = "specialrandom" // Special random endpoint
	WikiPage      EndPoint = "wikipage"      // Wikipedia page endpoint
)

// WikiPageRequest is the request for the wikipage endpoint
// It contains the either the subject of the Wikipedia page to be retrieved
// or the link to an asset on the Wikipedia website
type WikiPageRequest struct {
	Subject string `json:"subject"`
}

// WikiPageResponse is the response for the wikipage endpoint
// It contains either a (string) page HTML or a (binary) asset,
// and an error message if the page or asset could not be retrieved
type WikiPageResponse struct {
	Page  string `json:"page"`
	Error string `json:"error"`
}
