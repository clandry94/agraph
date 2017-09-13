package agraph

type FilterGraph struct {
	Nodes map[string]Node
}

func (g FilterGraph) Insert(node *Node) error {
	return nil
}

func (g FilterGraph) Remove(node *Node) error {
	return nil
}

func (g FilterGraph) Process(data []uint8) {

}

type Node struct {
	Source chan []byte
	Sink   chan []byte
	Next chan []byte
}

func (n Node) process() error {
	return nil
}