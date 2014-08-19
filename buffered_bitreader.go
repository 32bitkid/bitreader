package bitreader

import "bytes"

type bufferedReader32 struct {
	reader Reader32
	buffer *bytes.Buffer
	sync   chan interface{}
}

var goToken = struct{}{}

func NewBufferedReader32(newReader32Fn NewReader32Fn) BufferedReader32 {
	buffer := &bytes.Buffer{}
	reader := newReader32Fn(buffer)
	sync := make(chan interface{}, 1)
	return &bufferedReader32{reader, buffer, sync}
}

func (wbr *bufferedReader32) Write(data []byte) (n int, err error) {
	n, err = wbr.buffer.Write(data)
	select {
	case <-wbr.sync:
		wbr.sync <- goToken
	default:
		wbr.sync <- goToken
	}
	return
}

func (wbr *bufferedReader32) Len() int {
	return wbr.buffer.Len()
}

func (wbr *bufferedReader32) Trash(len uint) (err error) {
	err = wbr.reader.Trash(len)
	if err == ErrNotAvailable {
		<-wbr.sync
		err = wbr.reader.Trash(len)
	}
	return
}

func (wbr *bufferedReader32) Peek32(len uint) (val uint32, err error) {
	val, err = wbr.reader.Peek32(len)
	if err == ErrNotAvailable {
		<-wbr.sync
		val, err = wbr.reader.Peek32(len)
	}
	return
}

func (wbr *bufferedReader32) Read32(len uint) (val uint32, err error) {
	val, err = wbr.reader.Read32(len)
	if err == ErrNotAvailable {
		<-wbr.sync
		val, err = wbr.reader.Read32(len)
	}
	return
}

func (wbr *bufferedReader32) PeekBit() (val bool, err error) {
	val, err = wbr.reader.PeekBit()
	if err == ErrNotAvailable {
		<-wbr.sync
		val, err = wbr.reader.PeekBit()
	}
	return
}

func (wbr *bufferedReader32) ReadBit() (val bool, err error) {
	val, err = wbr.reader.ReadBit()
	if err == ErrNotAvailable {
		<-wbr.sync
		val, err = wbr.reader.ReadBit()
	}
	return
}

func (wbr *bufferedReader32) Read(p []byte) (n int, err error) {
	n, err = wbr.reader.Read(p)
	return
}
