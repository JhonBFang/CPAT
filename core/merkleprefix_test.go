package core

import (
	"testing"
)

//*******************************
// 对于Merkle prefix tree的测试
//*******************************

// 目前只是测试了Merkle tree，Merkle prefix tre中的prefix 在monitor的时候生成根据最后的epoch和自己的epoch生成。
func TestAppend(t *testing.T) {
	m := NewMerklePT(4)
	m.Append(3, 3, 27)
	m.Append(3, 3, 27)

	l1 := m.getLeafNode(0)
	l2 := m.getLeafNode(1)

	if !l1.isComplete() || !l2.isComplete() || !l1.getParent().isComplete() {
		t.Error()
	}

	m.Append(3, 3, 27)
	// m.Append([]byte("4"))

	l3 := m.getLeafNode(2)
	l4 := m.getLeafNode(3)

	if !l3.isComplete() || l4.isComplete() || l3.getParent().isComplete() {
		t.Error()
	}

	if m.Size != 3 {
		t.Error()
	}

	//此处 0.008， Merkle2是1.022
	m1 := NewMerklePT(20)
	numAppends := 1
	for i := 0; i < numAppends; i++ {
		m1.Append(8192, 1, 8192) // k, depth, numbers
	}

	if m1.Size != 1 {
		t.Error()
	}

}

// 写consistency proof
func TestGenerateExtensionProof(t *testing.T) {

	m0 := createTestingTree(7, 3)
	m1 := createTestingTree(15, 4)
	m2 := createTestingTree(3500, 16)

	tables := []struct {
		ms            *MerklePT
		oldSize       uint32
		requestedSize uint32
	}{
		{m0, 1, 7},
		{m1, 1, 8},
		{m1, 7, 15},
		{m1, 1, 15},
		{m2, 1000, 2000},
		{m2, 150, 2100},
		{m2, 350, 2200},
		{m2, 1000, 3500},
	}

	for _, table := range tables {
		oldDigest := table.ms.GetOldDigest(table.oldSize)
		newDigest := table.ms.GetOldDigest(table.requestedSize)

		proof := table.ms.GenerateConsistencyProof(table.oldSize, table.requestedSize)

		if !VerifyExtensionProof(oldDigest, newDigest, proof) {
			t.Log(table.ms.depth)
			t.Error()
		}
	}

}

func createTestingTree(size uint32, depth uint32) *MerklePT {
	m := NewMerklePT(depth)

	var i uint32
	for i = 0; i < size; i++ {
		m.Append(3, 3, 27)
	}
	return m
}

func TestComputeContentHash(t *testing.T) {
	m := ComputeContentHash([]byte("1"), 0)
	if m == nil {
		t.Error()
	}
}
