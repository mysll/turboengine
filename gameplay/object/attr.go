package object

import (
	"io"
)

// 属性接口
type Attr interface {
	Flag() uint32
	SetFlag(f uint32)
	ClearFlag(f uint32)
	Name() string
	Index() uint32
	Write(stream io.Writer) (uint32, error)
	Read(reader io.Reader) (uint32, error)
}

type AttrHolder struct {
	name  string
	index uint32
	flag  uint32
}

func (h *AttrHolder) Flag() uint32 {
	return h.flag
}

func (h *AttrHolder) SetFlag(f uint32) {
	h.flag |= f
}

func (h *AttrHolder) ClearFlag(f uint32) {
	h.flag &= ^f
}

func (h *AttrHolder) Name() string {
	return h.name
}

func (h *AttrHolder) Index() uint32 {
	return h.index
}
