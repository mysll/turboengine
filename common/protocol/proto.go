package protocol

type ProtoMsg struct {
	Id   int32
	Src  Mailbox
	Dest Mailbox
	Data any
}

type ProtoEncoder interface {
	Encode(*ProtoMsg) (*Message, error)
}

type ProtoDecoder interface {
	Decode([]byte) (*ProtoMsg, error)
}
