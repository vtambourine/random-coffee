package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"github.com/robfig/cron"
	"github.com/natefinch/lumberjack"
)

var scheduler chan string

func main() {
	log.SetOutput(&lumberjack.Logger{
		Filename:   "/home/vkleshchenko/random-coffee_logs/rc.log",
		MaxSize:    100, // megabytes
		MaxAge:     28, //days
	})

	log.Println("Random Coffee initialized")

	accessToken := os.Getenv("PAGE_ACCESS_TOKEN")
	verifyToken := os.Getenv("VERIFY_TOKEN")

	if accessToken == "" {
		panic("no PAGE_ACCESS_TOKEN environment variable")
	}
	if verifyToken == "" {
		panic("no VERIFY_TOKEN environment variable")
	}

	var port string
	if port = os.Getenv("PORT"); len(port) == 0 {
		port = "3000"
	}

	// Run server
	addr := fmt.Sprintf("127.0.0.1:%s", port)
	messenger := NewMessenger(accessToken, verifyToken)
	go messenger.Start(addr)
	log.Printf("started at %s\n", addr)

	// Create new data storage
	db := NewStorage()
	db.Init("./storage.db")
	defer db.Connection.Close()

	// Create new scheduler channel to trigger weekly events
	scheduler = make(chan string)

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
	for i := 0; i < 0; i++ {
		e = &Employee{
			ID:                fmt.Sprintf("id-%d", i),
			Active:            true,
			Name:              names[rand.Intn(len(names))],
			Oldie:             false,
			PreferredLocation: preferredLocations[rand.Intn(len(preferredLocations))],
		}
		roster.Add(e)
	}

	// Start listen to new messages ot the bot
	go func() {
		for {
			select {
			case m := <-messenger.C:
				go processMessage(m, messenger, roster, db)
			}
		}
	}()

	// Start listen to cheat codes
	go func() {
		for {
			select {
			case event := <-scheduler:
				switch event {
				case "TRIGGER_AVAILABILITY":
					checkAvailability(roster, db, messenger)

				case "TRIGGER_MATCH":
					notifyPairs(roster.GetMatches(), messenger)
				}
			}
		}
	}()

	// Schedule events
	c := cron.New()

	// Every Wednesday morning check person's availability
	c.AddFunc("0 01 10 * * WED", func() {
		checkAvailability(roster, db, messenger)
	})

	// Every Wednesday afternoon send match results
	c.AddFunc("0 01 13 * * WED", func() {
		notifyPairs(roster.GetMatches(), messenger)
	})

	// Block main process
	select {}
}

