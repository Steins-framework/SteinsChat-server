package main

import (
	"chat/lib"
	"chat/lib/pool"
	"chat/models"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	_ "net/http/pprof"
	"os"
)

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

	client.User = &user
	client.Status = lib.ClientNone

	pool.UserConnPool.Add(user.Id, client)
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

	c1, err := pool.UserConnPool.OfUser(message.Sender)
	if err == nil {
		_, _ = c1.Emit("message", message)
	}
	c2, err := pool.UserConnPool.OfUser(message.Receiver)
	if err == nil {
		_, _ = c2.Emit("message", message)
	}
}

func handleMatching(client *lib.Client, _ []byte) {
	if client.User.Id == "" {
		return
	}
	var other *models.User
	for {
		if pool.WaitMatchPool.Len() == 0 {
			pool.WaitMatchPool.Add(client.User)
			return
		}

		other = pool.WaitMatchPool.Pop()

		if other.Id == client.User.Id {
			continue
		}

		if pool.UserConnPool.IsOnline(other) {
			break
		}
	}

	room := &models.SingleRoom{
		U1:     client.User,
		U2:     other,
		Topic:  "",
		Status: 0,
	}

	c1, err := pool.UserConnPool.OfUser(room.U1)

	if err == nil {
		_, _ = c1.Emit("matched", room)
	}

	c2, err := pool.UserConnPool.OfUser(room.U2)

	if err == nil {
		_, _ = c2.Emit("matched", room)
	}

	pool.ChatRoomPool.AddRoom(room)
}

func handleLeave(client *lib.Client, _ []byte) {
	room := pool.ChatRoomPool.Leave(client.User)

	fmt.Print(room)
	if room == nil {
		return
	}

	c2, err := pool.UserConnPool.OfUser(room.Other(client.User))

	fmt.Println(err)
	if err == nil {
		fmt.Print(c2)
		_, _ = c2.Emit("leave", room)
	}
}

func heartbeat(client *lib.Client, ping []byte) {
	_ ,_ = client.Emit("heartbeat", string(ping))
}

func handleClose(client *lib.Client, _ []byte) {
	fmt.Printf("Connect close:%s\n", client.Conn.RemoteAddr().String())

	if client.User == nil {
		return
	}

	handleLeave(client, nil)

	pool.UserConnPool.Delete(client.User.Id)
}