package bitreader

import "io"

type Reader interface {
	ReadBit() bool
	PeekBit() bool
	Trash(uint)
}

type Reader32 interface {
	Reader
	Read32(uint) uint32
	Peek32(uint) uint32
}

func NewReader32(r io.Reader) Reader32 {
	return &simpleReader32{r, make([]byte, 4), 0, 0}
}
