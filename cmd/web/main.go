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

	app := &application{
		logger: logger,
	}

	db, err := pgxpool.New(context.Background(), *dsn)
	if err != nil {
		app.serverError(err)
	}
	defer db.Close()

	// TODO: set the DB pool in the app struct

	app.logger.Info("starting server", slog.String("port", serverAddr))
	err = http.ListenAndServe(serverAddr, app.routes())
	app.serverError(err)
}
