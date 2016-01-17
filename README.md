# bitreader

Provides basic interfaces to read and traverse an `io.Reader` as a stream of bits, rather than a stream of bytes.

[![GoDoc](https://godoc.org/github.com/32bitkid/bitreader?status.svg)](https://godoc.org/github.com/32bitkid/bitreader)

## Installation

```bash
$ go get github.com/32bitkid/bitreader
```

## Examples

Ever wanted to count the number of 0 bits from the start of a file?

```go

file, _ := os.Open("file")
br := bitreader.NewBitReader(file)
n := 0
for {
    val, err := br.ReadBit()
    if err != nil || val == true {
        break
    }
    n += 1
}
fmt.Printf("The file starts with %d off bits", n)
```

But seriously, this is used for parsing densely packed binary formats where data may not be byte aligned. For example, decoding values packed with [Huffman Coding](https://en.wikipedia.org/wiki/Huffman_coding).
