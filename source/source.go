package source

import (
	"fmt"
	"io"
	"math"
	"os"
	"unsafe"

	"golang.org/x/exp/mmap"
)

type Buffer struct {
	data   []byte
	mmaped bool
}

func shouldUseMmap(size int64) bool {
	return size > 4*4096 || int(size) >= os.Getpagesize()
}

func NewFromFile(filename string) (*Buffer, error) {
	f, err := os.OpenFile(filename, os.O_RDONLY, 0444)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		return nil, err
	}

	size := fi.Size()
	if size >= math.MaxUint32 {
		return nil, fmt.Errorf("file is over the 2GiB input limit (%d bytes)", size)
	}

	if shouldUseMmap(size) {
		m, err := mmap.Open(filename)
		if err != nil {
			return nil, err
		}
		sb := (*Buffer)(unsafe.Pointer(m))
		sb.mmaped = true
		return sb, nil
	}

	data, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}
	return &Buffer{data: data}, nil
}

func NewFromReader(r io.Reader) (*Buffer, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	return &Buffer{data: data}, nil
}

func (b *Buffer) Len() int {
	return len(b.data)
}

func (b *Buffer) At(i int) byte {
	return b.data[i]
}

func (b *Buffer) From(start int) []byte {
	return b.data[start:]
}

func (b *Buffer) Range(start, end int) []byte {
	return b.data[start:end]
}

func (b *Buffer) Bytes() []byte {
	return b.data
}
