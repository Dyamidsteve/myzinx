package ziface

import "net"

//定义链接模块的抽象接口
type IConnection interface {

	//启动链接，让当前链接开始工作
	Start()

	//停止连接，结束当前链接工作
	Stop()

	//获取当前链接绑定socket conn
	GetTCPConnection() *net.TCPConn

	//获取当前链接模块的链接ID
	GetConnID() uint32

	//获取远程客户端的TCP状态 IP.port
	RemoteAddr() net.Addr

	//发送数据给客户端
	SendMsg(msgId uint32, data []byte) error

	// 设置链接属性
	SetProperty(key string, value interface{})

	// 获取链接属性
	GetProperty(key string) (interface{}, error)

	// 删除链接属性
	RemoveProperty(key string)
}

//type HandleFunc func(*net.TCPConn, []byte, int) error
