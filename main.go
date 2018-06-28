package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"
)

var offices = []Office{AMS3, AMS9, AMS10, AMS11, AMS14, AMS15, AMS16, AMS17, AMS22}

var scheduler chan string

func main() {
	log.Println("Random Coffee initialized")

	db := NewStorage()
	db.Init("./storage.db")
	defer db.Connection.Close()

	scheduler = make(chan string)

	accessToken := os.Getenv("PAGE_ACCESS_TOKEN")
	verifyToken := os.Getenv("VERIFY_TOKEN")

	if accessToken == "" {
		panic("no PAGE_ACCESS_TOKEN environment variable")
	}
	if verifyToken == "" {
		panic("no VERIFY_TOKEN environment variable")
	}

	addr := "127.0.0.1:80"
	msngr := NewMessenger(accessToken, verifyToken)

	go msngr.Start(addr)

	log.Printf("started at %s\n", addr)

	roster := NewRoster(db)

	names := []string{
		"Shanel", "Al", "Sandee", "Jimmie",
		"Luella", "Neoma", "Edmond", "Marilee",
		"Kristy", "Deonna", "Margo", "Bethany",
		"Cuc", "Kathlene", "Mica", "Shanti",
		"Joycelyn", "Norbert", "Ardath", "Nichell",
	}
	preferredLocations := []OfficeGroup{Rembrandtplein, Vijzelstraat, PietHeinkade, Sloterdijk, Zuid}

	var e *Employee
	for i := 0; i <= 10; i++ {
		e = &Employee{
			ID:                fmt.Sprintf("id-%d", i),
			Active:            true,
			Name:              names[rand.Intn(len(names))],
			Oldie:             false,
			PreferredLocation: preferredLocations[rand.Intn(len(preferredLocations))],
		}
		roster.Add(e)
	}

	for _, emp := range roster.Employees {
		fmt.Printf("%v in %s\n", emp.Name, emp.PreferredLocation)
	}

	go func() {
		for {
			select {
			case m := <-msngr.C:
				go processMessage(m, msngr, roster, db)
			}
		}
	}()

	// Ticker

	//ticker := time.NewTicker(10 * time.Second)
	//defer ticker.Stop()

	for {
		select {
		case event := <-scheduler:
			log.Println(event)
			switch event {
			case "AVA":
				roster.SetAvailabilityAll(Unavailable)
				for _, e := range roster.Employees {
					if e.Availability == Unavailable {
						e.Availability = Unknown
						go msngr.Send(Messaging{
							Recipient: User{
								ID: e.ID,
							},
							Message: &Message{
								Text: fmt.Sprintf("Good morning %s! Are you available to grab a coffee with someone today?", e.Name),
								QuickReplies: &[]QuickReply{
									{
										ContentType: "text",
										Title:       "Yes",
										Payload:     "<AVAILABILITY:YES>",
									},
									{
										ContentType: "text",
										Title:       "Not today",
										Payload:     "<AVAILABILITY:NO>",
									},
								},
							},
						})
					}
				}

			case "PAIR":
				notifyPairs(roster.GetMatches(), msngr)
			}
		}
	}

	wednesdayMorning := time.After(3 * time.Second)
	//wednesdayAfternoon := time.NewTicker(20 * time.Second)

	for {
		select {
		//case <-done:
		//	fmt.Println("Done!")
		//	return
		//case <-ticker.C:
		case <-time.After(15 * time.Second):
			notifyPairs(roster.GetMatches(), msngr)

		// This should happen every Wednesday morning
		case <-wednesdayMorning:
			log.Println("\nROSTOER\n")
			roster.SetAvailabilityAll(Unavailable)

			for _, e := range roster.Employees {
				log.Printf("Sending to %s\n", e.ID)
				if e.Availability == Unavailable {
					e.Availability = Unknown
					go msngr.Send(Messaging{
						Recipient: User{
							ID: e.ID,
						},
						Message: &Message{
							Text: fmt.Sprintf("Good morning %s! Are you available to grab a coffee with someone today?", e.Name),
							QuickReplies: &[]QuickReply{
								{
									ContentType: "text",
									Title:       "Yes",
									Payload:     "<AVAILABILITY:YES>",
								},
								{
									ContentType: "text",
									Title:       "Not today",
									Payload:     "<AVAILABILITY:NO>",
								},
							},
						},
					})
				}
			}

		default:
			// do nothing
		}
	}
}

func processMessage(m Messaging, messenger *Messenger, roster *Roster, db *Storage) {
	senderID := m.Sender.ID

	if m.Postback != nil {
		if p := m.Postback.Payload; len(p) > 0 {
			switch p {
			case "AVA":
				fallthrough
			case "PAIR":
				scheduler <- p
				return
			}
		}
	}

	employee, ok := roster.GetByID(senderID)
	if !ok {
		m := messenger.GetMember(senderID)

		log.Printf("received memeber %#v", m)


		employee = &Employee{
			ID:           senderID,
			Name:         m.Name,
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
				Text: fmt.Sprintf("Hey %s, I hope you’re having a great day! I’m here to find a random colleague for you to grab a coffee with.", employee.Name),
			},
		})
		employee.Oldie = true
	}

	var qr *QuickReply

	if m.Message != nil {
		qr = m.Message.QuickReply
	}

	// Handle selection of preferred location
	if qr != nil {
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
			db.SaveEmployee(employee)

			messenger.Send(Messaging{
				Recipient: User{
					ID: senderID,
				},
				Message: &Message{
					Text: fmt.Sprintf("Great! I’m going to grind some beans and I’ll get back to you with a match shortly."),
				},
			})
		}
	}

	// Handle availability change
	if qr != nil {
		switch qr.Payload {
		case "<AVAILABILITY:YES>":
			(*employee).Availability = Available
			messenger.Send(Messaging{
				Recipient: User{
					ID: senderID,
				},
				Message: &Message{
					Text: "Perfect. I’ll get back to you around noon with your match. Have a good morning!",
				},
			})

		case "<AVAILABILITY:NO>":
			messenger.Send(Messaging{
				Recipient: User{
					ID: senderID,
				},
				Message: &Message{
					Text: "No problem, I’ll talk to you next Wednesday. Have a great week and weekend!",
					QuickReplies: &[]QuickReply{
						{
							ContentType: "text",
							Title:       "Sounds good!",
							Payload:     "<AVAILABILITY:POSTPONE>",
						},
						{
							ContentType: "text",
							Title:       "I’d like to unsubscribe",
							Payload:     "<AVAILABILITY:UNSUBSCRIBE>",
						},
					},
				},
			})

		case "<AVAILABILITY:UNSUBSCRIBE>":
			employee.Active = false
			messenger.Send(Messaging{
				Recipient: User{
					ID: senderID,
				},
				Message: &Message{
					Text: "YOU'VE BEEN UNSUBSCRBED",
				},
			})
		}
	}

	if qr != nil {
		return
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
		return
	}
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
