package main

import (
	"github.com/LinkinStars/golang-util/gu"
	"github.com/LinkinStars/simple-chatroom/v3"
	_ "net/http/pprof"
)

func main() {
	// 初始化日志（可以忽略，日志好看一些）
	gu.InitEasyZapDefault("simple-chatroom")
	// 启动！
	v3.StartChatRoom()
}
