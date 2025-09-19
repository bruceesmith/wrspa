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

func TestNewClient(t *testing.T) {
	tests := []struct {
		name    string
		wikiURL string
		want    string
	}{
		{
			name:    "default url",
			wikiURL: "",
			want:    defaultWikiURL,
		},
		{
			name:    "custom url",
			wikiURL: "http://localhost:8080",
			want:    "http://localhost:8080",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewClient(tt.wikiURL)
			if c.wikiURL != tt.want {
				t.Errorf("got %s, want %s", c.wikiURL, tt.want)
			}
		})
	}
}

func TestGet(t *testing.T) {
	// success and server failure cases
	t.Run("success and server failure", func(t *testing.T) {
		tests := []struct {
			name            string
			body            string
			want            string
			wantContentType string
			shouldFail      bool
		}{
			{
				name:            "success",
				body:            "Hello, client",
				want:            "Hello, client",
				wantContentType: "text/plain; charset=utf-8",
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
					w.Header().Set("Content-Type", tt.wantContentType)
					w.Write([]byte(tt.body))
				}))
				defer server.Close()

				c := NewClient(server.URL)
				body, contentType, err := c.Get("")
				if tt.shouldFail {
					if err == nil {
						t.Errorf("expected error, got nil")
					}
					return
				}

				if string(body) != tt.want {
					t.Errorf("got %s, want %s", string(body), tt.want)
				}

				if contentType != tt.wantContentType {
					t.Errorf("got content type %s, want %s", contentType, tt.wantContentType)
				}
			})
		}
	})

	// network error case
	t.Run("network error", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
		server.Close() // Close server immediately

		c := NewClient(server.URL)
		_, _, err := c.Get("")
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

		c := NewClient(server.URL)
		_, _, err := c.Get("")
		if err == nil {
			t.Fatal("expected an error reading the body, but got nil")
		}
	})
}

func TestGetRandom(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Location", "/wiki/Random_Page")
			w.WriteHeader(http.StatusFound)
		}))
		defer server.Close()

		c := NewClient(server.URL)
		path := c.GetRandom()
		if path != "Random_Page" {
			t.Errorf("got %s, want Random_Page", path)
		}
	})

	t.Run("server failure", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))
		defer server.Close()

		c := NewClient(server.URL)
		path := c.GetRandom()
		if path != "" {
			t.Errorf("expected empty path, got %s", path)
		}
	})

	t.Run("no location header", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusFound)
		}))
		defer server.Close()

		c := NewClient(server.URL)
		path := c.GetRandom()
		if path != "" {
			t.Errorf("expected empty path, got %s", path)
		}
	})

	t.Run("network error", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
		server.Close() // Close server immediately

		c := NewClient(server.URL)
		path := c.GetRandom()
		if path != "" {
			t.Fatalf("expected empty path on network error, but got '%s'", path)
		}
	})
}
