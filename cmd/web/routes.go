package main

import "net/http"

func (a *application) routes() *http.ServeMux {
	mux := http.NewServeMux()

	// Define all routes below
	mux.HandleFunc("/{$}", a.home)
	mux.HandleFunc("/collection", a.collection)
	mux.HandleFunc("/history", a.history)
	mux.HandleFunc("/mixes", a.mixes)

	return mux
}
