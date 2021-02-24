package chat

import (
	"chat/lib"
	"chat/models"
	"encoding/json"
)

type SingleChat struct {
	c1, c2 *lib.Client

	Status int
}

func NewSingleChat(c1, c2 *lib.Client) *SingleChat {
	return &SingleChat{
		c1:     c1,
		c2:     c2,
		Status: 0,
	}
}

func (s *SingleChat) Other(client *lib.Client) *lib.Client {
	if s.c1 == client {
		return s.c2
	}
	return s.c1
}

func (s *SingleChat) SendJoinMessage() {

	c1Response, _ := json.Marshal(lib.UnifiedDataFormat{
		Event: "matched",
		Data:  s.c1.User,
	})

	c2Response, _ := json.Marshal(lib.UnifiedDataFormat{
		Event: "matched",
		Data:  s.c1.User,
	})

	_, _ = s.c1.Conn.Write(c1Response)

	_, _ = s.c2.Conn.Write(c2Response)
}

func (s *SingleChat) SendMessage(message models.Message)  {

}