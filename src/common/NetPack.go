package common

const (
	PACK_HEADER_SIZE = 7 //packetType & Opcode & reqIdx
	INDEX_TYPE       = 0 //uint8
	INDEX_MSG_ID     = 1 //uint16
	INDEX_REQ_IDX    = 3 //uint32

	// INDEX_TYPE：写通用错误码
	Err_not_found = 255
)

type NetPack struct {
	*ByteBuffer
}

func NewNetPackCap(capacity int) *NetPack {
	if capacity < 16 {
		capacity = 16
	}
	self := &NetPack{NewByteBufferCap(capacity)}
	self.ReadPos = PACK_HEADER_SIZE
	self.buf = self.buf[:PACK_HEADER_SIZE] //len == PACK_HEADER_SIZE
	return self
}
func NewNetPackLen(length int) *NetPack {
	if length < PACK_HEADER_SIZE {
		return nil
	} else {
		self := &NetPack{NewByteBufferLen(length)}
		self.ReadPos = PACK_HEADER_SIZE
		return self
	}
}
func NewNetPack(data []byte) *NetPack {
	if len(data) < PACK_HEADER_SIZE {
		return nil
	} else {
		self := &NetPack{NewByteBuffer(data)}
		self.ReadPos = PACK_HEADER_SIZE
		return self
	}
}

func (self *NetPack) Body() []byte  { return self.buf[PACK_HEADER_SIZE:] }
func (self *NetPack) BodySize() int { return len(self.buf) - PACK_HEADER_SIZE }
func (self *NetPack) Clear() {
	self.buf = self.buf[:PACK_HEADER_SIZE]
	self.ReadPos = PACK_HEADER_SIZE
	self.SetMsgId(0)
}
func (self *NetPack) Reset(data []byte) bool {
	if len(data) < PACK_HEADER_SIZE {
		self.Clear()
		return false
	} else {
		self.buf = data
		self.ReadPos = PACK_HEADER_SIZE
		return true
	}
}
func (self *NetPack) ResetHead(other *NetPack) {
	self.buf = self.buf[:0]
	self.WriteBuf(other.buf[:PACK_HEADER_SIZE])
	self.ReadPos = PACK_HEADER_SIZE
}

//! head
func (self *NetPack) SetType(v uint8) { self.buf[INDEX_TYPE] = v }
func (self *NetPack) GetType() uint8  { return self.buf[INDEX_TYPE] }
func (self *NetPack) SetMsgId(id uint16) {
	self.buf[INDEX_MSG_ID] = byte(id)
	self.buf[INDEX_MSG_ID+1] = byte(id >> 8)
}
func (self *NetPack) GetMsgId() uint16 {
	return uint16(self.buf[INDEX_MSG_ID+1])<<8 | uint16(self.buf[INDEX_MSG_ID])
}
func (self *NetPack) SetReqIdx(idx uint32) {
	self.buf[INDEX_REQ_IDX] = byte(idx)
	self.buf[INDEX_REQ_IDX+1] = byte(idx >> 8)
	self.buf[INDEX_REQ_IDX+2] = byte(idx >> 16)
	self.buf[INDEX_REQ_IDX+3] = byte(idx >> 24)
}
func (self *NetPack) GetReqIdx() (ret uint32) {
	for i := 0; i < 4; i++ {
		ret |= uint32(self.buf[INDEX_REQ_IDX+i]) << uint(i*8)
	}
	return
}
func (self *NetPack) GetReqKey() uint64 {
	return uint64(self.GetMsgId())<<32 | uint64(self.GetReqIdx())
}
func (self *NetPack) SetReqKey(key uint64) {
	self.SetMsgId(uint16(key >> 32))
	self.SetReqIdx(uint32(0xFFFFFFFF & key))
}

//! Set
func (self *NetPack) SetPos(pos int, v uint32) { self.ByteBuffer.SetPos(PACK_HEADER_SIZE+pos, v) }
func (self *NetPack) GetPos(pos int) uint32    { return self.ByteBuffer.GetPos(PACK_HEADER_SIZE + pos) }
