package wrserver

import (
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	slog.SetDefault(slog.New(slog.NewTextHandler(io.Discard, nil)))
	os.Exit(m.Run())
}

func TestGet(t *testing.T) {
	tests := []struct {
		name       string
		body       string
		want       string
		shouldFail bool
	}{
		{
			name: "success",
			body: "Hello, client",
			want: "Hello, client",
		},
		{
			name:       "failure",
			shouldFail: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if tt.shouldFail {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				w.Write([]byte(tt.body))
			}))
			defer server.Close()

			wikiURL = server.URL

			c := &Client{}
			body, err := c.Get("")
			if tt.shouldFail {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
				return
			}

			if string(body) != tt.want {
				t.Errorf("got %s, want %s", string(body), tt.want)
			}
		})
	}
}

func TestGetRandom(t *testing.T) {
	tests := []struct {
		name       string
		path       string
		want       string
		shouldFail bool
	}{
		{
			name: "success",
			path: "/wiki/Special:Random",
			want: "Special:Random",
		},
		{
			name:       "failure",
			shouldFail: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if tt.shouldFail {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				w.Header().Set("Location", tt.path)
				w.WriteHeader(http.StatusFound)
			}))
			defer server.Close()

			wikiURL = server.URL

			c := &Client{}
			path := c.GetRandom()

			if tt.shouldFail {
				if path != "" {
					t.Errorf("expected empty path, got %s", path)
				}
				return
			}

			if path != tt.want {
				t.Errorf("got %s, want %s", path, tt.want)
			}
		})
	}
}
