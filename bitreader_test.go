package bitreader_test

import (
	"bytes"
	"io"
	"testing"

	"github.com/ysh86/bitreader"
)

type read32 func(uint) (uint32, error)
type read64 func(uint) (uint64, error)

func createReader(b ...byte) bitreader.BitReader {
	return bitreader.NewReader(bytes.NewReader(b))
}

func check32(t *testing.T, fn read32, len uint, expected uint32) {
	actual, err := fn(len)
	if err != nil {
		t.Fatal(err)
		return
	}
	if actual != expected {
		t.Fatalf("Expected %d, got %d", expected, actual)
		return
	}
}

func check64(t *testing.T, fn read64, len uint, expected uint64) {
	actual, err := fn(len)
	if err != nil {
		t.Fatal(err)
		return
	}
	if actual != expected {
		t.Fatalf("Expected %d, got %d", expected, actual)
		return
	}
}

func TestPeekingForZero(t *testing.T) {
	br := createReader(0, 0, 0, 0, 0, 0, 0, 0)
	for i := uint(1); i < 32; i++ {
		check32(t, br.Peek32, i, 0)
	}
}

func TestPeekingForOne(t *testing.T) {
	br := createReader(255, 255, 255, 255, 255, 255, 255, 255)
	for i := uint(1); i < 32; i++ {
		check32(t, br.Peek32, i, ^uint32(0)>>(32-i))
	}
}

func TestPeekingFor64Zero(t *testing.T) {
	br := createReader(0, 0, 0, 0, 0, 0, 0, 0)
	for i := uint(1); i < 64; i++ {
		check64(t, br.Peek64, i, 0)
	}
}

func TestPeekingFor64One(t *testing.T) {
	br := createReader(255, 255, 255, 255, 255, 255, 255, 255)
	for i := uint(1); i < 64; i++ {
		check64(t, br.Peek64, i, ^uint64(0)>>(64-i))
	}
}

func TestTrashingBits(t *testing.T) {
	br := createReader(1)
	br.Skip(7)
	check32(t, br.Peek32, 1, 1)
}

func TestReadingBits(t *testing.T) {
	// 0000 0001 0000 0001 0000 0010 0000 0100
	br := createReader(1, 1, 2, 4)
	check32(t, br.Read32, 7, 0)
	check32(t, br.Read32, 1, 1)
	check32(t, br.Read32, 7, 0)
	check32(t, br.Read32, 1, 1)
	check32(t, br.Read32, 6, 0)
	check32(t, br.Read32, 1, 1)
	check32(t, br.Read32, 6, 0)
	check32(t, br.Read32, 1, 1)
}

func TestPeekingBools(t *testing.T) {
	// 01 010 101
	br := createReader(0125)
	for i := 0; i < 4; i++ {
		val, err := br.Peek1()
		if val != false || err != nil {
			t.Fatal("Expected false")
		}
		err = br.Skip(1)
		if err != nil {
			t.Fatal("Unexpected error")
		}
		val, err = br.Peek1()
		if val != true || err != nil {
			t.Fatal("Expected true")
		}
		err = br.Skip(1)
		if err != nil {
			t.Fatal("Unexpected error")
		}
	}
}

func TestReadingBools(t *testing.T) {
	// 01 010 101
	br := createReader(0125)
	for i := 0; i < 4; i++ {
		val, err := br.Read1()
		if val != false || err != nil {
			t.Fatal("Expected false")
		}
		val, err = br.Read1()
		if val != true || err != nil {
			t.Fatal("Expected true")
		}
	}
}

func TestReadingLongStrings(t *testing.T) {
	data := []byte{0x48, 0xbb, 0xad, 0x83, 0xa6, 0xa4, 0xe1, 0x43, 0x25, 0xb, 0x19, 0xe2, 0xf5, 0x5d, 0x27, 0x2, 0x69, 0xf9, 0xd3, 0x50}
	br := createReader(data...)
	for _, val := range data {
		check32(t, br.Read32, 8, uint32(val))
	}
}

func TestPeekEOF(t *testing.T) {
	br := createReader(0x01)
	br.Skip(8)
	_, err := br.Peek32(8)
	if err != io.EOF {
		t.Fatalf("Expected %s but got %s\n", io.EOF, err)
	}
}

func TestPeek1EOF(t *testing.T) {
	br := createReader(0x01)
	br.Skip(8)
	_, err := br.Peek1()
	if err != io.EOF {
		t.Fatalf("Expected %s but got %s\n", io.EOF, err)
	}
}

func TestReadEOF(t *testing.T) {
	br := createReader(0x01)
	br.Skip(8)
	_, err := br.Read32(8)
	if err != io.ErrUnexpectedEOF {
		t.Fatalf("Expected %s error but got %s\n", io.ErrUnexpectedEOF, err)
	}
}

