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

	file, err := os.OpenFile("ringbackA.wav", os.O_RDWR, 066)
	if err != nil {
		fmt.Println(err)
	}

	reader, err := agraph.NewWaveReader(file)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(reader)

	start := time.Now()
	for {
		_, err := reader.ReadRawSample()
		if err != nil {
			fmt.Println(err)
			break
		}
		//fmt.Printf(" %v ", data)
	}
	end := time.Now()
	fmt.Println(end.Sub(start))



	/*
	firstNode, _ := agraph.NewNode(agraph.NopFilter)
	secondNode, _ := agraph.NewNode(agraph.VolumeFilter)
	thirdNode, _ := agraph.NewNode(agraph.NopFilter)

	firstNode.SetSink(secondNode.Source())
	secondNode.SetSink(thirdNode.Source())
	*/
}
