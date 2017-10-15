package main

import (
	"encoding/binary"
	"fmt"
	"github.com/clandry94/agraph"
	"os"
	"time"
)

func main() {
	_, err := agraph.New("ringbackA.wav")
	if err != nil {
		fmt.Println(err)
	}

	file, err := os.OpenFile("ringbackA.wav", os.O_RDWR, 066)
	if err != nil {
		fmt.Println(err)
	}

	reader, err := agraph.NewWaveReader(file)
	if err != nil {
		fmt.Println(err)
	}

	f, err := os.Create("output.wav")
	defer f.Close()
	if err != nil {
		fmt.Println(err)
	}

	writer, err := agraph.NewWaveWriter(f,
		agraph.NumChannels(int(reader.Fmt.Data.NumChannels)),
		agraph.SampleRate(int(reader.Fmt.Data.SampleRate)),
		agraph.BitsPerSample(int(reader.Fmt.Data.BitsPerSample)))

	firstNode, _ := agraph.NewNode(agraph.NopFilter, "nop1")
	secondNode, _ := agraph.NewNode(agraph.NopFilter, "nop2")

	firstNode.SetSink(secondNode.Source())
	secondNode.SetSink(make(chan []uint16, agraph.SOURCE_SIZE))

	go firstNode.Process()
	go secondNode.Process()

	start := time.Now()

	for {
		data, err := reader.ReadSampleInt16()
		if err != nil {
			fmt.Println(err)
			break
		}
		fmt.Println(data)

		modifiedAudio := make([]byte, 2)

		binary.LittleEndian.PutUint16(modifiedAudio, uint16(data[0]))

		//binary.LittleEndian.PutUint16(modifiedAudio, uint16(data[1]))

		//firstNode.Source() <- data

		//_ = <-secondNode.Sink()
		//filteredData = <-secondNode.Sink()

		//fmt.Println(filtered)
		//fmt.Printf(" %v ", data)

		writer.Write(modifiedAudio)
	}

	writer.Close()

	end := time.Now()
	fmt.Println(end.Sub(start))

}