func TestRead1EOF(t *testing.T) {
	br := createReader(0x01)
	br.Skip(8)
	_, err := br.Read1()
	if err != io.ErrUnexpectedEOF {
		t.Fatalf("Expected %s error but got %s\n", io.ErrUnexpectedEOF, err)
	}
}

func TestTrashEOF(t *testing.T) {
	br := createReader(0x01)
	br.Skip(8)
	err := br.Skip(8)
	if err != io.ErrUnexpectedEOF {
		t.Fatalf("Expected %s error but got %s\n", io.ErrUnexpectedEOF, err)
	}
}

func TestRead16UnexpectedEOF(t *testing.T) {
	br := createReader(0x01)
	_, err := br.Read16(16)
	if err != io.ErrUnexpectedEOF {
		t.Fatalf("Expected %s error but got %s\n", io.ErrUnexpectedEOF, err)
	}
}

func TestRead32UnexpectedEOF(t *testing.T) {
	br := createReader(0x01)
	_, err := br.Read32(32)
	if err != io.ErrUnexpectedEOF {
		t.Fatalf("Expected %s error but got %s\n", io.ErrUnexpectedEOF, err)
	}
}

func TestRead64UnexpectedEOF(t *testing.T) {
	br := createReader(0x01)
	_, err := br.Read64(64)
	if err != io.ErrUnexpectedEOF {
		t.Fatalf("Expected %s error but got %s\n", io.ErrUnexpectedEOF, err)
	}
}

func TestTrashUnexpectedEOF(t *testing.T) {
	br := createReader(0x01)
	err := br.Skip(32)
	if err != io.ErrUnexpectedEOF {
		t.Fatalf("Expected %s error but got %s\n", io.ErrUnexpectedEOF, err)
	}
}

func TestBasicReading(t *testing.T) {
	data := []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	br := createReader(data...)

	buffer := make([]byte, 5)
	_, err := io.ReadAtLeast(br, buffer, 5)

	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(buffer, data[:5]) {
		t.Fatalf("Expected %+v to equal %+v", buffer, data[:5])
	}
}

func TestReadingAfterBitOperation(t *testing.T) {
	data := []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	br := createReader(data...)

	br.Skip(8)

	buffer := make([]byte, 5)
	_, err := io.ReadAtLeast(br, buffer, 5)

	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(buffer, data[1:6]) {
		t.Fatalf("Expected %+v to equal %+v", buffer, data[1:6])
	}
}

func TestRealignmentReading(t *testing.T) {
	data := []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	br := createReader(data...)

	br.Skip(20)

	buffer := make([]byte, 5)
	_, err := io.ReadAtLeast(br, buffer, 5)

	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(buffer, data[3:8]) {
		t.Fatalf("Expected %+v to equal %+v", buffer, data[3:8])
	}
}

func TestReadingAfterPeeking(t *testing.T) {
	data := []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	br := createReader(data...)
	buffer := make([]byte, 3)

	// 1
	_, err := br.Peek1()
	if err != nil {
		t.Fatal(err)
	}
	_, err = br.Peek8(7)
	if err != nil {
		t.Fatal(err)
	}
	_, err = br.Peek8(3)
	if err != nil {
		t.Fatal(err)
	}
	_, err = br.Peek16(13)
	if err != nil {
		t.Fatal(err)
	}
	_, err = br.Peek8(8)
	if err != nil {
		t.Fatal(err)
	}

	n, err := io.ReadAtLeast(br, buffer, 3)

	if err != nil {
		t.Fatal(err)
	}

	if n != 3 {
		t.Fatalf("Expected %+v to equal %+v", n, 3)
	}
	if !bytes.Equal(buffer, data[0:3]) {
		t.Fatalf("Expected %+v to equal %+v", buffer, data[0:3])
	}

	// 2
	_, err = br.Peek64(49)
	if err != nil {
		t.Fatal(err)
	}

	n, err = io.ReadAtLeast(br, buffer, 3)

	if err != nil {
		t.Fatal(err)
	}

	if n != 3 {
		t.Fatalf("Expected %+v to equal %+v", n, 3)
	}
	if !bytes.Equal(buffer, data[3:6]) {
		t.Fatalf("Expected %+v to equal %+v", buffer, data[3:6])
	}
}

func TestByteAlignment(t *testing.T) {
	br := createReader(0, 255, 0, 0, 0)
	if br.IsAligned() != true {
		t.Fail()
	}
	br.Skip(1)
	if br.IsAligned() != false {
		t.Fail()
	}

	for !br.IsAligned() {
		br.Skip(1)
	}

	if val, err := br.Peek1(); val != true || err != nil {
		t.Fail()
	}

}
