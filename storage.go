package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
)

type Storage struct {
	filename   string
	Connection *sql.DB
}

func NewStorage() *Storage {
	return &Storage{
		filename: "",
	}
}

func (s *Storage) Init(filename string) {
	s.filename = filename
	var err error
	s.Connection, err = sql.Open("sqlite3", s.filename)
	if err != nil {
		fmt.Printf("[DATABASE ERROR] %v", err)
		return
	}
}

func (s *Storage) GetAllEmployees() map[string]*Employee {
	employees := make(map[string]*Employee)
	rows, err := s.Connection.Query("SELECT * FROM employee")
	if err != nil {
		fmt.Printf("[DATABASE ERROR] %v", err)
		return employees
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		var workplaceId string
		var name string
		var active int
		_ = rows.Scan(&id, &workplaceId, &name, &active)
		employees[workplaceId] = &Employee{
			ID:   workplaceId,
			Name: name,
		}
	}
	return employees
}
