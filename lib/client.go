package lib

import (
	"chat/models"
	"encoding/json"
	"fmt"
	"net"
)

type Client struct {
	Conn net.Conn
	User models.User
	Status int
	event map[string]func(*Client, []byte)
}

const (
	ClientNone = iota
	ClientWait
	ClientChat
	ClientMatch
)

func (c *Client) On(event string, closure func(*Client, []byte)) {
	if c.event == nil {
		c.event = make(map[string]func(*Client, []byte))
	}

	c.event[event] = closure
}

func (c *Client) Trigger(format UnifiedDataFormat) {
	bytes, err := format.getDataBytes()

	if err != nil {
		fmt.Println(err)
		return
	}

	if closure, ok := c.event[format.Event]; ok {
		closure(c, bytes)
	}
}

func (c *Client) Close() error{
	defer c.Trigger(UnifiedDataFormat{Event: "_close"})

	return c.Conn.Close()
}

func (c *Client) Listen() {
	for {
		receive := make([]byte,1024)

		l, err := c.Conn.Read(receive)
		if err != nil {
			_ = c.Conn.Close()
			c.Trigger(UnifiedDataFormat{Event: "_close"})
		}

		receive = receive[:l]

		if l < 2 {
			continue
		}

		fmt.Printf("receiv %d:%s\n", l, string(receive))

		format := UnifiedDataFormat{}

		err = json.Unmarshal(receive, &format)

		if err != nil {
			continue
		}

		c.Trigger(format)
	}
}