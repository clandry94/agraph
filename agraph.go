package agraph

import (
	"github.com/clandry94/agraph/filter"
)

type Graph struct {
	Meta     MetaData
	graphDef *GraphDef
	graph    *graph
}

func (g *Graph) Compile() error {
	meta := filter.MetaData{
		SampleRate:  g.Meta.SampleRate,
		NumChannels: g.Meta.NumChannels,
	}


	var prevFilter filter.Node

	for i, filterDef := range g.graphDef.filters {
		opts := filter.Options{
			VolumeMultiplier:    filterDef.fields["VolumeMultiplier"].(float32),
			Delay:               filterDef.fields["Delay"].(int),
			Decay:               filterDef.fields["Decay"].(float32),
			MovingAverageLength: filterDef.fields["MovingAverageLength"].(int),
			Angle:               filterDef.fields["Angle"].(float64),
		}

		filt, err := filter.NewNode(filterDef.filter,
			meta,
			"name not implemented",
			opts,
		)
		if err != nil {
			return err
		}

		if i == 0 {
			g.graph.sink = filt.Source()
			prevFilter = filt
			continue
		} else if i == len(g.graphDef.filters) {
			g.graph.sink = filt.Sink()
		} else if i != 0 && i != len(g.graphDef.filters) {
			filt.SetSource(prevFilter.Sink())
		}

		prevFilter = filt
	}

	return nil
}

type graph struct {
	source chan []uint16
	sink   chan []uint16
}

type Fields map[string]interface{}

// Meta is a type which has information about a file
type MetaData struct {
	SampleRate  uint32
	NumChannels uint16
}

type GraphDef struct {
	filters []*FilterDef
}

type FilterDef struct {
	filter filter.Type
	fields Fields
}

func NewFilter(filter filter.Type, fields Fields) *FilterDef {
	return &FilterDef{
		filter: filter,
		fields: fields,
	}
}

// New returns a new *Graph if file exists, otherwise returns an error
func NewGraph(graphDef *GraphDef, meta MetaData) (*Graph, error) {
	return &Graph{
		Meta:     meta,
		graphDef: graphDef,
		graph:    nil,
	}, nil
}
