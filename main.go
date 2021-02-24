package main

import (
	"chat/lib"
	"chat/models"
	"container/list"
	"encoding/json"
	"fmt"
	"net"
	"os"
)

var userConnMap = make(map[int]*lib.Client, 10)

var waitMatch = list.New()

func main() {
	listener, err := net.Listen("tcp", "0.0.0.0:65535")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("Start Server...")
	for {
		conn, err := listener.Accept()

		if err != nil {
			_ = conn.Close()
			fmt.Println("Connect close:" + conn.RemoteAddr().String())
			continue
		}

		fmt.Printf("New connect form: %s\n", conn.RemoteAddr().String())

		client := lib.Client{Conn: conn}

		client.On("register", handleRegister)
		client.On("message", handleMessage)
		client.On("matching", handleMatching)

		go client.Listen()
	}

}

func matching() {
	c1 := waitMatch.Back().Value.(lib.Client)
	c2 := waitMatch.Front().Value.(lib.Client)

	c1Response, _ := json.Marshal(lib.UnifiedDataFormat{
		Event: "matched",
		Data:  c2.User,
	})

	c2Response, _ := json.Marshal(lib.UnifiedDataFormat{
		Event: "matched",
		Data:  c1.User,
	})

	_, err := c1.Conn.Write(c1Response)

	if err != nil {
		waitMatch.PushBack(c2)
		return
	}

	_, err = c2.Conn.Write(c2Response)

	if err != nil {
		waitMatch.PushFront(c1)
		return
	}
}

func handleRegister(client *lib.Client, data []byte){
	user := models.User{}

	err := json.Unmarshal(data, &user)

	if err != nil {
		fmt.Println(err)
		return
	}

	client.User = user

	userConnMap[user.Id] = client
}

func handleMessage(client *lib.Client, data []byte)  {
	message := models.Message{}

	err := json.Unmarshal(data, &message)

	if err != nil {
		fmt.Println(err)
		return
	}

	receiver, exist := userConnMap[message.Receiver.Id]

	if ! exist {
		return
	}

	fmt.Printf("send: %s", message.Text)

	responseData, _ := json.Marshal(message)

	_, err = receiver.Conn.Write(responseData)

	if err != nil {
		_ = client.Conn.Close()
		fmt.Println("Connect Close:" + client.Conn.RemoteAddr().String())
	}

}

func handleMatching(client *lib.Client, _ []byte) {
	waitMatch.PushBack(client)
}