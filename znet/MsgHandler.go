package znet

import (
	"fmt"
	"strconv"
	"zinx-demo/utils"
	"zinx-demo/ziface"
)

// 消息处理模块实现
type MsgHandle struct {
	// 存放每个msgID对应的处理方法
	Apis map[uint32]ziface.IRouter

	//负责worker取任务的消息队列
	TaskQueue []chan ziface.IRequest

	// Worker池的消息队列数量
	WorkerPoolSize uint32
}

func NewMsgHandler() *MsgHandle {
	mh := &MsgHandle{
		Apis:           make(map[uint32]ziface.IRouter),
		WorkerPoolSize: utils.GlobalConf.WorkerPoolSize,
		TaskQueue:      make([]chan ziface.IRequest, utils.GlobalConf.WorkerPoolSize),
	}
	return mh
}

// 调度/执行对应的Router消息处理方法
func (mh *MsgHandle) DoMsgHandler(req ziface.IRequest) {
	// 从req中的msgID找到handle中对应的router
	router, ok := mh.Apis[req.GetMsgID()]
	if !ok {
		//不存在对应的msgID
		panic("msgID " + strconv.Itoa(int(req.GetMsgID())) + " not match in Apis")
	}
	// 处理消息
	router.PreHandle(req)
	router.MainHandle(req)
	router.PostHandle(req)
}

// 为消息添加具体的处理逻辑
func (mh *MsgHandle) AddRouter(msgId uint32, router ziface.IRouter) {
	if _, ok := mh.Apis[msgId]; ok {
		//id 已注册
		//panic停止当前协程(defer过的func仍然会调度)
		panic("repeat api,msgID=" + strconv.Itoa(int(msgId)))
	}
	mh.Apis[msgId] = router
	fmt.Println("Add router msgId:", msgId, "success")
}

// 启动一个Worker工作池
func (mh *MsgHandle) StartWorkerPool() {
	// 根据workerPoolSize分别开启Worker
	for i := 0; i < int(mh.WorkerPoolSize); i++ {

		//为每个worker对应的消息队列开辟空间
		mh.TaskQueue[i] = make(chan ziface.IRequest, utils.GlobalConf.MaxWorkerTaskLen)
		//启动当前worker,阻塞等待消息从chanel传递进来
		go mh.StartOneWorker(i, mh.TaskQueue[i])
	}
}

// 启动一个Worker工作流程
func (mh *MsgHandle) StartOneWorker(workerID int, taskQueue chan ziface.IRequest) {
	fmt.Println("Worker ID=", workerID, "started...")
	// 阻塞等待消息并处理
	for {
		select {
		case req := <-taskQueue:
			mh.DoMsgHandler(req)
		}
	}
}

// 将消息发给taskQueue，由worker进行处理
func (mh *MsgHandle) SendMsgToTaskQueue(req ziface.IRequest) {
	//将消息平均分配给不同的worker
	// 根据客户端建立的ConnID来分配达到平均分配
	workerID := req.GetConn().GetConnID() % mh.WorkerPoolSize
	fmt.Println("Add ConnID=", req.GetConn().GetConnID(),
		"req MsgID=", req.GetMsgID(), "to WorkerID=", workerID)
	//将消息发给对应worker
	mh.TaskQueue[workerID] <- req

}
