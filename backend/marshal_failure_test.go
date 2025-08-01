package wrserver

import (
	"errors"
	"testing"
)

func Test_marshalFailure(t *testing.T) {
	type simpleResponse struct {
		Field string
	}

	type complexResponse struct {
		ID   int
		Name string
	}

	tests := []struct {
		name     string
		function string
		err      error
		response any
		want     string
	}{
		{
			name:     "simple case with a struct response",
			function: "testFunc",
			err:      errors.New("test error"),
			response: simpleResponse{Field: "value"},
			want:     `{"msg": "unable to marshal API response", "function": "testFunc", "error": "test error", "response": "{Field:value}"}`,
		},
		{
			name:     "case with a nil response",
			function: "anotherFunc",
			err:      errors.New("another error"),
			response: nil,
			want:     `{"msg": "unable to marshal API response", "function": "anotherFunc", "error": "another error", "response": "<nil>"}`,
		},
		{
			name:     "case with a more complex struct response",
			function: "complexFunc",
			err:      errors.New("complex error"),
			response: complexResponse{ID: 123, Name: "test"},
			want:     `{"msg": "unable to marshal API response", "function": "complexFunc", "error": "complex error", "response": "{ID:123 Name:test}"}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := marshalFailure(tt.function, tt.err, tt.response); got != tt.want {
				t.Errorf("marshalFailure() = %v, want %v", got, tt.want)
			}
		})
	}
}
