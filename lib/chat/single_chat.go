package chat

import (
	"chat/lib"
	"chat/models"
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
	_, _ = s.c1.Emit("matched", s.c2.User)
	_, _ = s.c2.Emit("matched", s.c1.User)
}

func (s *SingleChat) SendMessage(message models.Message)  {

}