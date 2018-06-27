package main

import "time"

type Office int

const (
	AMS3 Office = iota
	AMS9
	AMS10
	AMS11
)

var office = [...]string{
	"AMS3",  // ??
	"AMS9",  // the bank
	"AMS10", // learning center
	"AMS11", // spaces
	"AMS17", // piet hein?
	"AMS19", // ??
}

func (o Office) String() string {
	return office[o]
}

type Frequency int

const (
	Weekly Frequency = iota
	Biweekly
	Triweekly
	Monthly
)

type Employee struct {
	ID                string
	Name              string
	Office            Office
	LastMatch         Match
	Availability      []time.Weekday
	Frequency         Frequency
	PreferredLocation string // Should be group of the offices
	Oldie            bool
}

func (e *Employee) wasMatchedToday() bool {
	return false
}

func (e *Employee) AddWeekday(w time.Weekday) {
	for _, a := range e.Availability {
		if a == w {
			return
		}
	}
	e.Availability = append(e.Availability, w)
}

func (e *Employee) RemoveWeekday(w time.Weekday) {
	for i, a := range e.Availability {
		if a == w {
			e.Availability = append(e.Availability[:i], e.Availability[i+1:]...)
		}
	}
}

func (e *Employee) SharedWeekdays(d *Employee) {

}

func (e *Employee) MatchedThisWeek() bool {
	return false
}

type Rating int

const (
	BAD Rating = iota - 1
	NONE
	GOOD
)

type Match struct {
	Employees []Employee
	Time      time.Time
	Reviewed  bool
	Rating    [2]Rating
}

//type EmployeeRoster struct {
//	roster map[string]*Employee
//}
//
//func NewEmployeeRoster() *EmployeeRoster {
//	return &EmployeeRoster{
//		roster: make(map[string]*Employee),
//	}
//}
//
//func (er *EmployeeRoster) Add(employee *Employee) {
//	if _, ok := er.roster[employee.ID]; !ok {
//		er.roster[employee.ID] = employee
//	}
//}
