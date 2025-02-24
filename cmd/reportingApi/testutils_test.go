package main

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"twitch_chat_analysis/internal/data"
	"twitch_chat_analysis/internal/data/mocks"
)

// Create a newTestApplication helper which returns an instance of our // application struct containing mocked dependencies.
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

// Create a newTestServer helper which initalizes and returns a new instance // of our custom testServer type.
func newTestServer(t *testing.T, h http.Handler) *testServer {
	ts := httptest.NewServer(h)
	return &testServer{ts}
}

// Implement a get() method on our custom testServer type. This makes a GET
// request to a given url path using the test server client, and returns the
// response status code, headers and body.
func (ts *testServer) get(t *testing.T, urlPath string) (int, http.Header, string) {
	rs, err := ts.Client().Get(ts.URL + urlPath)
	if err != nil {
		t.Fatal(err)
	}
	defer rs.Body.Close()
	body, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}
	bytes.TrimSpace(body)
	return rs.StatusCode, rs.Header, string(body)
}
