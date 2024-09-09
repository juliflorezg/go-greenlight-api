package main

import (
	"testing"

	"github.com/juliflorezg/greenlight/internal/assert"
)

func TestShowMovie(t *testing.T) {
	app := NewTestApplication(t)

	ts := NewTestServer(t, app.routes())
	defer ts.Close()

	// tests table
	tt := map[string]struct {
		urlPath  string
		wantCode int
		wantBody string
	}{
		"ID zero": {
			urlPath:  "/v1/movies/0",
			wantCode: 404,
			wantBody: "404 page not found",
		},
		"Negative ID": {
			urlPath:  "/v1/movies/-234",
			wantCode: 404,
			wantBody: "404 page not found",
		},
		"ID string": {
			urlPath:  "/v1/movies/asd",
			wantCode: 404,
			wantBody: "404 page not found",
		},
	}

	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			code, _, body := ts.get(t, tc.urlPath)

			assert.Equal(t, code, tc.wantCode)
			if tc.wantBody != "" {
				assert.StringContains(t, body, tc.wantBody)
			}
		})
	}
}
