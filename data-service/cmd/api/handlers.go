package main

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
)

var deck = map[string][]jsonResponse{
	"fruit": []jsonResponse{
		{ID: 1, Term: "jahoda", Definition: "strawberry"},
		{ID: 2, Term: "banan", Definition: "banana"},
		{ID: 3, Term: "jablko", Definition: "apple"},
		{ID: 4, Term: "broskev", Definition: "peach"},
		{ID: 5, Term: "malina", Definition: "raspberry"},
	},
	"vegetables": []jsonResponse{
		{ID: 1, Term: "rajce", Definition: "tomato"},
		{ID: 2, Term: "brokolice", Definition: "brokoli"},
		{ID: 3, Term: "okurek", Definition: "cucumber"},
		{ID: 4, Term: "paprika", Definition: "paprika"},
		{ID: 5, Term: "cibule", Definition: "onion"},
	},
}

func (app *Config) ListDecks(w http.ResponseWriter, t *http.Request) {
	var resp jsonResponse
	for k, _ := range deck {
		resp.Name = k
		app.writeJSON(w, http.StatusOK, resp)
	}
}

func (app *Config) GetDeck(w http.ResponseWriter, r *http.Request) {
	deckParam := chi.URLParam(r, "deck")
	idParam := r.URL.Query().Get("id")
	termParam := r.URL.Query().Get("term")

	if idParam != "" {
		id, err := strconv.Atoi(idParam)
		if err != nil {
			app.errorJSON(w, errors.New("Can't parse id to int"))
			return
		}

		if len(deck[deckParam]) < id-1 || id < 0 {
			app.errorJSON(w, errors.New("No item matching requested id..."), http.StatusNotFound)
			return
		}

		app.writeJSON(w, http.StatusOK, deck[deckParam][id-1])
	} else if termParam != "" {
		for _, item := range deck[deckParam] {
			if strings.ToLower(item.Term) == strings.ToLower(termParam) {
				app.writeJSON(w, http.StatusOK, item)
				return
			}
		}

		app.errorJSON(w, errors.New("No item matching requested term..."), http.StatusNotFound)
	} else {
		for _, item := range deck[deckParam] {
			app.writeJSON(w, http.StatusOK, item)
		}
	}
}

func (app *Config) UpdateItem(w http.ResponseWriter, r *http.Request) {
	deckParam := chi.URLParam(r, "deck")
	idParam := chi.URLParam(r, "id")

	var item jsonResponse

	id, err := strconv.Atoi(idParam)
	if err != nil {
		app.errorJSON(w, errors.New("Can't parse id to int"))
		return
	}

	err = app.readJSON(w, r, &item)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	if item.Term != "" {
		deck[deckParam][id-1].Term = item.Term
	}

	if item.Definition != "" {
		deck[deckParam][id-1].Definition = item.Definition
	}
}
