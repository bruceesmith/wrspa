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
	// success and server failure cases
	t.Run("success and server failure", func(t *testing.T) {
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
				name:       "http status failure",
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
	})

	// network error case
	t.Run("network error", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
		wikiURL = server.URL
		server.Close() // Close server immediately

		c := &Client{}
		_, err := c.Get("")
		if err == nil {
			t.Fatal("expected a network error, but got nil")
		}
	})

	// read body error case
	t.Run("read body error", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			hj, ok := w.(http.Hijacker)
			if !ok {
				http.Error(w, "webserver doesn't support hijacking", http.StatusInternalServerError)
				return
			}
			conn, _, err := hj.Hijack()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			// Write a partial response and then close the connection to cause io.ReadAll to fail
			conn.Write([]byte("HTTP/1.1 200 OK\r\nContent-Length: 100\r\n\r\nHello"))
			conn.Close()
		}))
		defer server.Close()

		wikiURL = server.URL
		c := &Client{}
		_, err := c.Get("")
		if err == nil {
			t.Fatal("expected an error reading the body, but got nil")
		}
	})
}

func TestGetRandom(t *testing.T) {
	// success and server failure cases
	t.Run("success and server failure", func(t *testing.T) {
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
	})

	// network error case
	t.Run("network error", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
		wikiURL = server.URL
		server.Close() // Close server immediately

		c := &Client{}
		path := c.GetRandom()
		if path != "" {
			t.Fatalf("expected empty path on network error, but got '%s'", path)
		}
	})
}
