package api

type SpecialRandomResponse struct {
	Start string `json:"start"`
	Goal  string `json:"goal"`
}

type EndPoint string

const (
	SpecialRandom EndPoint = "specialrandom"
	WikiPage      EndPoint = "wikipage"
)

type WikiPageRequest struct {
	Subject string `json:"subject"`
}

type WikiPageResponse struct {
	Page  string `json:"page"`
	Error string `json:"error"`
}
