package core

import (
	"encoding/binary"
	"fmt"

	crypto "github.com/ucbrise/MerkleSquare/lib/crypto"
)

// 先规定为8个字节，用于计算proof的大小。这只是hash的，后面可能会变。
const pointerSizeInBytes = 8

// 基础node struct
type node struct {
	hash      []byte
	acc       []byte //Accumulator
	parent    MerkleNode
	isRight   bool
	completed bool
	index     index
	// Verkle tree的root
}

// 中间节点
type InternalNode struct {
	node
	leftChild  MerkleNode
	rightChild MerkleNode
}

// 叶子节点
type LeafNode struct {
	node
	contentHash []byte //epoch 和acc的hash

	acc []byte // accumulator

	// accVerkle KaryTree
}

type index struct {
	depth uint32 //树的深度
	shift uint32 //树的偏移量也就是从左到右第几个节点
}

// MerkleNode interface for leaf/internal nodes
type MerkleNode interface {
	isLeafNode() bool
	setParent(MerkleNode)
	getHash() []byte
	getAcc() []byte //得到accumulator
	complete()      //prefix完成才能生成merkle,我这里不需要
	isComplete() bool
	isRightChild() bool
	getParent() MerkleNode
	createRightChild() MerkleNode
	createLeftChild() MerkleNode
	getRightChild() MerkleNode
	getLeftChild() MerkleNode
	getDepth() uint32
	print()
	// getPrefixTree() *prefixTree
	getShift() uint32
	getIndex() index
	getSibling() Sibling
	getContentHash() []byte
	// getPrefix() []byte
	// serialize() ([]byte, error) //在后面引用
	getSize() int
}

// 创建叶子节点, 这里的acc先使用数字代替，后面补上
func (node *LeafNode) completeLeaf(acc []byte, epo uint32) {

	contentHash := ComputeContentHash(acc, epo)
	// 添加verkle tree
	node.contentHash = contentHash
	node.hash = contentHash
	node.acc = acc //叶子节点的accumulator
	node.completed = true
}

// 创建叶子节点,还没有添加 accumulator
func createLeafNode(parent MerkleNode, isRight bool, shift uint32) *LeafNode {
	return &LeafNode{
		node: node{
			parent:  parent,
			isRight: isRight,
			// completed: false,
			index: index{
				depth: 0,
				shift: shift,
			},
		},
	}
}

// 创建中间节点 还没有添加accumulator
func createInternalNode(parent MerkleNode, depth uint32, isRight bool, shift uint32) *InternalNode {
	return &InternalNode{
		node: node{
			parent:  parent,
			isRight: isRight,
			// completed: false,
			index: index{
				depth: depth,
				shift: shift,
			},
		},
	}
}

// 创建root节点
func createRootNode(depth uint32) MerkleNode {
	return &InternalNode{
		node: node{
			completed: false,
			index: index{
				depth: depth,
				shift: 0,
			},
		},
	}
}

func (node *InternalNode) complete() {
	// hashVal := crypto.Hash(node.leftChild.getHash(), node.rightChild.getHash(), []byte("1"))
	hashVal := crypto.Hash(node.leftChild.getHash(), node.rightChild.getHash())
	node.hash = hashVal
	// node.acc = []byte("1")
	node.completed = true
}

func (node *InternalNode) createRightChild() MerkleNode {

	var newNode MerkleNode
	if node.getDepth() == 1 {
		newNode = createLeafNode(node, true, node.getShift()*2+1)
	} else {
		newNode = createInternalNode(node, node.getDepth()-1, true, node.getShift()*2+1)
	}
	node.rightChild = newNode
	return newNode
}

func (node *InternalNode) createLeftChild() MerkleNode {
	var newNode MerkleNode
	if node.getDepth() == 1 {
		newNode = createLeafNode(node, false, node.getShift()*2)
	} else {
		newNode = createInternalNode(node, node.getDepth()-1, false, node.getShift()*2)
	}

	node.leftChild = newNode

	return newNode
}

// 中间节点的兄弟节点
func (node *InternalNode) getSibling() Sibling {

	var hash []byte
	// var acc []byte
	if node.isRightChild() {
		hash = node.getParent().getLeftChild().getHash()
		// acc = node.getParent().getLeftChild().getAcc()
	} else {
		hash = node.getParent().getRightChild().getHash()
		// acc = node.getParent().getRightChild().getAcc()
	}

	return Sibling{
		Hash: hash,
		//还需要添加accumulator
		// Acc: acc,
	}
}

