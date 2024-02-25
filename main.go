package main

import (
	"fmt"
	"github.com/mikelangelon/dutch-words/core"
	"github.com/mikelangelon/dutch-words/db"
	"github.com/mikelangelon/dutch-words/server"
	"github.com/mikelangelon/dutch-words/services"
	"log/slog"
)

func main() {
	// Setup
	store := db.NewStore()
	store.Insert(&core.Word{Dutch: "hond", English: "dog"})
	store.Insert(&core.Word{Dutch: "paard", English: "horse"})
	sv := server.New(services.NewService(store))

	fmt.Println("Server is listening on http://localhost:8080")

	err := sv.ListenAndServe()
	if err != nil {
		slog.Error("unexpected error on server", "error", err)
	}
}
