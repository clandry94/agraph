package agraph

import (
	"os"
)

type Graph struct {
	Meta   Meta
	source *[]float64
	sink   *[]float64
	root   *Node
}

type Meta struct {
	Filepath string
	Size     int64
}

func New(filepath string) (*Graph, error) {
	_, err := os.OpenFile(filepath, os.O_RDWR, 0666)
	if err != nil {
		return nil, err
	}

	meta := Meta{
		Filepath: filepath,
	}

	/*
		reader, err := NewWaveReader(file)
		if err != nil {
			return &Graph{}, err
		}
	*/

	/*
		codec := Mp3{
			file: file,
		}
	*/

	return &Graph{
		Meta:   meta,
		source: nil,
		sink:   nil,
		root:   nil,
	}, nil
}
