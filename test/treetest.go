package test

type Node struct {
	Value int
	Left  *Node
	Right *Node
}

// 插入节点
func (node *Node) Insert(value int) *Node {
	if node == nil {
		return &Node{Value: value}
	}
	if value < node.Value {
		node.Left = node.Left.Insert(value)
	} else {
		node.Right = node.Right.Insert(value)
	}
	return node
}

// 查找节点
func (node *Node) Search(value int) bool {
	if node == nil {
		return false
	}
	if node.Value == value {
		return true
	}
	if value < node.Value {
		return node.Left.Search(value)
	}
	return node.Right.Search(value)
}

// func (node *Node) InOrderTraversal() []int {
// 	var result []int
// 	if node == nil {
// 		return result
// 	}
// 	result = append(node.Left.InOrderTraversal(), node.Value)
// 	result = append(result, node.Right.InOrderTraversal()...)
// 	return result
// }
