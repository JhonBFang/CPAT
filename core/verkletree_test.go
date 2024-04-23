package core

import (
	"testing"
)

func TestAddLeaf(t *testing.T) {
	v := NewKaryTree(3, 3)
	numTotal := 27
	for i := 0; i < numTotal; i++ {
		v.AddLeaf(uint32(i))
	}
	v.CalculateHashes(v.Root)
}
