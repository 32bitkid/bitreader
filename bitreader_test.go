package bitreader_test

import "testing"
import "bytes"
import "github.com/32bitkid/bitreader"

func createReader(b ...byte) bitreader.Bitreader {
	return bitreader.Create(bytes.NewReader(b))
}

func TestPeekingForZero(t *testing.T) {
	data := []byte{0, 0, 0, 0, 0, 0, 0, 0}
	reader := bytes.NewReader(data)
	br := bitreader.Create(reader)

	for i := uint(1); i < 64; i++ {
		actual := br.Peek(i)
		expected := uint32(0)
		if actual != expected {
			t.Fatalf("Expected %d, got %d", expected, actual)
		}
	}
}

func TestPeekingForOne(t *testing.T) {
	data := []byte{255, 255, 255, 255, 255, 255, 255, 255}
	reader := bytes.NewReader(data)
	br := bitreader.Create(reader)

	for i := uint(1); i < 64; i++ {
		actual := br.Peek(i)
		expected := uint32(1<<i - 1)
		if actual != expected {
			t.Fatalf("Expected %d, got %d", expected, actual)
		}
	}
}

func TestTrashingBits(t *testing.T) {
	br := createReader(1)
	br.Trash(7)
	if br.Peek(1) != 1 {
		t.Fatal("Expected to read one")
	}
}

func check(t *testing.T, br bitreader.Bitreader, len uint, expected uint32) {
	if actual := br.Read(len); actual != expected {
		t.Fatalf("Expected %d, got %d", expected, actual)
	}
}

func TestReadingBits(t *testing.T) {
	// 0000 0001 0000 0001 0000 0010 0000 0100
	br := createReader(1, 1, 2, 4)
	check(t, br, 7, 0)
	check(t, br, 1, 1)
	check(t, br, 7, 0)
	check(t, br, 1, 1)
	check(t, br, 6, 0)
	check(t, br, 1, 1)
	check(t, br, 6, 0)
	check(t, br, 1, 1)
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
