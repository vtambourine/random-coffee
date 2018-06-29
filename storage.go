package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"time"
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
		log.Printf("[DATABASE ERROR] %v", err)
		return
	}
}

func (s *Storage) GetEmployeeId(ID string) int {
	var id int
	err := s.Connection.QueryRow("SELECT id FROM employee WHERE workplace_id=?", ID).Scan(&id)
	if err != nil {
		return 0
	} else {
		return id
	}
}

func (s *Storage) SaveEmployee(employee *Employee) {
	stmt, err := s.Connection.Prepare("INSERT OR REPLACE INTO employee (workplace_id, name, preferred_location, availability, active) VALUES(?, ?, ?, ?, ?)")
	if err != nil {
		log.Printf("[DATABASE ERROR] %v", err)
	}
	_, err = stmt.Exec(employee.ID, employee.Name, employee.PreferredLocation, employee.Availability, employee.Active)
	if err != nil {
		log.Printf("[DATABASE ERROR] %v", err)
	}
}

func (s *Storage) GetAllEmployees() map[string]*Employee {
	employees := make(map[string]*Employee)
	rows, err := s.Connection.Query("SELECT * FROM employee")
	if err != nil {
		log.Printf("[DATABASE ERROR] %v", err)
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
			Oldie:             true,
		}
	}
	return employees
}

func (s *Storage) SaveMatch(match *Match) {
	stmt, err := s.Connection.Prepare("INSERT INTO matches (match_id_1, match_id_2, created_at, happened) VALUES(?, ?, ?, ?)")
	if err != nil {
		log.Printf("[DATABASE ERROR] %v", err)
	}
	_, err = stmt.Exec(s.GetEmployeeId(match.Pair[0].ID), s.GetEmployeeId(match.Pair[1].ID), match.Time, match.Happened)
	if err != nil {
		log.Printf("[DATABASE ERROR] %v", err)
	}
}

func (s *Storage) SaveAllMatches(matches [][]*Employee) {
	for _, pair := range matches {
		match := &Match{
			Pair:     pair,
			Time:     time.Now(),
			Happened: MatchUnknown,
		}
		s.SaveMatch(match)
		s.SaveEmployee(match.Pair[0])
		s.SaveEmployee(match.Pair[1])
	}
}
