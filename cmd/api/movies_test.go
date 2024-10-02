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

func TestListMoviesHandler(t *testing.T) {
	app := NewTestApplication(t)
	ts := NewTestServer(t, app.routes())
	defer ts.Close()

	tt := map[string]struct {
		queryString        string
		expectedResponse   string
		expectedCodeStatus int
	}{
		// "no query string(def values)": {
		// 	queryString: "",
		// },
		"incorrect page": {
			queryString: "/v1/movies?title=moana&genres=action,adventure&page=abc&page_size=10&sort=year",
			expectedResponse: `{
  "error": {
    "page": "must be an integer value"
  },
  "status_code": 422
}`,
			expectedCodeStatus: 422,
		},
		"incorrect page (negative)": {
			queryString: "/v1/movies?title=moana&genres=action,adventure&page=-10&page_size=10&sort=year",
			expectedResponse: `{
  "error": {
    "page": "must be greater than 0"
  },
  "status_code": 422
}`,
			expectedCodeStatus: 422,
		},
		"incorrect page (out of bounds)": {
			queryString: "/v1/movies?title=moana&genres=action,adventure&page=20000000&page_size=10&sort=year",
			expectedResponse: `{
  "error": {
    "page": "must be less than or equal to 10 million"
  },
  "status_code": 422
}`,
			expectedCodeStatus: 422,
		},
		"incorrect page_size": {
			queryString: "/v1/movies?title=moana&genres=action,adventure&page=1&page_size=abc&sort=year",
			expectedResponse: `{
  "error": {
    "page_size": "must be an integer value"
  },
  "status_code": 422
}`,
			expectedCodeStatus: 422,
		},
		"incorrect page_size (negative)": {
			queryString: "/v1/movies?title=moana&genres=action,adventure&page=1&page_size=-15&sort=year",
			expectedResponse: `{
  "error": {
    "page_size": "must be greater than 0"
  },
  "status_code": 422
}`,
			expectedCodeStatus: 422,
		},
		"incorrect page_size (out of bounds)": {
			queryString: "/v1/movies?title=moana&genres=action,adventure&page=1&page_size=200&sort=year",
			expectedResponse: `{
  "error": {
    "page_size": "must be less than or equal to 100"
  },
  "status_code": 422
}`,
			expectedCodeStatus: 422,
		},
		"incorrect sort value (typo)": {
			queryString: "/v1/movies?title=moana&genres=action,adventure&page=1&page_size=10&sort=yearr",
			expectedResponse: `{
  "error": {
    "sort": "must be a valid sort value"
  },
  "status_code": 422
}`,
			expectedCodeStatus: 422,
		},
		"incorrect sort value (value not recognized)": {
			queryString: "/v1/movies?title=moana&genres=action,adventure&page=1&page_size=10&sort=rating",
			expectedResponse: `{
  "error": {
    "sort": "must be a valid sort value"
  },
  "status_code": 422
}`,
			expectedCodeStatus: 422,
		},
		"valid values for title, genres, page, page_size & sort": {
			queryString: "/v1/movies?title=moana&genres=action,adventure&page=1&page_size=10&sort=year",
			expectedResponse: `{
  "movies": [
    {
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
  ]
}`,
			expectedCodeStatus: 200,
		},
	}

	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			status, _, body := ts.get(t, tc.queryString)

			assert.Equal(t, status, tc.expectedCodeStatus)
			assert.Equal(t, body, tc.expectedResponse)
		})
	}
}
