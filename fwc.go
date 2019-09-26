package fwc

import (
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/timdrysdale/hub"
)

// pass in the messaging hub as a parameter
// assume it is already running
func New(messages *hub.Hub) *Hub {

	h := &Hub{
		Messages: messages,
		Clients:  make(map[string]*Client), //map Id string to Client
		Rules:    make(map[string]Rule),    //map Id string to Rule
		Add:      make(chan Rule),
		Delete:   make(chan string), //Id string
	}

	return h
}

func (h *Hub) Run(closed chan struct{}) {

	for {
		select {
		case <-closed:
			return
		case rule := <-h.Add:

			// Delete any pre-existing client for this rule.Id
			// because it just became superseded
			if client, ok := h.Clients[rule.Id]; ok {
				client = h.Clients[rule.Id]
				close(client.Stopped) //stop writing to file
				delete(h.Clients, rule.Id)
			}
			if _, ok := h.Rules[rule.Id]; ok {
				delete(h.Rules, rule.Id)
			}

			//record the new rule for later convenience in reporting
			h.Rules[rule.Id] = rule

			// create client to handle stream messages
			messageClient := &hub.Client{Hub: h.Messages,
				Name:  rule.Filename,
				Topic: rule.Stream,
				Send:  make(chan hub.Message),
				Stats: hub.NewClientStats()}

			client := &Client{Hub: h,
				Messages: messageClient,
				Stopped:  make(chan struct{}),
				Filename: rule.Filename,
			}

			h.Clients[rule.Id] = client

			h.Messages.Register <- client.Messages //register for messages from hub

			go client.Writer()

		case ruleId := <-h.Delete:
			if client, ok := h.Clients[ruleId]; ok {
				close(client.Stopped)
				delete(h.Clients, ruleId)
			}
			if _, ok := h.Rules[ruleId]; ok {
				delete(h.Rules, ruleId)
			}
		}
	}
}

func (c *Client) Writer() {

	// open file
	f, err := os.OpenFile(c.Filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.WithFields(log.Fields{"filename": c.Filename, "error": err}).Error("Failed opening file for writing")
		return
	}
	defer f.Close()

	// write messages to file until we are stopped
	for {
		select {
		case <-c.Stopped:
			break
		case msg := <-c.Messages.Send:
			_, err := f.Write(msg.Data)
			if err != nil {
				log.WithFields(log.Fields{"filename": c.Filename, "error": err}).Error("Failed writing to file")
				return
			}
		}
	}
}
