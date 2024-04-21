package main

import (
	_ "embed"
	"io"
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	slog.Info("starting server")
	webTemplateFile, err := os.OpenFile("templates/static/web.html", os.O_RDONLY, 0644)
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
	webTemplateBytes, err := io.ReadAll(webTemplateFile)
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write(webTemplateBytes)
	})

	r.Post("/wystaw", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("<p>OK</p><a href='/'>Powrót</a>"))
	})

	http.ListenAndServe(":8080", r)
}
