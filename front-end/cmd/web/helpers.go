package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"text/template"
)

//go:embed templates
var templateFS embed.FS

func (app *Config) getDecks() ([]jsonResponse, error) {
	decksURL := fmt.Sprintf("%s/decks", apiURL)

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
	deckURL := fmt.Sprintf("%s/deck/%s", apiURL, deckParam)

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
		"templates/base.layout.gohtml",
		"templates/header.partial.gohtml",
		"templates/footer.partial.gohtml",
	}

	var templateSlice []string
	templateSlice = append(templateSlice, fmt.Sprintf("templates/%s", t))

	for _, x := range partials {
		templateSlice = append(templateSlice, x)
	}

	tmpl, err := template.ParseFS(templateFS, templateSlice...)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return nil
	}

	return tmpl
}

func (app *Config) checkDeck(deck string) bool {
	deckItems, err := app.getDecks()
	if err != nil {
		log.Fatal(err)
	}
	for _, deck := range deckItems {
		app.Decks = append(app.Decks, deck.Deck)
	}

	for _, d := range app.Decks {
		if deck == d {
			return true
		}
	}

	return false
}
