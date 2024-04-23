package core

import (
	"bytes"
	"encoding/binary"
	"math/bits"

	crypto "github.com/ucbrise/MerkleSquare/lib/crypto"
)

// Merkle prefix tree
type MerklePT struct {
	Roots   []MerkleNode
	root    MerkleNode
	next    MerkleNode
	Size    uint32
	depth   uint32
	accroot []byte //pre-compute中历史root的acc
}

// MerkleConsistency proof contains an existence proof and subset proof 对于一个特定的leafnode
type MerkleConsistencyProof struct {
	Siblings []Sibling //只需要使用兄弟节点和subset proof

	// Acc []byte //加上prefix上的中间节点的accumulator
}

type Sibling struct {
	Hash []byte
	// Accumulator
	Acc []byte
}

// digest 对于当前Merkle prefix tree的状态
type Digest struct {
	Roots [][]byte
	Acc   []byte //root 中的accumulator
	Size  uint32
}

// 叶子节点的hash
type LeafHash struct {
	NodeContentHash []byte
	Acc             []byte
}

// 添加元素到Merkle prefix tree，hash(epo和acc),acc
func (m *MerklePT) Append(k uint32, depth uint32, numverkle uint32) {
	if m.isFull() {
		return
	}
	node := m.next.(*LeafNode)
	tree := NewKaryTree(k, depth)
	for i := 0; i < int(numverkle); i++ {
		tree.AddLeaf(uint32(i))
	}
	nodeAcc := tree.CalculateHashes(tree.Root)

	node.completeLeaf(nodeAcc, m.Size)
	m.Size++
	p := m.next

	//如果节点是右节点，那么合并，合并的时候要取出来一个旧root，将新的root添加进去
	for p.isRightChild() {
		p = p.getParent()
		p.complete()
		m.pop()
	}
	m.addRoot(p)          //左节点作为新的root添加到森林中
	m.addAccs(p.getAcc()) //pre-compute部分先不写，最后加1试试

	//查看tree是否是满的
	if m.isFull() {
		return
	}
	_ = p.getParent().createRightChild()
	p = p.getParent().getRightChild()
	for p.getDepth() > 0 {
		_ = p.createLeftChild()
		p = p.getLeftChild()
	}
	m.next = p
}

func (m *MerklePT) getLeafNode(epo uint32) MerkleNode {
	node := m.root

	for node.getDepth() > 0 {
		shift := node.getDepth() - 1
		if epo&(1<<shift)>>shift == 1 {
			node = node.getRightChild()
		} else {
			node = node.getLeftChild()
		}
	}
	return node
}

//存在证明应该在verkle tree中给出来

// GenerateExistenceProof 为给定的key/高度对生成存在证明

// 给一个digest生成consistency proof
func (m *MerklePT) GenerateConsistencyProof(oldSize uint32, requestedSize uint32) *MerkleConsistencyProof {

	roots := m.getOldRoots(requestedSize)
	oldDigestRoots := m.getOldRoots(oldSize)
	res := &MerkleConsistencyProof{}

	for i, root := range oldDigestRoots {
		if !bytes.Equal(root.getHash(), roots[i].getHash()) {
			lastNode := oldDigestRoots[len(oldDigestRoots)-1]

			generateConsistencyProof(lastNode, res, roots[i].getDepth())
			break
		}
	}
	return res
}

func generateConsistencyProof(node MerkleNode, proof *MerkleConsistencyProof, depth uint32) {
	siblings := []Sibling{}

	for node.getDepth() != depth { // if we want size param: remove isComplete() and pass in correct depth

		if !node.isRightChild() {
			sibling := node.getSibling()
			siblings = append(siblings, sibling)
		}

		node = node.getParent()
	}

	proof.Siblings = siblings
	//没有中间节点的acc 或者前一个版本的acc,将前一个版本的acc全部换位1,然后使用对比
	// proof.Acc = []byte("1")
}

func (m *MerklePT) getOldRoots(oldSize uint32) []MerkleNode {
	Roots := []MerkleNode{}
	var totalKeys uint32 = 0
	var mask uint32 = 1 << m.depth
	for mask > 0 {

		if bits.OnesCount32(mask&oldSize) == 1 {

			depth := bits.TrailingZeros32(mask)
			shift := totalKeys >> bits.TrailingZeros32(mask)

			Roots = append(Roots, m.getNode(uint32(depth), shift))

			totalKeys += mask
		}
		mask = mask >> 1
	}

	return Roots

}

func (m *MerklePT) getNode(depth uint32, shift uint32) MerkleNode {
	index := index{
		depth: depth,
		shift: shift,
	}

	return m.getNodeFromIndex(index)
}

