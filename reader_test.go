package iff

import (
	"io/ioutil"
	"testing"
)

func TestReadIFF(t *testing.T) {
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

		for _, chunk := range iff.Chunks {
			t.Logf("Chunk ID : %s", string(chunk.ChunkID[:]))
		}

		if len(iff.Chunks) != testFile.ChunkSize {
			t.Fatalf("Invalid length of chunks")
		}

		if iff.FileSize != testFile.FileSize {
			t.Fatalf("File size if invalid: %d", iff.FileSize)
		}

		if string(iff.FileType[:]) != testFile.FileType {
			t.Fatalf("File type is invalid: %s", iff.FileType)
		}

		data, err := ioutil.ReadAll(iff.Chunks[0])
		if err != nil {
			t.Fatalf("Can't read data from chunk")
		}

		t.Logf("Length of the first chunk: %d", len(data))
		file.Close()
	}
}

func TestReadNotIFF(t *testing.T) {
	file, err := fixtureFile("../reader.go")
	if err != nil {
		t.Fatalf("Failed to open fixture file")
	}

	reader := NewReader(file)
	_, err = reader.Read()

	if err.Error() != "Given bytes is not a IFF format" {
		t.Fatal("Non-IFF file should not be read")
	}
}
