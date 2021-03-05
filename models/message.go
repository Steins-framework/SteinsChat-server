package models

type Message struct {
	Sender User `json:"sender"`
	Receiver User `json:"receiver"`
	Time int `json:"time"`
	Text string `json:"text"`
	Key string `json:"key"`
}
