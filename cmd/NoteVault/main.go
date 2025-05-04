package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/LoL-KeKovich/NoteVault/internal/config"
	"github.com/LoL-KeKovich/NoteVault/internal/service"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func main() {
	cfg := config.Load()

	log := SetupLogger(cfg.Env)
	log.Info("Starting NoteVault at", "address", cfg.HTTPServer.Address)

	noteService := service.NoteService{}

	router := chi.NewRouter()
	router.Use(middleware.Logger)

	router.Route("/api/v1", func(router chi.Router) {
		router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("NoteVault is OK!"))
		})

		router.Get("/notes/{id}", noteService.GetNoteByID)
		router.Get("/notes", noteService.GetNotes)
		router.Post("/notes", noteService.CreateNote)
		router.Put("/notes/{id}", noteService.UpdateNote)
		router.Delete("/notes/{id}", noteService.DeleteNote)
	})

	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Error("Failed to start NoteVault!!!")
	}
}

func SetupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case "local":
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case "dev":
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case "prod":
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	default:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}
