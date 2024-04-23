package core

import (
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
)

type Node struct {
	Children []*Node // 子节点
	Hash     []byte  // 当前节点的哈希
}

type KaryTree struct {
	Root  *Node  // 树的根节点
	K     uint32 // 分叉因子
	Depth uint32 // 树的高度
}

// 创建新的K叉树
func NewKaryTree(k uint32, depth uint32) *KaryTree {
	if depth < 1 {
		panic("树的深度至少为1")
	}
	return &KaryTree{
		Root:  &Node{},
		K:     k,
		Depth: depth,
	}
}

// 为树添加叶子节点
func (t *KaryTree) AddLeaf(pos uint32) {
	posAsByte := make([]byte, 4)
	binary.LittleEndian.PutUint32(posAsByte, pos)
	leaf := &Node{Hash: posAsByte}
	if !t.addLeaf(t.Root, leaf, 1) {
		panic("无法添加更多叶子节点：树已满")
	}
}

// 递归地添加叶子节点
func (t *KaryTree) addLeaf(current *Node, leaf *Node, depth uint32) bool {
	// 如果达到树的最大深度，则添加叶子节点
	if depth == t.Depth {
		if len(current.Children) < int(t.K) {
			current.Children = append(current.Children, leaf)
			return true
		}
		return false
	}

	// 如果不是在最大深度，则需要在中间节点中添加
	for _, child := range current.Children {
		if t.addLeaf(child, leaf, depth+1) {
			return true
		}
	}

	// 如果所有子节点都已检查并且都满了，则添加一个新的中间节点（如果可能）
	if len(current.Children) < int(t.K) {
		newMiddle := &Node{}
		current.Children = append(current.Children, newMiddle)
		return t.addLeaf(newMiddle, leaf, depth+1)
	}

	return false
}

// 计算哈希值，叶子节点为传入值，中间节点为子节点哈希的组合
func (t *KaryTree) CalculateHashes(node *Node) []byte {
	if len(node.Children) == 0 {
		// 叶子节点已经有了哈希
		return []byte{}
	}
	// 中间节点的哈希是其所有子节点的哈希组合
	var hashes []byte
	for _, child := range node.Children {
		t.CalculateHashes(child)
		// hashes += child.Hash
		hashes = append(hashes, child.Hash...)
		println("中间节点hashes:", hashes)
	}
	// node.Hash = t.calculateHash(hashes)
	node.Hash = hashes
	println("最终的hashes:", node.Hash)
	return node.Hash
}

// 计算字符串的SHA-256哈希值
func (t *KaryTree) calculateHash(input string) string {
	hashBytes := sha256.Sum256([]byte(input))
	return hex.EncodeToString(hashBytes[:])
}
