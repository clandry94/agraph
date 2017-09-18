package main

import (
	"fmt"
	"github.com/clandry94/agraph"
)

func main() {
	ar, err := agraph.New("ringbackA.wav")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(ar.Meta)

	/*
	firstNode, _ := agraph.NewNode(agraph.NopFilter)
	secondNode, _ := agraph.NewNode(agraph.VolumeFilter)
	thirdNode, _ := agraph.NewNode(agraph.NopFilter)

	firstNode.SetSink(secondNode.Source())
	secondNode.SetSink(thirdNode.Source())
	*/
}
