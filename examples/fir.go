package main

import (
	"encoding/binary"
	"fmt"
	"os"
	"time"

	"github.com/clandry94/agraph"
)

func main() {
	file, err := os.OpenFile("buttons.wav", os.O_RDWR, 066)
	if err != nil {
		fmt.Println(err)
	}

	reader, err := agraph.NewWaveReader(file)
	if err != nil {
		fmt.Println(err)
	}

	f, err := os.Create("buttons_fir.wav")
	if err != nil {
		fmt.Println(err)
	}
	defer f.Close()

	writer, err := agraph.NewWaveWriter(f,
		agraph.NumChannels(int(reader.Fmt.Data.NumChannels)),
		agraph.SampleRate(int(reader.Fmt.Data.SampleRate)),
		agraph.BitsPerSample(int(reader.Fmt.Data.BitsPerSample)))

	firNode, _ := agraph.NewNode(agraph.FIRFilter,
		"moving_average",
		agraph.Taps(10))

	firNode.SetSink(make(chan []uint16, 0))

	go firNode.Process()
	start := time.Now()

	for {
		data, err := reader.ReadSampleInt16()
		if err != nil {
			fmt.Println(err)
			break
		}

		firNode.Source() <- data
		modifiedData := <-firNode.Sink()

		modifiedDataAsBytes := make([]byte, 2)
		binary.LittleEndian.PutUint16(modifiedDataAsBytes, modifiedData[0])

		writer.Write(modifiedDataAsBytes)
	}

	writer.Close()

	end := time.Now()
	fmt.Println(end.Sub(start))
}
