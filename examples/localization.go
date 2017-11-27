package main

import (
	"encoding/binary"
	"fmt"
	"github.com/clandry94/agraph"
	"os"
	"time"
)

func main() {
	file, err := os.OpenFile("M1F1-Alaw-AFsp.wav", os.O_RDWR, 066)
	if err != nil {
		fmt.Println(err)
	}

	reader, err := agraph.NewWaveReader(file)
	if err != nil {
		fmt.Println(err)
	}

	f, err := os.Create("M1F1-Alaw-AFsp_localized.wav")
	defer f.Close()
	if err != nil {
		fmt.Println(err)
	}

	writer, err := agraph.NewWaveWriter(f,
		agraph.NumChannels(int(reader.Fmt.Data.NumChannels)),
		agraph.SampleRate(int(reader.Fmt.Data.SampleRate)),
		agraph.BitsPerSample(int(reader.Fmt.Data.BitsPerSample)))

	localizationNode, err := agraph.NewNode(agraph.LocalizationFilter,
		"localization",
		agraph.Angle(-45))

	if err != nil {
		panic(err)
	}

	localizationNode.SetSink(make(chan []uint16, 0))

	go localizationNode.Process()
	start := time.Now()

	for {
		data, err := reader.ReadSampleInt16()
		if err != nil {
			fmt.Println(err)
			break
		}


		localizationNode.Source() <- data
		modifiedData := <-localizationNode.Sink()
		//fmt.Printf("Localized: %v\n", modifiedData)

		leftByte := make([]byte, 2)
		rightByte := make([]byte, 2)
		binary.LittleEndian.PutUint16(leftByte, modifiedData[0])
		binary.LittleEndian.PutUint16(rightByte, modifiedData[1])

		fmt.Print("Packet Info: \n")
		fmt.Printf(" - Original [%v]uint16 = %v \n", len(data), data)
		// fmt.Printf(" - Modified [%v]uint16 = %v \n", len(data), data)
		fmt.Printf(" - L: [%v]byte = %v \n - R: [%v]byte = %v \n", len(leftByte), leftByte, len(rightByte), rightByte)


		// pack all the bytes into the correct ordering
		fullByte := make([]byte, 4)

		fullByte[0] = leftByte[0]
		fullByte[1] = leftByte[1]
		fullByte[2] = rightByte[0]
		fullByte[3] = rightByte[1]

		// fmt.Printf(" - Stereo: [%v]byte = %v \n", len(fullByte), fullByte)


		// fmt.Println()
		writer.Write(fullByte)
	}

	writer.Close()

	end := time.Now()
	fmt.Println(end.Sub(start))
}
