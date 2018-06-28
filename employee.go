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
	Rembrandtplein   OfficeGroup = "Rembrandtplein"
	Vijzelstraat     OfficeGroup = "Vijzelstraat"
	PietHeinkade     OfficeGroup = "Piet Heinkade"
	Sloterdijk       OfficeGroup = "Sloterdijk"
	Zuid             OfficeGroup = "Zuid"
	OfficeGroupEmpty OfficeGroup = ""
)

type Frequency int

const (
	Weekly Frequency = iota
	Biweekly
	Triweekly
	Monthly
)

type Availability int

const (
	Unknown Availability = iota
	Unavailable// Unavailable for matching this week
	Available // Available for matching this week
	Matched // Already matched this week
)

type ConversationState int

const (
	ConversationStateUnknown ConversationState = iota
	ConversationStateLocation
	ConversationStateAvailability
	ConversationStateFeedback
)

type Employee struct {
	ID                string
	Name              string
	Active            bool
	Matches           Matches
	Availability      Availability
	PreferredLocation OfficeGroup // Preferred office group
	Oldie             bool        // Already talked to the bot and received introduction message
	ConversationState ConversationState
}

func (e *Employee) wasMatchedToday() bool {
	return false
}

func (e *Employee) MatchedThisWeek() bool {
	return false
}

type MatchStatus int

const (
	MatchUnknown  MatchStatus = iota // 0
	MatchHappened                    // 1
	MatchMissed                      // 2
)

type Match struct {
	Pair     [2]*Employee
	Time     time.Time
	Happened MatchStatus // Default to true
}

type Matches []Match

func (ms Matches) wasMatchedWithBefore(e *Employee) bool {
	for _, m := range ms {
		if m.Pair[0] == e || m.Pair[1] == e {
			return true
		}
	}
	return false
}
