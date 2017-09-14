package agraph

import (
	"fmt"
	id3 "github.com/clandry94/id3-go"
	"os"
)

type Graph struct {
	MetaData    MetaData
	Codec       Codec
	source      *[]byte
	sink        *[]byte
	filterGraph *FilterGraph
}

type MetaData struct {
	Id3      id3.Tagger
	Filepath string
	Size     int64
}

func New(filepath string) (*Graph, error) {
	fileInfo, err := os.Stat(filepath)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	file, err := os.OpenFile(filepath, os.O_RDWR, 0666)
	if err != nil {
		return nil, err
	}

	id3Data, err := id3.Parse(file)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	meta := MetaData{
		Id3:      id3Data.Tagger,
		Filepath: filepath,
		Size:     fileInfo.Size(),
	}

	codec := Mp3{
		file: file,
	}

	return &Graph{
		MetaData:    meta,
		Codec:       codec,
		source:      nil,
		sink:        nil,
		filterGraph: nil,
	}, nil
}
