package ziface

// 封装了消息ID、消息长度、数据的GET和SET方法的接口
type IMessage interface {
	GetMsgID() uint32

	SetMsgID(uint32)

	GetMsgLen() uint32

	SetMsgLen(uint32)

	GetData() []byte

	SetData([]byte)
}
