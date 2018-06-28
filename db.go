package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"strconv"
)

func db_test() {
	database, _ := sql.Open("sqlite3", "./storage.db")
	statement, _ := database.Prepare("INSERT INTO employee (email_address) VALUES (?)")
	statement.Exec("test@booking.com")
	rows, _ := database.Query("SELECT * FROM employee")
	fmt.Println("%v", rows)
	var id int
	var email string
	var active int
	for rows.Next() {
		rows.Scan(&id, &email, &active)
		fmt.Println(strconv.Itoa(id) + " : " + email + " : " + strconv.Itoa(active))
	}
}
