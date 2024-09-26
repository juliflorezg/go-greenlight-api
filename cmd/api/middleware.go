package main

import (
	"fmt"
	"net/http"
)

func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")

				app.serverErrorResponse(w, r, fmt.Errorf("%s", err))
			}
		}()

		next.ServeHTTP(w, r)
	})

}

func (app *application) enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// todo: restrict when app gets to prod
		// w.Header().Set("Access-Control-Allow-Origin", "*")

		w.Header().Set("Vary", "Origin")

		origin := w.Header().Get("Origin")

		if origin != "" {
			for _, v := range app.config.cors.trustedOrigins {
				if origin == v {
					w.Header().Set("Access-Control-Allow-Origin", v)
					break
				}
			}
		}

		next.ServeHTTP(w, r)
	})
}
