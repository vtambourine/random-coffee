package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
		"time"
)

type EmployeeRoster map[string]*Employee

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

	offices := []Office{AMS9, AMS10, AMS11}
	names := []string{
		"Shanel", "Al", "Sandee", "Jimmie",
		"Luella", "Neoma", "Edmond", "Marilee",
		"Kristy", "Deonna", "Margo", "Bethany",
		"Cuc", "Kathlene", "Mica", "Shanti",
		"Joycelyn", "Norbert", "Ardath", "Nichell",
	}
	preferredLocations := []string{"AMS3:AMS11", "AMS9:AMS10", "AMS17:AMS19"}
	//frequencies := []Frequency{Weekly, Biweekly, Triweekly}

	var e *Employee
	for i := 0; i <= 20; i++ {
		e = &Employee{
			ID:     string(i),
			Name:   names[rand.Intn(len(names))],
			Office: offices[rand.Intn(len(offices))],
			Oldie:  false,
			//Frequency: frequencies[rand.Intn(len(frequencies))],
			PreferredLocation: preferredLocations[rand.Intn(len(preferredLocations))],
		}
		//roster[e.ID] = e
		roster.Add(e)

		fmt.Printf("%v\n", (*e).Name)
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

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

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
			roster.GetMatches()
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
				Text: fmt.Sprintf("GREETING"),
			},
		})
		employee.Oldie = true
	}

	// If user doesn;t have preferred location
	if len(employee.PreferredLocation) == 0 {
		messenger.Send(Messaging{
			Recipient: User{
				ID: senderID,
			},
			Message: &Message{
				Text: "SELECT PREFERRED LOCATION",
				QuickReplies: &[]QuickReply{
					{
						ContentType: "text",
						Title:       "rembrandt",
						Payload:     "AMS9:AMS10",
					},
					{
						ContentType: "text",
						Title:       "vijzel",
						Payload:     "AMS3:AMS11",
					},
					{
						ContentType: "text",
						Title:       "piethein",
						Payload:     "AMS17:AMS19",
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
				Text: fmt.Sprintf("PREFERRED LOCATION RECEIVED"),
			},
		})
		messenger.Send(Messaging{
			Recipient: User{
				ID: senderID,
			},
			Message: &Message{
				Text: fmt.Sprintf("TOT ZIENS"),
			},
		})
	}

	// Process Postbacks

	// Handle Office selection
	if qr := m.Message.QuickReply; qr != nil {
		switch qr.Payload {
		case "AMS9:AMS10":
			fallthrough
		case "AMS3:AMS11":
			fallthrough
		case "AMS17:AMS19":
			//for _, o := range strings.Split(qr.Payload, ":") {
			//	//log.Printf("employee %#v wants to meet at %s", employee, o)
			//}
			employee.PreferredLocation = qr.Payload
		}
	}

	log.Printf("after process:\n%#v\n", *employee)
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
