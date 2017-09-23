package main

import (
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

	file, err := os.OpenFile("tone.wav", os.O_RDWR, 066)
	if err != nil {
		fmt.Println(err)
	}

	reader, err := agraph.NewWaveReader(file)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(reader)

	firstNode, _ := agraph.NewNode(agraph.NopFilter, "nop1")
	secondNode, _ := agraph.NewNode(agraph.NopFilter, "nop2")

	firstNode.SetSink(secondNode.Source())
	secondNode.SetSink(make(chan []float64, agraph.SOURCE_SIZE))

	fmt.Println(firstNode.Source())
	fmt.Println(firstNode.Sink())
	fmt.Println(secondNode.Source())
	fmt.Println(secondNode.Sink())

	go firstNode.Process()
	go secondNode.Process()

	start := time.Now()

	for {
		data, err := reader.ReadSampleFloat()
		if err != nil {
			fmt.Println(err)
			break
		}

		firstNode.Source() <- data

		//_ = <-secondNode.Sink()
		filtered := <-secondNode.Sink()

		fmt.Println(filtered)
		//fmt.Printf(" %v ", data)
	}

	end := time.Now()
	fmt.Println(end.Sub(start))



}