func processMessage(m Messaging, messenger *Messenger, roster *Roster, db *Storage) {
	senderID := m.Sender.ID
	log.Printf("recieved message from: %s", senderID)

	employee, ok := roster.GetByID(senderID)
	if !ok {
		m := messenger.GetMember(senderID)

		employee = &Employee{
			ID:           senderID,
			Name:         m.Name,
			FirstName:    m.FirstName,
			Availability: Unavailable,
			Active:       true,
		}
		roster.Add(employee)
	}

	messenger.Send(Messaging{
		Recipient: User{
			ID: senderID,
		},
		SendingAction: typingOn,
	})

	defer messenger.Send(Messaging{
		Recipient: User{
			ID: senderID,
		},
		SendingAction: typingOff,
	})

	if m.Postback != nil {
		if p := m.Postback.Payload; len(p) > 0 {

			// Handle cheat codes
			switch p {
			case "TRIGGER_MATCH":
				fallthrough
			case "TRIGGER_AVAILABILITY":
				log.Printf("Force %s", p)
				if employee.IsAdmin() {
					scheduler <- p
				} else {
					messenger.SendMessage(Messaging{
						Recipient: User{
							ID: senderID,
						},
						Message: &Message{
							Text: "DANGER ZONE. ADMINS ONLY",
						},
					})
				}
				return
			}

			// Handle other payload
			switch p {
			case "SUBSCRIBE_PAYLOAD":
				if employee.Active {
					messenger.SendMessage(Messaging{
						Recipient: User{
							ID: senderID,
						},
						Message: &Message{
							Text: "You already subscribed.",
						},
					})
					return
				}
				employee.Active = true

				fallthrough
			case "GET_STARTED_PAYLOAD":
				log.Printf("%s:%s - subscribed", senderID, employee.Name)
				messenger.SendMessage(Messaging{
					Recipient: User{
						ID: senderID,
					},
					Message: &Message{
						Text: fmt.Sprintf("Hey %s, welcome to Random Coffee, I hope you’re having a great day! Each Wednesday at noon I’ll find a random colleague for you to grab a coffee with. All you have to do is choose which area your office is closest to.", employee.FirstName),
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

			case "UNSUBSCRIBE_PAYLOAD":
				employee.Active = false
				db.SaveEmployee(employee)
				log.Printf("%s:%s - unsubscribed %s", senderID, employee.Name)
				if employee.Active {
					messenger.SendMessage(Messaging{
						Recipient: User{
							ID: senderID,
						},
						Message: &Message{
							Text: "Sorry to see you go, but if you change your mind you can click preferences and then subscribe again.",
						},
					})
				} else {
					messenger.SendMessage(Messaging{
						Recipient: User{
							ID: senderID,
						},
						Message: &Message{
							Text: "You already unsubscribed.",
						},
					})
				}
				return

			case "CHANGE_LOCATION_PAYLOAD":
				messenger.SendMessage(Messaging{
					Recipient: User{
						ID: senderID,
					},
					Message: &Message{
						Text: "Which area office is closest to?",
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
		}
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
			log.Printf("%s:%s - prefers %s", senderID, employee.Name, qr.Payload)

			messenger.SendMessage(Messaging{
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
			employee.Availability = Available
			db.SaveEmployee(employee)
			log.Printf("%s:%s - is available", senderID, employee.Name)
			messenger.SendMessage(Messaging{
				Recipient: User{
					ID: senderID,
				},
				Message: &Message{
					Text: "Perfect. I’ll get back to you around noon with your match. Have a good morning!",
				},
			})

		case "<AVAILABILITY:NO>":
			log.Printf("%s:%s - is unavailable", senderID, employee.Name)
			messenger.SendMessage(Messaging{
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
			db.SaveEmployee(employee)
			log.Printf("%s:%s - unsubscribed", senderID, employee.Name)
			messenger.SendMessage(Messaging{
				Recipient: User{
					ID: senderID,
				},
				Message: &Message{
					Text: "Sorry to see you go, but if you change your mind you can click preferences and then subscribe again.",
				},
			})
		}
	}
}

func checkAvailability(roster *Roster, db *Storage, messenger *Messenger) {
	roster.SetAvailabilityAll(Unavailable)
	log.Printf("Running checkAvailability")
	for _, employee := range roster.Employees {
		if employee.Availability == Unavailable {
			employee.Availability = Unknown
			db.SaveEmployee(employee)
			go messenger.SendMessage(Messaging{
				MessagingType: "UPDATE",
				Recipient: User{
					ID: employee.ID,
				},
				Message: &Message{
					Text: fmt.Sprintf("Good morning %s! Are you available to grab a coffee with someone today?", employee.FirstName),
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
}

// Send notifications to the pairs
func notifyPairs(matches [][]*Employee, messenger *Messenger) {
	log.Printf("Running notifyPairs")
	for _, pair := range matches {
		log.Printf("%s:%s - matched with - %s:%s", pair[0].ID, pair[0].Name, pair[1].ID, pair[1].Name)

		go messenger.SendMessage(Messaging{
			MessagingType: MessagingTypeUpdate,
			Recipient: User{
				ID: pair[0].ID,
			},
			Message: &Message{
				Text: fmt.Sprintf("Hi %s. Your match this week is %s. Shoot them a message on Workplace and organize a time to meet!", pair[0].FirstName, pair[1].Name),
			},
		})

		go messenger.SendMessage(Messaging{
			Recipient: User{
				ID: pair[1].ID,
			},
			Message: &Message{
				Text: fmt.Sprintf("Hi %s. Your match this week is %s. Shoot them a message on Workplace and organize a time to meet!", pair[1].FirstName, pair[0].Name),
			},
		})
	}
}
