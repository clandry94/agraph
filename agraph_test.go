package agraph

import (
	"testing"
	"github.com/clandry94/agraph/filter"
)

func TestNew(t *testing.T) {
	nop := NewFilter(filter.NopFilter, Fields{
		"nop" : "nop",
	})

	meta := MetaData {
		SampleRate: 0,
		NumChannels: 0,
	}

	graphDef := &GraphDef{ filters: []*FilterDef{nop}}
	_, err := New(graphDef, meta)

	if err != nil {
		t.Errorf(err.Error())
	}

}
