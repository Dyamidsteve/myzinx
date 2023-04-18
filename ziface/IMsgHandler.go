package ziface

// 消息处理模块接口
type IMsgHandle interface {
	//调度/执行对应的Router消息处理方法
	DoMsgHandler(req IRequest)

	// 为消息添加具体的处理逻辑
	AddRouter(msgID uint32, router IRouter)

	// 启动一个Worker工作池
	StartWorkerPool()

	// 将消息发给taskQueue，由worker进行处理
	SendMsgToTaskQueue(req IRequest)
}
