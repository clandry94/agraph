package agraph

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

/*
	Implementatiton specific to WAVE
*/

const (
	maxFileSize             = 2 << 31
	fmtChunkSize            = 16

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
type Data struct {
	ID   []byte
	Size uint32
	Data []byte
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

	dataSource  int64
	SampleCount uint32
	ReadSample  int32
	SampleTime  int
}

func NewWaveReader(fp *os.File) (*WaveReader, error) {
	defer fp.Close()
	fStat, err := fp.Stat()
	if err != nil {
		return nil, err
	}

	if fStat.Size() > maxFileSize {
		return nil, fmt.Errorf("File size (%v bytes) is too large", fStat.Size())
	}

	data, err := ioutil.ReadAll(fp)
	if err != nil {
		return nil, err
	}

	fmt.Printf("File Size: %v\n", fStat.Size())
	reader := new(WaveReader)
	reader.size = fStat.Size()
	reader.in = bytes.NewReader(data)

	err = reader.parseRiffChunk()
	if err != nil {
		return nil, err
	}

	err = reader.parseFmtChunk()
	if err != nil {
		return nil, err
	}

	err = reader.parseDataChunk()
	if err != nil {
		return nil, err
	}

	reader.SampleCount = reader.Data.Size / uint32(reader.Fmt.Data.BlockAlign)
	reader.SampleTime = int(reader.SampleCount / reader.Fmt.Data.SampleRate)

	return reader, nil
}

/*
	TODO: Clean these parsers up
*/
func (r *WaveReader) parseRiffChunk() error {
	chunkId := make([]byte, 4)

	// Read the RIFF token from the
	err := binary.Read(r.in, binary.BigEndian, chunkId)
	fmt.Printf("Chunk ID in bytes: %v\n", chunkId)
	fmt.Printf("Chunk ID as string: %v\n", string(chunkId))
	if err != nil {
		return err
	}

	if string(chunkId[:]) != riffChunkToken {
		return fmt.Errorf("Not a WAVE file. File is of type %s", string(chunkId[:]))
	}

	// Grab the 16 bit fmt chunk size
	// This is the size of the
	// entire file in bytes minus 8 bytes for the
	// two fields not included in this count:
	// ChunkID and ChunkSize.
	chunkSizeBytes := make([]byte, 4)
	err = binary.Read(r.in, binary.LittleEndian, chunkSizeBytes)
	if err != nil {
		return err
	}

	fmt.Printf("Chunk Size in Bytes: %v\n", chunkSizeBytes)
	chunkSize := binary.LittleEndian.Uint32(chunkSizeBytes)
	fmt.Printf("Chunk size as decimal %v\n", chunkSize)

	if chunkSize != uint32(r.size)-8 {
		return fmt.Errorf("RIFF Chunk Size %v must == file size-8 bytes %v", chunkSize, r.size-8)
	}

	format := make([]byte, 4)
	err = binary.Read(r.in, binary.BigEndian, format)
	if err != nil {
		return err
	}

	if string(format[:]) != wavFormatToken {
		return fmt.Errorf("File is not a WAVE file. It is %s", string(format[:]))
	}

	riff := Riff{
		ChunkID:   chunkId,
		ChunkSize: chunkSize,
		Format:    format,
	}

	r.Riff = &riff

	return nil
}

func (r *WaveReader) parseFmtChunk() error {
	subChunk1ID := make([]byte, 4)
	err := binary.Read(r.in, binary.BigEndian, subChunk1ID)
	if err != nil {
		return err
	}

	if string(subChunk1ID) != fmtChunkToken {
		return fmt.Errorf("invalid data chunk %s", string(subChunk1ID))
	}

	subChunk1Bytes := make([]byte, 4)
	err = binary.Read(r.in, binary.LittleEndian, subChunk1Bytes)
	if err != nil {
		return err
	}
	subChunk1Size := binary.LittleEndian.Uint32(subChunk1Bytes)
	if subChunk1Size != fmtChunkSize {
		return fmt.Errorf("Fmt Chunk Size %v must == %v", subChunk1ID, fmtChunkSize)
	}

	audioFormatByte := make([]byte, 2)
	err = binary.Read(r.in, binary.LittleEndian, audioFormatByte)
	if err != nil {
		return err
	}
	audioFormat := binary.LittleEndian.Uint16(audioFormatByte)

	numChannelsByte := make([]byte, 2)
	err = binary.Read(r.in, binary.LittleEndian, numChannelsByte)
	if err != nil {
		return err
	}
	numChannels := binary.LittleEndian.Uint16(numChannelsByte)

	sampleRateByte := make([]byte, 4)
	err = binary.Read(r.in, binary.LittleEndian, sampleRateByte)
	if err != nil {
		return err
	}
	sampleRate := binary.LittleEndian.Uint32(sampleRateByte)

	byteRateByte := make([]byte, 4)
	err = binary.Read(r.in, binary.LittleEndian, byteRateByte)
	if err != nil {
		return err
	}
	byteRate := binary.LittleEndian.Uint32(byteRateByte)

	blockAlignByte := make([]byte, 2)
	err = binary.Read(r.in, binary.LittleEndian, blockAlignByte)
	if err != nil {
		return err
	}
	blockAlign := binary.LittleEndian.Uint16(blockAlignByte)

	bitsPerSampleByte := make([]byte, 2)
	err = binary.Read(r.in, binary.LittleEndian, bitsPerSampleByte)
	if err != nil {
		return err
	}
	bitsPerSample := binary.LittleEndian.Uint16(bitsPerSampleByte)

	fmt.Printf("AudioFormat: %v\n", audioFormat)
	fmt.Printf("NumChannels: %v\n", numChannels)
	fmt.Printf("SampleRate: %v\n", sampleRate)
	fmt.Printf("ByteRate: %v\n", byteRate)
	fmt.Printf("BlockAlign: %v\n", blockAlign)
	fmt.Printf("BitsPerSample: %v\n", bitsPerSample)


	r.Fmt = &Fmt{
		ID:   subChunk1ID,
		Size: subChunk1Size,
		Data: &WavFmtData {
			AudioFormat: audioFormat,
			NumChannels: numChannels,
			SampleRate: sampleRate,
			ByteRate: byteRate,
			BlockAlign: blockAlign,
			BitsPerSample: bitsPerSample,
		},
	}
	
	return nil
}

func (r *WaveReader) parseDataChunk() error {
	subChunk2ID := make([]byte, 4)
	err := binary.Read(r.in, binary.BigEndian, subChunk2ID)
	if err != nil {
		return err
	}

	if string(subChunk2ID) != dataChunkToken {
		return fmt.Errorf("invalid data chunk %s", string(subChunk2ID))
	}

	subChunk2Bytes := make([]byte, 4)
	err = binary.Read(r.in, binary.LittleEndian, subChunk2Bytes)
	if err != nil {
		return err
	}

	subChunk2Size := binary.LittleEndian.Uint32(subChunk2Bytes)

	data := make([]byte, subChunk2Size)

	err = binary.Read(r.in, binary.LittleEndian, data)
	if err != nil {
		return err
	}

	r.Data = &Data{
		ID: subChunk2ID,
		Size: subChunk2Size,
		Data: data,
	}

	return nil
}