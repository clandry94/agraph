package agraph

import "fmt"

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
type Node interface{
	process()
}

type NopFilter struct {
	Source chan []byte
	Sink   chan []byte
	Next chan []byte
}

func (n NopFilter) process() error {
	for {
		select {
			case data := <-n.Source:
				fmt.Println("found data")
				var filteredData, err = n.do(&data)

				if err != nil {
					panic("Could not filter!")
				}
				n.Sink <- filteredData
		}
	}
	return nil
}

func (n NopFilter) do(data *[]byte) ([]byte, error) {
	return []byte("Hello, world!"), nil
}