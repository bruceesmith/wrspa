package wrserver

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bruceesmith/terminator"
	"github.com/bruceesmith/wrspa/backend/wrserver/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func Test_newClientAdapter(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "Test newClientAdapter",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := newClientAdapter()
			assert.NotNil(t, got, "newClientAdapter() should not return nil")
		})
	}
}

func Test_clientAdapter_Get(t *testing.T) {
	// Create a mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/success" {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`success`))
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	// Hijack the wikiURL to point to our mock server
	wikiURL = server.URL

	type fields struct {
		client ClientInterface
	}
	type args struct {
		path string
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		wantBody []byte
		wantErr  bool
	}{
		{
			name: "Test Get success",
			fields: fields{
				client: &Client{},
			},
			args: args{
				path: "/success",
			},
			wantBody: []byte(`success`),
			wantErr:  false,
		},
		{
			name: "Test Get failure",
			fields: fields{
				client: &Client{},
			},
			args: args{
				path: "/failure",
			},
			wantBody: nil,
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ca := &clientAdapter{
				client: tt.fields.client,
			}
			gotBody, err := ca.Get(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("clientAdapter.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.wantBody, gotBody, "clientAdapter.Get() gotBody = %v, want %v", gotBody, tt.wantBody)
		})
	}
}

func Test_newServerAdapter(t *testing.T) {
	type args struct {
		port   string
		static string
		client ClientInterface
	}
	tests := []struct {
		name    string
		args    args
		wantS   *serverAdapter
		wantErr bool
	}{
		{
			name: "Test newServerAdapter success",
			args: args{
				port:   "8080",
				static: "/tmp",
				client: mocks.NewMockClientInterface(gomock.NewController(t)),
			},
			wantErr: false,
		},
		{
			name: "Test newServerAdapter failure - bad port",
			args: args{
				port:   "bad port",
				static: "/tmp",
				client: mocks.NewMockClientInterface(gomock.NewController(t)),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotS, err := newServerAdapter(tt.args.port, tt.args.static, tt.args.client)
			if (err != nil) != tt.wantErr {
				t.Errorf("newServerAdapter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				assert.NotNil(t, gotS, "newServerAdapter() should not return nil on success")
			}
		})
	}
}

func Test_serverAdapter_Methods(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockServer := mocks.NewMockServerInterface(ctrl)
	sa := &serverAdapter{server: mockServer}

	tests := []struct {
		name  string
		setup func()
		act   func()
	}{
		{
			name: "API",
			setup: func() {
				mockServer.EXPECT().API(gomock.Any(), gomock.Any()).Times(1)
			},
			act: func() {
				sa.API(nil, nil)
			},
		},
		{
			name: "MarshalFailure",
			setup: func() {
				mockServer.EXPECT().MarshalFailure("test", errors.New("test error"), nil).Times(1)
			},
			act: func() {
				sa.MarshalFailure("test", errors.New("test error"), nil)
			},
		},
		{
			name: "Serve",
			setup: func() {
				mockServer.EXPECT().Serve(gomock.Any()).Times(1)
			},
			act: func() {
				sa.Serve(terminator.New())
			},
		},
		{
			name: "Settings",
			setup: func() {
				mockServer.EXPECT().Settings(gomock.Any(), gomock.Any()).Times(1)
			},
			act: func() {
				sa.Settings(nil, nil)
			},
		},
		{
			name: "SPAFile",
			setup: func() {
				mockServer.EXPECT().SPAFile(gomock.Any(), gomock.Any()).Times(1)
			},
			act: func() {
				sa.SPAFile(nil, nil)
			},
		},
		{
			name: "SpecialRandom",
			setup: func() {
				mockServer.EXPECT().SpecialRandom(gomock.Any(), gomock.Any()).Times(1)
			},
			act: func() {
				sa.SpecialRandom(nil, nil)
			},
		},
		{
			name: "WikiPage",
			setup: func() {
				mockServer.EXPECT().WikiPage(gomock.Any(), gomock.Any()).Times(1)
			},
			act: func() {
				sa.WikiPage(nil, nil)
			},
		},
		{
			name: "WikipediaFile",
			setup: func() {
				mockServer.EXPECT().WikipediaFile(gomock.Any(), gomock.Any()).Times(1)
			},
			act: func() {
				sa.WikipediaFile(nil, nil)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			tt.act()
		})
	}
}
