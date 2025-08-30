package wrserver

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/bruceesmith/terminator"
	"github.com/bruceesmith/wrspa/backend/wrserver/mocks"
	"go.uber.org/mock/gomock"
)

func TestNewServer(t *testing.T) {
	tests := []struct {
		name       string
		port       string
		static     string
		client     ClientInterface
		shouldFail bool
	}{
		{
			name:   "success",
			port:   "8080",
			static: "testdata",
			client: &Client{},
		},
		{
			name:       "bad port",
			port:       "-1",
			static:     "testdata",
			shouldFail: true,
		},
		{
			name:       "non-numeric port",
			port:       "a",
			static:     "testdata",
			shouldFail: true,
		},
		{
			name:       "static folder does not exist",
			port:       "8080",
			static:     "non-existent-folder",
			shouldFail: true,
		},
		{
			name:       "static is a file",
			port:       "8080",
			static:     "testdata/test.txt",
			shouldFail: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewServer(tt.port, tt.static, tt.client)
			if tt.shouldFail {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			}
		})
	}
}

func TestAPI(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mocks.NewMockClientInterface(ctrl)

	tests := []struct {
		name           string
		method         string
		function       string
		body           any
		statusCode     int
		expectedHeader map[string]string
		mockSetup      func()
	}{
		{
			name:           "settings",
			method:         http.MethodGet,
			function:       "settings",
			statusCode:     http.StatusOK,
			expectedHeader: map[string]string{"Content-Type": "application/json"},
		},
		{
			name:           "specialrandom",
			method:         http.MethodGet,
			function:       "specialrandom",
			statusCode:     http.StatusOK,
			expectedHeader: map[string]string{"Content-Type": "application/json"},
			mockSetup: func() {
				mockClient.EXPECT().GetRandom().Return("start")
				mockClient.EXPECT().GetRandom().Return("goal")
			},
		},
		{
			name:           "wikipage",
			method:         http.MethodPost,
			function:       "wikipage",
			body:           WikiPageRequest{Subject: "/wiki/test"},
			statusCode:     http.StatusOK,
			expectedHeader: map[string]string{"Content-Type": "text/html"},
			mockSetup: func() {
				mockClient.EXPECT().Get("/wiki/test").Return([]byte("<html><body><p>test</p></body></html>"), "text/html", nil)
			},
		},
		{
			name:       "wikipage invalid subject",
			method:     http.MethodPost,
			function:   "wikipage",
			body:       WikiPageRequest{Subject: "test"},
			statusCode: http.StatusBadRequest,
		},
		{
			name:       "wikipage read body error",
			method:     http.MethodPost,
			function:   "wikipage",
			body:       errorReader{},
			statusCode: http.StatusInternalServerError,
		},
		{
			name:       "wikipage unmarshal error",
			method:     http.MethodPost,
			function:   "wikipage",
			body:       "not json",
			statusCode: http.StatusBadRequest,
		},
		{
			name:       "wikipage client get error",
			method:     http.MethodPost,
			function:   "wikipage",
			body:       WikiPageRequest{Subject: "/wiki/test"},
			statusCode: http.StatusNotFound,
			mockSetup: func() {
				mockClient.EXPECT().Get("/wiki/test").Return(nil, "", errors.New("not found"))
			},
		},
		{
			name:       "wikipage no body tag",
			method:     http.MethodPost,
			function:   "wikipage",
			body:       WikiPageRequest{Subject: "/wiki/test"},
			statusCode: http.StatusInternalServerError,
			mockSetup: func() {
				mockClient.EXPECT().Get("/wiki/test").Return([]byte("<html></html>"), "text/html", nil)
			},
		},
		{
			name:       "bad method",
			method:     http.MethodDelete,
			function:   "settings",
			statusCode: http.StatusOK, // This will be the default status code from the recorder
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mockSetup != nil {
				tt.mockSetup()
			}

			s, _ := NewServer("8080", "testdata", mockClient)

			var body io.Reader
			if br, ok := tt.body.(io.Reader); ok {
				body = br
			} else {
				bodyBytes, _ := json.Marshal(tt.body)
				body = bytes.NewReader(bodyBytes)
			}

			req := httptest.NewRequest(tt.method, "/api/"+tt.function, body)
			w := httptest.NewRecorder()

			s.API(w, req)

			if w.Code != tt.statusCode {
				t.Errorf("got status code %d, want %d", w.Code, tt.statusCode)
			}

			for key, value := range tt.expectedHeader {
				if w.Header().Get(key) != value {
					t.Errorf("got header %s: %s, want %s", key, w.Header().Get(key), value)
				}
			}
		})
	}
}

