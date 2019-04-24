package v3

import (
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"net/http"
	"time"
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
	register    chan *Client
	unregister  chan *Client
	clientsPool map[*Client]bool
	send        chan Message
}

// 处理所有websocket请求
func chatRoomHandle(w http.ResponseWriter, r *http.Request) {
	log := zap.S()
	conn, err := ug.Upgrade(w, r, nil)
	if err != nil {
		log.Error(err)
		return
	}

	// 创建客户端
	c := &Client{
		conn: conn,
		send: make(chan Message, 128),
	}

	go c.ReadMessage()
	go c.SendMessage()
	chatRoom.register <- c
}

// 处理所有管道任务
func (room *Room) ProcessTask() {
	log := zap.S()
	log.Info("启动处理任务")
	for {
		select {
		case c := <-room.register:
			log.Info("当前有客户端进行注册")
			room.clientsPool[c] = true
		case c := <-room.unregister:
			log.Info("当前有客户端离开")
			delete(room.clientsPool, c)
		case m := <-room.send:
			time.Sleep(3 * time.Second)
			for c := range room.clientsPool {
				c.send <- m
			}
		default:
			break
		}
	}
}

// 启动聊天室
func StartChatRoom() {
	log := zap.S()
	log.Info("聊天室启动....")
	chatRoom = &Room{
		register:    make(chan *Client),
		unregister:  make(chan *Client),
		clientsPool: map[*Client]bool{},
		send:        make(chan Message),
	}
	http.HandleFunc("/chatroom", chatRoomHandle)
	go chatRoom.ProcessTask()
	if err := http.ListenAndServe(":8081", nil); err != nil {
		panic(err)
	}
}
