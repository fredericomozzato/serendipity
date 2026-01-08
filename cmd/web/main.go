package main

import (
	"context"
	"flag"
	"log/slog"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

type application struct {
	logger *slog.Logger
}

func main() {
	port := flag.String("port", "4000", "Network port to access the app. Defaults to 4000")
	dsn := flag.String("dsn", "", "Data Source Name")

	flag.Parse()

	serverAddr := ":" + *port

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		// keeping the default options explicit
		Level:     slog.LevelInfo,
		AddSource: false,
	}))

	db, err := pgxpool.New(context.Background(), *dsn)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	defer db.Close()

	app := &application{
		logger: logger,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/{$}", app.home)
	mux.HandleFunc("/collection", app.collection)
	mux.HandleFunc("/history", app.history)
	mux.HandleFunc("/mixes", app.mixes)

	app.logger.Info("starting server", slog.String("port", serverAddr))
	err = http.ListenAndServe(serverAddr, mux)
	app.logger.Error(err.Error())
	os.Exit(1)
}
