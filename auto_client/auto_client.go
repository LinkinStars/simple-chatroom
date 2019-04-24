package main

import (
	"fmt"
	"github.com/LinkinStars/simple-chatroom/v3"
	"github.com/gorilla/websocket"
	"math/rand"
	"net/url"
	"strconv"
	"time"
)

type Sender struct {
	conn *websocket.Conn
	send chan v3.Message
}

func main() {
	clients := createClients(10)
	process(clients)
}

// 创建一定数量的客户端
func createClients(amount int) []*Sender {
	var clients []*Sender
	u := url.URL{Scheme: "ws", Host: "127.0.0.1:8081", Path: "/chatroom"}
	for i := 0; i < amount; i++ {
		c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
		if err != nil {
			panic(err)
		}
		sender := &Sender{
			conn: c,
			send: make(chan v3.Message, 128),
		}
		go sender.loopSendMessage()
		clients = append(clients, sender)
	}
	return clients
}

// 随机的选择一些用户去发送消息
func process(clients []*Sender) {
	// 设置随机种子
	rand.Seed(time.Now().UnixNano())

	clientsAmount := len(clients) - 1

	flag := 0
	for {
		time.Sleep(time.Millisecond * 1)
		go func() {
			randIndex := rand.Intn(clientsAmount)
			clients[randIndex].sendMessage(strconv.Itoa(flag))
			flag++
		}()
	}
}

// 写入消息
func (sender *Sender) sendMessage(str string) {
	message := v3.Message{
		Token:   time.Now().Format(time.RFC3339),
		Content: str,
	}
	sender.send <- message
}

// 循环发送消息
func (sender *Sender) loopSendMessage() {
	for {
		m := <-sender.send
		if err := sender.conn.WriteJSON(m); err != nil {
			fmt.Println(err)
		}
		fmt.Println("发送消息", m)
	}
}
