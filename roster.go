package main

import "sort"

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
			if e.Availability != Available {
				continue
			}
			for _, e2 := range g[i:] {
				if e2.Availability == Available && !(e2.Matches.wasMatchedWithBefore(e)) {
					e.Availability = Matched
					e2.Availability = Matched
					pair := []*Employee{e, e2}
					matches = append(matches, pair)
					continue
				}
			}
		}
	}

	return matches
}
