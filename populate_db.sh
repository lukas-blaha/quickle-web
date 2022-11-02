#!/bin/bash

# Create tables if not exists
docker exec -ti quickle-postgres-1 psql -U postgres -d quickle -c "create table if not exists quickle(id serial primary key, deck text, term text, definition text)" >/dev/null 2>&1


# Populate DB with example data
echo "Sending example data to db..."
while read line; do
	curl http://localhost:8888/deck/fruit -X POST -d "$line"
done < ./fruit.list
while read line; do
	curl http://localhost:8888/deck/vegetables -X POST -d "$line"
done < ./vegetables.list
echo "Done"
