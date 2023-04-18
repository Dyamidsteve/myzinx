package znet

import (
	"errors"
	"fmt"
	"io"
	"net"
	"sync"
	"zinx-demo/utils"
	"zinx-demo/ziface"
)

type Connection struct {
	//当前Conn隶属于哪个Server(可以扩展，分布式，多Server业务等)
	TcpServer ziface.IServer

	//当前链接到达socket TCP
	Conn *net.TCPConn

	//链接ID
	ConnID uint32

	//链接状态
	isClose bool

	//用于通知链接退出的channel
	ExitChan chan bool

	//无缓冲管道，用于读写Goroutine之间消息通信
	msgChan chan []byte

	//该连接处理的消息处理模块
	MsgHandle ziface.IMsgHandle

	//链接属性集合
	property map[string]interface{}

	//链接属性map锁
	propertyLock sync.RWMutex
}

// 初始化链接模块
func NewConnection(server ziface.IServer, connection *net.TCPConn, ID uint32, msgHandler ziface.IMsgHandle) *Connection {
	conn := &Connection{
		TcpServer: server,
		Conn:      connection,
		ConnID:    ID,
		isClose:   false,
		msgChan:   make(chan []byte),
		MsgHandle: msgHandler,
		ExitChan:  make(chan bool, 1),
		property:  make(map[string]interface{}),
	}

	//将Conn加入ConnManager中
	conn.TcpServer.GetConnManager().Add(conn)
	return conn
}

func (c *Connection) StartReader() {
	fmt.Println("Reader Start...")

	defer fmt.Println("Reader exit,connID=", "c.ConnID,Reader is exit,remote Addr=", c.RemoteAddr())
	defer c.Stop()

	for {
		//创建包处理器
		dp := NewDataPack()
		headBuffer := make([]byte, dp.GetHeadLen())
		//第一次读取头(ReadFull可保证完整获取)
		if _, err := io.ReadFull(c.Conn, headBuffer); err != nil {
			fmt.Println("io read conn error:", err)
			break
		}

		//解包获取头信息
		msgHead, err := dp.UnPack(headBuffer)
		if err != nil {
			fmt.Println("dataPack UnPack error:", err)
			break
		}
		msg := msgHead.(*Message) //强制转化实现类
		//存在消息内容
		if msgHead.GetMsgLen() > 0 {
			//第二次读取消息内容
			//开辟数据空间
			msg.Data = make([]byte, msg.GetMsgLen())
			//读取数据
			if _, err := io.ReadFull(c.Conn, msg.Data); err != nil {
				fmt.Println("io read conn error:", err)
				break
			}
			fmt.Println("-->Recv MsgID:", msg.Id, "datalen=", msg.DataLen, "data:", string(msg.Data))

		}

		// 封装请求连接和消息
		req := Request{
			conn: c,
			data: msg,
		}

		if utils.GlobalConf.WorkerPoolSize > 0 {
			//已经开启了工作池，将消息发送给工作池即可
			c.MsgHandle.SendMsgToTaskQueue(&req)
		} else {
			//根据绑定的msgID找到对应的方法
			go c.MsgHandle.DoMsgHandler(&req)
		}
	}

}

// 写消息模块，专用于将管道消息发送给客户端
func (c *Connection) StartWriter() {
	fmt.Println("Writer Goroutine Start...")

	defer fmt.Println(c.RemoteAddr().String(), "[conn Writer exit]")

	for {
		select {
		case data := <-c.msgChan:
			if _, err := c.Conn.Write(data); err != nil {
				fmt.Println("Conn write error:", err)
				return
			}
		case <-c.ExitChan:
			//可读代表Reader退出，Writer也要退出
			return
		}
	}
}

// 启动链接，让当前链接开始工作
func (c *Connection) Start() {
	fmt.Println("Connection Start... ConnID=", c.ConnID)
	//启动当前连接读数据的业务
	go c.StartReader()

	go c.StartWriter()

	//调用开发者传来的Server端的Hook-OnConnStart方法
	c.TcpServer.CallOnConnStart(c)

}

// 停止连接，结束当前链接工作
func (c *Connection) Stop() {
	fmt.Println("Connection Stop... ConnID=", c.ConnID)

	if c.isClose {
		return
	}

	c.isClose = true //链接状态关闭

	c.Conn.Close() //结束链接

	//将Server端的manager中的该连接删除
	c.TcpServer.GetConnManager().Remove(c)

	//调用开发者传来的Server端的Hook-OnConnStop方法
	c.TcpServer.CallOnConnStop(c)

	//关闭通知管道，节省内存
	close(c.ExitChan)
	close(c.msgChan)
}

// 获取当前链接绑定socket conn
func (c *Connection) GetTCPConnection() *net.TCPConn {

	return c.Conn
}

// 获取当前链接模块的链接ID
func (c *Connection) GetConnID() uint32 {

	return c.ConnID
}

// 获取远程客户端的TCP状态 IP.port
func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

// 发送数据给客户端
func (c *Connection) SendMsg(msgId uint32, date []byte) error {
	if c.isClose {
		return errors.New("Connection closed when send msg")
	}
	//创建包处理器
	dp := NewDataPack()
	//封装消息
	msg := NewMessage(msgId, date)
	//打包消息
	buffer, err := dp.Pack(msg)
	if err != nil {
		return err
	}

	//发送消息给管道
	c.msgChan <- buffer
	fmt.Println("msg Chan sended")
	return nil
}

// 设置链接属性
func (c *Connection) SetProperty(key string, value interface{}) {
	c.propertyLock.Lock()
	c.property[key] = value
	c.propertyLock.Unlock()
}

// 获取链接属性
func (c *Connection) GetProperty(key string) (interface{}, error) {
	c.propertyLock.Lock()
	//若放在后面则可能直接return了轮不到
	//因此对于可能提前return的语句，应该用defer保证解锁
	defer c.propertyLock.Unlock()
	if val, ok := c.property[key]; ok {
		return val, nil
	} else {
		return nil, fmt.Errorf("Property Not Found")
	}

}

// 删除链接属性
func (c *Connection) RemoveProperty(key string) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()
	// 删除属性
	delete(c.property, key)
}
