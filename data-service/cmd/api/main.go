package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/lukas-blaha/quickle-web/data-service/data"
)

const webPort = "8888"

var counts int

type Config struct {
	DB     *sql.DB
	Models data.Models
}

func main() {
	conn := connectToDB()
	if conn == nil {
		log.Fatal("Can't connect to Postgres!")
	}

	app := Config{
		DB:     conn,
		Models: data.New(conn),
	}

	log.Printf("Starting data-service on port %s\n", webPort)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	err := srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func connectToDB() *sql.DB {
	dsn := os.Getenv("DSN")

	for {
		conn, err := openDB(dsn)
		if err != nil {
			log.Println("Postgres not yet ready...")
			counts++
		} else {
			log.Println("Connected to Postgres!")
			return conn
		}

		if counts > 10 {
			log.Println(err)
			return nil
		}

		log.Println("Backing off for two seconds...")
		time.Sleep(2 * time.Second)
	}
}
