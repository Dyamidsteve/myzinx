package znet

import "zinx-demo/ziface"

// 实现router时，先嵌入该BaseRouter积累，根据该基类方法重写子类方法即可
type BaseRouter struct{}

// 处理业务前方法
func (br *BaseRouter) PreHandle(req ziface.IRequest) {}

// 处理业务主方法
func (br *BaseRouter) MainHandle(req ziface.IRequest) {}

// 处理业务后方法
func (br *BaseRouter) PostHandle(req ziface.IRequest) {}
