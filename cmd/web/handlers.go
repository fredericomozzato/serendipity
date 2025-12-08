package main

import (
	"log"
	"net/http"
	"text/template"
)

func home(w http.ResponseWriter, r *http.Request) {
	files := []string{
		"./ui/html/layout.html",
		"./ui/html/partials/navbar.html",
		"./ui/html/pages/home.html",
		"./ui/html/partials/player.html",
	}

	template, err := template.ParseFiles(files...)
	if err != nil {
		handleInternalError(w, err)
	}

	err = template.ExecuteTemplate(w, "layout", nil)
	if err != nil {
		handleInternalError(w, err)
	}
}

func collection(w http.ResponseWriter, r *http.Request) {
	files := []string{
		"./ui/html/layout.html",
		"./ui/html/partials/navbar.html",
		"./ui/html/pages/collection.html",
	}

	template, err := template.ParseFiles(files...)
	if err != nil {
		handleInternalError(w, err)
	}

	err = template.ExecuteTemplate(w, "layout", nil)
	if err != nil {
		handleInternalError(w, err)
	}
}

func history(w http.ResponseWriter, r *http.Request) {
	files := []string{
		"./ui/html/layout.html",
		"./ui/html/partials/navbar.html",
		"./ui/html/pages/history.html",
	}

	template, err := template.ParseFiles(files...)
	if err != nil {
		handleInternalError(w, err)
	}

	err = template.ExecuteTemplate(w, "layout", nil)
	if err != nil {
		handleInternalError(w, err)
	}
}

func handleInternalError(w http.ResponseWriter, err error) {
	log.Println(err)
	http.Error(w, err.Error(), http.StatusInternalServerError)
}
