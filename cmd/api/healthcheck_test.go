package main

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/juliflorezg/greenlight/internal/assert"
)

func TestHealthCheck(t *testing.T) {

	// expectedBody := "status: available\nenvironment: development\nversion: 1.0.0"
	expectedBody := `{"status": "available", "environment": "development", "version": "1.0.0"}`
	expectedHeader := "application/json"

	app := NewTestApplication(t)

	ts := NewTestServer(t, app.routes())
	defer ts.Close()

	status, header, body := ts.get(t, "/v1/healthcheck")
	fmt.Printf("header map: %+v\n", header)
	fmt.Printf("content type for healthcheck test: %q\n", header.Get("Content-Type"))

	assert.Equal(t, status, http.StatusOK)
	assert.Equal(t, header.Get("Content-Type"), expectedHeader)
	assert.Equal(t, body, expectedBody)

}
