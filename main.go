package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"
	"strings"
)

func reply(m Messaging, msngr *Messenger) {
	//sender := m.Sender.ID

	// Check, if sender is already exists
}
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

	employees := make(map[string]*Employee)
	var sender string

	//employeesChans := make(map[string]chan struct{})

	for {
		select {
		case m := <-msngr.C:
			sender = m.Sender.ID

			log.Printf("recieved message from: %s", sender)

			//if e, ok := employees[sender]; ok {
			//	log.Printf("found by %s, %s", sender, e.Name)
			//} else {
			//	log.Printf("Not found by %s", sender)
			//	e := &Employee{
			//		Name: "Not found",
			//	}
			//	employees[sender] = e
			//	//(*e).Name = "Not found"
			//}

			//continue

			// If postback is not empty, the message is one of the quick reply responses
			if m.Postback != nil {
				log.Printf("Got postback: %#v", m.Postback)

				employee, ok := employees[sender]
				if !ok {
					continue
				}

				//employeeChans, _ := employeesChans[sender]

				switch m.Postback.Payload {
				case "weekely":
					employee.Frequency = Weekly

				case "biweekely":
					employee.Frequency = Biweekly

				}

				msngr.Send(Messaging{
					Recipient: User{
						ID: sender,
					},
					Message: &Message{
						Text: fmt.Sprintf("You will be invited %n", employee.Frequency),
					},
				})

				continue
			}

			// Handle Office selection
			if qr := m.Message.QuickReply; qr != nil {

				employee, ok := employees[sender]
				if !ok {
					continue
				}

				switch qr.Payload {
				case "AMS9:AMS10":
					fallthrough
				case "AMS3:AMS11":
					fallthrough
				case "AMS17:AMS19":
					for _, o := range strings.Split(qr.Payload, ":") {
						log.Printf("employee %#v wants to meet at %s", employee, o)
					}
					employee.PreferredLocation = qr.Payload
				}

				//continue
			}


			if e, ok := employees[sender]; ok {
				//msngr.Send(Messaging{
				//	Recipient: User{
				//		ID: sender,
				//	},
				//	Message: &Message{
				//		Text: fmt.Sprintf("REPEAT GREETING"),
				//	},
				//})
			} else {
				e = &Employee{
					Name: "John Doe",
				}

				employees[sender] = e

				log.Println("send greeting ")
				msngr.Send(Messaging{
					Recipient: User{
						ID: sender,
					},
					Message: &Message{
						Text: fmt.Sprintf("GREETING"),
					},
				})
			}

			employee := *employees[sender]

			// Person didn't confirm location
			log.Println(employee.PreferredLocation)
			//time.Sleep(2 * time.Second)
			if len(employee.PreferredLocation) == 0 {
				log.Println("send office ")
				msngr.Send(Messaging{
					Recipient: User{
						ID: sender,
					},
					Message: &Message{
						Text: "Pam Pam",
						QuickReplies: &[]QuickReply{
							{
								ContentType: "text",
								Title: "rembrandt",
								Payload: "AMS9:AMS10",
							},
							{
								ContentType: "text",
								Title: "vijzel",
								Payload: "AMS3:AMS11",
							},
							{
								ContentType: "text",
								Title: "piethein",
								Payload: "AMS17:AMS19",
							},
						},
					},
				})
			} else {
				msngr.Send(Messaging{
					Recipient: User{
						ID: sender,
					},
					Message: &Message{
						Text: fmt.Sprintf("PREFERRED LOCATION RECEIVED"),
					},
				})
				msngr.Send(Messaging{
					Recipient: User{
						ID: sender,
					},
					Message: &Message{
						Text: fmt.Sprintf("TOT ZIENS"),
					},
				})
			}

			//if employee.Frequency == 0 {
			//	msngr.Send(Messaging{
			//		Recipient: User{
			//			ID: m.Sender.ID,
			//		},
			//		Message: &Message{
			//			Attachment: &Attachment{
			//				Type: "template",
			//				Payload: Payload{
			//					TemplateType: "button",
			//					Text: "How often you want two be invited?",
			//					Buttons: &[]Button {
			//						{
			//							Title: "Every one week",
			//							Type: "postback",
			//							Payload: "weekely",
			//						},
			//						{
			//							Title: "Every two week",
			//							Type: "postback",
			//							Payload: "biweekely",
			//						},
			//					},
			//				},
			//			},
			//		},
			//	})
			//} else {
			//	msngr.Send(Messaging{
			//		Recipient: User{
			//			ID: sender,
			//		},
			//		Message: &Message{
			//			Text: fmt.Sprintf("You will be invited %n", employee.Frequency),
			//		},
			//	})
			//}
		}
	}

	offices := []Office{AMS9, AMS10, AMS11}
	names := []string{
		"Shanel", "Al", "Sandee", "Jimmie",
		"Luella", "Neoma", "Edmond", "Marilee",
		"Kristy", "Deonna", "Margo", "Bethany",
		"Cuc", "Kathlene", "Mica", "Shanti",
		"Joycelyn", "Norbert", "Ardath", "Nichell",
	}
	frequencies := []Frequency{Weekly, Biweekly, Triweekly}

	var e Employee
	for i := 0; i <= 20; i++ {
		e = Employee{
			Name:      names[rand.Intn(len(names))],
			Office:    offices[rand.Intn(len(offices))],
			Frequency: frequencies[rand.Intn(len(frequencies))],
		}
		employees[string(i)] = &e

		fmt.Printf("%v\n", e)
	}

	//msgr := NewMatcher()
	//msgr.Add(alice)
	//msgr.Add(bob)

	// Restore employees from database

	// From all employees choose who

	//fmt.Printf("%v", msgr.GetMatches())

	//ticker := time.NewTicker(10 * time.Second)
	//defer ticker.Stop()
	//
	//done := make(chan bool)
	//go func() {
	//	time.Sleep(10 * time.Second)
	//	//done <- true
	//}()
	//match(&employees)
	//for {
	//	select {
	//	case <-done:
	//		fmt.Println("Done!")
	//		return
	//	case <-ticker.C:
	//		match(&employees)
	//	}
	//}
}

func match(employees *[]Employee) {
	fmt.Print("MATCHING ")
	fmt.Println(time.Now().Format(time.UnixDate))

	match := NewMatcher()
	for _, e := range *employees {
		match.Add(&e)
	}

	fmt.Printf("Match: %v\n\n", match.GetMatches())
}
