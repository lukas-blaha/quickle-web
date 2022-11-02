package data

import (
	"context"
	"database/sql"
	"log"
	"time"
)

const dbTimeout = 3 * time.Second

var db *sql.DB

// Models contains all types we want to be available to our application
type Models struct {
	Item Item
}

func New(dbPool *sql.DB) Models {
	db = dbPool

	return Models{
		Item: Item{},
	}
}

type Item struct {
	ID         int    `json:"id,omitempty"`
	Deck       string `json:"deck,omitempty"`
	Term       string `json:"term,omitempty"`
	Definition string `json:"definition,omitempty"`
}

func (i *Item) GetDecks() ([]*Item, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `select distinct deck from quickle`

	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*Item

	for rows.Next() {
		var item Item
		err = rows.Scan(&item.Deck)
		if err != nil {
			log.Println("Error scanning", err)
			return nil, err
		}

		items = append(items, &item)
	}

	return items, nil
}

func (i *Item) GetAll(deck string) ([]*Item, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt := `select id, deck, term, definition from quickle where deck = $1 order by id`

	rows, err := db.QueryContext(ctx, stmt, deck)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*Item

	for rows.Next() {
		var item Item
		err = rows.Scan(
			&item.ID,
			&item.Deck,
			&item.Term,
			&item.Definition,
		)

		if err != nil {
			log.Println("Error scanning", err)
			return nil, err
		}

		items = append(items, &item)
	}

	return items, nil
}

func (i *Item) GetByID(deck, id string) (*Item, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `select id, deck, term, definition from quickle where deck = $1 and id = $2`

	var item Item
	row := db.QueryRowContext(ctx, query, deck, id)
	err := row.Scan(
		&item.ID,
		&item.Deck,
		&item.Term,
		&item.Definition,
	)

	if err != nil {
		return nil, err
	}

	return &item, nil
}

func (i *Item) GetByTerm(deck, term string) (*Item, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `select id, deck, term, definition from quickle where deck = $1 and term = $2`

	var item Item
	row := db.QueryRowContext(ctx, query, deck, term)
	err := row.Scan(
		&item.ID,
		&item.Deck,
		&item.Term,
		&item.Definition,
	)

	if err != nil {
		return nil, err
	}

	return &item, nil
}

func (i *Item) Update(def, term, id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	var stmt string

	if def != "" && term != "" {
		stmt = `update quickle set term = $2, definition = $3 where id = $1`
	} else if term != "" {
		stmt = `update quickle set term = $2 where id = $1`
	} else if def != "" {
		stmt = `update quickle set definition = $3 where id = $1`
	}

	_, err := db.ExecContext(ctx, stmt, id, term, def)
	if err != nil {
		return err
	}

	return nil
}

func (i *Item) Delete(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt := `delete from quickle where id = $1`

	_, err := db.ExecContext(ctx, stmt, id)
	if err != nil {
		return err
	}

	return nil
}

func (i *Item) DeleteDeck(deck string) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt := `delete from quickle where deck = $1`

	_, err := db.ExecContext(ctx, stmt, deck)
	if err != nil {
		return err
	}

	return nil
}

func (i *Item) Insert(deck, term, def string) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	var newID int
	stmt := `insert into quickle (deck, term, definition) values ($1, $2, $3) returning id`

	err := db.QueryRowContext(ctx, stmt, deck, term, def).Scan(&newID)
	if err != nil {
		return 0, err
	}

	return newID, nil
}
