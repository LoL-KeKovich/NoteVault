package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/LoL-KeKovich/NoteVault/internal/config"
	"github.com/LoL-KeKovich/NoteVault/internal/repository/mongodb"
	"github.com/LoL-KeKovich/NoteVault/internal/service"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	cfg := config.Load()

	log := SetupLogger(cfg.Env)
	log.Info("Starting NoteVault at", "address", cfg.HTTPServer.Address)

	mongoClient, ctx := mongoConnect(cfg, log)
	defer mongoClient.Disconnect(ctx)

	noteCollection := mongoClient.Database("NoteVault").Collection("notes")

	noteService := service.NoteService{
		DBClient: mongodb.MongoClient{
			Client: *noteCollection,
		},
	}

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

func mongoConnect(cfg *config.Config, log *slog.Logger) (*mongo.Client, context.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(cfg.StoragePath)

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		panic(err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		panic(err)
	}

	log.Info("Successfully connected to mongo!")

	return client, ctx
}
