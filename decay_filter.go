package agraph

import (
	"fmt"
)

/*
	Changes volume amount
*/
type Reverb struct {
	source     chan []uint16
	sink       chan []uint16
	Delay int // something such as 1.2, 0.3, etc
	Decay float32
	i int
	prevSamples []uint16
	delayBuffer []float32
}

func newReverb(delay int, decay float32) (Node, error) {
	return &Reverb{
		source:     make(chan []uint16, SOURCE_SIZE),
		sink:       nil,
		Delay: delay,
		i: 0,
		delayBuffer: make([]float32, delay),
		Decay: decay,
	}, nil
}

func (n *Reverb) SetSink(c chan []uint16) {
	n.sink = c
}

func (n *Reverb) Source() chan []uint16 {
	return n.source
}

func (n *Reverb) Sink() chan []uint16 {
	return n.sink
}

func (n *Reverb) Process() error {
	fmt.Println("Starting up!")
	//delayIndex := n.i - n.Delay

	for {
		select {
		case data := <-n.source:
			sample := data[0]
			fmt.Println(n.i)
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

func (n *Reverb) do(data []uint16) ([]uint16, error) {
	sample := data[0]

	data[0] = sample
	return data, nil
}
