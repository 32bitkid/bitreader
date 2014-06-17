package bitreader_test

import "testing"
import "bytes"
import "github.com/32bitkid/bitreader"

type read32 func(uint) uint32

func createReader(b ...byte) bitreader.Reader32 {
	return bitreader.NewReader32(bytes.NewReader(b))
}

func check32(t *testing.T, fn read32, len uint, expected uint32) {
	if actual := fn(len); actual != expected {
		t.Fatalf("Expected %d, got %d", expected, actual)
	}
}

func TestPeekingForZero(t *testing.T) {
	br := createReader(0, 0, 0, 0, 0, 0, 0, 0)
	for i := uint(1); i < 64; i++ {
		check32(t, br.Peek32, i, 0)
	}
}

func TestPeekingForOne(t *testing.T) {
	br := createReader(255, 255, 255, 255, 255, 255, 255, 255)
	for i := uint(1); i < 64; i++ {
		check32(t, br.Peek32, i, uint32(1<<i-1))
	}
}

func TestTrashingBits(t *testing.T) {
	br := createReader(1)
	br.Trash(7)
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
		if br.PeekBit() != false {
			t.Fatal("Expected false")
		}
		br.Trash(1)
		if br.PeekBit() != true {
			t.Fatal("Expected true")
		}
		br.Trash(1)
	}
}

func TestReadingBools(t *testing.T) {
	// 01 010 101
	br := createReader(0125)
	for i := 0; i < 4; i++ {
		if br.ReadBit() != false {
			t.Fatal("Expected false")
		}
		if br.ReadBit() != true {
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
