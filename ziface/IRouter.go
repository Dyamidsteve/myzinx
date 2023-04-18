package ziface

// 路由抽象接口，路由内数据都是IRequest
type IRouter interface {
	//处理业务前方法
	PreHandle(req IRequest)

	//处理业务主方法
	MainHandle(req IRequest)

	//处理业务后方法
	PostHandle(req IRequest)
}

//接口也能做到指针的性能，需要解引用
