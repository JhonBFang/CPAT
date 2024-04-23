package test

import (
	"testing"
)

// 测试用例
func TestInsert(t *testing.T) {
	tree := &Node{}

	//插入节点
	values := []int{4, 2, 7, 1, 3, 6, 9}
	for _, v := range values {
		tree.Insert(v)
	}

	//检查中序遍历是否正确
	// inOrderResult := tree.InOrderTraversal()
	// expectedInOrder := []int{1, 2, 3, 4, 6, 7, 9}
	// for i, v := range inOrderResult {
	// 	if v != expectedInOrder[i] {
	// 		t.Errorf("InOrderTranversal error, expected %v at index %d,got %v", expectedInOrder[i], i, v)
	// 	}
	// }

	//查找节点
	found := tree.Search(4)
	if !found {
		t.Errorf("Search error, expected to find 4")
	}
	found = tree.Search(5)
	if found {
		t.Errorf("Search error, did not expect to find 5")
	}
}
