package v3

import (
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"time"
)

// 客户端
type Client struct {
	conn *websocket.Conn
	send chan Message
}

// 读取消息
func (c *Client) ReadMessage() {
	preMessageTime := int64(0)
	for {
		// 接收消息
		message := &Message{}
		if err := c.conn.ReadJSON(message); err != nil {
			c.conn.Close()
			chatRoom.unregister <- c
			return
		}

		// 限制用户发送消息频率，每1秒只能发送一条消息
		curMessageTime := time.Now().Unix()
		if curMessageTime-preMessageTime < 1 {
			zap.S().Warn("1秒内无法再次发送")
			continue
		}
		preMessageTime = curMessageTime

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
			return
		}
	}
}
