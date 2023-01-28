package main

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCreateMessage(t *testing.T) {
	// Create a new instance of our application struct which uses the mocked
	// dependencies.
	app := newTestApplication(t)
	// Establish a new test server for running end-to-end tests.
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	tests := []struct {
		name      string
		wantCode  int
		wantBody  string
		sender    string
		receiver  string
		message   string
		createdAt time.Time
	}{
		{
			name:      "Valid Submission",
			wantCode:  200,
			wantBody:  "OK",
			sender:    "john",
			receiver:  "doe",
			message:   "Hello",
			createdAt: time.Now(),
		},
		{
			name:      "Missing required field",
			wantCode:  400,
			wantBody:  "required",
			sender:    "john",
			receiver:  "",
			message:   "Hello",
			createdAt: time.Now(),
		},
		{
			name:      "System Error",
			wantCode:  400,
			wantBody:  "Bad Request",
			sender:    "error",
			receiver:  "doe",
			message:   "Hello",
			createdAt: time.Now(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			payload := fmt.Sprintf(`{
				"sender":   "%s",
				"receiver": "%s",
				"message":  "%s"
    		}`, tt.sender, tt.receiver, tt.message)
			code, _, body := ts.post(t, "/message", payload)
			assert.Equal(t, tt.wantCode, code)
			if tt.wantBody != "" {
				assert.Contains(t, body, tt.wantBody)
			}
		})
	}
}
