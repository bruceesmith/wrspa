package wrserver

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
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
			static: "/tmp",
			client: &Client{},
		},
		{
			name:       "bad port",
			port:       "-1",
			shouldFail: true,
		},
		{
			name:       "non-numeric port",
			port:       "a",
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
		name       string
		method     string
		function   string
		body       any
		statusCode int
		mockSetup  func()
	}{
		{
			name:       "settings",
			method:     http.MethodGet,
			function:   "settings",
			statusCode: http.StatusOK,
		},
		{
			name:       "specialrandom",
			method:     http.MethodGet,
			function:   "specialrandom",
			statusCode: http.StatusOK,
			mockSetup: func() {
				mockClient.EXPECT().GetRandom().Return("start")
				mockClient.EXPECT().GetRandom().Return("goal")
			},
		},
		{
			name:       "wikipage",
			method:     http.MethodPost,
			function:   "wikipage",
			body:       WikiPageRequest{Subject: "test"},
			statusCode: http.StatusOK,
			mockSetup: func() {
				mockClient.EXPECT().Get("/test").Return([]byte("<body class=\"test\">test</body>"), nil)
			},
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
			body:       WikiPageRequest{Subject: "test"},
			statusCode: http.StatusNotFound,
			mockSetup: func() {
				mockClient.EXPECT().Get("/test").Return(nil, errors.New("not found"))
			},
		},
		{
			name:       "wikipage regex mismatch",
			method:     http.MethodPost,
			function:   "wikipage",
			body:       WikiPageRequest{Subject: "test"},
			statusCode: http.StatusInternalServerError,
			mockSetup: func() {
				mockClient.EXPECT().Get("/test").Return([]byte("no body tag"), nil)
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

			s, _ := NewServer("8080", "/tmp", mockClient)

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
		})
	}
}

func TestMarshalFailure(t *testing.T) {
	s := &Server{}
	json := s.MarshalFailure("test", errors.New("test error"), "test response")
	if json != `{"msg": "unable to marshal API response", "function": "test", "error": "test error", "response": "test response"}` {
		t.Errorf("unexpected json: %s", json)
	}
}

func TestServe(t *testing.T) {
	s, err := NewServer("8081", "/tmp", &Client{})
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
		name       string
		path       string
		statusCode int
		mockSetup  func()
	}{
		{
			name:       "success",
			path:       "/w/test",
			statusCode: http.StatusOK,
			mockSetup: func() {
				mockClient.EXPECT().Get("/w/test").Return([]byte("test"), nil)
			},
		},
		{
			name:       "not found",
			path:       "/w/notfound",
			statusCode: http.StatusNotFound,
			mockSetup: func() {
				mockClient.EXPECT().Get("/w/notfound").Return(nil, errors.New("not found"))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mockSetup != nil {
				tt.mockSetup()
			}

			s, _ := NewServer("8080", "/tmp", mockClient)

			req := httptest.NewRequest(http.MethodGet, tt.path, nil)
			w := httptest.NewRecorder()

			s.WikipediaFile(w, req)

			if w.Code != tt.statusCode {
				t.Errorf("got status code %d, want %d", w.Code, tt.statusCode)
			}
		})
	}
}

// errorReader is a helper type that implements io.Reader and always returns an error.
type errorReader struct{}

func (errorReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("forced read error")
}
