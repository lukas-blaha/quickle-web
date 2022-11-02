package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func (app *Config) routes() http.Handler {
	mux := chi.NewRouter()

	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://*", "https://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	mux.Use(middleware.Heartbeat("/ping"))
	mux.Get("/decks", app.ListDecks)
	mux.Get("/deck/{deck}", app.GetItems)
	mux.Post("/deck/{deck}", app.AddItem)
	mux.Patch("/deck/{deck}/{id}", app.UpdateItem)
	mux.Delete("/deck/{deck}/{id}", app.RemoveItem)
	mux.Delete("/deck/{deck}", app.RemoveDeck)

	return mux
}
