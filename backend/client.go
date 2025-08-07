package wrserver

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/bruceesmith/logger"
)

const defaultWikiURL = "https://en.wikipedia.org"

type Client struct {
	wikiURL string
}

func NewClient(wikiURL string) *Client {
	if wikiURL == "" {
		wikiURL = defaultWikiURL
	}
	return &Client{wikiURL: wikiURL}
}

// Get fetches a Wikipedia URL (either a static file or a dynamic page)
func (c *Client) Get(path string) (body []byte, err error) {
	var resp *http.Response
	url := c.wikiURL + path
	logger.TraceID("client", "get", "URL", url)
	resp, err = http.Get(url)
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
	url := c.wikiURL + "/wiki/Special:Random"
	resp, err := http.Head(url)
	if err != nil {
		logger.Error("error fetching Special:Random", "error", err.Error())
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusFound {
		logger.Error("unexpected status getting random", "status", resp.Status)
		return
	}

	path = strings.TrimPrefix(resp.Request.URL.Path, "/wiki/")
	return
}
