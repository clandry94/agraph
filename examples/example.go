package main

import (
	"github.com/clandry94/agraph"
	"fmt"
)

func main() {
	ag, err := agraph.New("example.mp3");
	if err != nil {
		fmt.Println(err)
	}

	firstNode := agraph.NewNode(agraph.NopFilter)
	secondNode := agraph.NewNode(agraph.VolumeFilter)
	thirdNode := agraph.NewNode(agraph.NopFilter)

	firstNode.Sink = secondNode.Source
	secondNode.Sink = thirdNode.Source

}
