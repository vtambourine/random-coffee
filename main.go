package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"
)

func reply(m Messaging, msngr *Messenger) {
	//sender := m.Sender.ID

	// Check, if sender is already exists
}
func main() {
	log.Println("Random Coffee initialized")

	addr := "localhost:3031"
	msngr := NewMessenger("EAAGUCR82yskBAEuLKH3iufeRtiSV8vSZBOp0EmZBcUbtkNoJxyt2vOjqT87ZChQaAdHzRgSntGaZAnTMgtbNIBc7BDggF3zRlZB3NS50U8v0xA8krYdFi8TSllZA0geZCr2iSuS4FrbnhmP6dtZBnfF5LTc9YnWIb46GL5jAsp5ZAzwZDZD", os.Getenv("VERIFY_TOKEN"))
	go msngr.Start(addr)
	log.Printf("started at %s\n", addr)

	employees := make(map[string]*Employee)
	var sender string

	for {
		select {
		case m := <-msngr.C:
			sender = m.Sender.ID
			log.Printf("recieved message from: %s", sender)
			//
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

			if m.Postback != nil {
				log.Printf("Got postback: %#v", m.Postback)
				employee, ok := employees[sender]
				if !ok {
					continue
				}

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

			if e, ok := employees[sender]; ok {
				msngr.Send(Messaging{
					Recipient: User{
						ID: sender,
					},
					Message: &Message{
						Text: fmt.Sprintf("Hello, %s", e.Name),
					},
				})
			} else {
				e = &Employee{
					Name: "John Doe",
				}

				employees[sender] = e

				msngr.Send(Messaging{
					Recipient: User{
						ID: sender,
					},
					Message: &Message{
						Text: fmt.Sprintf("Hello, %s", e.Name),
					},
				})
			}

			log.Printf("%#v", employees)

			employee := *employees[sender]

			if employee.Frequency == 0 {
				msngr.Send(Messaging{
					Recipient: User{
						ID: m.Sender.ID,
					},
					Message: &Message{
						Attachment: &Attachment{
							Type: "template",
							Payload: Payload{
								TemplateType: "button",
								Text: "How often you want two be invited?",
								Buttons: &[]Button {
									{
										Title: "Every one week",
										Type: "postback",
										Payload: "weekely",
									},
									{
										Title: "Every two week",
										Type: "postback",
										Payload: "biweekely",
									},
								},
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
						Text: fmt.Sprintf("You will be invited %n", employee.Frequency),
					},
				})
			}
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
