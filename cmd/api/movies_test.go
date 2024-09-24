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
			wantBody: `{
  "error": "The requested resource could not be found",
  "status_code": 404
}`,
		},
		"Negative ID": {
			urlPath:  "/v1/movies/-234",
			wantCode: 404,
			wantBody: `{
  "error": "The requested resource could not be found",
  "status_code": 404
}`,
		},
		"ID string": {
			urlPath:  "/v1/movies/asd",
			wantCode: 404,
			wantBody: `{
  "error": "The requested resource could not be found",
  "status_code": 404
}`,
		},
		"valid ID": {
			urlPath:  "/v1/movies/1",
			wantCode: 200,
			wantBody: `{
  "movie": {
    "id": 1,
    "title": "The Hunger Games",
    "year": 2012,
    "runtime": "142 mins",
    "genres": [
      "dystopian sci-fi",
      "action",
      "adventure"
    ],
    "version": 1
  }
}`,
		},
	}

	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			code, _, body := ts.get(t, tc.urlPath)

			assert.Equal(t, code, tc.wantCode)
			if tc.wantBody != "" {
				assert.Equal(t, body, tc.wantBody)
			}
		})
	}
}

func TestCreateMovieHandler(t *testing.T) {
	app := NewTestApplication(t)
	ts := NewTestServer(t, app.routes())
	defer ts.Close()

	tt := map[string]struct {
		input            string
		expectedOutput   string
		expectedResCode  int
		checkHeader      bool
		expectedLocation string
	}{
		"invalid xml": {
			input: `xml version="1.0" encoding="UTF-8"?><note><to>Alex</to></note>`,
			expectedOutput: `{
  "error": "body contains a badly-formed JSON (at character 1)",
  "status_code": 400
}`, // bc of the addition of \n in response
			expectedResCode: http.StatusBadRequest,
			checkHeader:     false,
		},
		"malformed json": {
			input: `{"title": "Moana", }`,
			expectedOutput: `{
  "error": "body contains a badly-formed JSON (at character 20)",
  "status_code": 400
}`,
			expectedResCode: http.StatusBadRequest,
			checkHeader:     false,
		},
		"array instead of object": {
			input: `["foo", "bar"]`,
			expectedOutput: `{
  "error": "request JSON body contains an incorrect type (at character 1)",
  "status_code": 400
}`,
			expectedResCode: http.StatusBadRequest,
			checkHeader:     false,
		},
		"numeric title value": {
			input: `{"title": 123}`,
			expectedOutput: `{
  "error": "request JSON body could not be parsed due to an incorrect type \"number\" for field \"title\" (type: string)",
  "status_code": 400
}`,
			expectedResCode: http.StatusBadRequest,
		},
		"empty body": {
			input: ``,
			expectedOutput: `{
  "error": "body must not be empty",
  "status_code": 400
}`,
			expectedResCode: http.StatusBadRequest,
			checkHeader:     false,
		},
		"valid input": {
			input: `{"title":"Deadpool","year":2016, "runtime":"108 mins","genres":["action","comedy"]}`,
			expectedOutput: `{
  "movie": {
    "id": 1,
    "title": "Deadpool",
    "year": 2016,
    "runtime": "108 mins",
    "genres": [
      "action",
      "comedy"
    ],
    "version": 1
  }
}`,
			expectedResCode:  http.StatusCreated,
			checkHeader:      true,
			expectedLocation: "/v1/movies/1",
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

			statusCode, headers, body := ts.post(t, "/v1/movies", tc.input)

			assert.Equal(t, statusCode, tc.expectedResCode)

			// assert.Equal(t, strings.ReplaceAll(body, " ", "-"), tc.expectedOutput)
			assert.Equal(t, body, tc.expectedOutput)

			if tc.checkHeader {
				assert.Equal(t, headers.Get("Location"), tc.expectedLocation)
			}
		})
	}
}
