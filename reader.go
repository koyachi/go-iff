package iff

import (
	"errors"
	"io"
)

type IFFReader interface {
	io.Reader
	io.ReaderAt
}

type Reader struct {
	IFFReader
}

type IFFChunk struct {
	FileSize uint32
	FileType []byte
	Chunks   []*Chunk
}

type Chunk struct {
	ChunkID   []byte
	ChunkSize uint32
	IFFReader
}

func NewReader(r IFFReader) *Reader {
	return &Reader{r}
}

func (r *Reader) Read() (chunk *IFFChunk, err error) {
	chunk, err = readIFFChunk(r)

	return
}

func readIFFChunk(r *Reader) (chunk *IFFChunk, err error) {
	bytes := newByteReader(r)

	chunkId := bytes.readBytes(4)
	if string(chunkId[:]) != "FORM" {
		err = errors.New("Given bytes is not a IFF format")
		return
	}

	fileSize := bytes.readBEUint32()
	fileType := bytes.readBytes(4)

	chunk = &IFFChunk{fileSize, fileType, make([]*Chunk, 0)}

	for bytes.offset < fileSize {
		chunkId = bytes.readBytes(4)
		chunkSize := bytes.readBEUint32()
		offset := bytes.offset

		if chunkSize%2 == 1 {
			chunkSize += 1
		}

		bytes.offset += chunkSize

		chunk.Chunks = append(
			chunk.Chunks,
			&Chunk{
				chunkId,
				chunkSize,
				io.NewSectionReader(r, int64(offset), int64(chunkSize))})
	}

	return
}
