package wrserver

import (
	"context"
	"testing"
	"time"

	"github.com/bruceesmith/terminator"
	"github.com/bruceesmith/wrspa/backend/wrserver/mocks"
	"github.com/urfave/cli/v3"
	"go.uber.org/mock/gomock"
)

func Test_daemon(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "should call server Serve method",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockServer := mocks.NewMockServerInterface(ctrl)
			serveCalled := make(chan struct{})

			mockServer.EXPECT().Serve(gomock.Any()).Do(func(term *terminator.Terminator) {
				close(serveCalled)
			})

			term := terminator.New()
			go daemon(mockServer, term)

			select {
			case <-serveCalled:
				// Test passed
			case <-time.After(time.Second):
				t.Fatal("Serve was not called within a second")
			}

			// Stop the terminator to avoid leaking the goroutine
			term.Stop()
		})
	}
}

func TestDaemon_ErrorHandling(t *testing.T) {
	cmd := &cli.Command{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "port",
				Value: "invalid-port",
			},
			&cli.StringFlag{
				Name:  "static",
				Value: "/tmp",
			},
		},
	}

	err := Daemon(context.Background(), cmd)
	if err == nil {
		t.Fatal("expected an error, got nil")
	}
}
