package znet

import (
	"fmt"
	"io"
	"net"
	"testing"
)

func TestDataPack(t *testing.T) {
	//服务器监听
	listenner, err := net.Listen("tcp", "127.0.0.1:8999")
	if err != nil {
		fmt.Println("server listen error:", err)
		return
	}

	go func() {
		//接收客户端请求
		conn, err := listenner.Accept()
		if err != nil {
			fmt.Println("listenner accept error:", err)
			return
		}
		go func(conn net.Conn) {
			//处理客户端请求
			dp := NewDataPack()
			for {
				//第一次从conn中读数据，先将包的head(Len,ID)读出来
				//根据头长度创建缓冲
				header := make([]byte, dp.GetHeadLen())
				_, err := io.ReadFull(conn, header)
				if err != nil {
					fmt.Println("read head error")
					break //没有数据则跳出循环
				}

				msgHead, err := dp.UnPack(header)
				if err != nil {
					fmt.Println("server unPack error:", err)
					return
				}

				if msgHead.GetMsgLen() > 0 {
					//msg有数据，需要第二次读取
					msg := msgHead.(*Message) //将接口强制转为struct
					msg.Data = make([]byte, msg.GetMsgLen())

					//根据dataLen再次从io读取
					_, err := io.ReadFull(conn, msg.Data)
					if err != nil {
						fmt.Println("conn readFull error:", err)
						return
					}

					fmt.Println("-->Recv MsgID:", msg.Id, "datalen=", msg.DataLen, "data:", string(msg.Data))

				}
			}
		}(conn)
	}()

	//模拟客户端发送数据
	conn, err1 := net.Dial("tcp", "127.0.0.1:8999")
	if err1 != nil {
		fmt.Println("net dial error:", err1)
		return
	}

	//模拟粘包过程发两次包
	dp := NewDataPack()

	//第一次封包
	msg1 := &Message{
		Id:      1,
		DataLen: 5,
		Data:    []byte("HELLO"),
	}
	data1, err2 := dp.Pack(msg1)
	if err2 != nil {
		fmt.Println("datapack error:", err2)
		return
	}
	//第二次封包
	msg2 := &Message{
		Id:      2,
		DataLen: 5,
		Data:    []byte("hello"),
	}
	data2, err3 := dp.Pack(msg2)
	if err3 != nil {
		fmt.Println("datapack error:", err3)
		return
	}
	//将包粘在一起
	sendData := append(data1, data2...)

	conn.Write(sendData)

	//客户端阻塞
	select {}
}
