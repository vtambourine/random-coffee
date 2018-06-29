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
	log.Printf("[STORAGE] Initialising storage with file %s", filename)
	s.filename = filename
	var err error
	s.Connection, err = sql.Open("sqlite3", s.filename)
	if err != nil {
		log.Printf("[DATABASE ERROR] %v", err)
		return
	}
}

func (s *Storage) GetEmployee(rowId int) *Employee {
	var id int
	var workplaceId string
	var name string
	var firstName string
	var preferredLocation string
	var availability int
	var active int
	err := s.Connection.QueryRow("SELECT * FROM employee WHERE id=?", rowId).
		Scan(&id, &workplaceId, &name, &firstName, &preferredLocation, &availability, &active)
	if err != nil {
		return &Employee{}
	} else {
		currentEmployee := &Employee{
			ID:                workplaceId,
			Name:              name,
			FirstName:         firstName,
			PreferredLocation: OfficeGroup(preferredLocation),
			Availability:      Availability(availability),
			Active:            active != 0,
			Oldie:             true,
		}
		return currentEmployee
	}
}

func (s *Storage) GetEmployeeId(workplaceId string) int {
	var id int
	err := s.Connection.QueryRow("SELECT id FROM employee WHERE workplace_id=?", workplaceId).Scan(&id)
	if err != nil {
		return 0
	} else {
		return id
	}
}

func (s *Storage) SaveEmployee(employee *Employee) {
	log.Printf("[STORAGE] Saving employee (%s) %s to the database", employee.ID, employee.Name)
	id := s.GetEmployeeId(employee.ID)
	if id != 0 {
		stmt, err := s.Connection.Prepare("UPDATE employee SET workplace_id = ?, name = ?, first_name = ?, preferred_location = ?, availability = ?, active = ? WHERE id = ?")
		if err != nil {
			log.Printf("[DATABASE ERROR] %v", err)
		}
		_, err = stmt.Exec(employee.ID, employee.Name, employee.FirstName, employee.PreferredLocation, employee.Availability, employee.Active, id)
		if err != nil {
			log.Printf("[DATABASE ERROR] %v", err)
		}
	} else {
		stmt, err := s.Connection.Prepare("INSERT INTO employee (workplace_id, name, first_name, preferred_location, availability, active) VALUES(?, ?, ?, ?, ?, ?)")
		if err != nil {
			log.Printf("[DATABASE ERROR] %v", err)
		}
		_, err = stmt.Exec(employee.ID, employee.Name, employee.FirstName, employee.PreferredLocation, employee.Availability, employee.Active)
		if err != nil {
			log.Printf("[DATABASE ERROR] %v", err)
		}

	}
}

func (s *Storage) GetAllEmployees() map[string]*Employee {
	log.Printf("[STORAGE] Fetching all employees from database file %s", s.filename)
	employees := make(map[string]*Employee)
	dbEmployees, err := s.Connection.Query("SELECT * FROM employee")
	if err != nil {
		log.Printf("[DATABASE ERROR] %v", err)
		return employees
	}
	defer dbEmployees.Close()
	for dbEmployees.Next() {
		var id int
		var workplaceId string
		var name string
		var firstName string
		var preferredLocation string
		var availability int
		var active int
		_ = dbEmployees.Scan(&id, &workplaceId, &name, &firstName, &preferredLocation, &availability, &active)
		currentEmployee := &Employee{
			ID:                workplaceId,
			Name:              name,
			FirstName:         firstName,
			PreferredLocation: OfficeGroup(preferredLocation),
			Availability:      Availability(availability),
			Active:            active != 0,
			Oldie:             true,
		}
		log.Printf("[STORAGE] Fetching all previous matches from employee with id %d (%s)", id, workplaceId)
		dbMatchesForEmployee, err := s.Connection.Query("SELECT match1_id, match2_id, created_at, happened FROM matches WHERE match1_id = ? OR match2_id = ?", id, id)
		if err != nil {
			log.Printf("[DATABASE ERROR] %v", err)
		}
		var matches Matches
		for dbMatchesForEmployee.Next() {
			var match1Id int
			var match2Id int
			var createdAt time.Time
			var happened MatchStatus
			_ = dbEmployees.Scan(&match1Id, &match2Id, &createdAt, &happened)
			colleagueId := match1Id
			if id == match1Id {
				colleagueId = match2Id
			}
			colleague := s.GetEmployee(colleagueId)

			pair := []*Employee{currentEmployee, colleague}
			match := Match{
				Pair:     pair,
				Time:     createdAt,
				Happened: MatchStatus(happened),
			}
			matches = append(matches, match)
		}
		currentEmployee.Matches = matches
		employees[workplaceId] = currentEmployee
	}
	return employees
}

func (s *Storage) SaveMatch(match *Match) {
	log.Printf("[STORAGE] Saving match between (%s) %s and (%s) %s to the database", match.Pair[0].ID, match.Pair[0].Name, match.Pair[1].ID, match.Pair[1].Name)
	stmt, err := s.Connection.Prepare("INSERT INTO matches (match1_id, match2_id, created_at, happened) VALUES(?, ?, ?, ?)")
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
