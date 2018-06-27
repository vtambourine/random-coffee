package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
)

var schema = `
CREATE TABLE person (
	id text,
	name text
);

CREATE TABLE office (
	id integer
	name text
)`

type Person struct {
	ID   string `db:"name"`
	Name string `db:"name"`
}

type Office struct {
	ID   int
	Name string
}

func main() {
	log.Println("Hello, World")

	os.Remove("./coffee.db")

	db, err := sql.Open("sqlite3", "./coffee.db")
	if err != nil {
		log.Fatal(err)

	}
	defer db.Close()

	_, err = db.Exec(schema)
}
