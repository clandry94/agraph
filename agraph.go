package agraph

import (
	"github.com/clandry94/agraph/filter"
)

type Graph struct {
	Meta   MetaData
	graph  *graph
}

type graph struct {

}

type Fields map[string]interface{}

// Meta is a type which has information about a file
type MetaData struct {
	SampleRate uint32
	NumChannels uint16
}

type GraphDef struct {
	filters []*FilterDef
}

type FilterDef struct {
	filter filter.Type
	fields Fields
}

func NewFilter(filter filter.Type, fields Fields) *FilterDef{
	return &FilterDef {
		filter: filter,
		fields: fields,
	}
}

// New returns a new *Graph if file exists, otherwise returns an error
func New(graphDefinition *GraphDef, meta MetaData) (*Graph, error) {

	return &Graph{
		Meta:   meta,
		graph: nil,
	}, nil
}
