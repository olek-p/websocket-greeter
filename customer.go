package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

type Customer struct {
	id     string
	urls   []string
	fields map[string]string
	ws     *websocket.Conn
}
type Customers struct {
	customers map[string]map[string]Customer
}

func (c *Customers) Init() {
	c.customers = make(map[string]map[string]Customer)
}

func (c *Customers) Get(licence string) (map[string]Customer, error) {
	if customers, ok := c.customers[licence]; ok {
		return customers, nil
	}
	return nil, fmt.Errorf("No customers found for licence %s", licence)
}

func (c *Customers) Add(licence, id, url string, fields map[string]string, ws *websocket.Conn) error {
	customers, err := c.Get(licence)
	if err != nil {
		c.customers[licence] = make(map[string]Customer)
		customers = c.customers[licence]
	}
	if customer, ok := customers[id]; ok {
		customer.urls = append(customer.urls, url)
		for k, v := range fields {
			customer.fields[k] = v
		}
		customer.CheckGreetable()
	} else {
		customers[id] = Customer{id, []string{url}, fields, ws}
		customers[id].CheckGreetable()
	}

	return nil
}

func (c Customer) CheckGreetable() {
	for _, url := range c.urls {
		if strings.Contains(url, "test=1") {
			c.AddGreeting()
		}
	}
}

func (c Customer) AddGreeting() {
	go func() {
		<-time.After(time.Second * 5)
		fmt.Println("Sending a greeting")
		sendToWs(c.ws, "Greetings!")
	}()
}
