package main

import (
	"log/slog"
	"net/http"
	"runtime/debug"
)

func (a *application) serverError(w http.ResponseWriter, r *http.Request, err error) {
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
