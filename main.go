package main

import (
	"fmt"
	"zinx-demo/ziface"
	"zinx-demo/znet"
)

type HelloRouter struct {
	znet.BaseRouter
}

func (prouter *HelloRouter) PreHandle(req ziface.IRequest) {
	fmt.Println("HelloRouter prehandle...")
	fmt.Println("recv from client:msgID = ", req.GetMsgID(),
		",data=", string(req.GetData()))

	err := req.GetConn().SendMsg(1, []byte("hihihihi"))
	if err != nil {
		fmt.Println("send msg error:", err)
		return
	}
}

type PingRouter struct {
	znet.BaseRouter
}

func (prouter *PingRouter) PreHandle(req ziface.IRequest) {
	fmt.Println("PingRouter prehandle...")
}

func (prouter *PingRouter) MainHandle(req ziface.IRequest) {
	fmt.Println("PingRouter MainHandle...")
	//显示要处理的数据
	fmt.Println("recv from client:msgID = ", req.GetMsgID(),
		",data=", string(req.GetData()))

	err := req.GetConn().SendMsg(0, []byte("ping..ping...ping"))
	if err != nil {
		fmt.Println("send msg error:", err)
		return
	}
}

func (prouter *PingRouter) PostHandle(req ziface.IRequest) {
	fmt.Println("PingRouter PostHandle...")
	err := req.GetConn().SendMsg(0, []byte("asd"))
	if err != nil {
		fmt.Println("call back after ping error")
	}
}

var (
	name string
)

func init() {
	name = "znix-v1.0"
}

func main() {
	//创建服务器
	server := znet.NewServer(name)

	//注册链接Hook钩子函数
	server.SetOnConnStart(func(conn ziface.IConnection) {
		fmt.Println("===> Hook Func OnConnStart Called...")
		if err := conn.SendMsg(5, []byte("Connection BEGIN")); err != nil {
			fmt.Println("Conn Send Msg error:", err)
			return
		}

		conn.SetProperty("Name", "aristal")
		conn.SetProperty("Password", "12345")
	})
	server.SetOnConnStop(func(conn ziface.IConnection) {
		fmt.Println("===> Hook Func OnConnStop Called...")

		if val, err := conn.GetProperty("Name"); err == nil {
			fmt.Println("GetProperty:Key=", "Name", "val=", val)
		}
		if val, err := conn.GetProperty("Password"); err == nil {
			fmt.Println("GetProperty:Key=", "Password", "val=", val)
		}
	})

	//添加封装了消息处理方法的路由
	server.AddRouter(0, &PingRouter{})
	server.AddRouter(1, &HelloRouter{})

	server.Serve() //启动服务
}
