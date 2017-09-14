package agraph

import "fmt"

/*
	Null operation filter. Does nothing.
*/
type Nop struct {
	source chan []byte
	sink   chan []byte
}

func newNop() (Node, error) {
	return &Nop{
		source: make(chan []byte, SOURCE_SIZE),
		sink:   nil,
	}, nil
}

func (n *Nop) SetSink(c chan []byte) {
	n.sink = c
}

func (n *Nop) Source() chan []byte {
	return n.source
}

func (n *Nop) Process() error {
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

func (n *Nop) do(data []byte) ([]byte, error) {
	return data, nil
}
