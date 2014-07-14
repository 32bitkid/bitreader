package bitreader

import "io"
import "errors"

type Reader interface {
	ReadBit() (bool, error)
	PeekBit() (bool, error)
	Trash(uint) error
	io.Reader
}

type Reader32 interface {
	Reader
	Read32(uint) (uint32, error)
	Peek32(uint) (uint32, error)
}

func NewReader32(r io.Reader) Reader32 {
	return &simpleReader32{r, make([]byte, 4), 0, 0}
}

var ErrNotAvailable = errors.New("not available")
