package agraph

import (
	"bytes"
	"io"
)

const (
	maxFileSize  = 2 << 31
	fmtChunkSize = 16

	riffChunkToken = "RIFF"
	wavFormatToken = "WAVE"
	fmtChunkToken  = "fmt "
	// listChunkToken = "LIST"
	dataChunkToken = "data"
)

// WAVE File Format
//0		 ---------------
//		| ChunkID       | 4 bytes
//4		|---------------|
//		| ChunkSize     | 4
//8		|---------------|
//		| Format 		| 4
//12	|---------------\
//		| Subchunk1ID   | 4
//16	|---------------|`
//		| Subchunk1Size | 4
//20	|---------------|
//		| AudioFormat   | 2
//22	|---------------|
//		| NumChannels   | 2
//24	|---------------|
//		| SampleRate    | 4
//28	|---------------|
//		| ByteRate      | 4
//32	|---------------|
//		| BlockAlign    | 2
//34	|---------------|
//		| BitsPerSample | 2
//36	|---------------|
//		| Subchunk2ID   | 4
//40	|---------------|
//		| Subchunk2Size | 4
//44	|---------------|
//		|				|
//		|				|
//		|     data		| Subchunk2Size
//		|		        |
//		|               |
//		 ---------------

// Riff Chunk
type Riff struct {
	ChunkID   []byte
	ChunkSize uint32
	Format    []byte
}

// Fmt Chunk
type Fmt struct {
	ID   []byte
	Size uint32
	Data *WavFmtData
}

type WavFmtData struct {
	AudioFormat   uint16
	NumChannels   uint16
	SampleRate    uint32 /// per second
	ByteRate      uint32
	BlockAlign    uint16 // data block size in bytes aka numChannels * BitsPerSample /8
	BitsPerSample uint16
}

// Data Chunk
type data struct {
	ID   []byte
	Size uint32

	// Holds the data of the wave file. Shouldn't be accessed directly since it has a reader
	data DataReader
}

type DataReader interface {
	io.Reader
	io.ReaderAt
}

type ReadSeeker interface {
	io.Reader
	io.Seeker
	io.ReaderAt
}

type DataWriterChunk struct {
	ID   []byte
	Size uint32
	Data *bytes.Buffer
}
