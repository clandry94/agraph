package agraph


/*
	Null operation filter. Does nothing.

	Assuming a simple, spherical head and that the sound sources are
	infinitely far away s.t. the sound reaches the ears
    in a straight line
*/

type Nop struct {
	source chan []uint16
	sink   chan []uint16
	Name   string
}

func newNop(name string) (Node, error) {

	return &Nop{
		source: make(chan []uint16, SOURCE_SIZE),
		sink:   nil,
		Name:   name,
	}, nil
}

func (n *Nop) SetSink(c chan []uint16) {
	n.sink = c
}

func (n *Nop) Source() chan []uint16 {
	return n.source
}

func (n *Nop) Sink() chan []uint16 {
	return n.sink
}

func (n *Nop) Process() error {
	for {
		select {
		case data := <-n.source:
			var filteredData, err = n.do(data)
			//fmt.Printf("Data processed from %v, here it is: %v\n", n.Name, filteredData)
			if err != nil {
				panic("Could not filter!")
			}
			n.sink <- filteredData
		}
	}
	return nil
}

func (n *Nop) do(data []uint16) ([]uint16, error) {
	return data, nil
}

