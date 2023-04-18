package znet

// Message包含Id、len、data三个属性，并实现get和set方法
type Message struct {
	Id uint32

	DataLen uint32

	Data []byte
}

// 初始化Message
func NewMessage(msgId uint32, data []byte) *Message {
	msg := &Message{
		Id:      msgId,
		DataLen: uint32(len(data)),
		Data:    data,
	}
	return msg
}

func (m *Message) GetMsgID() uint32 {
	return m.Id
}

func (m *Message) SetMsgID(ID uint32) {
	m.Id = ID
}

func (m *Message) GetMsgLen() uint32 {
	return m.DataLen
}

func (m *Message) SetMsgLen(Len uint32) {
	m.DataLen = Len
}

func (m *Message) GetData() []byte {
	return m.Data
}

func (m *Message) SetData(data []byte) {
	m.Data = data
}
