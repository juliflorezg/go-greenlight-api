package main

import (
	"encoding/json"
	"net/http"
)

// Declare a handler which writes a plain-text response with information about the
// application status, operating environment and version.
func (app *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	// fmt.Fprintln(w, "status: available")
	// fmt.Fprintf(w, "environment: %s\n", app.config.env)
	// fmt.Fprintf(w, "version: %s\n", version)
	data := map[string]string{
		"status":      "available",
		"environment": app.config.env,
		"version":     version,
	}

	js, err := json.Marshal(data)
	if err != nil {
		app.logger.Error(err.Error())
		http.Error(w, "The server encountered an error and could not process your request", http.StatusInternalServerError)
	}

	// js = append(js, '\n')

	w.Header().Set("Content-Type", "application/json")

	w.Write(js)

}
