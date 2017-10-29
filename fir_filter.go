package agraph


/*
	FIR moving average filter
*/
type FIR struct {
	source chan []uint16
	sink   chan []uint16
	Name   string
	maBuffer [][]uint16 // maybe needs to be a pointer
	maSize int
}

func newFIR(name string, maSize int) (Node, error) {
	return &FIR{
		source: make(chan []uint16, SOURCE_SIZE),
		sink:   nil,
		Name:   name,
		maBuffer: make([][]uint16, 0),
		maSize: maSize,
	}, nil
}

func (n *FIR) SetSink(c chan []uint16) {
	n.sink = c
}

func (n *FIR) Source() chan []uint16 {
	return n.source
}

func (n *FIR) Sink() chan []uint16 {
	return n.sink
}

func (n *FIR) Process() error {
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

func (n *FIR) do(data []uint16) ([]uint16, error) {
	var movAvg float64
	movAvg = 0.0
	if len(n.maBuffer) < n.maSize {
		n.maBuffer = append(n.maBuffer, data)
	} else {
		x := (n.maBuffer)[:n.maSize-1][0]
		z := make([][]uint16, 1)
		val := make([]uint16, 1)
		z[0] = val
		n.maBuffer = append(z, x)
	}

	for _, sample := range n.maBuffer {
		movAvg += float64(sample[0]) * (1 / float64(n.maSize))
	}

	data[0] = uint16(movAvg)

	return data, nil
}



