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

	noteCollection := mongoClient.Database("NoteVault").Collection("notes")         //Захардкожено
	noteBookCollection := mongoClient.Database("NoteVault").Collection("notebooks") //Захардкожено
	tagCollection := mongoClient.Database("NoteVault").Collection("tags")           //Захардкожено

	noteService := service.NoteService{
		DBClient: mongodb.MongoClient{
			Client: *noteCollection,
		},
		HelperNoteBookClient: mongodb.MongoClient{
			Client: *noteBookCollection,
		},
	}

	noteBookService := service.NoteBookService{
		DBClient: mongodb.MongoClient{
			Client: *noteBookCollection,
		},
		HelperNoteClient: mongodb.MongoClient{
			Client: *noteCollection,
		},
	}

	tagService := service.TagService{
		DBClient: mongodb.MongoClient{
			Client: *tagCollection,
		},
	}

	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(middleware.SetHeader("CONTENT-TYPE", "application/json"))

	router.Route("/api/v1", func(router chi.Router) {
		router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("NoteVault is OK!"))
		})

		router.Get("/notes/{id}", noteService.HandleGetNoteByID)
		router.Get("/notes", noteService.HandleGetNotes)
		router.Get("/notes/trash", noteService.HandleGetTrashedNotes)
		router.Get("/notes/trash/{id}", noteService.HandleRestoreNoteFromTrash)
		router.Get("/notes/group/{id}", noteService.HandleGetNotesByNoteBookID)
		router.Post("/notes", noteService.HandleCreateNote)
		router.Put("/notes/{id}", noteService.HandleUpdateNote)
		router.Put("/notes/notebook/{id}", noteService.HandleUpdateNoteNoteBook)
		router.Delete("/notes/{id}", noteService.HandleDeleteNote)
		router.Delete("/notes/trash/{id}", noteService.HandleMoveNoteToTrash)

		router.Get("/notebooks/{id}", noteBookService.HandleGetNoteBookByID)
		router.Get("/notebooks", noteBookService.HandleGetNoteBooks)
		router.Post("/notebooks", noteBookService.HandleCreateNoteBook)
		router.Put("/notebooks/{id}", noteBookService.HandleUpdateNoteBook)
		router.Delete("/notebooks/{id}", noteBookService.HandleDeleteNoteBook)

		router.Get("/tags/{id}", tagService.HandleGetTagByID)
		router.Get("/tags", tagService.HandleGetTags)
		router.Post("/tags", tagService.HandleCreateTag)
		router.Put("/tags", tagService.HandleUpdateTag)
		router.Delete("/tags", tagService.HandleDeleteTag)
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
