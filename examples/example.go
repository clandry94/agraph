package main

import (
	"fmt"
	"github.com/clandry94/agraph"
)

func main() {
	ar, err := agraph.New("test.mp3")
	if err != nil {
		fmt.Println(err)
	}

	firstNode, _ := agraph.NewNode(agraph.NopFilter)
	secondNode, _ := agraph.NewNode(agraph.VolumeFilter)
	thirdNode, _ := agraph.NewNode(agraph.NopFilter)

	firstNode.SetSink(secondNode.Source())
	secondNode.SetSink(thirdNode.Source())

	fmt.Println(ar.MetaData.Id3.Title())
	fmt.Println(ar.MetaData.Id3.Artist())
	fmt.Println(ar.MetaData.Id3.Year())
	fmt.Println(ar.MetaData.Id3.Genre())

}
