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
func (c *Client) Get(path string) (body []byte, contentType string, err error) {
	var resp *http.Response
	url := c.wikiURL + path
	logger.TraceID("client", "get", "URL", url)

	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		logger.Error("error creating http.Request", "error", err.Error())
		return
	}
	req.Header.Add("Accept", "*/*")
	req.Header.Set("User-Agent", "wrspa/1.0")

	resp, err = client.Do(req)
	if err != nil {
		logger.Error("error fetching "+path, "error", err.Error())
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("unexpected status: %s", resp.Status)
		return
	}

	contentType = resp.Header.Get("Content-Type")

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
	logger.TraceID("client", "getrandom", "URL", url)

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	req, err := http.NewRequest("HEAD", url, nil)
	if err != nil {
		logger.Error("error creating http.Request", "error", err.Error())
		return
	}
	req.Header.Add("Accept", "*/*")
	req.Header.Set("User-Agent", "wrspa/1.0")

	resp, err := client.Do(req)
	if err != nil {
		logger.Error("error fetching Special:Random", "error", err.Error())
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusFound {
		logger.Error("unexpected status getting random", "status", resp.Status)
		return
	}

	l, err := resp.Location()
	if err != nil {
		logger.Error("error getting location", "error", err.Error())
		return
	}
	path = strings.TrimPrefix(l.Path, "/wiki/")
	return
}
