package main

import (
	"fmt"
	"github.com/clandry94/agraph"
	"os"
	"time"
)

const outputFileName = "two_channel_noise_filtered.wav"

func main() {
	file, err := os.OpenFile("M1F1-Alaw-AFsp.wav", os.O_RDWR, 066)
	if err != nil {
		fmt.Println(err)
	}

	reader, err := agraph.NewWaveReader(file)
	if err != nil {
		fmt.Println(err)
	}

	f, err := os.Create(outputFileName)
	defer f.Close()
	if err != nil {
		fmt.Println(err)
	}

	writer, err := agraph.NewWaveWriter(f,
		agraph.NumChannels(int(reader.Fmt.Data.NumChannels)),
		agraph.SampleRate(int(reader.Fmt.Data.SampleRate)),
		agraph.BitsPerSample(int(reader.Fmt.Data.BitsPerSample)))

	meta := fmt.Sprintf("File info: \n" +
		"- Filename: %v \n" +
		"- Audio Format: %v \n" +
		"- Number of Channels: %v \n" +
		"- Sample Rate: %v \n" +
		"- Bit Depth: %v \n",
		f.Name(),
		reader.Fmt.Data.AudioFormat,
		reader.Fmt.Data.NumChannels,
		reader.Fmt.Data.SampleRate,
		reader.Fmt.Data.BitsPerSample)

	fmt.Println(meta)

	start := time.Now()

	i := 0
	for {

		data, err := reader.ReadSampleRaw()
		if err != nil {
			fmt.Println(err)
			break
		}

		data[1] = 0
		data[0] = 0

		writer.Write(data)

		// [FL][FL][FR][FR]
		//fmt.Println(data)

		//binary.LittleEndian.PutUint16(modifiedDataAsBytes[0:1], data[])

		//fmt.Print("not writing ")

		i++
	}

	writer.Close()

	end := time.Now()
	fmt.Println(end.Sub(start))

	writtenMeta := fmt.Sprintf("File info: \n" +
		"- Filename: %v \n" +
		"- Audio Format: %v \n" +
		"- Number of Channels: %v \n" +
		"- Sample Rate: %v \n" +
		"- Bit Depth: %v \n",
		outputFileName,
		writer.Fmt.Data.AudioFormat,
		writer.Fmt.Data.NumChannels,
		writer.Fmt.Data.SampleRate,
		writer.Fmt.Data.BitsPerSample)

	fmt.Println(writtenMeta)
}
