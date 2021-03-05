package pool

import (
	"chat/models"
)

var ChatRoomPool = ChatRoom{
	pool: map[string]*models.SingleRoom{},
}

type ChatRoom struct {
	pool map[string] *models.SingleRoom
}

func (c *ChatRoom) AddRoom(room *models.SingleRoom) {
	c.pool[room.U1.Id] = room
	c.pool[room.U2.Id] = room
}

func (c *ChatRoom) Leave(user *models.User) *models.SingleRoom {
	room, exist := c.pool[user.Id]

	if ! exist {
		return nil
	}
	delete(c.pool, room.U1.Id)
	delete(c.pool, room.U2.Id)

	return room
}