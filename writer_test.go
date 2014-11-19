package iff

import (
	"io"
	"io/ioutil"
	"os"
	"testing"
)

func TestWriteIFF(t *testing.T) {
	testFiles := []testFile{
		testFile{
			"a.aiff",
			3,
			247784,
			"AIFF",
		},
		//		testFile{
		//		},
	}
	for _, testFile := range testFiles {
		file, err := fixtureFile(testFile.Name)
		if err != nil {
			t.Fatalf("Failed to open fixture file")
		}

		reader := NewReader(file)
		iff, err := reader.Read()
		if err != nil {
			t.Fatal(err)
		}

		outfile, err := ioutil.TempFile("/tmp", "outfile")
		if err != nil {
			t.Fatal(err)
		}
		defer func() {
			outfile.Close()
			os.Remove(outfile.Name())
		}()

		writer := NewWriter(outfile, iff.FileType, iff.FileSize)
		for _, chunk := range iff.Chunks {
			writer.WriteChunk(chunk.ChunkID, chunk.ChunkSize, func(w io.Writer) {
				buf, err := ioutil.ReadAll(chunk)
				if err != nil {
					t.Fatal(err)
				}

				w.Write(buf)
			})
		}
		outfile.Close()

		// reopen to check file content
		file, err = os.Open(outfile.Name())
		if err != nil {
			t.Fatal(err)
		}

		reader = NewReader(file)
		iff, err = reader.Read()
		if err != nil {
			t.Fatal(err)
		}

		for _, chunk := range iff.Chunks {
			t.Logf("Chunk ID: %s", string(chunk.ChunkID[:]))
		}

		if len(iff.Chunks) != testFile.ChunkSize {
			t.Fatalf("Invalid length of chunks")
		}

		if iff.FileSize != testFile.FileSize {
			t.Fatalf("File size is invalid: %d", iff.FileSize)
		}

		if string(iff.FileType[:]) != testFile.FileType {
			t.Fatalf("File type is invalid: %s", iff.FileType)
		}
	}
}
