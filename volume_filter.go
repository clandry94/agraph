package agraph

import "fmt"

/*
	Changes volume amount
*/
type Volume struct {
	source     chan []float64
	sink       chan []float64
	Multiplier int
}

func newVolume(multiplier int) (Node, error) {
	return &Volume{
		source:     make(chan []float64, SOURCE_SIZE),
		sink:       nil,
		Multiplier: multiplier,
	}, nil
}

func (n *Volume) SetSink(c chan []float64) {
	n.sink = c
}

func (n *Volume) Source() chan []float64 {
	return n.source
}

func (n *Volume) Sink() chan []float64 {
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

func (n *Volume) do(data []float64) ([]float64, error) {
	return data, nil
}
