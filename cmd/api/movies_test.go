package main

import (
	"net/http"
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

func TestCreateMovieHandler(t *testing.T) {
	app := NewTestApplication(t)
	ts := NewTestServer(t, app.routes())
	defer ts.Close()

	tt := map[string]struct {
		input          string
		expectedOutput string
	}{
		"invalid xml": {
			input: `xml version="1.0" encoding="UTF-8"?><note><to>Alex</to></note>`,
			expectedOutput: `{
  "error": "body contains a badly-formed JSON (at character 1)",
  "status_code": 400
}
`, // bc of the addition of \n in response
		},
		"malformed json": {
			input: `{"title": "Moana", }`,
			expectedOutput: `{
  "error": "body contains a badly-formed JSON (at character 20)",
  "status_code": 400
}
`,
		},
		"array instead of object": {
			input: `["foo", "bar"]`,
			expectedOutput: `{
  "error": "request JSON body contains an incorrect type (at character 1)",
  "status_code": 400
}
`,
		},
		"numeric title value": {
			input: `{"title": 123}`,
			expectedOutput: `{
  "error": "request JSON body could not be parsed due to an incorrect type string for field title",
  "status_code": 400
}
`,
		},
		"empty body": {
			input: ``,
			expectedOutput: `{
  "error": "body must not be empty",
  "status_code": 400
}
`,
		},
	}

	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			// res, err := ts.Client().Post(ts.URL+"/v1/movies", "application/json", strings.NewReader(tc.input))
			// if err != nil {
			// 	t.Fatal(err)
			// }
			// defer res.Body.Close()
			// body, _ := io.ReadAll(res.Body)

			statusCode, _, body := ts.post(t, "/v1/movies", tc.input)

			assert.Equal(t, statusCode, http.StatusBadRequest)

			assert.Equal(t, body, tc.expectedOutput)
		})
	}
}
