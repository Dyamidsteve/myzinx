package test

import (
	"fmt"
	"io"
	"net"
	"testing"
	"time"
	"zinx-demo/znet"
)

func TestClient(t *testing.T) {
	fmt.Println("[Clien]Client Start...")

	time.Sleep(time.Second * 1)
	conn, err := net.Dial("tcp4", "127.0.0.1:8999")
	if err != nil {
		fmt.Println("Net Dial Error:", err)
		return
	}
	//最后关闭连接
	defer conn.Close()
	for {
		//生成消息
		msg := znet.NewMessage(0, []byte("hello i'm client"))
		//fmt.Println("****请输入内容****")
		//fmt.Scanln(&msg)

		//创建包处理器
		dp := znet.NewDataPack()

		data, err := dp.Pack(msg)
		if err != nil {
			fmt.Println("dp pack error:", err)
			break
		}

		_, err2 := conn.Write(data)
		if err2 != nil && err != io.EOF {
			fmt.Println("Conn write error:", err)
			continue
		}

		//接收消息
		headBuffer := make([]byte, dp.GetHeadLen())
		//第一次读取头(ReadFull可保证完整获取)
		if _, err := io.ReadFull(conn, headBuffer); err != nil {
			fmt.Println("io read conn error:", err)
			break
		}

		//解包获取头信息
		msgHead, err := dp.UnPack(headBuffer)
		if err != nil {
			fmt.Println("dataPack UnPack error:", err)
			break
		}
		msg = msgHead.(*znet.Message) //强制转化实现类
		//存在消息内容
		if msgHead.GetMsgLen() > 0 {
			//第二次读取消息内容
			//开辟数据空间
			msg.Data = make([]byte, msg.GetMsgLen())
			//读取数据
			if _, err := io.ReadFull(conn, msg.Data); err != nil {
				fmt.Println("io read conn error:", err)
				break
			}
			fmt.Println("-->Recv MsgID:", msg.Id, "datalen=", msg.DataLen, "data:", string(msg.Data))

		}

		time.Sleep(time.Second * 1)

	}
}
