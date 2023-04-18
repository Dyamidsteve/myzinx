package znet

import (
	"fmt"
	_ "io"
	"net"

	"zinx-demo/utils"
	"zinx-demo/ziface"
)

type Server struct {
	IPVersion   string                        //IP版本号
	Name        string                        //服务器名称
	Ip          string                        //Ip
	Port        int                           //端口
	MsgHandler  ziface.IMsgHandle             //消息管理模块
	ConnManager ziface.IConnManager           //链接管理模块
	OnConnStart func(conn ziface.IConnection) //连接开始处理Hook函数
	OnConnStop  func(conn ziface.IConnection) //连接停止处理Hook函数
}

// 启动服务
func (s *Server) Start() {
	fmt.Printf("[Start:%s]Server start listener at Ip:%s port:%d\n", s.Name, s.Ip, s.Port)
	// 开启工作池
	s.MsgHandler.StartWorkerPool()

	//获取TCP的Addr
	addr, _ := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.Ip, s.Port))
	//tpc4地址监听
	listener, err := net.ListenTCP(s.IPVersion, addr)
	if err != nil {
		fmt.Println("Listen error:", err)
		return
	}
	var cid uint32
	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			fmt.Println("Listen accepct error:", err)
			continue //出现error跳过该请求
		}
		//连接数超出则取消连接
		if s.ConnManager.GetConnLen() >= utils.GlobalConf.MaxConn {
			// 给客户端相应一个超出最大连接的消息包
			fmt.Println("Warning:Exceed the MaxConnSize:", utils.GlobalConf.MaxConn)
			conn.Close()
			continue
		}

		//将链接的路由和本身的路由进行绑定
		go NewConnection(s, conn, cid, s.MsgHandler).Start()
		cid++

	}

}

// 停止服务
func (s *Server) Stop() {

	fmt.Println("Server Stop...")
	s.ConnManager.Clear()
}

// 运行服务
func (s *Server) Serve() {
	//异步启动服务器
	go s.Start()

	//启动服务器后业务

	//死锁
	select {}
}

// 添加路由
func (s *Server) AddRouter(msgID uint32, router ziface.IRouter) {
	s.MsgHandler.AddRouter(msgID, router)
	fmt.Println("Add Router success...")
}

// 获取当前ConnMagager
func (s *Server) GetConnManager() ziface.IConnManager {
	return s.ConnManager
}

// 注册OnConnStart Hook函数的方法
func (s *Server) SetOnConnStart(hookFunc func(conn ziface.IConnection)) {
	s.OnConnStart = hookFunc
}

// 注册OnConnStop Hook函数的方法
func (s *Server) SetOnConnStop(hookFunc func(conn ziface.IConnection)) {
	s.OnConnStop = hookFunc
}

// 调用OnConnStart Hook函数的方法
func (s *Server) CallOnConnStart(conn ziface.IConnection) {
	if s.OnConnStart != nil {
		fmt.Println("-----> Call OnConnStart()...")
		s.OnConnStart(conn)
	}
}

// 调用OnConnStop Hook函数的方法
func (s *Server) CallOnConnStop(conn ziface.IConnection) {
	if s.OnConnStop != nil {
		fmt.Println("-----> Call OnConnStart()...")
		s.OnConnStop(conn)
	}
}

// 初始化Serveer
func NewServer(name string) ziface.IServer {
	server := &Server{
		IPVersion:   "tcp4",
		Name:        utils.GlobalConf.Name,
		Ip:          utils.GlobalConf.Host,
		Port:        utils.GlobalConf.TcpPort,
		MsgHandler:  NewMsgHandler(),
		ConnManager: NewConnManager(),
		
	}

	return server
}
