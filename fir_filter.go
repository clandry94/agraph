package agraph

import "container/list"

/*
	FIR moving average filter
*/
type FIR struct {
	source    chan []uint16
	sink      chan []uint16
	Name      string
	tapBuffer *list.List // maybe needs to be a pointer
	tapCount  int
}

func newFIR(name string, maSize int) (Node, error) {
	return &FIR{
		source:    make(chan []uint16, SOURCE_SIZE),
		sink:      nil,
		Name:      name,
		tapBuffer: list.New(),
		tapCount:  maSize,
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

	if n.tapBuffer.Len() >= n.tapCount {
		n.tapBuffer.Remove(n.tapBuffer.Back())
	}
	n.tapBuffer.PushFront(data)

	p := n.tapBuffer.Front()
	for p.Next() != nil {
		movAvg += float64(p.Value.([]uint16)[0]) * (1 / float64(n.tapCount))
		p = p.Next()
	}

	data[0] = uint16(movAvg)

	return data, nil
}



