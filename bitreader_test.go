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
