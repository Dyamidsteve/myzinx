package ziface

// IRequst接口
// 将客户端请求的链接信息和请求数据包装到一个Request当中
type IRequest interface {
	//得到当前链接
	GetConn() IConnection

	//得到请求数据
	GetData() []byte

	//得到请求消息的ID
	GetMsgID() uint32
}
