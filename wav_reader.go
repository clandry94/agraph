package agraph

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"os"
	"io"
)

/*
	Implementation specific to WAVE
*/
type WaveReader struct {
	in   ReadSeeker
	size int64

	Riff *Riff
	Fmt  *Fmt
	data *data

	dataSource  int64
	SampleCount uint32
	SampleNum   int32
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

	reader.SampleCount = reader.data.Size / uint32(reader.Fmt.Data.BlockAlign)
	reader.SampleTime = int(reader.SampleCount / reader.Fmt.Data.SampleRate)

	return reader, nil
}

// Reads into a byte slice
func (r *WaveReader) Read(b []byte) (int, error) {
	return r.data.data.Read(b)
}

// Returns a byte slice the size of one sample
func (r *WaveReader) ReadSampleRaw() ([]byte, error) {
	size := r.Fmt.Data.BlockAlign
	b := make([]byte, size)
	n, err := r.Read(b)

	if err == nil {
		r.SampleNum += 1
	} else {
		return b, err
	}

	if uint16(n) != size {
		return b, fmt.Errorf("sample read is not the "+
			"same size as the sample size: %v != %v", n, size)
	}

	return b, err
}

func (r *WaveReader) ReadSampleInt16() ([]uint16, error) {
	rawSample, err := r.ReadSampleRaw()
	if err != nil {
		return nil, err
	}

	sample := make([]uint16, int(r.Fmt.Data.NumChannels))
	length := len(rawSample) / int(r.Fmt.Data.NumChannels)

	for i := 0; i < int(r.Fmt.Data.NumChannels); i++ {
		sample[i] = bytesToInt(rawSample[length * i : length * (i + 1)])
		if err != nil && err != io.EOF {
			return sample, err
		}
	}

	return sample, nil
}


func bytesToInt(b []byte) uint16 {
	var ret uint16
	switch len(b) {
	case 1:
		// 0 ~ 128 ~ 255
		ret = uint16(b[0])
	case 2:
		// -32768 ~ 0 ~ 32767
		ret = uint16(b[0]) + uint16(b[1])<<8
		//	fmt.Printf("%08b %08b ", b[1], b[0])
		//	fmt.Printf("%016b => %d\n", ret, ret)
	case 3:
		// HiReso / DVDAudio
		ret = uint16(b[0]) + uint16(b[1])<<8 + uint16(b[2])<<16
	default:
		ret = 0
	}
	return ret
}

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

	r.Fmt = &Fmt{
		ID:   subChunk1ID,
		Size: subChunk1Size,
		Data: &WavFmtData{
			AudioFormat:   audioFormat,
			NumChannels:   numChannels,
			SampleRate:    sampleRate,
			ByteRate:      byteRate,
			BlockAlign:    blockAlign,
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

	b := make([]byte, subChunk2Size)

	err = binary.Read(r.in, binary.LittleEndian, b)
	if err != nil {
		return err
	}

	dataReader := bytes.NewReader(b)

	r.data = &data{
		ID:   subChunk2ID,
		Size: subChunk2Size,
		data: dataReader,
	}

	return nil
}
