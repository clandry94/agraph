package agraph

import (
	"fmt"
	"os"
)

type Graph struct {
	Meta Meta
	Codec    Codec
	source   *[]byte
	sink     *[]byte
	root     *Node
}

type Meta struct {
	Filepath string
	Size     int64
}

func New(filepath string) (*Graph, error) {
	file, err := os.OpenFile(filepath, os.O_RDWR, 0666)
	if err != nil {
		return nil, err
	}

	meta := Meta{
		Filepath: filepath,
	}

	reader, err := NewWaveReader(file)
	if err != nil {
		return &Graph{}, err
	}

	fmt.Print(reader)


	/*
	codec := Mp3{
		file: file,
	}
	*/

	return &Graph{
		Meta: meta,
		Codec:    nil,
		source:   nil,
		sink:     nil,
		root:     nil,
	}, nil
}
