package ziface

type IServer interface {

	//启动
	Start()

	//停止
	Stop()

	//运行
	Serve()

	//添加路由
	AddRouter(msgID uint32, router IRouter)

	//获取当前ConnManager对象
	GetConnManager() IConnManager

	// 注册OnConnStart Hook函数的方法
	SetOnConnStart(func(conn IConnection))

	// 注册OnConnStop Hook函数的方法
	SetOnConnStop(func(conn IConnection))

	// 调用OnConnStart Hook函数的方法
	CallOnConnStart(conn IConnection)

	// 调用OnConnStop Hook函数的方法
	CallOnConnStop(conn IConnection)

}
