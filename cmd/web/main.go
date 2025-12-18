package main

import (
	"context"
	"log"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	dsn := "postgres://serendipity:password@localhost:5432/serendipity"
	db, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		log.Fatal(err)
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
