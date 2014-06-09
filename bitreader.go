package bitreader

import "io"

type Bitreader interface {
	ReadBit() bool
	PeekBit() bool
	Trash(uint)
}

type Bitreader32 interface {
	Bitreader
	Read(uint) uint32
	Peek(uint) uint32
}

type simpleBitreader32 struct {
	source     io.Reader
	readBuffer []byte
	buffer     uint64
	bitsLeft   uint
}

func (b *simpleBitreader32) Peek(len uint) uint32 {
	b.check(len)
	shift := (64 - len)
	var mask uint64 = (1 << (len + 1)) - 1
	return uint32(b.buffer & (mask << shift) >> shift)
}

func (b *simpleBitreader32) Trash(len uint) {
	b.check(len)
	b.buffer <<= len
	b.bitsLeft -= len
}

func (b *simpleBitreader32) Read(len uint) uint32 {
	defer b.Trash(len)
	return b.Peek(len)
}

func (b *simpleBitreader32) PeekBit() bool {
	return b.Peek(1) == 1
}

func (b *simpleBitreader32) ReadBit() bool {
	defer b.Trash(1)
	return b.PeekBit()
}

func (b *simpleBitreader32) check(len uint) {
	if b.bitsLeft < len {
		b.fill()
	}
}

func (b *simpleBitreader32) fill() {
	len, err := b.source.Read(b.readBuffer)
	if err != nil {
		panic(err)
	}
	for i := 0; i < len; i++ {
		b.buffer = b.buffer | uint64(b.readBuffer[i])<<(64-8-b.bitsLeft)
		b.bitsLeft += 8
	}
}

func Create(r io.Reader) Bitreader32 {
	return &simpleBitreader32{r, make([]byte, 4), 0, 0}
}
