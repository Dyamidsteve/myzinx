package znet

import "zinx-demo/ziface"

type Request struct {
	//已经和客户端建立好的连接
	conn ziface.IConnection

	//客户端请求数据
	data ziface.IMessage
}

// 得到连接
func (r *Request) GetConn() ziface.IConnection {
	return r.conn
}

// 得到请求数据
func (r *Request) GetData() []byte {
	return r.data.GetData()
}

// 得到请求消息的ID
func (r *Request) GetMsgID() uint32 {
	return r.data.GetMsgID()
}
