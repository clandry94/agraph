package agraph

import (
	"fmt"
)

/*
	Changes volume amount
*/
type Volume struct {
	source     chan []uint16
	sink       chan []uint16
	Multiplier float32 // something such as 1.2, 0.3, etc
}

func newVolume(multiplier float32) (Node, error) {
	return &Volume{
		source:     make(chan []uint16, SOURCE_SIZE),
		sink:       nil,
		Multiplier: multiplier,
	}, nil
}

func (n *Volume) SetSink(c chan []uint16) {
	n.sink = c
}

func (n *Volume) Source() chan []uint16 {
	return n.source
}

func (n *Volume) Sink() chan []uint16 {
	return n.sink
}

func (n *Volume) Process() error {
	fmt.Println("Starting up!")

	for {
		select {
		case data := <-n.source:
			filteredData, err := n.do(data)
			if err != nil {
				panic("Could not filter!")
			}
			n.sink <- filteredData
		}
	}
	return nil
}

func (n *Volume) do(data []uint16) ([]uint16, error) {
	sample := data[0]
	sample = sample + 1000
	data[0] = sample
	return data, nil
}
