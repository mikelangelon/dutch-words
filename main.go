package main

import (
	"fmt"
	"log/slog"

	"github.com/mikelangelon/dutch-words/config"
	"github.com/mikelangelon/dutch-words/db"
	"github.com/mikelangelon/dutch-words/server"
	"github.com/mikelangelon/dutch-words/services"
)

func main() {
	// Parse config
	cfg, err := config.Parse()
	if err != nil {
		slog.Error("problem parsing dependencies", "error", err)
	}
	// Setup dependencies
	mongoStore, err := db.NewMongoStore(cfg)
	if err != nil {
		slog.Error("problem parsing dependencies", "error", err)
	}

	sv := server.New(
		services.NewService(mongoStore),
		services.NewSentencesService(mongoStore),
		services.NewGameService(mongoStore),
		services.NewScoreService(mongoStore),
	)

	fmt.Println("Server is listening on http://localhost:8080")

	err = sv.ListenAndServe()
	if err != nil {
		slog.Error("unexpected error on server", "error", err)
	}
}
