package main

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	_ "github.com/lib/pq"
)

func main() {
	// Open a connection to the database.
	db, err := openDatabase()

	if err != nil {
		panic(err)
	}

	defer db.Close()

	// Create the runs table and insert a new row to track the run.
	if _, err := db.Exec(`
CREATE TABLE IF NOT EXISTS runs (
	id SERIAL PRIMARY KEY,
	date TIMESTAMP WITHOUT TIME ZONE NOT NULL
);

INSERT INTO runs(date) VALUES (now());
	`); err != nil {
		panic(err)
	}

	serve(":8080", func(w http.ResponseWriter, r *http.Request) {
		var runsCount uint

		if err := db.
			QueryRowContext(r.Context(), "SELECT COUNT(id) FROM runs").
			Scan(&runsCount); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		fmt.Fprintf(w, "Total runs: %d", runsCount)
	})
}

func openDatabase() (*sql.DB, error) {
	// Get the DSN from the environment variable.
	dsn := os.Getenv("DATABASE_URL")

	// Open a connection to the database.
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	// Ping the database to verify the connection.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		panic(err)
	}

	return db, nil
}

func serve(addr string, handler http.HandlerFunc) {
	// Create a channel to receive an interrupt signal.
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	server := &http.Server{Addr: addr, Handler: handler}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()

	// Wait for an interrupt signal.
	<-stop

	// Shutdown the server gracefully.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		panic(err)
	}
}
