package znet

import (
	"errors"
	"fmt"
	"sync"
	"zinx-demo/ziface"
)

type ConnManager struct {
	// 管理的链接模块
	connections map[uint32]ziface.IConnection
	// 链接模块读写锁
	connLock sync.RWMutex
}

func NewConnManager() *ConnManager {
	return &ConnManager{
		connections: make(map[uint32]ziface.IConnection),
	}
}

// 添加链接
func (cm *ConnManager) Add(conn ziface.IConnection) {
	cm.connLock.Lock()
	defer cm.connLock.Unlock()
	cm.connections[conn.GetConnID()] = conn

	fmt.Println("Add connection ID=", conn.GetConnID())
}

// 删除链接
func (cm *ConnManager) Remove(conn ziface.IConnection) {
	cm.connLock.Lock()
	defer cm.connLock.Unlock()

	delete(cm.connections, conn.GetConnID())
	fmt.Println("delete connection ID=", conn.GetConnID())
}

// 根据connID获取链接
func (cm *ConnManager) Get(connID uint32) (ziface.IConnection, error) {
	cm.connLock.Lock()
	defer cm.connLock.Unlock()
	conn, ok := cm.connections[connID]
	if ok {
		return conn, nil
	} else {
		return nil, errors.New("connection not FOUND")
	}
}

// 获取当前连接总数
func (cm *ConnManager) GetConnLen() int {
	return len(cm.connections)
}

// 清除并关闭所有链接
func (cm *ConnManager) Clear() {
	cm.connLock.Lock()
	defer cm.connLock.Unlock()

	// 删除conn并停止conn的工作
	for connID, conn := range cm.connections {
		conn.Stop() //停止链接

		delete(cm.connections, connID) //删除链接
	}
	fmt.Println("Clear All connections current Length=", len(cm.connections))
}
