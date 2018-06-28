package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"
)

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
			Active:            true,
			Name:              names[rand.Intn(len(names))],
			Oldie:             false,
			PreferredLocation: preferredLocations[rand.Intn(len(preferredLocations))],
		}
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

	// Ticker

	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	wednesdayMorning := time.NewTicker(10 * time.Second)
	//wednesdayAfternoon := time.NewTicker(20 * time.Second)

	// This should happen every Wednesday morning
	roster.SetAvailabilityAll(Unavailable)

	for {
		select {
		//case <-done:
		//	fmt.Println("Done!")
		//	return
		case <-ticker.C:
			notifyPairs(roster.GetMatches(), msngr)

		// This should happen every Wednesday morning
		case <-wednesdayMorning.C:
			for _, e := range roster.Employees {
				if e.Availability == Unavailable {
					e.Availability = Uncertain
					go msngr.Send(Messaging{
						Recipient: User{
							ID: e.ID,
						},
						Message: &Message{
							Text: "Good morning {{NAME}}! Are you available to grab a coffee with someone today?",
							QuickReplies: &[]QuickReply{
								{
									ContentType: "text",
									Title:       "Yes",
									Payload:     "<AVAILABILITY:YES>",
								},
								{
									ContentType: "text",
									Title:       "<AVAILABILITY:NO>",
								},
							},
						},
					})
				}
			}
		}
	}
}

func processMessage(m Messaging, messenger *Messenger, roster *Roster) {
	senderID := m.Sender.ID

	employee, ok := roster.GetByID(senderID)
	if !ok {
		employee = &Employee{
			ID:           senderID,
			Name:         "New Name",
			Availability: Unavailable,
		}
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

	// If user have preferred location (and other condition might apply)
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

// Send notifications to the pairs
func notifyPairs(matches [][]*Employee, messenger *Messenger) {
	for _, pairs := range matches {
		fmt.Println("== Pair ==")

		go messenger.Send(Messaging{
			Recipient: User{
				ID: pairs[0].ID,
			},
			Message: &Message{
				Text: fmt.Sprintf("Perfect. Your match this week is %s. Shoot them a message on Workplace and organize a time to meet!", pairs[1].Name),
			},
		})

		go messenger.Send(Messaging{
			Recipient: User{
				ID: pairs[1].ID,
			},
			Message: &Message{
				Text: fmt.Sprintf(" Perfect. Your match this week is %s. Shoot them a message on Workplace and organize a time to meet!", pairs[0].Name),
			},
		})

		for _, p := range pairs {
			fmt.Printf("%s : %v\n", p.Name, p.ID)

		}

		fmt.Println()
	}
}
