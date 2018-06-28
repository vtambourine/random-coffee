package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
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
func (s *Storage) SaveEmployee(employee *Employee) {
	stmt, err := s.Connection.Prepare("INSERT INTO employee (workplace_id, name, preferred_location, availability, active) VALUES(?, ?, ?, ?, ?)")
	if err != nil {
		log.Print(err)
	}
	_, err = stmt.Exec(employee.ID, employee.Name, employee.PreferredLocation, employee.Availability, employee.Active)
	if err != nil {
		log.Print(err)
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
		var preferredLocation string
		var availability int
		var active int
		_ = rows.Scan(&id, &workplaceId, &name, &preferredLocation, &availability, &active)
		employees[workplaceId] = &Employee{
			ID:                workplaceId,
			Name:              name,
			PreferredLocation: OfficeGroup(preferredLocation),
			Availability:      Availability(availability),
			Active:            active != 0,
		}
	}
	return employees
}
