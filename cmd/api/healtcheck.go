package main

import (
	"net/http"
)

func (app *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	// fmt.Fprintln(w, "status: available")
	// fmt.Fprintf(w, "environment: %s\n", app.config.env)
	// fmt.Fprintf(w, "version: %s\n", version)
	data := envelope{
		"status": "available",
		"status_info": map[string]string{
			"environment": app.config.env,
			"version":     version,
		},
	}

	err := app.writeJSON(w, http.StatusOK, data, nil)
	if err != nil {
		// app.logger.Error(err.Error())
		// http.Error(w, "The server encountered an error and could not process your request", http.StatusInternalServerError)
		app.serverErrorResponse(w, r, err)
	}
}
