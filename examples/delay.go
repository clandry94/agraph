package main

import (
	"encoding/binary"
	"fmt"
	"github.com/clandry94/agraph"
	"os"
	"time"
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

	f, err := os.Create("delay_march.wav")
	defer f.Close()
	if err != nil {
		fmt.Println(err)
	}

	writer, err := agraph.NewWaveWriter(f,
		agraph.NumChannels(int(reader.Fmt.Data.NumChannels)),
		agraph.SampleRate(int(reader.Fmt.Data.SampleRate)),
		agraph.BitsPerSample(int(reader.Fmt.Data.BitsPerSample)))

	delayNode, _ := agraph.NewNode(agraph.DelayFilter,
		"delay1",
		agraph.DelayLength(220),
		agraph.Decay(1.0))

	delayNode.SetSink(make(chan []uint16, 0))

	go delayNode.Process()
	start := time.Now()

	for {
		data, err := reader.ReadSampleInt16()
		if err != nil {
			fmt.Println(err)
			break
		}

		delayNode.Source() <- data
		modifiedData := <-delayNode.Sink()

		modifiedDataAsBytes := make([]byte, 2)
		binary.LittleEndian.PutUint16(modifiedDataAsBytes, modifiedData[0])

		writer.Write(modifiedDataAsBytes)
	}

	writer.Close()

	end := time.Now()
	fmt.Println(end.Sub(start))

}
