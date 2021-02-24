package models

type Message struct {
	Sender User
	Receiver User
	Time string
	Text string
	Key string
}
