package main

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"twitch_chat_analysis/internal/data"
	"twitch_chat_analysis/internal/data/mocks"
)

// Create a newTestApplication helper which returns an instance of our
// application struct containing mocked dependencies.
func newTestApplication(t *testing.T) *application {
	return &application{
		models: data.Models{
			Message: mocks.MessageModel{},
		},
	}
}

// Define a custom testServer type which embeds a httptest.Server instance.
type testServer struct {
	*httptest.Server
}

// Create a newTestServer helper which initalizes and returns a new instance
// of our custom testServer type.
func newTestServer(t *testing.T, h http.Handler) *testServer {
	ts := httptest.NewServer(h)
	return &testServer{ts}
}

// Create a post method for sending POST requests to the test server.
func (ts *testServer) post(t *testing.T, urlPath string, payload string) (int, http.Header, string) {
	rs, err := ts.Client().Post(ts.URL+urlPath, "application/json", strings.NewReader(payload))
	if err != nil {
		t.Fatal(err)
	}
	// Read the response body from the test server.
	defer rs.Body.Close()
	body, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}
	bytes.TrimSpace(body)
	// Return the response status, headers and body.
	return rs.StatusCode, rs.Header, string(body)
}
