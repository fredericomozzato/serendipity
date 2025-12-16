package main

import (
	"database/sql"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	dsn := "web:password@tcp(localhost:3306)/serendipity?parseTime=true"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
		db.Close()
	}

	defer db.Close()

	mux := http.NewServeMux()
	mux.HandleFunc("/{$}", home)
	mux.HandleFunc("/collection", collection)
	mux.HandleFunc("/history", history)
	mux.HandleFunc("/mixes", mixes)

	log.Println("starting server on :4000")
	err = http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}
