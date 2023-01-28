package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetMessages(t *testing.T) {
	// Create a new instance of our application struct which uses the mocked
	// dependencies.
	app := newTestApplication(t)
	// Establish a new test server for running end-to-end tests.
	ts := newTestServer(t, app.routes())

	defer ts.Close()

	tests := []struct {
		name     string
		urlPath  string
		wantCode int
		wantBody string
	}{
		{
			name:     "Valid parameters sent",
			urlPath:  "/message/list/john/doe",
			wantCode: 200,
			wantBody: "john",
		},
		{
			name:     "Incomplete parameter",
			urlPath:  "/message/list/john",
			wantCode: 404,
			wantBody: "",
		},
		{
			name:     "Missing Key",
			urlPath:  "/message/list/not/found",
			wantCode: 500,
			wantBody: "redis: nil",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, _, body := ts.get(t, tt.urlPath)
			assert.Equal(t, tt.wantCode, code)
			if tt.wantBody != "" {
				assert.Contains(t, body, tt.wantBody)
			}
		})
	}
}
