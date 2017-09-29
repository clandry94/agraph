package agraph

import (
	"bytes"
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
