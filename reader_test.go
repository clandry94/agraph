package agraph

import (
	"os"
	"testing"
)

func TestNewWaveReader(t *testing.T) {
	file, err := os.OpenFile("examples/ringbackA.wav", os.O_RDWR, 066)
	if err != nil {
		t.Error(err)
	}

	reader, err := NewWaveReader(file)
	if err != nil {
		t.Error(err)
	}

	expectedID := "RIFF"
	expectedSize := uint32(160038) // take off the 8 bits??
	expectedFormat := "WAVE"

	if string(reader.Riff.ChunkID) != expectedID {
		t.Errorf("Actual chunk ID %v is not equal to expected chunk ID %v", expectedID, reader.Riff.ChunkID)
	}

	if reader.Riff.ChunkSize != expectedSize {
		t.Errorf("Actual chunk size %v is not equal to expected chunk size %v", expectedSize, reader.Riff.ChunkSize)
	}
	if string(reader.Riff.Format) != expectedFormat {
		t.Errorf("Actual format %v is not equal to expected formmat %v", expectedFormat, reader.Riff.Format)
	}
}
