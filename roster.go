package main

import (
	"log"
	"fmt"
)

type Roster struct {
	employees map[string]*Employee
}

func NewRoster() *Roster {
	return &Roster{
		employees: make(map[string]*Employee),
	}
}

func (r *Roster) Add(employee *Employee) {
	if _, ok := r.employees[employee.ID]; ok {
		return
	}
	r.employees[employee.ID] = employee
}

func (r *Roster) Has(employee *Employee) bool {
	_, ok := r.employees[employee.ID]
	return ok
}

func (r *Roster) GetByID(id string) (*Employee, bool) {
	e, ok := r.employees[id]
	if ok {
		return e, true
	} else {
		return nil, false
	}
}

func (r *Roster) GetMatches() [][]*Employee {
	groups := make(map[OfficeGroup][]*Employee)

	for _, e := range r.employees {
		og := e.Office.GetGroup()
		groups[og] = append(groups[og], e)
	}

	matches := [][]*Employee{}
	for n, g := range groups {
		log.Printf("Group from %s", n)
		m := []*Employee{}
		for _, e := range g {
			fmt.Print(e.Name, " ")
			m = append(m, e)
		}

		fmt.Println(" ")
		matches = append(matches, m)
	}

	return matches
}
