package fwc

import "github.com/timdrysdale/hub"

type Hub struct {
	Messages *hub.Hub
	Clients  map[string]*Client //map Id string to client
	Rules    map[string]Rule    //map Id string to Rule
	Add      chan Rule
	Delete   chan string //Id string
}

type Rule struct {
	Id       string
	Stream   string
	Filename string
}

type Client struct {
	Hub      *Hub
	Messages *hub.Client
	Stopped  chan struct{}
	Filename string
}
