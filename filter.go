package agraph

import "fmt"

type FilterType int

const (
	// Filters
	NopFilter    FilterType = 1
	VolumeFilter FilterType = 2

	SOURCE_SIZE = 512
)

type FilterGraph struct {
	head *Node
}

func (g FilterGraph) Insert(node *Node) error {
	return nil
}

func (g FilterGraph) Remove(node *Node) error {
	return nil
}

func (g FilterGraph) Process(data []uint8) {

}

// Filters are implemented as structs which implement the type Node. Filters
// are initialized with their source channel created and sink channel as nil.
// In order to connect filters, set the sink of each node to the next node's source.
//
// Example:
//     ----------   ----------   ----------
//  ->|  Node 1 |->|  Node 1 |->|  Node 1 |->
//    ----------   ----------   ----------
//
//	ag, err := agraph.New("example.mp3");
//	if err != nil {
//		fmt.Println(err)
//	}
//
//	firstNode := agraph.NewNode(agraph.NopFilter)
//	secondNode := agraph.NewNode(agraph.VolumeFilter)
//	thirdNode := agraph.NewNode(agraph.NopFilter)
//
//	firstNode.Sink = secondNode.Source
//	secondNode.Sink = thirdNode.Source

type Node interface {
	process() error
	do() ([]byte, error)
}

func NewNode(t FilterType) interface{} {
	switch t {
	case NopFilter:
		return newNop()
	case VolumeFilter:
		return newVolume(3) // increase multipler
	default:
		return newNop()
	}
}

/*
	Null operation filter. Does nothing.
*/
type Nop struct {
	Source chan []byte
	Sink   chan []byte
}

func newNop() interface{} {
	return Nop{
		Source: make(chan []byte, SOURCE_SIZE),
		Sink:   nil,
	}
}

func (n Nop) process() error {
	for {
		select {
		case data := <-n.Source:
			fmt.Println("found data")
			var filteredData, err = n.do(data)

			if err != nil {
				panic("Could not filter!")
			}
			n.Sink <- filteredData
		}
	}
	return nil
}

func (n Nop) do(data []byte) ([]byte, error) {
	return data, nil
}

/*
	Changes volume amount
*/
type Volume struct {
	Source     chan []byte
	Sink       chan []byte
	Multiplier int
}

func newVolume(multiplier int) interface{} {
	return Volume{
		Source:     make(chan []byte, SOURCE_SIZE),
		Sink:       nil,
		Multiplier: multiplier,
	}
}

func (n Volume) process() error {
	for {
		select {
		case data := <-n.Source:
			fmt.Println("found data")
			var filteredData, err = n.do(data)

			if err != nil {
				panic("Could not filter!")
			}
			n.Sink <- filteredData
		}
	}
	return nil
}

func (n Volume) do(data []byte) ([]byte, error) {
	return data, nil
}
