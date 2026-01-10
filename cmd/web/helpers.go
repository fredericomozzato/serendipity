package main

import (
	"log/slog"
	"net/http"
	"os"
	"runtime/debug"
)

func (a *application) serverError(err error) {
	a.logger.Error(err.Error(), slog.String("trace", string(debug.Stack())))
	os.Exit(1)
}

func (a *application) internalServerError(w http.ResponseWriter, r *http.Request, err error) {
	a.logger.Error(
		err.Error(),
		slog.String("method", r.Method),
		slog.String("url", r.URL.RequestURI()),
		slog.String("trace", string(debug.Stack())),
	)

	http.Error(
		w,
		http.StatusText(http.StatusInternalServerError),
		http.StatusInternalServerError,
	)
}

// TODO: implement clientError to handle user input errors
