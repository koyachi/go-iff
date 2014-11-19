package iff

import (
	"io"
)

type byteReader struct {
	offset uint32
	io.ReaderAt
}

func newByteReader(r io.ReaderAt) (bytes *byteReader) {
	bytes = &byteReader{0, r}

	return
}

func (bytes *byteReader) readBEUint32() uint32 {
	offset := bytes.offset
	data := make([]byte, 4)

	n, err := bytes.ReadAt(data, int64(offset))
	if err != nil || n < 4 {
		panic("Can't read bytes")
	}

	defer func() {
		bytes.offset += 4
	}()

	return uint32(data[0])<<24 +
		uint32(data[1])<<16 +
		uint32(data[2])<<8 +
		uint32(data[3])
}

func (bytes *byteReader) readBEUint16() uint16 {
	offset := bytes.offset
	data := make([]byte, 2)

	n, err := bytes.ReadAt(data, int64(offset))
	if err != nil || n < 2 {
		panic("Can't read bytes")
	}

	defer func() {
		bytes.offset += 2
	}()

	return uint16(data[0])<<8 + uint16(data[1])
}

func (bytes *byteReader) readBEInt16() int16 {
	offset := bytes.offset
	data := make([]byte, 2)

	n, err := bytes.ReadAt(data, int64(offset))

	if err != nil || n < 2 {
		panic("Can't read bytes")
	}

	defer func() {
		bytes.offset += 2
	}()

	return int16(data[offset+0])<<8 + int16(data[offset+1])
}

func (bytes *byteReader) readBytes(size uint32) []byte {
	offset := bytes.offset
	data := make([]byte, size)

	n, err := bytes.ReadAt(data, int64(offset))
	if err != nil || n < int(size) {
		panic("Can't read bytes")
	}

	defer func() {
		bytes.offset += size
	}()

	return data
}
