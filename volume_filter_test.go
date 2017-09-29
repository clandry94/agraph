package agraph

import (
	"testing"
)

func TestVolumeFilterCreation(t *testing.T) {
	_, err := NewNode(NopFilter, "test1")
	if err != nil {
		t.Error(err)
	}
}

func TestVolumeProcess(t *testing.T) {

}

func TestVolumeDo(t *testing.T) {

}
