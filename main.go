package main

import (
	"chat/lib"
	"chat/lib/chat"
	"chat/models"
	"container/list"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"time"
)

var userConnMap = make(map[int]*lib.Client, 10)

var waitMatch = list.New()

var chatRooms = make([]chat.SingleChat, 10)

var other = models.User{
	Id:     9999,
	Sex:    1,
	Name:   "Mitsuha",
	Avatar: "",
}

func main() {
	listener, err := net.Listen("tcp", "0.0.0.0:65535")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("Start Server...")

	go matching()
	go func() {
		for {
			time.Sleep(5*time.Second)
			fmt.Printf("Now client connect: %d\n", len(userConnMap))
		}
	}()

	go func() {
		for {
			time.Sleep(5 * time.Second)

			for _, client := range userConnMap {
				if client == nil {
					continue
				}
				response, _ := json.Marshal(lib.UnifiedDataFormat{
					Event: "leave",
				})

				client.Conn.Write(response)
			}
		}
	}()
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
	for {
		if waitMatch.Len() < 2 {
			continue
		}
		c1 := waitMatch.Remove(waitMatch.Back()).(*lib.Client)
		c2 := waitMatch.Remove(waitMatch.Front()).(*lib.Client)

		room := chat.NewSingleChat(c1,c2)

		room.SendJoinMessage()

		c1.Status = lib.ClientChat
		c2.Status = lib.ClientChat
	}
}

func handleRegister(client *lib.Client, data []byte){
	user := models.User{}

	err := json.Unmarshal(data, &user)

	if err != nil{
		fmt.Println(err)
		return
	}

	if user.Id == 0 {
		return
	}

	fmt.Println(user)

	client.User = user
	client.Status = lib.ClientNone

	userConnMap[user.Id] = client
}

func handleMessage(client *lib.Client, data []byte)  {
	message := models.Message{}

	err := json.Unmarshal(data, &message)

	if err != nil {
		fmt.Println(err)
		return
	}

	sender := message.Sender
	message.Sender = message.Receiver
	message.Receiver = sender

	receiver, exist := userConnMap[message.Receiver.Id]

	if ! exist {
		return
	}

	fmt.Printf("send %s to %d\n", message.Text, message.Receiver.Id)

	responseData, _ := json.Marshal(lib.UnifiedDataFormat{
		Event: "message",
		Data:  message,
	})

	_, err = receiver.Conn.Write(responseData)

	if err != nil {
		_ = client.Conn.Close()
		fmt.Println("Connect Close:" + client.Conn.RemoteAddr().String())
	}

}

func handleMatching(client *lib.Client, _ []byte) {
	//if client.Status == lib.ClientWait {
	//	return
	//}
	//waitMatch.PushBack(client)
	//client.Status = lib.ClientMatch

	c2Response, _ := json.Marshal(lib.UnifiedDataFormat{
		Event: "matched",
		Data:  other,
	})

	client.Conn.Write(c2Response)
}

func handleLeave(client *lib.Client, _ []byte) {

}