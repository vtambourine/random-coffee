package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"
)

type EmployeeRoster map[string]*Employee

var offices = []Office{AMS3, AMS9, AMS10, AMS11, AMS14, AMS15, AMS16, AMS17, AMS22}

func main() {
	log.Println("Random Coffee initialized")

	log.Println("Random Coffee initialized")

	accessToken := os.Getenv("PAGE_ACCESS_TOKEN")
	verifyToken := os.Getenv("VERIFY_TOKEN")

	if accessToken == "" {
		panic("no PAGE_ACCESS_TOKEN environment variable")
	}
	if verifyToken == "" {
		panic("no VERIFY_TOKEN environment variable")
	}

	addr := "127.0.0.1:3001"
	msngr := NewMessenger(accessToken, verifyToken)

	go msngr.Start(addr)

	log.Printf("started at %s\n", addr)

	//roster := make(EmployeeRoster)
	roster := NewRoster()

	names := []string{
		"Shanel", "Al", "Sandee", "Jimmie",
		"Luella", "Neoma", "Edmond", "Marilee",
		"Kristy", "Deonna", "Margo", "Bethany",
		"Cuc", "Kathlene", "Mica", "Shanti",
		"Joycelyn", "Norbert", "Ardath", "Nichell",
	}
	preferredLocations := []OfficeGroup{Rembrandtplein, Vijzelstraat, PietHeinkade, Sloterdijk, Zuid}

	var e *Employee
	for i := 0; i <= 20; i++ {
		e = &Employee{
			ID:                fmt.Sprintf("id-%d", i),
			Name:              names[rand.Intn(len(names))],
			Office:            offices[rand.Intn(len(offices))],
			Oldie:             false,
			PreferredLocation: preferredLocations[rand.Intn(len(preferredLocations))],
		}
		//roster[e.ID] = e
		roster.Add(e)

		fmt.Printf("%v in %s\n", (*e).Name, (*e).PreferredLocation)
	}

	go func() {
		for {
			select {
			case m := <-msngr.C:
				go processMessage(m, msngr, roster)
			}
		}
	}()

	//msgr := NewMatcher()
	//msgr.Add(alice)
	//msgr.Add(bob)

	// Restore employees from database

	// From all employees choose who

	//fmt.Printf("%v", msgr.GetMatches())

	ticker := time.NewTicker(100 * time.Second)
	defer ticker.Stop()

	notifyPairs(roster.GetMatches(), msngr)
	//done := make(chan bool)
	//go func() {
	//time.Sleep(1 * time.Second)
	//done <- true
	//}()

	//match(&employees)

	for {
		select {
		//case <-done:
		//	fmt.Println("Done!")
		//	return
		case <-ticker.C:
			//match(&roster)
			notifyPairs(roster.GetMatches(), msngr)
		}
	}
}

func processMessage(m Messaging, messenger *Messenger, roster *Roster) {
	senderID := m.Sender.ID

	//employee, ok := employees[senderID]
	//if !ok {
	//	employee = &Employee{
	//		ID:   senderID,
	//		Name: "New Name",
	//	}
	//	employees[senderID] = employee
	//}

	employee, ok := roster.GetByID(senderID)
	if !ok {
		employee = &Employee{
			ID:   senderID,
			Name: "New Name",
		}
		//employees[senderID] = employee
		roster.Add(employee)
	}

	log.Printf("recieved message from: %s", senderID)
	log.Printf("before process:\n%#v", *employee)


	// If user contact bot for the first time, greet him
	if !employee.Oldie {
		messenger.Send(Messaging{
			Recipient: User{
				ID: senderID,
			},
			Message: &Message{
				Text: fmt.Sprintf("Hey {{NAME}}, I hope you’re having a great day! I’m here to find a random colleague for you to grab a coffee with."),
			},
		})
		employee.Oldie = true
	}

	// Handle selection of preferred location
	if qr := m.Message.QuickReply; qr != nil {
		switch qr.Payload {
		case string(Rembrandtplein):
			fallthrough
		case string(Vijzelstraat):
			fallthrough
		case string(PietHeinkade):
			fallthrough
		case string(Sloterdijk):
			fallthrough
		case string(Zuid):
			(*employee).PreferredLocation = OfficeGroup(qr.Payload)
		}
	}

	// If user doesn;t have preferred location
	if len(employee.PreferredLocation) == 0 {
		messenger.Send(Messaging{
			Recipient: User{
				ID: senderID,
			},
			Message: &Message{
				Text: "Which office are you in?",
				QuickReplies: &[]QuickReply{
					{
						ContentType: "text",
						Title:       string(Rembrandtplein),
						Payload:     string(Rembrandtplein),
					},
					{
						ContentType: "text",
						Title:       string(Vijzelstraat),
						Payload:     string(Vijzelstraat),
					},
					{
						ContentType: "text",
						Title:       string(PietHeinkade),
						Payload:     string(PietHeinkade),
					},
					{
						ContentType: "text",
						Title:       string(Sloterdijk),
						Payload:     string(Sloterdijk),
					},
					{
						ContentType: "text",
						Title:       string(Zuid),
						Payload:     string(Zuid),
					},
				},
			},
		})
	}

	// If user have preferred location (and other condition moght apply)
	if len(employee.PreferredLocation) > 0 {
		messenger.Send(Messaging{
			Recipient: User{
				ID: senderID,
			},
			Message: &Message{
				Text: fmt.Sprintf("Great! I’m going to grind some beans and I’ll get back to you with a match shortly."),
			},
		})
	}

	log.Printf("after process:\n%#v\n", *employee)
}

func notifyPairs(matches [][]*Employee, messenger *Messenger) {
	for _, employees := range matches {
		fmt.Println("== Group: ==")
		for _, e := range employees {
			fmt.Printf("%s : %v\n", e.Name, e.ID)
		}
		fmt.Println()
	}
}

func match(employees *EmployeeRoster) {
	fmt.Print("MATCHING ")
	fmt.Println(time.Now().Format(time.UnixDate))

	//match := NewMatcher()
	for _, e := range *employees {
		log.Println(e.Name)
		//match.Add(e)
	}

	//fmt.Printf("Match: %v\n\n", match.GetMatches())
}
