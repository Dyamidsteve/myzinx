package ziface

type IConnManager interface {
	//添加链接
	Add(conn IConnection)

	//删除链接
	Remove(conn IConnection)

	//根据connID获取链接
	Get(connID uint32) (IConnection, error)

	//获取当前连接总数
	GetConnLen() int

	//清除并关闭所有链接
	Clear()
}
