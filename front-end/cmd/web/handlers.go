package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

var apiURL = "https://quickle-api.fikos.cz"

func (app *Config) ListDecks(w http.ResponseWriter, r *http.Request) {
	tmpl := app.render(w, "main.page.gohtml")

	items, err := app.getDecks()
	if err != nil {
		log.Println(err)
	}

	if err := tmpl.Execute(w, items); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (app *Config) ListItems(w http.ResponseWriter, r *http.Request) {
	deckParam := chi.URLParam(r, "deck")

	if !app.checkDeck(deckParam) {
		http.Error(w, "Study set does not exist!", http.StatusBadRequest)
		return
	}

	tmpl := app.render(w, "deck.page.gohtml")

	items, err := app.getItems(deckParam)
	if err != nil {
		log.Println(err)
	}

	_ = tmpl.Execute(w, items)
}

func (app *Config) EditItems(w http.ResponseWriter, r *http.Request) {
	deckParam := chi.URLParam(r, "deck")

	if !app.checkDeck(deckParam) {
		http.Error(w, "Study set does not exist!", http.StatusBadRequest)
		return
	}

	tmpl := app.render(w, "edit.page.gohtml")

	items, err := app.getItems(deckParam)
	if err != nil {
		log.Println(err)
	}

	_ = tmpl.Execute(w, items)
}

func (app *Config) Cards(w http.ResponseWriter, r *http.Request) {
	deckParam := chi.URLParam(r, "deck")

	if !app.checkDeck(deckParam) {
		http.Error(w, "Study set does not exist!", http.StatusBadRequest)
		return
	}

	tmpl := app.render(w, "cards.page.gohtml")

	if err := tmpl.Execute(w, deckParam); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (app *Config) WriteMode(w http.ResponseWriter, r *http.Request) {
	deckParam := chi.URLParam(r, "deck")

	if !app.checkDeck(deckParam) {
		http.Error(w, "Study set does not exist!", http.StatusBadRequest)
		return
	}

	tmpl := app.render(w, "write.page.gohtml")

	if err := tmpl.Execute(w, deckParam); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
