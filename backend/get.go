package wrserver

import (
	"io"
	"net/http"
	"strings"

	"github.com/bruceesmith/logger"
)

func getRandom() (path string) {
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			path = strings.TrimPrefix(req.URL.Path, "/wiki/")
			return http.ErrUseLastResponse
		},
	}
	r, err := client.Get("https://en.wikipedia.org/wiki/Special:Random")
	if err != nil {
		logger.Error("error fetching Special:Random", "error", err.Error())
		return
	}
	defer r.Body.Close()
	return
}

func get(path string) (body []byte, err error) {
	var resp *http.Response
	logger.TraceID("server", "get", "URL", "https://en.wikipedia.org"+path)
	resp, err = http.Get("https://en.wikipedia.org" + path)
	if err != nil {
		logger.Error("error fetching "+path, "error", err.Error())
		return
	}
	defer resp.Body.Close()
	body, err = io.ReadAll(resp.Body)
	if err != nil {
		logger.Error("error reading response to GET("+path+")", "error", err.Error())
		return
	}
	return
}

func getString(path string) (body string, err error) {
	var b []byte
	b, err = get(path)
	if err != nil {
		return
	}
	body = string(b)
	return
}
