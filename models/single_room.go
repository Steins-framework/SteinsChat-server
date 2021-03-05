package models

type SingleRoom struct {
	U1 *User `json:"u1"`
	U2 *User `json:"u2"`
	Topic string `json:"topic"`
	Status int `json:"status"`
}

func (s *SingleRoom) Other(user *User) *User {
	if s.U1.Id == user.Id {
		return s.U2
	}
	return s.U1
}
