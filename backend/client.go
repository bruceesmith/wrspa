package wrserver

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/bruceesmith/logger"
)

type Client struct {
}

var wikiURL = "https://en.wikipedia.org"

// Get fetches a Wikipedia URL (either a static file or a dynamic page)
func (c *Client) Get(path string) (body []byte, err error) {
	var resp *http.Response
	logger.TraceID("client", "get", "URL", wikiURL+path)
	resp, err = http.Get(wikiURL + path)
	if err != nil {
		logger.Error("error fetching "+path, "error", err.Error())
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("unexpected status: %s", resp.Status)
		return
	}

	body, err = io.ReadAll(resp.Body)
	if err != nil {
		logger.Error("error reading response to GET("+path+")", "error", err.Error())
		return
	}
	return
}

// GetRandom fetches the URL path but not the page content returned
// by fetching the path /wiki/SpecialRandom but not redirecting
func (c *Client) GetRandom() (path string) {
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			path = strings.TrimPrefix(req.URL.Path, "/wiki/")
			return http.ErrUseLastResponse
		},
	}
	r, err := client.Get(wikiURL + "/wiki/Special:Random")
	if err != nil {
		logger.Error("error fetching Special:Random", "error", err.Error())
		return
	}
	defer r.Body.Close()

	if r.StatusCode != http.StatusFound {
		logger.Error("unexpected status getting random", "status", r.Status)
		return
	}

	return
}
