package docker

import (
	"testing"
)

func TestFillNew(t *testing.T) {

	t.Log("===========================")

	reg := new(Registry)
	// fill registry with all available docker machines
	reg.fillNew()

	t.Log("===========================")

}
