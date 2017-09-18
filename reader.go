package agraph

import (
	"io"
	"os"
	"fmt"
	"io/ioutil"
	"bytes"
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

func NewWaveReader(fp *os.File) (reader *WaveReader, err error) {
	defer fp.Close()
	fStat, err := fp.Stat()
	if err != nil {
		return reader, err
	}

	if fStat.Size() > maxFileSize {
		return reader, fmt.Errorf("File size (%v bytes) is too large", fStat.Size())
	}

	data, err := ioutil.ReadAll(fp)
	if err != nil {
		return reader, err
	}

	reader.size = fStat.Size()
	reader.in = bytes.NewReader(data)

	/*
		TODO: Need to parse riff, fmt, list, and data chunks here
	 */

	return reader, nil
}

//func WavFormatReader(r io.Reader, n int64) io.Reader { return &WaveReader{r, n }}

func (w WaveReader) Read(p []byte) (n int, err error) { return 0, nil }