// 叶子节点的兄弟节点
func (node *LeafNode) getSibling() Sibling {

	var hash []byte
	// var acc []byte
	if node.isRightChild() {
		hash = node.getParent().getLeftChild().getHash()
		// acc = node.getParent().getLeftChild().getAcc()
	} else {
		hash = node.getParent().getRightChild().getHash()
		// acc = node.getParent().getRightChild().getAcc()
	}

	return Sibling{
		Hash: hash,
		// Acc:  acc,
	}
}

// 叶子节点的hash长度
func (node *LeafNode) getSize() int {

	// pointer to parent 8 bytes
	total := pointerSizeInBytes

	// hash, isRight, isComplete sizes
	total += binary.Size(node.hash) + binary.Size(node.isRight) + binary.Size(node.isComplete)

	// size of index
	total += binary.Size(node.index.depth) + binary.Size(node.index.shift)

	return total
}

// 中间节点的hash长度
func (node *InternalNode) getSize() int {
	// pointer to parent 8 bytes + left and right child pointers + prefix tree pointer
	total := pointerSizeInBytes * 4

	// hash, isRight, isComplete sizes
	total += binary.Size(node.hash) + binary.Size(node.isRight) + binary.Size(node.isComplete)

	// size of index
	total += binary.Size(node.index.depth) + binary.Size(node.index.shift)

	// right child
	if node.getRightChild() != nil {
		total += node.getRightChild().getSize()
	}

	if node.getLeftChild() != nil {
		total += node.getLeftChild().getSize()
	}

	return total
}

func (node *InternalNode) isComplete() bool            { return node.completed }
func (node *InternalNode) isRightChild() bool          { return node.isRight }
func (node *InternalNode) getParent() MerkleNode       { return node.parent }
func (node *InternalNode) isLeafNode() bool            { return false }
func (node *InternalNode) setParent(parent MerkleNode) { node.parent = parent }
func (node *InternalNode) getHash() []byte             { return node.hash }
func (node *InternalNode) getAcc() []byte              { return node.acc }
func (node *InternalNode) getRightChild() MerkleNode   { return node.rightChild }
func (node *InternalNode) getLeftChild() MerkleNode    { return node.leftChild }
func (node *InternalNode) getDepth() uint32            { return node.index.depth }
func (node *InternalNode) print()                      { fmt.Print(node.isComplete()) }

// func (node *InternalNode) getPrefixTree() *prefixTree  { return node.prefixTree }
func (node *InternalNode) getShift() uint32       { return node.index.shift }
func (node *InternalNode) getIndex() index        { return node.index }
func (node *InternalNode) getContentHash() []byte { return []byte("") }

// func (node *InternalNode) getPrefix() []byte      { return []byte("") }

func (node *LeafNode) isComplete() bool             { return node.completed }
func (node *LeafNode) isRightChild() bool           { return node.isRight }
func (node *LeafNode) getParent() MerkleNode        { return node.parent }
func (node *LeafNode) isLeafNode() bool             { return true }
func (node *LeafNode) setParent(parent MerkleNode)  { node.parent = parent }
func (node *LeafNode) getHash() []byte              { return node.hash }
func (node *LeafNode) getAcc() []byte               { return node.acc }
func (node *LeafNode) complete()                    {}
func (node *LeafNode) createLeftChild() MerkleNode  { return &LeafNode{} }
func (node *LeafNode) createRightChild() MerkleNode { return &LeafNode{} }
func (node *LeafNode) getRightChild() MerkleNode    { return &LeafNode{} }
func (node *LeafNode) getLeftChild() MerkleNode     { return &LeafNode{} }
func (node *LeafNode) getDepth() uint32             { return 0 }
func (node *LeafNode) print()                       { fmt.Print(node.isComplete()) }

// func (node *LeafNode) getPrefixTree() *prefixTree   { return NewPrefixTree() }
func (node *LeafNode) getShift() uint32       { return node.index.shift }
func (node *LeafNode) getIndex() index        { return node.index }
func (node *LeafNode) getContentHash() []byte { return node.contentHash }

// func (node *LeafNode) getPrefix() []byte      { return makePrefixFromKey(node.key) }
