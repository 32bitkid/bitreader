package bitreader

import "io"

type simpleBitReader32 struct {
	source     io.Reader
	readBuffer []byte
	buffer     uint64
	bitsLeft   uint
}

func (b *simpleBitReader32) Peek32(len uint) uint32 {
	b.check(len)
	shift := (64 - len)
	var mask uint64 = (1 << (len + 1)) - 1
	return uint32(b.buffer & (mask << shift) >> shift)
}

func (b *simpleBitReader32) Trash(len uint) {
	b.check(len)
	b.buffer <<= len
	b.bitsLeft -= len
}

func (b *simpleBitReader32) Read32(len uint) uint32 {
	defer b.Trash(len)
	return b.Peek32(len)
}

func (b *simpleBitReader32) PeekBit() bool {
	return b.Peek32(1) == 1
}

func (b *simpleBitReader32) ReadBit() bool {
	defer b.Trash(1)
	return b.PeekBit()
}

func (b *simpleBitReader32) check(len uint) {
	if b.bitsLeft < len {
		b.fill()
	}
}

func (b *simpleBitReader32) fill() {
	len, err := b.source.Read(b.readBuffer)

	// TODO propagate
	if err != nil {
		panic(err)
	}

	for i := 0; i < len; i++ {
		b.buffer = b.buffer | uint64(b.readBuffer[i])<<(64-8-b.bitsLeft)
		b.bitsLeft += 8
	}
}
