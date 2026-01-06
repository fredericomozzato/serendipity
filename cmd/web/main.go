package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	addr := flag.String("addr", "4000", "Network port to access the app. Defaults to 4000")
	dsn := flag.String("dsn", "", "Data Source Name")

	flag.Parse()

	serverAddr := ":" + *addr

	db, err := pgxpool.New(context.Background(), *dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	mux := http.NewServeMux()
	mux.HandleFunc("/{$}", home)
	mux.HandleFunc("/collection", collection)
	mux.HandleFunc("/history", history)
	mux.HandleFunc("/mixes", mixes)

	log.Println(fmt.Sprintf("starting server at %s", serverAddr))
	err = http.ListenAndServe(serverAddr, mux)
	log.Fatal(err)
}
