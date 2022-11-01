package main

import (
	"errors"
	"net/http"
	"strconv"

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
	items, err := app.Models.Item.GetTables()
	if err != nil {
		app.errorJSON(w, err)
	}

	for _, i := range items {
		app.writeJSON(w, http.StatusOK, i)
	}
}

func (app *Config) GetItems(w http.ResponseWriter, r *http.Request) {
	deckParam := chi.URLParam(r, "deck")
	idParam := r.URL.Query().Get("id")
	termParam := r.URL.Query().Get("term")

	if idParam != "" {
		_, err := app.checkID(w, deckParam, idParam)
		if err != nil {
			return
		}

		item, err := app.Models.Item.GetByID(deckParam, idParam)
		if err != nil {
			app.errorJSON(w, err)
			return
		}

		app.writeJSON(w, http.StatusOK, item)
	} else if termParam != "" {
		item, err := app.Models.Item.GetByTerm(deckParam, termParam)
		if err != nil {
			app.errorJSON(w, err)
			return
		}

		app.writeJSON(w, http.StatusOK, item)
	} else {
		items, err := app.Models.Item.GetAll(deckParam)
		if err != nil {
			app.errorJSON(w, err)
			return
		}

		for _, i := range items {
			app.writeJSON(w, http.StatusOK, i)
		}
	}
}

func (app *Config) UpdateItem(w http.ResponseWriter, r *http.Request) {
	deckParam := chi.URLParam(r, "deck")
	idParam := chi.URLParam(r, "id")

	var item jsonResponse

	_, err := app.checkID(w, deckParam, idParam)
	if err != nil {
		return
	}

	err = app.readJSON(w, r, &item)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	err = app.Models.Item.Update(deckParam, item.Definition, item.Term, idParam)
	if err != nil {
		app.errorJSON(w, err)
		return
	}
}

func (app *Config) RemoveItem(w http.ResponseWriter, r *http.Request) {
	deckParam := chi.URLParam(r, "deck")
	idParam := chi.URLParam(r, "id")

	_, err := app.checkID(w, deckParam, idParam)
	if err != nil {
		return
	}

	err = app.Models.Item.Delete(deckParam, idParam)
	if err != nil {
		app.errorJSON(w, err)
		return
	}
}

func (app *Config) AddItem(w http.ResponseWriter, r *http.Request) {
	deckParam := chi.URLParam(r, "deck")

	var item jsonResponse

	err := app.readJSON(w, r, &item)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	_, err = app.Models.Item.Insert(deckParam, item.Term, item.Definition)
	if err != nil {
		app.errorJSON(w, err)
		return
	}
}

func (app *Config) checkID(w http.ResponseWriter, deckParam, idParam string) (int, error) {
	id, err := strconv.Atoi(idParam)
	if err != nil {
		app.errorJSON(w, errors.New("Can't parse id to int"))
		return 0, err
	}

	// if len(deck[deckParam]) < id-1 || id < 0 {
	// 	err = errors.New("No item matching requested id...")
	// 	app.errorJSON(w, err, http.StatusNotFound)
	// 	return 0, err
	// }

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