// 从索引获取MerkleNode
func (m *MerklePT) getNodeFromIndex(index index) MerkleNode {

	node := m.root

	for node.getDepth() != index.depth {

		if isRightOf(node.getIndex(), index) {
			node = node.getRightChild()
		} else {
			node = node.getLeftChild()
		}

	}

	return node
}

// 如果index2在index1的右侧, 则返回true
func isRightOf(index1 index, index2 index) bool {

	heightDiff := index1.depth - index2.depth

	relativeWidth := uint32(1 << heightDiff)
	startIndex := relativeWidth * index1.shift

	return index2.shift >= startIndex+(relativeWidth>>1)
}

func (m *MerklePT) isFull() bool {
	return m.Size == 1<<m.depth
}

// 从forests中取出一个元素并返回
func (m *MerklePT) pop() MerkleNode {
	numTrees := len(m.Roots)
	node := m.Roots[numTrees-1]
	m.Roots = m.Roots[:numTrees-1]

	return node
}

// 向森林中添加一个元素
func (m *MerklePT) addRoot(node MerkleNode) {
	m.Roots = append(m.Roots, node)
}

// 向森林中添加acc, pre-compute
func (m *MerklePT) addAccs(acc []byte) {
	m.accroot = append(m.accroot, acc...)
}

// 将uint32转换为[]byte用来计算Hash
func ComputeContentHash(acc []byte, pos uint32) []byte {
	posAsByte := make([]byte, 4)
	binary.LittleEndian.PutUint32(posAsByte, pos) //使用小端序序列化，处理的更快

	contentHash := crypto.Hash(acc, posAsByte)

	return contentHash
}

// NewMerklePT是构造MerklePT对象的工厂方法
func NewMerklePT(depth uint32) *MerklePT {
	m := &MerklePT{
		Roots: []MerkleNode{},
		Size:  0,
		depth: depth,
	}

	next := createRootNode(depth)
	m.root = next

	for next.getDepth() > 0 {
		_ = next.createLeftChild()
		next = next.getLeftChild()
	}
	m.next = next
	return m
}

// VerifyExtensionProof verifies an ExtensionProof
func VerifyExtensionProof(oldDigest *Digest, newDigest *Digest, proof *MerkleConsistencyProof) bool {

	for i, oldRoot := range oldDigest.Roots {
		if bytes.Equal(oldRoot, newDigest.Roots[i]) {
			continue
		}

		p := len(oldDigest.Roots) - 2
		hash := oldDigest.Roots[p+1]
		// acc := oldDigest.Acc
		//acc

		lastRootDepth := GetOldDepth(oldDigest.Size-1, oldDigest.Size)
		newRootDepth := GetOldDepth(oldDigest.Size-1, newDigest.Size)
		shift := oldDigest.Size - 1
		siblingIndex := 0

		for j := 0; uint32(j) < newRootDepth; j++ {
			if uint32(j) >= lastRootDepth && isRight(shift) {
				hash = crypto.Hash(oldDigest.Roots[p], hash)
				// acc = []byte("1")
				p = p - 1
			} else if uint32(j) >= lastRootDepth {
				hash = crypto.Hash(hash, proof.Siblings[siblingIndex].Hash)
				// acc = []byte("1")
				siblingIndex++
			}

			shift = shift / 2
		}

		return bytes.Equal(hash, newDigest.Roots[i])
	}

	return true
}

// GetOldDepth given a position and size for an old forest, returns the depth of the tree pos belongs to
func GetOldDepth(pos uint32, size uint32) uint32 {

	index := getRootIndex(pos, size)
	leadingZeros := bits.LeadingZeros32(size)
	mask := 1 << (31 - leadingZeros)

	for index > 0 {

		index = index - 1
		mask = mask >> 1

		for bits.OnesCount(uint(mask&int(size))) == 0 {
			mask = mask >> 1
		}
	}

	return uint32(bits.TrailingZeros(uint(mask)))
}

// GetOldDigest returns a digest of the a MerkleSquare instance
// when it only contained oldSize keys.
func (m *MerklePT) GetOldDigest(oldSize uint32) *Digest {
	Roots := [][]byte{}

	for _, root := range m.getOldRoots(oldSize) {
		Roots = append(Roots, root.getHash())
	}

	return &Digest{
		Roots: Roots, //全是hash
		Size:  oldSize,
		//加上acc?
		Acc: []byte("1"),
	}
}

// Returns the root index that pos belongs to given the forest Size
func getRootIndex(pos uint32, Size uint32) int {

	xor := pos ^ Size
	leadingZeros := bits.LeadingZeros32(xor)

	forestSize := int(Size) >> (32 - leadingZeros)

	return bits.OnesCount32(uint32(forestSize))
}

func isRight(shift uint32) bool {
	return shift%2 != 0
}
