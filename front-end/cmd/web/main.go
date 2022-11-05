package main

import (
	"fmt"
	"log"
	"net/http"
)

const webPort = "80"

type jsonResponse struct {
	ID         int    `json:"id"`
	Deck       string `json:"deck"`
	First      int    `json:"first"`
	Last       int    `json:"last"`
	Term       string `json:"term"`
	Definition string `json:"definition"`
}

type Config struct {
	Decks []string
}

func main() {
	app := Config{}

	fmt.Printf("Starting front-end service on port %s\n", webPort)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	err := srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
