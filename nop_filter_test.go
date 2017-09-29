package agraph

import (
	"testing"
)

func TestNopFilterCreation(t *testing.T) {
	_, err := NewNode(NopFilter, "test1")
	if err != nil {
		t.Error(err)
	}
}

func TestNopProcess(t *testing.T) {

}

func TestNopDo(t *testing.T) {

}
