package main

import (
	"github.com/clandry94/agraph"
	"fmt"
)

func main() {
	_, err := agraph.New("example.mp3");
	if err != nil {
		fmt.Println(err)
	}

	firstNode, _ := agraph.NewNode(agraph.NopFilter)
	secondNode, _ := agraph.NewNode(agraph.VolumeFilter)
	thirdNode, _ := agraph.NewNode(agraph.NopFilter)

	firstNode.SetSink(secondNode.Source())
	secondNode.SetSink(thirdNode.Source())

}
