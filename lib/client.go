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
	event map[string]func(*Client, []byte)
}

func (c *Client) On(event string, closure func(*Client, []byte)) {
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

func (c *Client) Listen() {
	for {
		receive := make([]byte,1024)

		l, err := c.Conn.Read(receive)
		if err != nil {
			_ = c.Conn.Close()
			c.Trigger(UnifiedDataFormat{Event: "_close"})
		}

		receive = receive[:l]

		if l == 1 {
			continue
		}

		fmt.Printf("receiv %d:%s\n", l, string(receive))

		format := UnifiedDataFormat{}

		err = json.Unmarshal(receive, &format)

		if err != nil {
			fmt.Println(err)
			continue
		}

		c.Trigger(format)
	}
}