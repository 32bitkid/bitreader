package bitreader

import "io"
import "errors"

type Reader interface {
	ReadBit() (bool, error)
	PeekBit() (bool, error)
	Trash(uint) error
	IsByteAligned() bool
	io.Reader
}

type Reader32 interface {
	Reader
	Read32(uint) (uint32, error)
	Peek32(uint) (uint32, error)
}

type BufferedReader32 interface {
	Reader32
	Len() int
	io.Writer
}

type NewReaderFn func(io.Reader) Reader
type NewReader32Fn func(io.Reader) Reader32
type NewBuffered32Fn func() BufferedReader32

var ErrNotAvailable = errors.New("not available")
