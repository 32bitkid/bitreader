package bitreader

import "io"

type BitReader interface {
	ReadBit() bool
	PeekBit() bool
	Trash(uint)
}

type BitReader32 interface {
	BitReader
	Read32(uint) uint32
	Peek32(uint) uint32
}

func NewBitReader32(r io.Reader) BitReader32 {
	return &simpleBitReader32{r, make([]byte, 4), 0, 0}
}
