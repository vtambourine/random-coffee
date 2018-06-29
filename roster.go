package main

import (
	"fmt"
	"log"
	"sort"
	"time"
)

type Roster struct {
	Employees map[string]*Employee
	db        *Storage
}

func NewRoster(db *Storage) *Roster {
	existingEmployees := db.GetAllEmployees()
	return &Roster{
		Employees: existingEmployees,
		db:        db,
	}
}

func (r *Roster) Add(employee *Employee) {
	if _, ok := r.Employees[employee.ID]; ok {
		return
	}
	r.db.SaveEmployee(employee)
	r.Employees[employee.ID] = employee
}

func (r *Roster) Has(employee *Employee) bool {
	_, ok := r.Employees[employee.ID]
	return ok
}

func (r *Roster) SetAvailabilityAll(a Availability) {
	for _, e := range r.Employees {
		if e.Active {
			e.Availability = a
			r.db.SaveEmployee(e)
		}
	}
}

func (r *Roster) GetByID(id string) (*Employee, bool) {
	e, ok := r.Employees[id]
	if ok {
		return e, true
	} else {
		return nil, false
	}
}

func (r *Roster) GetMatches() [][]*Employee {
	log.Printf("Running GetMatches")
	groups := make(map[OfficeGroup][]*Employee)

	for _, e := range r.Employees {
		if og := e.PreferredLocation; len(og) > 0 {
			groups[og] = append(groups[og], e)
		}
	}

	matches := [][]*Employee{}
	for _, g := range groups {
		sort.Slice(g, func(i, j int) bool {
			return len(g[i].Matches) > len(g[j].Matches)
		})

		for i, e := range g {
			fmt.Printf("current user %v\n", e)
			if e.Availability != Available {
				continue
			}
			for _, e2 := range g[i+1:] {
				if e2.Availability == Available && !(e2.Matches.wasMatchedWithBefore(e)) {
					fmt.Printf("%v match %v\n", e, e2)
					e.Availability = Matched
					r.db.SaveEmployee(e)
					e2.Availability = Matched
					r.db.SaveEmployee(e2)
					pair := []*Employee{e, e2}
					matches = append(matches, pair)

					match := Match{
						Pair:     pair,
						Time:     time.Now(),
						Happened: MatchUnknown,
					}
					r.db.SaveMatch(&match)

					pair[0].Matches = append(pair[0].Matches, match)
					pair[1].Matches = append(pair[1].Matches, match)

					break
				}
			}
		}
	}

	return matches
}
