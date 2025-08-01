package wrserver

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestFileHandler_ServeHTTP(t *testing.T) {
	// Create a temporary directory to act as the root for static files.
	tempDir, err := os.MkdirTemp("", "test-static-")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	// Clean up the temporary directory after the test finishes.
	defer os.RemoveAll(tempDir)

	// Create a dummy index.html file.
	indexContent := "<html><body><h1>Index</h1></body></html>"
	if err := os.WriteFile(filepath.Join(tempDir, "index.html"), []byte(indexContent), 0644); err != nil {
		t.Fatalf("Failed to write index.html: %v", err)
	}

	// Create another dummy file.
	appJSContent := "console.log('hello world');"
	if err := os.WriteFile(filepath.Join(tempDir, "app.js"), []byte(appJSContent), 0644); err != nil {
		t.Fatalf("Failed to write app.js: %v", err)
	}

	handler := fileHandler{root: tempDir}

	tests := []struct {
		name           string
		target         string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "request for root should serve index.html",
			target:         "/",
			expectedStatus: http.StatusOK,
			expectedBody:   indexContent,
		},
		{
			name:           "request for a specific file should serve that file",
			target:         "/app.js",
			expectedStatus: http.StatusOK,
			expectedBody:   appJSContent,
		},
		{
			name:           "request for a non-existent file should return 404",
			target:         "/nonexistent.css",
			expectedStatus: http.StatusNotFound,
			expectedBody:   "404 page not found\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a request to pass to our handler.
			req, err := http.NewRequest("GET", tt.target, nil)
			if err != nil {
				t.Fatal(err)
			}

			// We use a ResponseRecorder to record the response.
			rr := httptest.NewRecorder()

			// Call the handler's ServeHTTP method directly.
			handler.ServeHTTP(rr, req)

			// Check the status code.
			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tt.expectedStatus)
			}

			// Check the response body.
			body, err := io.ReadAll(rr.Body)
			if err != nil {
				t.Fatalf("Could not read response body: %v", err)
			}

			if string(body) != tt.expectedBody {
				t.Errorf("handler returned unexpected body: got %q want %q",
					string(body), tt.expectedBody)
			}
		})
	}
}
