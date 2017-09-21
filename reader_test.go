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
	expectedSubChunk1ID := "fmt "
	expectedSubChunk1Size := uint32(16)
	expectedAudioFormat := uint16(1)
	expectedNumChannels := uint16(1)
	expectedSampleRate := uint32(8000)
	expectedByteRate := uint32(16000)
	expectedBlockAlign := uint16(2)
	expectedBitsPerSamp := uint16(16)

	if string(reader.Riff.ChunkID) != expectedID {
		t.Errorf("Actual chunk ID %v is not equal to expected chunk ID %v", expectedID, reader.Riff.ChunkID)
	}

	if reader.Riff.ChunkSize != expectedSize {
		t.Errorf("Actual chunk size %v is not equal to expected chunk size %v", expectedSize, reader.Riff.ChunkSize)
	}
	if string(reader.Riff.Format) != expectedFormat {
		t.Errorf("Actual format %v is not equal to expected formmat %v", expectedFormat, reader.Riff.Format)
	}

	if string(reader.Fmt.ID) != expectedSubChunk1ID {
		t.Errorf("Actual subChunk1ID %v is not equal to expected subChunk1 ID %v", string(reader.Fmt.ID), expectedSubChunk1ID)
	}

	if reader.Fmt.Size != expectedSubChunk1Size {
		t.Errorf("Actual subChunk1Size %v is not equal to expected subchunk1Size %v", reader.Fmt.Size, expectedSubChunk1Size)
	}

	if reader.Fmt.Data.AudioFormat != expectedAudioFormat {
		t.Errorf("Actual audio-format %v is not equal to expected audio-format %v", reader.Fmt.Data.AudioFormat, expectedAudioFormat)
	}

	if reader.Fmt.Data.NumChannels != expectedNumChannels {
		t.Errorf("Actual num-channels %v is not equal to expected num-channels %v", reader.Fmt.Data.NumChannels, expectedNumChannels)
	}

	if reader.Fmt.Data.SampleRate != expectedSampleRate {
		t.Errorf("Actual sample-rate %v is not equal to expected sample-rate %v", reader.Fmt.Data.SampleRate, expectedSampleRate)
	}

	if reader.Fmt.Data.ByteRate != expectedByteRate {
		t.Errorf("Actual byte-rate %v is not equal to expected byte-rate %v", reader.Fmt.Data.ByteRate, expectedByteRate)
	}

	if reader.Fmt.Data.BlockAlign != expectedBlockAlign {
		t.Errorf("Actual block-align %v is not equal to expected block-align %v", reader.Fmt.Data.BlockAlign, expectedBlockAlign)
	}

	if reader.Fmt.Data.BitsPerSample != expectedBitsPerSamp {
		t.Errorf("Actual bits-per-sample %v is not equal to expected bits-per-sample %v", reader.Fmt.Data.BitsPerSample, expectedBitsPerSamp)
	}
}
