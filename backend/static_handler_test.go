package wrserver

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestStaticHandler_ServeHTTP(t *testing.T) {
	// Store original wikiURL and restore it after the test.
	// This allows us to redirect the `get` function's requests to our test server.
	originalWikiURL := wikiURL
	defer func() { wikiURL = originalWikiURL }()

	t.Run("should return file content when get is successful", func(t *testing.T) {
		// Arrange
		expectedBody := "svg data"
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(expectedBody))
		}))
		defer server.Close()
		wikiURL = server.URL

		handler := staticHandler{}
		req := httptest.NewRequest("GET", "/static/image.svg", nil)
		rr := httptest.NewRecorder()

		// Act
		handler.ServeHTTP(rr, req)

		// Assert
		if rr.Code != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", rr.Code, http.StatusOK)
		}

		body, err := io.ReadAll(rr.Body)
		if err != nil {
			t.Fatalf("Could not read response body: %v", err)
		}
		if string(body) != expectedBody {
			t.Errorf("handler returned unexpected body: got %q want %q", string(body), expectedBody)
		}
	})

	t.Run("should return 404 when get fails", func(t *testing.T) {
		// Arrange
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
		// Immediately close the server to simulate a network error for the `get` function.
		server.Close()
		wikiURL = server.URL

		handler := staticHandler{}
		req := httptest.NewRequest("GET", "/static/nonexistent.png", nil)
		rr := httptest.NewRecorder()

		// Act
		handler.ServeHTTP(rr, req)

		// Assert
		if rr.Code != http.StatusNotFound {
			t.Errorf("handler returned wrong status code: got %v want %v", rr.Code, http.StatusNotFound)
		}
	})
}
