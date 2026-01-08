package main

import (
	"net/http"
	"text/template"
)

func (a *application) home(w http.ResponseWriter, r *http.Request) {
	files := []string{
		"./ui/html/layout.html",
		"./ui/html/partials/navbar.html",
		"./ui/html/pages/home.html",
		"./ui/html/partials/player.html",
	}

	template, err := template.ParseFiles(files...)
	if err != nil {
		a.handleInternalError(w, err)
	}

	err = template.ExecuteTemplate(w, "layout", nil)
	if err != nil {
		a.handleInternalError(w, err)
	}
}

func (a *application) collection(w http.ResponseWriter, r *http.Request) {
	files := []string{
		"./ui/html/layout.html",
		"./ui/html/partials/navbar.html",
		"./ui/html/pages/collection.html",
	}

	template, err := template.ParseFiles(files...)
	if err != nil {
		a.handleInternalError(w, err)
	}

	err = template.ExecuteTemplate(w, "layout", nil)
	if err != nil {
		a.handleInternalError(w, err)
	}
}

func (a *application) history(w http.ResponseWriter, r *http.Request) {
	files := []string{
		"./ui/html/layout.html",
		"./ui/html/partials/navbar.html",
		"./ui/html/pages/history.html",
	}

	template, err := template.ParseFiles(files...)
	if err != nil {
		a.handleInternalError(w, err)
	}

	err = template.ExecuteTemplate(w, "layout", nil)
	if err != nil {
		a.handleInternalError(w, err)
	}
}

func (a *application) mixes(w http.ResponseWriter, r *http.Request) {
	files := []string{
		"./ui/html/layout.html",
		"./ui/html/partials/navbar.html",
		"./ui/html/pages/mixes.html",
	}

	template, err := template.ParseFiles(files...)
	if err != nil {
		a.handleInternalError(w, err)
	}

	err = template.ExecuteTemplate(w, "layout", nil)
	if err != nil {
		a.handleInternalError(w, err)
	}
}

func (a *application) handleInternalError(w http.ResponseWriter, err error) {
	a.logger.Error(err.Error())
	http.Error(w, err.Error(), http.StatusInternalServerError)
}
