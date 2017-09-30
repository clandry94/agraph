package agraph

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

/*
	Implementation specific to WAVE
*/
type WaveWriter struct {
	Out                 io.WriteCloser
	SamplesWrittenCount int

	Riff *Riff
	Fmt  *Fmt
	Data *DataWriterChunk
}

type Options struct {
	NumChannels   int
	SampleRate    int
	BitsPerSample int
}

type Option func(*Options)

func NumChannels(channels int) Option {
	return func(args *Options) {
		args.NumChannels = channels
	}
}

func SampleRate(sampleRate int) Option {
	return func(args *Options) {
		args.SampleRate = sampleRate
	}
}

func BitsPerSample(bitsPerSample int) Option {
	return func(args *Options) {
		args.BitsPerSample = bitsPerSample
	}
}

func NewWaveWriter(out io.WriteCloser, opts ...Option) (*WaveWriter, error) {
	args := &Options{
		NumChannels:   1,
		SampleRate:    44100,
		BitsPerSample: 8,
	}

	for _, opt := range opts {
		opt(args)
	}

	writer := &WaveWriter{}
	writer.Out = out

	blockAlign := uint16(args.BitsPerSample*args.NumChannels) / 8
	byteRate := uint32(int(blockAlign) * args.SampleRate)

	writer.Riff = &Riff{
		ChunkID: []byte(riffChunkToken),
		Format:  []byte(wavFormatToken),
	}

	writer.Fmt = &Fmt{
		ID:   []byte(fmtChunkToken),
		Size: uint32(fmtChunkSize),
	}

	writer.Fmt.Data = &WavFmtData{
		AudioFormat:   uint16(1), // pulse code modulation (PCM)
		NumChannels:   uint16(args.NumChannels),
		SampleRate:    uint32(args.SampleRate),
		ByteRate:      byteRate,
		BlockAlign:    blockAlign,
		BitsPerSample: uint16(args.BitsPerSample),
	}

	writer.Data = &DataWriterChunk{
		ID:   []byte(dataChunkToken),
		Data: bytes.NewBuffer([]byte{}),
	}

	return writer, nil
}

func (w *WaveWriter) Write(b []byte) (int, error) {
	blockAlign := int(w.Fmt.Data.BlockAlign)
	if len(b) < blockAlign {
		return 0, fmt.Errorf("Size of b %v < BlockAlign %v", len(b), blockAlign)
	}

	if len(b)%blockAlign != 0 {
		return 0, fmt.Errorf("Sizez of b %v must be a multiple of BlockAlign %v", len(b), blockAlign)
	}

	numBytesWritten := len(b) / blockAlign

	n, err := w.Data.Data.Write(b)
	if err != nil {
		w.SamplesWrittenCount += numBytesWritten
	}

	return n, err
}

func (w *WaveWriter) Close() error {
	data := w.Data.Data.Bytes()
	length := uint32(len(data))

	// possibly shouldn't have the (8 + length) on the end
	w.Riff.ChunkSize = uint32(len(w.Riff.ChunkID)) + (8 + w.Fmt.Size) + (8 + length)
	w.Data.Size = length

	// Write the riff chunk
	err := binary.Write(w.Out, binary.BigEndian, w.Riff.ChunkID)
	err = binary.Write(w.Out, binary.LittleEndian, w.Riff.ChunkSize)
	err = binary.Write(w.Out, binary.BigEndian, w.Riff.Format)

	// Write the fmt chunk
	err = binary.Write(w.Out, binary.BigEndian, w.Fmt.ID)
	err = binary.Write(w.Out, binary.LittleEndian, w.Fmt.Size)
	err = binary.Write(w.Out, binary.LittleEndian, w.Fmt.Data)

	// Write the data chunk
	err = binary.Write(w.Out, binary.BigEndian, w.Data.ID)
	err = binary.Write(w.Out, binary.LittleEndian, w.Data.Size)
	_, err = w.Out.Write(data)

	err = w.Out.Close()

	if err != nil {
		return err
	}

	return nil
}
