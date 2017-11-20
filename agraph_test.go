package agraph

import (
	"os"
	"testing"
)

func TestNew(t *testing.T) {
	t.Run("Should return a new graph", func(t *testing.T) {
		_, err := New("README.md")
		if err != nil {
			t.Errorf("should pass without error, got: %v", err)
		}
	})

	t.Run("Should return an IsNotExist error", func(t *testing.T) {
		g, err := New("File which does not exist")
		if err == nil {
			t.Error("should return an error")
		}
		if !os.IsNotExist(err) {
			t.Errorf("should return isNotExist error, got: %v", err)
		}

		if g != nil {
			t.Errorf("returned graph should be nil, got: %v", g)
		}
	})
}
