package main

import (
	"log"
	"net/http"
	"text/template"
)

func home(w http.ResponseWriter, r *http.Request) {
	// load the template. The relative path is from the projects root
	template, err := template.ParseFiles("./ui/html/pages/home.html")
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	// execute the template
	err = template.Execute(w, nil)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
