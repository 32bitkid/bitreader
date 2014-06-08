package bitreader

import "io"

type Bitreader interface {
	Read(uint) uint32
	ReadBit() bool

	Peek(uint) uint32
	PeekBit() bool

	Trash(uint)
}

type bitreader struct {
	source     io.Reader
	readBuffer []byte
	buffer     uint64
	bitsLeft   uint
}

func (b *bitreader) Peek(len uint) uint32 {
	b.check(len)
	shift := (64 - len)
	var mask uint64 = (1 << (len + 1)) - 1
	return uint32(b.buffer & (mask << shift) >> shift)
}

func (b *bitreader) Trash(len uint) {
	b.check(len)
	b.buffer <<= len
}

func (b *bitreader) Read(len uint) uint32 {
	defer b.Trash(len)
	return b.Peek(len)
}

func (b *bitreader) PeekBit() bool {
	return b.Peek(1) == 1
}

func (b *bitreader) ReadBit() bool {
	defer b.Trash(1)
	return b.PeekBit()
}

func (b *bitreader) check(len uint) {
	if b.bitsLeft < len {
		b.fill()
	}
}

func (b *bitreader) fill() {
	len, err := b.source.Read(b.readBuffer)
	if err != nil {
		panic(err)
	}
	for i := 0; i < len; i++ {
		b.buffer = b.buffer | uint64(b.readBuffer[i])<<(64-8-b.bitsLeft)
		b.bitsLeft += 8
	}
}

func Create(r io.Reader) Bitreader {
	return &bitreader{r, make([]byte, 4), 0, 0}
}
