package filter

type Type int

const (
	SOURCE_SIZE = 512

	// Filters
	NopFilter    	   Type = 1
	VolumeFilter	   Type = 2
	DelayFilter  	   Type = 3
	FIRFilter    	   Type = 4
	LocalizationFilter Type = 5
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
	do(data []uint16) ([]uint16, error)
	SetSink(c chan []uint16)
	SetSource(c chan []uint16)
	Source() chan []uint16
	Sink() chan []uint16
}

type MetaData struct {
	SampleRate uint32
	NumChannels uint16
}

type Options struct {
	VolumeMultiplier float32
	Delay            int
	Decay            float32
	MovingAverageLength int
	Angle			 float64
}

func NewNode(t Type, meta MetaData, name string, options Options) (Node, error) {

	switch t {
	case NopFilter:
		return newNop(name, meta)
	case VolumeFilter:
		return newVolume(name, meta, options.VolumeMultiplier) // increase multiplier
	case DelayFilter:
		return newDelay(name, meta, options.Delay, options.Decay)
	case FIRFilter:
		return newFIR(name, meta, options.MovingAverageLength)
	case LocalizationFilter:
		return newLocalization(name, meta, options.Angle)
	default:
		return newNop("default", meta)
	}
}
