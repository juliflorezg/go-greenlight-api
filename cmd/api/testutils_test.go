package main

import (
	"bytes"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
)

type testServer struct {
	*httptest.Server
}

func NewTestApplication(t *testing.T) *application {

	return &application{
		logger: slog.New(slog.NewTextHandler(io.Discard, nil)),
		config: config{
			port: 4000,
			env:  "development",
		},
	}
}

// Create a newTestServer helper which initializes and returns a new instance
// of our custom testServer type.
func NewTestServer(t *testing.T, h http.Handler) *testServer {
	ts := httptest.NewServer(h)

	// Disable redirect-following for the test server client by setting a custom
	// CheckRedirect function. This function will be called whenever a 3xx response is
	// received by the client, and by always returning a http.ErrUseLastResponse error
	// it forces the client to immediately return the received response.
	ts.Client().CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	return &testServer{ts}
}

// This method makes a GET request to a given url path
// using the test server client, and returns the
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

	body = bytes.TrimSpace(body)

	return rs.StatusCode, rs.Header, string(body)
}
