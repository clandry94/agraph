package main

import (
	"encoding/binary"
	"fmt"
	"os"
	"time"

	"github.com/clandry94/agraph"
)

func main() {
	file, err := os.OpenFile("imperial_march.wav", os.O_RDWR, 066)
	if err != nil {
		fmt.Println(err)
	}

	reader, err := agraph.NewWaveReader(file)
	if err != nil {
		fmt.Println(err)
	}

	f, err := os.Create("volume_increase.wav")
	if err != nil {
		fmt.Println(err)
	}
	defer f.Close()

	writer, err := agraph.NewWaveWriter(f,
		agraph.NumChannels(int(reader.Fmt.Data.NumChannels)),
		agraph.SampleRate(int(reader.Fmt.Data.SampleRate)),
		agraph.BitsPerSample(int(reader.Fmt.Data.BitsPerSample)))

	volumeNode, _ := agraph.NewNode(agraph.VolumeFilter,
		"volume1",
		agraph.VolumeMultiplier(0.5))

	volumeNode.SetSink(make(chan []uint16, 0))

	go volumeNode.Process()
	start := time.Now()

	for {
		data, err := reader.ReadSampleInt16()
		if err != nil {
			fmt.Println(err)
			break
		}

		volumeNode.Source() <- data
		modifiedData := <-volumeNode.Sink()

		modifiedDataAsBytes := make([]byte, 2)
		binary.LittleEndian.PutUint16(modifiedDataAsBytes, modifiedData[0])

		writer.Write(modifiedDataAsBytes)
	}

	writer.Close()

	end := time.Now()
	fmt.Println(end.Sub(start))

}
