package main

import "time"

type Office string

const (
	AMS3  Office = "AMS3"  // Vijzelstraat
	AMS9  Office = "AMS9"  // The Bank
	AMS10 Office = "AMS10" // Learning Center
	AMS11 Office = "AMS11" // Prins & Keizer
	AMS14 Office = "AMS14" // Sloterdijk
	AMS15 Office = "AMS15" // Piet Hein
	AMS16 Office = "AMS16" // Bloomhouse
	AMS17 Office = "AMS17" // UP
	AMS22 Office = "AMS22" // Atrium
)

//var office = [...]string{
//	"AMS3",  // ??
//	"AMS9",  // the bank
//	"AMS10", // learning center
//	"AMS11", // spaces
//	"AMS17", // piet hein?
//	"AMS19", // ??
//}

//func (o Office) String() string {
//	return office[o]
//}

var officeToGroup = map[Office]OfficeGroup{
	AMS3:  Vijzelstraat,
	AMS9:  Rembrandtplein,
	AMS10: Rembrandtplein,
	AMS11: Vijzelstraat,
	AMS14: Sloterdijk,
	AMS15: PietHeinkade,
	AMS16: Rembrandtplein,
	AMS17: PietHeinkade,
	AMS22: Zuid,
}

func (o Office) GetGroup() OfficeGroup {
	return officeToGroup[o]
}

type OfficeGroup string

const (
	Rembrandtplein OfficeGroup = "Rembrandtplein"
	Vijzelstraat   OfficeGroup = "Vijzelstraat"
	PietHeinkade   OfficeGroup = "Piet Heinkade"
	Sloterdijk     OfficeGroup = "Sloterdijk"
	Zuid           OfficeGroup = "Zuid"
)

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
	PreferredLocation OfficeGroup
	Oldie             bool
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
