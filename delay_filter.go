package agraph


/*
	Changes volume amount
*/
type Delay struct {
	source      chan []uint16
	sink        chan []uint16
	Delay       int // something such as 1.2, 0.3, etc
	Decay       float32
	i           int
	prevSamples []uint16
	delayBuffer []float32
}

func newDelay(delay int, decay float32) (Node, error) {
	return &Delay{
		source:      make(chan []uint16, SOURCE_SIZE),
		sink:        nil,
		Delay:       delay,
		i:           0,
		delayBuffer: make([]float32, delay),
		Decay:       decay,
	}, nil
}

func (n *Delay) SetSink(c chan []uint16) {
	n.sink = c
}

func (n *Delay) Source() chan []uint16 {
	return n.source
}

func (n *Delay) Sink() chan []uint16 {
	return n.sink
}

func (n *Delay) Process() error {
	for {
		select {
		case data := <-n.source:
			sample := data[0]
			if n.i < 0 {
				n.i += n.Delay
			}

			delayedSample := n.delayBuffer[n.i]
			filteredData := uint16(delayedSample)

			n.delayBuffer[n.i] = (delayedSample * n.Decay) + float32(sample)
			n.i++
			if n.i >= n.Delay {
				n.i -= n.Delay
			}

			data[0] = filteredData
			n.sink <- data
		}
	}
	return nil
}

func (n *Delay) do(data []uint16) ([]uint16, error) {
	sample := data[0]

	data[0] = sample
	return data, nil
}
