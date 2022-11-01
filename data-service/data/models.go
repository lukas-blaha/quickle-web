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
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Term       string `json:"term"`
	Definition string `json:"definition"`
}

func (i *Item) GetAll(deck string) ([]*Item, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `select id, term, definition from $1 order by id`

	rows, err := db.QueryContext(ctx, query, deck)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*Item

	for rows.Next() {
		var item Item
		err = rows.Scan(
			&item.ID,
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

	query := `select id, term, definition from $1 where id = $2`

	var item Item
	row := db.QueryRowContext(ctx, query, deck, id)
	err := row.Scan(
		&item.ID,
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

	query := `select id, term, definition from $1 where term = $2`

	var item Item
	row := db.QueryRowContext(ctx, query, deck, term)
	err := row.Scan(
		&item.ID,
		&item.Term,
		&item.Definition,
	)

	if err != nil {
		return nil, err
	}

	return &item, nil
}

func (i *Item) Update(deck, def, term, id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	var stmt string

	if def != "" && term != "" {
		stmt = `update $1 set term = $3, definition = $4 where id = $2`
	} else if term != "" {
		stmt = `update $1 set term = $3 where id = $2`
	} else if def != "" {
		stmt = `update $1 set definition = $4 where id = $2`
	}

	_, err := db.ExecContext(ctx, stmt, deck, id, term, def)
	if err != nil {
		return err
	}

	return nil
}

func (i *Item) Delete(deck, id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt := `delete from $1 where id = $2`

	_, err := db.ExecContext(ctx, stmt, deck, id)
	if err != nil {
		return err
	}

	return nil
}

func (i *Item) Insert(deck, term, def string) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	var newID int
	stmt := `insert into $1 (term, definition)
		values ($2, $3) returning id`

	err := db.QueryRowContext(ctx, stmt, deck, term, def).Scan(&newID)
	if err != nil {
		return 0, err
	}

	return newID, nil
}

func (i *Item) GetTables() ([]*Item, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `select tablename from pg_catalog.pg_tables where schemaname = 'public'`

	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*Item

	for rows.Next() {
		var item Item
		err = rows.Scan(&item.Name)
		if err != nil {
			log.Println("Error scanning", err)
			return nil, err
		}

		items = append(items, &item)
	}

	return items, nil
}
