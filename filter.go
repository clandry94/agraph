package agraph

type FilterType int

const (
	// Filters
	NopFilter    FilterType = 1
	VolumeFilter FilterType = 2

	SOURCE_SIZE = 512
)

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
	Process() error
	do(data []byte) ([]byte, error)
	SetSink(c chan []byte)
	Source() chan []byte
}

func NewNode(t FilterType) (Node, error) {
	switch t {
	case NopFilter:
		return newNop()
	case VolumeFilter:
		return newVolume(3) // increase multipler
	default:
		return newNop()
	}
}
