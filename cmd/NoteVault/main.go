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
	"github.com/go-chi/cors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	cfg := config.Load()

	log := SetupLogger(cfg.Env)
	log.Info("Starting NoteVault at", "address", cfg.HTTPServer.Address)

	mongoClient, ctx := mongoConnect(cfg, log)
	defer mongoClient.Disconnect(ctx)

	noteCollection := mongoClient.Database(cfg.Database).Collection(cfg.Collections.Notes)
	noteBookCollection := mongoClient.Database(cfg.Database).Collection(cfg.Collections.NoteBooks)
	tagCollection := mongoClient.Database(cfg.Database).Collection(cfg.Collections.Tags)
	userCollection := mongoClient.Database(cfg.Database).Collection(cfg.Collections.Users)

	indexEmail := mongo.IndexModel{
		Keys:    bson.D{{Key: "email", Value: 1}},
		Options: options.Index().SetUnique(true),
	}

	_, err := userCollection.Indexes().CreateOne(context.Background(), indexEmail)
	if err != nil {
		log.Error("Failed to create unique index for email", slog.String("error", err.Error()))
	}

	noteService := service.NoteService{
		DBClient: mongodb.MongoClient{
			Client: *noteCollection,
		},
		HelperNoteBookClient: mongodb.MongoClient{
			Client: *noteBookCollection,
		},
		HelperTagClient: mongodb.MongoClient{
			Client: *tagCollection,
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
		HelperNoteClient: mongodb.MongoClient{
			Client: *noteCollection,
		},
	}

	userService := service.UserService{
		DBClient: mongodb.MongoClient{
			Client: *userCollection,
		},
	}

	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(middleware.SetHeader("CONTENT-TYPE", "application/json"))
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	router.Route("/api/v1", func(router chi.Router) {
		router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("NoteVault is OK!"))
		})

		router.Get("/notes/{id}", noteService.HandleGetNoteByID)
		router.Get("/notes", noteService.HandleGetNotes)
		router.Get("/notes/trash", noteService.HandleGetTrashedNotes)
		router.Get("/notes/archive", noteService.HandleGetArchivedNotes)
		router.Get("/notes/trash/{id}", noteService.HandleRestoreNoteFromTrash)
		router.Get("/notes/archive/{id}", noteService.HandleRestoreNoteFromArchive)
		router.Get("/notes/group/{id}", noteService.HandleGetNotesByNoteBookID)
		router.Post("/notes", noteService.HandleCreateNote)
		router.Post("/notes/tag", noteService.HandleGetNotesByTags)
		router.Put("/notes/{id}", noteService.HandleUpdateNote)
		router.Put("/notes/notebook/{id}", noteService.HandleUpdateNoteNoteBook)
		router.Put("/notes/tag/{id}", noteService.HandleAddTagToNote)
		router.Patch("/notes/tag/{id}", noteService.HandleRemoveTagFromNote)
		router.Delete("/notes/{id}", noteService.HandleDeleteNote)
		router.Delete("/notes/trash/{id}", noteService.HandleMoveNoteToTrash)
		router.Delete("/notes/archive/{id}", noteService.HandleMoveNoteToArchive)

		router.Get("/notebooks/{id}", noteBookService.HandleGetNoteBookByID)
		router.Get("/notebooks", noteBookService.HandleGetNoteBooks)
		router.Post("/notebooks", noteBookService.HandleCreateNoteBook)
		router.Put("/notebooks/{id}", noteBookService.HandleUpdateNoteBook)
		router.Delete("/notebooks/{id}", noteBookService.HandleDeleteNoteBook)

		router.Get("/tags/{id}", tagService.HandleGetTagByID)
		router.Get("/tags", tagService.HandleGetTags)
		router.Post("/tags", tagService.HandleCreateTag)
		router.Put("/tags/{id}", tagService.HandleUpdateTag)
		router.Delete("/tags/{id}", tagService.HandleDeleteTag)

		router.Post("/users/login", userService.HandleLoginUser)
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
