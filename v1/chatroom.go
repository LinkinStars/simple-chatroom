package v1

import (
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"net/http"
)

var (
	chatRoom *Room
	ug       = websocket.Upgrader{
		// 允许跨域
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

// 聊天室配置
type Room struct {
	Connections []*websocket.Conn // 连接池，保存所有连接用户
}

// 处理所有websocket请求
func chatRoomHandle(w http.ResponseWriter, r *http.Request) {
	log := zap.S()
	conn, err := ug.Upgrade(w, r, nil)
	if err != nil {
		log.Error(err)
		return
	}
	defer conn.Close()

	// 每次有连接进来就加入连接池
	chatRoom.Connections = append(chatRoom.Connections, conn)

	for {
		// 接收消息
		message := &Message{}
		if err := conn.ReadJSON(message); err != nil {
			log.Error(err)
			return
		}

		// 处理消息
		chatRoom.batchSendMessage(*message)
	}
}

// 群发消息
func (*Room) batchSendMessage(message Message) {
	log := zap.S()
	for _, conn := range chatRoom.Connections {
		if err := conn.WriteJSON(message); err != nil {
			log.Error(err)
			conn.Close()
		}
	}
}

// 启动聊天室
func StartChatRoom() {
	log := zap.S()
	log.Info("聊天室启动....")
	chatRoom = &Room{}
	http.HandleFunc("/chatroom", chatRoomHandle)
	http.ListenAndServe(":8081", nil)
}
