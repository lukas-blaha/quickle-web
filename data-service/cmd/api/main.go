package main

import (
	"fmt"
	"log"
	"net/http"
)

const webPort = "8888"

type Config struct{}

func main() {
	app := Config{}

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	log.Printf("Starting data-service on port %s\n", webPort)
	err := srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
