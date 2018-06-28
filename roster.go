package main


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
		if og := e.PreferredLocation; len(og) > 0 {
			groups[og] = append(groups[og], e)
		}
	}

	matches := [][]*Employee{}
	for _, g := range groups {
		m := []*Employee{}
		for _, e := range g {
			m = append(m, e)
		}
		matches = append(matches, m)
	}

	return matches
}
