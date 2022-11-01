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

func (app *Config) GetItems(w http.ResponseWriter, r *http.Request) {
	deckParam := chi.URLParam(r, "deck")
	idParam := r.URL.Query().Get("id")
	termParam := r.URL.Query().Get("term")

	err := app.checkDeck(w, deckParam)
	if err != nil {
		return
	}

	if idParam != "" {
		id, err := app.checkID(w, deckParam, idParam)
		if err != nil {
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

	err := app.checkDeck(w, deckParam)
	if err != nil {
		return
	}

	id, err := app.checkID(w, deckParam, idParam)
	if err != nil {
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

func (app *Config) RemoveItem(w http.ResponseWriter, r *http.Request) {
	deckParam := chi.URLParam(r, "deck")
	idParam := chi.URLParam(r, "id")

	err := app.checkDeck(w, deckParam)
	if err != nil {
		return
	}

	id, err := app.checkID(w, deckParam, idParam)
	if err != nil {
		return
	}

	deck[deckParam] = append(deck[deckParam][:id-1], deck[deckParam][id:]...)
}

func (app *Config) RemoveDeck(w http.ResponseWriter, r *http.Request) {
	deckParam := chi.URLParam(r, "deck")

	err := app.checkDeck(w, deckParam)
	if err != nil {
		return
	}

	delete(deck, deckParam)
}

func (app *Config) AddItem(w http.ResponseWriter, r *http.Request) {
	deckParam := chi.URLParam(r, "deck")

	var item jsonResponse

	err := app.checkDeck(w, deckParam)
	if err != nil {
		return
	}

	err = app.readJSON(w, r, &item)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	deck[deckParam] = append(deck[deckParam], item)
}

func (app *Config) AddDeck(w http.ResponseWriter, r *http.Request) {
	deckParam := chi.URLParam(r, "deck")

	_, exists := deck[deckParam]
	if exists {
		app.errorJSON(w, errors.New("Deck already exists..."))
		return
	}

	if deckParam != "" {
		deck[deckParam] = []jsonResponse{}
	}
}

func (app *Config) checkID(w http.ResponseWriter, deckParam, idParam string) (int, error) {
	id, err := strconv.Atoi(idParam)
	if err != nil {
		app.errorJSON(w, errors.New("Can't parse id to int"))
		return 0, err
	}

	if len(deck[deckParam]) < id-1 || id < 0 {
		err = errors.New("No item matching requested id...")
		app.errorJSON(w, err, http.StatusNotFound)
		return 0, err
	}

	return id, nil
}

func (app *Config) checkDeck(w http.ResponseWriter, deckParam string) error {
	for d, _ := range deck {
		if d == deckParam {
			return nil
		}
	}

	err := errors.New("Deck not found...")
	app.errorJSON(w, err, http.StatusNotFound)
	return err
}
