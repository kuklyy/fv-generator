package main

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func main() {
	slog.Info("starting server")
	r := chi.NewRouter()

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello from fv server"))
	})

	http.ListenAndServe(":8080", r)
}
