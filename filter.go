package agraph

type FilterType int

const (
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
	do(data []uint16) ([]uint16, error)
	SetSink(c chan []uint16)
	Source() chan []uint16
	Sink() chan []uint16
}

type NodeInitOptions struct {
	VolumeMultiplier float32
	Delay            int
	Decay            float32
	MovingAverageLength int
	Angle			 float64
}

type NodeInitOption func(*NodeInitOptions)

func VolumeMultiplier(m float32) NodeInitOption {
	return func(args *NodeInitOptions) {
		args.VolumeMultiplier = m
	}
}

// delay in milliseconds
func DelayLength(m int) NodeInitOption {
	return func(args *NodeInitOptions) {
		args.Delay = m
	}
}

func Decay(m float32) NodeInitOption {
	return func(args *NodeInitOptions) {
		args.Decay = m
	}
}

func Taps(m int) NodeInitOption {
	return func(args *NodeInitOptions) {
		args.MovingAverageLength = m
	}
}

func Angle(m float64) NodeInitOption {
	return func(args *NodeInitOptions) {
		args.Angle = m
	}
}

func NewNode(t FilterType, name string, options ...NodeInitOption) (Node, error) {
	args := &NodeInitOptions{
		VolumeMultiplier: 0,
		Delay:            0,
	}

	for _, option := range options {
		option(args)
	}

	switch t {
	case NopFilter:
		return newNop(name)
	case VolumeFilter:
		return newVolume(name, args.VolumeMultiplier) // increase multiplier
	case DelayFilter:
		return newDelay(name, args.Delay, args.Decay)
	case FIRFilter:
		return newFIR(name, args.MovingAverageLength)
	case LocalizationFilter:
		return newLocalization(name, args.Angle)
	default:
		return newNop("default")
	}
}
