package wrserver

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func Test_getRandom(t *testing.T) {
	// Keep original wikiURL and restore it after the test
	originalWikiURL := wikiURL
	defer func() { wikiURL = originalWikiURL }()

	t.Run("successful redirect", func(t *testing.T) {
		expectedPath := "Test_Path"
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// The client should not follow this redirect.
			// The path is extracted from the Location header.
			w.Header().Set("Location", "/wiki/"+expectedPath)
			w.WriteHeader(http.StatusFound) // 302
		}))
		defer server.Close()

		wikiURL = server.URL

		path := getRandom()
		if path != expectedPath {
			t.Errorf("getRandom() = %v, want %v", path, expectedPath)
		}
	})

	t.Run("http get error", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
		// Close server immediately to simulate a connection error
		server.Close()

		wikiURL = server.URL

		path := getRandom()
		if path != "" {
			t.Errorf("getRandom() should return an empty string on error, but got %v", path)
		}
	})
}

func Test_get(t *testing.T) {
	// Keep original wikiURL and restore it after the test
	originalWikiURL := wikiURL
	defer func() { wikiURL = originalWikiURL }()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/good":
			fmt.Fprint(w, "hello")
		case "/bad":
			w.WriteHeader(http.StatusNotFound)
		default:
			w.WriteHeader(http.StatusInternalServerError)
		}
	}))
	defer server.Close()
	wikiURL = server.URL

	type args struct {
		path string
	}
	tests := []struct {
		name     string
		args     args
		wantBody []byte
		wantErr  bool
	}{
		{"good path", args{"/good"}, []byte("hello"), false},
		{"bad path", args{"/bad"}, []byte(""), false},
		{"server error", args{"/error"}, []byte(""), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotBody, err := get(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotBody, tt.wantBody) {
				t.Errorf("get() = %s, want %s", gotBody, tt.wantBody)
			}
		})
	}

	t.Run("http get error", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
		server.Close()
		wikiURL = server.URL

		_, err := get("/somepath")

		if err == nil {
			t.Error("get() expected an error, but got nil")
		}
	})
}

func Test_getString(t *testing.T) {
	originalWikiURL := wikiURL
	defer func() { wikiURL = originalWikiURL }()

	t.Run("successful get", func(t *testing.T) {
		expectedBody := "this is the body"
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, expectedBody)
		}))
		defer server.Close()
		wikiURL = server.URL

		body, err := getString("/somepath")

		if err != nil {
			t.Errorf("getString() returned unexpected error: %v", err)
		}
		if body != expectedBody {
			t.Errorf("getString() body = %q, want %q", body, expectedBody)
		}
	})

	t.Run("http get error", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
		server.Close()
		wikiURL = server.URL

		_, err := getString("/somepath")

		if err == nil {
			t.Error("getString() expected an error, but got nil")
		}
	})

	t.Run("body read error", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Length", "100")
			fmt.Fprint(w, "short body")
		}))
		defer server.Close()
		wikiURL = server.URL

		_, err := getString("/somepath")

		if err == nil {
			t.Error("getString() expected an error on body read, but got nil")
		}
	})
}
