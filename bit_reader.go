// Package bitreader provides basic interfaces to read and traverse
// an io.Reader as a stream of bits, rather than a stream of bytes.
package bitreader

import "io"

// Reader is the interface that wraps the basic ReadBit method
//
// ReadBit will return true or false depending on whether or not
// the next bit in the bit stream is set, then advance one bit
// forward in the bitstream.
//
// ReadBit() is the equivalent to PeekBit() followed by Trash(1)
type Reader interface {
	ReadBit() (bool, error)
}

// Peeker is the interface that wraps the basic PeekBit method
//
// PeekBit will return true or false depending on whether or not
// the next bit in the bit stream is set; it does not advance
// the stream any bits.
type Peeker interface {
	PeekBit() (bool, error)
}

// Trasher is the interface that wraps the basic Trash method
//
// Trash will advance the bit stream by n bits.
type Trasher interface {
	Trash(n uint) error
}

// Aligner is the interface that allows for byte realignment.
//
// IsByteAligned() returns true if the bit stream is currently
// aligned to a byte boundary.
//
// ByteAlign() will trash the necessary bit to realign the bit
// stream to a byte boundary. It returns the number of bits trashed
// (0 <= n < 8) during realignment.
type Aligner interface {
	IsByteAligned() bool
	ByteAlign() (n uint, err error)
}

// Reader32 is the interface that wraps the basic Read32 method.
//
// Read32 allows for reading multiple bits, where (1 <= n <= 32) as a uint32
// from the bit stream. Then advancing the bit stream by n bits.
//
// Read32(n) is equivalent to Peek32(n) followed by Trash(n)
type Reader32 interface {
	Read32(n uint) (uint32, error)
}

// Peeker32 is the interface that wraps the basic Peek32 method.
//
// Peek32 allows for reading multiple bits, where (1 <= n <= 32) as a uint32
// from the bit stream; it does not advance
// the bit stream any bits.
type Peeker32 interface {
	Peek32(n uint) (uint32, error)
}

// BitReader is the interface that groups together the methods required to
// perform basic bit operations on a bit stream
type BitReader interface {
	io.Reader
	Reader
	Peeker
	Trasher
	Aligner
	Reader32
	Peeker32
}

// NewBitReader returns the default implementation of a BitReader
func NewBitReader(r io.Reader) BitReader {
	return NewSimpleBitReader(r)
}
