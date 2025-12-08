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
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	err = template.ExecuteTemplate(w, "layout", nil)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	err = template.ExecuteTemplate(w, "layout", nil)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
