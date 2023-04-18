package znet

import (
	"bytes"
	"encoding/binary"
	"errors"
	"zinx-demo/utils"
	"zinx-demo/ziface"
)

type DataPack struct {
}

func NewDataPack() *DataPack {
	return &DataPack{}
}

func (dp *DataPack) GetHeadLen() uint32 {
	return 8 //Len4个字节+ID4个字节
}

func (dp *DataPack) Pack(msg ziface.IMessage) ([]byte, error) {
	//创建一个存放bytes字节的缓存
	Buffer := bytes.NewBuffer([]byte{})

	//将dataLen写入Buffer中
	if err := binary.Write(Buffer, binary.LittleEndian, msg.GetMsgLen()); err != nil {
		return nil, err
	}

	//将msgId写入Buffer中
	if err := binary.Write(Buffer, binary.LittleEndian, msg.GetMsgID()); err != nil {
		return nil, err
	}

	//将data写入Buffer中
	if err := binary.Write(Buffer, binary.LittleEndian, msg.GetData()); err != nil {
		return nil, err
	}

	return Buffer.Bytes(), nil
}

func (dp *DataPack) UnPack(data []byte) (ziface.IMessage, error) {
	//创建一个从输入二进制数据的ioReader
	Buffer := bytes.NewReader(data)

	//只解压head信息，得到dataLen和msgID
	msg := &Message{}

	//读取dataLen
	if err := binary.Read(Buffer, binary.LittleEndian, &msg.DataLen); err != nil {
		return nil, err
	}

	//读取ID
	if err := binary.Read(Buffer, binary.LittleEndian, &msg.Id); err != nil {
		return nil, err
	}

	// //读取data
	// if err := binary.Read(Buffer, binary.LittleEndian, &msg.Data); err != nil {
	// 	return nil, err
	// }

	//判断dataLen是否超出允许的最大包长度
	if msg.DataLen > utils.GlobalConf.MaxPackageSize {
		return nil, errors.New("too large msg data recv")
	}
	return msg, nil
}
