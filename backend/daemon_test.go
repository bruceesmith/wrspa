package wrserver

import (
	"testing"
	"time"

	"github.com/bruceesmith/wrspa/backend/wrserver/mocks"
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

			mockServer.EXPECT().Serve(gomock.Any()).Do(func(interface{}) {
				close(serveCalled)
			})

			// daemon() is a blocking function that waits for a signal.
			// We run it in a goroutine to avoid blocking the test.
			// This means the goroutine will be leaked, which is a tradeoff
			// for testing this kind of function without refactoring.
			go daemon(mockServer)

			select {
			case <-serveCalled:
				// Test passed
			case <-time.After(time.Second):
				t.Fatal("Serve was not called within a second")
			}
		})
	}
}
