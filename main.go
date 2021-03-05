package main

import (
	"chat/lib"
	"chat/lib/chat"
	"chat/models"
	"container/list"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	_ "net/http/pprof"
	"os"
)

var userConnMap = make(map[string]*lib.Client, 10)

var waitMatch = list.New()

var chatRooms = make(map[string]*chat.SingleChat) // user_id => SingleChat

var other = models.User{
	Id:     "AAAAAAAAA",
	Sex:    1,
	Name:   "Mitsuha",
	Avatar: "",
}

func main() {
	go func() {
		_ = http.ListenAndServe("0.0.0.0:9999", nil)
	}()

	listener, err := net.Listen("tcp", "0.0.0.0:9966")
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
		client.On("leave", handleLeave)
		//client.On("input", handleInput)
		//client.On("reading", handleRead)
		client.On("heartbeat", heartbeat)
		client.On("_close", handleClose)

		go client.Listen()
	}

}

func handleRegister(client *lib.Client, data []byte){
	user := models.User{}

	err := json.Unmarshal(data, &user)

	if err != nil{
		fmt.Println(err)
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

	if client.User.Id != message.Sender.Id {
		return
	}
	///////////
	//message.Sender, message.Receiver = message.Receiver , message.Sender
	//////////////

	receiver, exist := userConnMap[message.Receiver.Id]

	if ! exist {
		return
	}

	fmt.Printf("send %s to %s\n", message.Text, message.Receiver.Id)

	_, _ = receiver.Emit("message", message)
	_, _ = client.Emit("message", message)
}

func handleMatching(client *lib.Client, _ []byte) {

	//client.Emit("matched", other)
	//
	//return

	if client.User.Id == "" {
		return
	}
	var other *lib.Client
	for {
		if waitMatch.Len() == 0 {
			waitMatch.PushBack(client)
			client.Status = lib.ClientMatch
			return
		}

		other = waitMatch.Remove(waitMatch.Front()).(*lib.Client)

		if other.User.Id == client.User.Id {
			continue
		}

		if other.IsConnect() {
			break
		}
	}

	room := chat.NewSingleChat(client,other)

	room.SendJoinMessage()

	chatRooms[client.User.Id] = room
	chatRooms[other.User.Id] = room

	client.Status = lib.ClientChat
	other.Status = lib.ClientChat
}

func handleLeave(client *lib.Client, _ []byte) {
	room, exist := chatRooms[client.User.Id]
	if ! exist {
		return
	}
	other := room.Other(client)

	_, _ = other.Emit("leave", room)

	client.Status = lib.ClientWait
	other.Status = lib.ClientWait
	delete(chatRooms, client.User.Id)
	delete(chatRooms, other.User.Id)
}

func heartbeat(client *lib.Client, ping []byte) {
	_ ,_ = client.Emit("heartbeat", string(ping))
}

func handleClose(client *lib.Client, _ []byte) {
	fmt.Printf("Connect close:%s\n", client.Conn.RemoteAddr().String())

	handleLeave(client, nil)

	delete(userConnMap, client.User.Id)
	delete(chatRooms, client.User.Id)
}