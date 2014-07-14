package bitreader

import "io"
import "bytes"

type BufferedBitreader interface {
	Reader32
	Len() int
	io.Writer
}

type bufferedBitreader struct {
	reader Reader32
	buffer *bytes.Buffer
	sync   chan interface{}
}

var goToken = struct{}{}

func NewBufferedBitreader() BufferedBitreader {
	buffer := &bytes.Buffer{}
	reader := NewReader32(buffer)
	sync := make(chan interface{}, 1)
	return &bufferedBitreader{reader, buffer, sync}
}

func (wbr *bufferedBitreader) Write(data []byte) (n int, err error) {
	n, err = wbr.buffer.Write(data)
	select {
	case <-wbr.sync:
		wbr.sync <- goToken
	default:
		wbr.sync <- goToken
	}
	return
}

func (wbr *bufferedBitreader) Len() int {
	return wbr.buffer.Len()
}

func (wbr *bufferedBitreader) Trash(len uint) (err error) {
	err = wbr.reader.Trash(len)
	if err == ErrNotAvailable {
		<-wbr.sync
		err = wbr.reader.Trash(len)
	}
	return
}

func (wbr *bufferedBitreader) Peek32(len uint) (val uint32, err error) {
	val, err = wbr.reader.Peek32(len)
	if err == ErrNotAvailable {
		<-wbr.sync
		val, err = wbr.reader.Peek32(len)
	}
	return
}

func (wbr *bufferedBitreader) Read32(len uint) (val uint32, err error) {
	val, err = wbr.reader.Read32(len)
	if err == ErrNotAvailable {
		<-wbr.sync
		val, err = wbr.reader.Read32(len)
	}
	return
}

func (wbr *bufferedBitreader) PeekBit() (val bool, err error) {
	val, err = wbr.reader.PeekBit()
	if err == ErrNotAvailable {
		<-wbr.sync
		val, err = wbr.reader.PeekBit()
	}
	return
}

func (wbr *bufferedBitreader) ReadBit() (val bool, err error) {
	val, err = wbr.reader.ReadBit()
	if err == ErrNotAvailable {
		<-wbr.sync
		val, err = wbr.reader.ReadBit()
	}
	return
}

func (wbr *bufferedBitreader) Read(p []byte) (n int, err error) {
	n, err = wbr.reader.Read(p)
	return
}
