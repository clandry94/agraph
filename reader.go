package agraph

import (
	"io"
)

/*
	Implementatiton specific to WAVE
*/

const (
	maxFileSize             = 2 << 31
	riffChunkSize           = 12
	listChunkOffset         = 36
	riffChunkSizeBaseOffset = 36
	fmtChunkSize            = 16

	riffChunkToken = "RIFF"
	wavFormatToken = "WAVE"
	fmtChunkToken  = "fmt"
	listChunkToken = "LIST"
	dataChunkToken = "data"
)

// WAVE File Format
//0		 ---------------
//		| ChunkID       | 4 bytes
//4		|---------------|
//		| ChunkSize     | 4
//8		|---------------|
//		| Format 		| 4
//12	|---------------|
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
type Riff struct {
	ID         []byte
	Size       uint32
	FormatType []byte
}

type Fmt struct {
	ID   []byte
	Size uint32
	Data *WavFmtData
}

type WavFmtData struct {
	AudioFormat   uint16
	NumChannels   uint16
	SampleRate    uint32
	ByteRate      uint32
	BlockAlign    uint16 // data block size in bytes
	BitsPerSample uint16
}

type Data struct {
	ID   []byte
	Size uint32
	Data Reader
}

type Reader interface {
	io.Reader
	io.ReaderAt
}

type ReadSeeker interface {
	io.Reader
	io.Seeker
	io.ReaderAt
}

type WaveReader struct {
	in   ReadSeeker
	size int64

	Riff *Riff
	Fmt  *Fmt
	Data *Data

	dataSource int64
	NumSamples uint32
	ReadSample int32
	SampleTime int
}

//func WavFormatReader(r io.Reader, n int64) io.Reader { return &WaveReader{r, n }}

func (w WaveReader) Read(p []byte) (n int, err error) {
	return 0, nil
}
