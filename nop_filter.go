package agraph

/*
	Null operation filter. Does nothing.
*/
type Nop struct {
	source chan []float64
	sink   chan []float64
	Name string
}

func newNop(name string) (Node, error) {
	return &Nop{
		source: make(chan []float64, SOURCE_SIZE),
		sink:   nil,
		Name: name,
	}, nil
}

func (n *Nop) SetSink(c chan []float64) {
	n.sink = c
}

func (n *Nop) Source() chan []float64 {
	return n.source
}

func (n *Nop) Sink() chan []float64 {
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

func (n *Nop) do(data []float64) ([]float64, error) {
	return data, nil
}
