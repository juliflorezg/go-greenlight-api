package main

import (
	"net/http"
	"testing"

	"github.com/juliflorezg/greenlight/internal/assert"
)

func TestHealthCheck(t *testing.T) {

	expectedBody := "status: available\nenvironment: development\nversion: 1.0.0"

	app := NewTestApplication(t)

	ts := NewTestServer(t, app.routes())
	defer ts.Close()

	status, _, body := ts.get(t, "/v1/healthcheck")

	assert.Equal(t, status, http.StatusOK)
	assert.Equal(t, body, expectedBody)

}
