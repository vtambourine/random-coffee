package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

const (
	markSeen  = "mark_seen"
	typingOn  = "typing_on"
	typingOff = "typing_off"
)

const apiTmpl = "https://graph.facebook.com/v2.12"

type Messenger struct {
	C      chan Messaging // The channel on which the messages are delivered
	secret string         // Bot application key
	token  string         // Verification token
	api    string
	stream chan struct{}
}

type Event struct {
	Object string `json:"object,omitempty"`
	Entry  []struct {
		ID        string      `json:"id,omitempty"`
		Time      int         `json:"time,omitempty"`
		Messaging []Messaging `json:"messaging,omitempty"`
	} `json:"entry,omitempty"`
}

type Messaging struct {
	Sender        User      `json:"sender,omitempty"`
	Recipient     User      `json:"recipient,omitempty"`
	Timestamp     int       `json:"timestamp,omitempty"`
	Message       *Message  `json:"message,omitempty"`
	SendingAction string    `json:"sender_action,omitempty"`
	Postback      *Postback `json:"postback,omitempty"`
}

type User struct {
	ID string `json:"id,omitempty"`
}

type Message struct {
	MID          string        `json:"mid,omitempty"`
	Text         string        `json:"text,omitempty"`
	QuickReply   *QuickReply   `json:"quick_reply,omitempty"`
	QuickReplies *[]QuickReply `json:"quick_replies,omitempty"`
	Attachment   *Attachment   `json:"attachment,omitempty"`
	Attachments  *[]Attachment `json:"attachments,omitempty"`
}

type Postback struct {
	Title   string `json:"title,omitempty"`
	Payload string `json:"payload,omitempty"`
}

type AttachmentType string

const (
	audio    AttachmentType = "audio"
	fallback AttachmentType = "fallback"
	file     AttachmentType = "file"
	image    AttachmentType = "image"
	location AttachmentType = "location"
	video    AttachmentType = "video"
	template AttachmentType = "template"
)

type Attachment struct {
	Type    AttachmentType `json:"type,omitempty"`
	Payload Payload        `json:"payload,omitempty"`
}

type Payload struct {
	URL          string       `json:"url,omitempty"`
	Text         string       `json:"text,omitempty"`
	Coordinates  *Coordinates `json:"coordinates,omitempty"`
	TemplateType string       `json:"template_type,omitempty"`
	Elements     *[]Element   `json:"element,omitempty"`
	Buttons      *[]Button    `json:"buttons,omitempty"`
}

type Coordinates struct {
	Lat  int `json:"lat,omitempty"`
	Long int `json:"long,omitempty"`
}

type Button struct {
	Type    string `json:"type,omitempty"`
	Title   string `json:"title,omitempty"`
	Payload string `json:"payload,omitempty"`
}

type Element struct {
	Title    string `json:"title,omitempty"`
	Subtitle string `json:"subtitle,omitempty"`
}

type QuickReply struct {
	ContentType string `json:"content_type,omitempty"`
	Title       string `json:"title,omitempty"`
	Payload     string `json:"payload,omitempty"`
}

func NewMessenger(secret, token string) *Messenger {
	c := make(chan Messaging, 1)
	m := &Messenger{
		C:      c,
		secret: secret,
		token:  token,
		api:    apiTmpl,
		stream: make(chan struct{}, 50),
	}
	return m
}

// webhookHandler handle the Messenger platform requests
func (m *Messenger) webhookHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		m.verificationEndpoint(w, r)
	case http.MethodPost:
		m.messagesEndpoint(w, r)
	}
}

func (m *Messenger) verificationEndpoint(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	mode := q.Get("hub.mode")
	token := q.Get("hub.verify_token")
	challenge := q.Get("hub.challenge")

	if mode != "" && token == m.token {
		w.WriteHeader(200)
		w.Write([]byte(challenge))
	} else {
		w.WriteHeader(403)
		w.Write([]byte("Wrong verification token"))
	}

}

func (m *Messenger) messagesEndpoint(w http.ResponseWriter, r *http.Request) {
	var event Event
	json.NewDecoder(r.Body).Decode(&event)
	if event.Object == "page" {
		for _, entry := range event.Entry {
			for _, message := range entry.Messaging {
				m.C <- message
				m.Send(Messaging{
					Recipient: User{
						ID: message.Sender.ID,
					},
					SendingAction: markSeen,
				})
			}
		}
	}
}

// Start create the server and register webhooks handler
func (m *Messenger) Start(addr string) {
	http.HandleFunc("/", m.webhookHandler)
	log.Fatal(http.ListenAndServe(addr, nil))
}

func (m *Messenger) Send(message Messaging) {
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(&message)
	url := fmt.Sprintf("%s/me/messages?access_token=%s", m.api, m.secret)
	req, err := http.NewRequest(http.MethodPost, url, body)
	req.Header.Add("Content-Type", "application/json")
	if err != nil {
		log.Fatal(err)
	}

	client := &http.Client{}

	//go func() {
	//	<-m.stream

	resp, err := client.Do(req)
	defer resp.Body.Close()

	if err != nil {
		log.Fatal(err)
	}
	//}()
}

type Member struct {
	ID    string `json:"id,omitempty"`
	Name  string `json:"name,omitempty"`
	Email string `json:"email,omitempty"`
}

func (m *Messenger) GetMember(id string) Member {
	url := fmt.Sprintf("%s/%s?fields=name,email&access_token=%s", m.api, id, m.secret)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	req.Header.Add("Content-Type", "application/json")
	if err != nil {
		log.Fatal(err)
	}

	client := &http.Client{}

	resp, err := client.Do(req)
	defer resp.Body.Close()

	if err != nil {
		log.Fatal(err)
	}

	member := new(Member)
	err = json.NewDecoder(resp.Body).Decode(member)

	if err != nil {
		log.Fatal(err)
	}

	return *member
}
