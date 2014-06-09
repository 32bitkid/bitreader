package bitreader_test

import "testing"
import "bytes"
import "github.com/32bitkid/bitreader"

type read32 func(uint) uint32

func createReader(b ...byte) bitreader.Bitreader32 {
	return bitreader.Create(bytes.NewReader(b))
}

func check32(t *testing.T, fn read32, len uint, expected uint32) {
	if actual := fn(len); actual != expected {
		t.Fatalf("Expected %d, got %d", expected, actual)
	}
}

func TestPeekingForZero(t *testing.T) {
	br := createReader(0, 0, 0, 0, 0, 0, 0, 0)
	for i := uint(1); i < 64; i++ {
		check32(t, br.Peek, i, 0)
	}
}

func TestPeekingForOne(t *testing.T) {
	br := createReader(255, 255, 255, 255, 255, 255, 255, 255)
	for i := uint(1); i < 64; i++ {
		check32(t, br.Peek, i, uint32(1<<i-1))
	}
}

func TestTrashingBits(t *testing.T) {
	br := createReader(1)
	br.Trash(7)
	check32(t, br.Peek, 1, 1)
}

func TestReadingBits(t *testing.T) {
	// 0000 0001 0000 0001 0000 0010 0000 0100
	br := createReader(1, 1, 2, 4)
	check32(t, br.Read, 7, 0)
	check32(t, br.Read, 1, 1)
	check32(t, br.Read, 7, 0)
	check32(t, br.Read, 1, 1)
	check32(t, br.Read, 6, 0)
	check32(t, br.Read, 1, 1)
	check32(t, br.Read, 6, 0)
	check32(t, br.Read, 1, 1)
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