func TestExtractBody(t *testing.T) {
	s := &Server{}

	tests := []struct {
		name    string
		html    string
		expect  string
		wantErr bool
	}{
		{
			name:   "simple body",
			html:   "<html><body><p>Hello</p></body></html>",
			expect: "<p>Hello</p>",
		},
		{
			name:    "no body",
			html:    "<html><head></head></html>",
			wantErr: true,
		},
		{
			name:   "empty body",
			html:   "<html><body></body></html>",
			expect: "",
		},
		{
			name:    "malformed html",
			html:    "<html",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, err := s.extractBody([]byte(tt.html))

			if (err != nil) != tt.wantErr {
				t.Errorf("extractBody() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && string(body) != tt.expect {
				t.Errorf("extractBody() = %v, want %v", string(body), tt.expect)
			}
		})
	}
}

func TestMarshalFailure(t *testing.T) {
	s := &Server{}
	json := s.MarshalFailure("test", errors.New("test error"), "test response")
	if !strings.Contains(json, `"msg": "unable to marshal API response"`) ||
		!strings.Contains(json, `"function": "test"`) ||
		!strings.Contains(json, `"error": "test error"`) ||
		!strings.Contains(json, `"response": "test response"`) {
		t.Errorf("unexpected json: %s", json)
	}
}

func TestServe(t *testing.T) {
	s, err := NewServer("8081", "testdata", &Client{})
	if err != nil {
		t.Fatalf("failed to create server: %v", err)
	}

	term := terminator.New()
	go s.Serve(term)

	// Give the server a moment to start
	time.Sleep(100 * time.Millisecond)

	// Send shutdown signal
	term.Stop()
}

func TestSPAFile(t *testing.T) {
	s, _ := NewServer("8080", "testdata", &Client{})

	tests := []struct {
		name       string
		path       string
		statusCode int
	}{
		{
			name:       "index",
			path:       "/",
			statusCode: http.StatusOK,
		},
		{
			name:       "other file",
			path:       "/test.txt",
			statusCode: http.StatusOK,
		},
		{
			name:       "not found",
			path:       "/notfound",
			statusCode: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tt.path, nil)
			w := httptest.NewRecorder()

			s.SPAFile(w, req)

			if w.Code != tt.statusCode {
				t.Errorf("got status code %d, want %d", w.Code, tt.statusCode)
			}
		})
	}
}

func TestWikipediaFile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mocks.NewMockClientInterface(ctrl)

	tests := []struct {
		name            string
		path            string
		statusCode      int
		wantContentType string
		mockSetup       func()
	}{
		{
			name:            "success",
			path:            "/w/test",
			statusCode:      http.StatusOK,
			wantContentType: "image/png",
			mockSetup: func() {
				mockClient.EXPECT().Get("/w/test").Return([]byte("test"), "image/png", nil)
			},
		},
		{
			name:       "not found",
			path:       "/w/notfound",
			statusCode: http.StatusNotFound,
			mockSetup: func() {
				mockClient.EXPECT().Get("/w/notfound").Return(nil, "", errors.New("not found"))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mockSetup != nil {
				tt.mockSetup()
			}

			s, _ := NewServer("8080", "testdata", mockClient)

			req := httptest.NewRequest(http.MethodGet, tt.path, nil)
			w := httptest.NewRecorder()

			s.WikipediaFile(w, req)

			if w.Code != tt.statusCode {
				t.Errorf("got status code %d, want %d", w.Code, tt.statusCode)
			}

			if w.Header().Get("Content-Type") != tt.wantContentType {
				t.Errorf("got content type %s, want %s", w.Header().Get("Content-Type"), tt.wantContentType)
			}
		})
	}
}

// errorReader is a helper type that implements io.Reader and always returns an error.
type errorReader struct{}

func (errorReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("forced read error")
}
