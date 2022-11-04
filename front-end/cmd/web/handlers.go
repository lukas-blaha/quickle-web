package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"text/template"

	"github.com/go-chi/chi/v5"
)

var baseApiURL = "http://data-service:8888"

func (app *Config) getDecks() ([]jsonResponse, error) {
	decksURL := fmt.Sprintf("%s/decks", baseApiURL)

	req, err := http.Get(decksURL)
	if err != nil {
		return nil, err
	}
	defer req.Body.Close()

	var items []jsonResponse

	body, err := io.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}

	json.Unmarshal(body, &items)

	return items, err
}

func (app *Config) getItems(deckParam string) ([]jsonResponse, error) {
	deckURL := fmt.Sprintf("%s/deck/%s", baseApiURL, deckParam)

	req, err := http.Get(deckURL)
	if err != nil {
		return nil, err
	}
	defer req.Body.Close()

	var items []jsonResponse

	body, err := io.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}

	json.Unmarshal(body, &items)

	return items, err
}

func (app *Config) render(w http.ResponseWriter, t string) *template.Template {
	partials := []string{
		"./cmd/web/templates/base.layout.gohtml",
		"./cmd/web/templates/header.partial.gohtml",
		"./cmd/web/templates/footer.partial.gohtml",
	}

	var templateSlice []string
	templateSlice = append(templateSlice, fmt.Sprintf("./cmd/web/templates/%s", t))

	for _, x := range partials {
		templateSlice = append(templateSlice, x)
	}

	tmpl, err := template.ParseFiles(templateSlice...)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return nil
	}

	return tmpl
}

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

	tmpl := app.render(w, "deck.page.gohtml")

	items, err := app.getItems(deckParam)
	if err != nil {
		log.Println(err)
	}

	if err := tmpl.Execute(w, items); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (app *Config) Cards(w http.ResponseWriter, r *http.Request) {
	deckParam := chi.URLParam(r, "deck")

	tmpl := app.render(w, "cards.page.gohtml")

	if err := tmpl.Execute(w, deckParam); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (app *Config) WriteMode(w http.ResponseWriter, r *http.Request) {
	deckParam := chi.URLParam(r, "deck")

	tmpl := app.render(w, "write.page.gohtml")

	if err := tmpl.Execute(w, deckParam); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
