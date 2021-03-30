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
