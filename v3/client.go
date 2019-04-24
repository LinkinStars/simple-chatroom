package v3

import (
	"github.com/gorilla/websocket"
)

// 客户端
type Client struct {
	conn *websocket.Conn
	send chan Message
}

// 读取消息
func (c *Client) ReadMessage() {
	for {
		// 接收消息
		message := &Message{}
		if err := c.conn.ReadJSON(message); err != nil {
			c.conn.Close()
			chatRoom.unregister <- c
			return
		}

		chatRoom.send <- *message
	}
}

// 发送消息
func (c *Client) SendMessage() {
	for {
		m := <-c.send
		if err := c.conn.WriteJSON(m); err != nil {
			c.conn.Close()
			chatRoom.unregister <- c
		}
	}
}
