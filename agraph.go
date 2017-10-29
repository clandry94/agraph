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

// Meta is a type which has information about a file
type Meta struct {
	Filepath string
	Size     int64
}

// New returns a new *Graph if file exists, otherwise returns an error
func New(filepath string) (*Graph, error) {
	f, err := os.Stat(filepath)
	if err != nil {
		return nil, err
	}

	meta := Meta{
		Filepath: filepath,
		Size:     f.Size(),
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
