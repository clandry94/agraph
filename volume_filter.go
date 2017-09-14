package agraph

import "fmt"

/*
	Changes volume amount
*/
type Volume struct {
	source     chan []byte
	sink       chan []byte
	Multiplier int
}

func newVolume(multiplier int) (Node, error) {
	return &Volume{
		source:     make(chan []byte, SOURCE_SIZE),
		sink:       nil,
		Multiplier: multiplier,
	}, nil
}

func (n *Volume) SetSink(c chan []byte) {
	n.sink = c
}

func (n *Volume) Source() chan []byte {
	return n.source
}

func (n *Volume) Process() error {
	for {
		select {
		case data := <-n.source:
			fmt.Println("found data")
			var filteredData, err = n.do(data)

			if err != nil {
				panic("Could not filter!")
			}
			n.sink <- filteredData
		}
	}
	return nil
}

func (n *Volume) do(data []byte) ([]byte, error) {
	return data, nil
}
