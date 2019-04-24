package v3

// 消息结构
type Message struct {
	Token   string `json:"token"`   //token用于表示发起用户
	Content string `json:"content"` //content表示消息内容，这里简化别的消息
}
