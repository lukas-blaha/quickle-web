package main

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func (app *Config) ListDecks(w http.ResponseWriter, t *http.Request) {
	items, err := app.Models.Item.GetDecks()
	if err != nil {
		app.errorJSON(w, err)
	}

	app.writeJSON(w, http.StatusOK, items)
}

func (app *Config) DeckRange(w http.ResponseWriter, r *http.Request) {
	deckParam := chi.URLParam(r, "deck")

	item, err := app.Models.Item.GetDeckRange(deckParam)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	app.writeJSON(w, http.StatusOK, item)
}

func (app *Config) GetItems(w http.ResponseWriter, r *http.Request) {
	deckParam := chi.URLParam(r, "deck")
	idParam := r.URL.Query().Get("id")
	termParam := r.URL.Query().Get("term")

	if idParam != "" {
		_, err := app.checkID(w, idParam)
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

		app.writeJSON(w, http.StatusOK, items)
	}
}

func (app *Config) UpdateItem(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")

	var item jsonResponse

	_, err := app.checkID(w, idParam)
	if err != nil {
		return
	}

	err = app.readJSON(w, r, &item)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	err = app.Models.Item.Update(item.Definition, item.Term, idParam)
	if err != nil {
		app.errorJSON(w, err)
		return
	}
}

func (app *Config) RemoveItem(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")

	_, err := app.checkID(w, idParam)
	if err != nil {
		return
	}

	err = app.Models.Item.Delete(idParam)
	if err != nil {
		app.errorJSON(w, err)
		return
	}
}

func (app *Config) RemoveDeck(w http.ResponseWriter, r *http.Request) {
	deckParam := chi.URLParam(r, "deck")

	err := app.Models.Item.DeleteDeck(deckParam)
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

func (app *Config) checkID(w http.ResponseWriter, idParam string) (int, error) {
	id, err := strconv.Atoi(idParam)
	if err != nil {
		app.errorJSON(w, errors.New("Can't parse id to int"))
		return 0, err
	}

	return id, nil
}
