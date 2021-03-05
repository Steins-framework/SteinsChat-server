package models

type User struct {
	Id string `json:"id"`
	Sex int `json:"sex"`
	Name string `json:"name"`
	Avatar string `json:"avatar"`
	Coordinate []string `json:"coordinate"`
	Tags []string `json:"tags"`
}