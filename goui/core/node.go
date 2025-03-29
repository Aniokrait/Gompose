package core

import (
	"fmt"
)

// NodeType は UIノードの種類を表す型です
type NodeType string

const (
	// 基本的なノードタイプ
	TextNodeType    NodeType = "Text"
	ContainerNodeType NodeType = "Container"
	RowNodeType     NodeType = "Row"
	ColumnNodeType  NodeType = "Column"
	BoxNodeType     NodeType = "Box"
	CustomNodeType  NodeType = "Custom"
)

// Node は UI ツリーの基本要素です
type Node struct {
	Type      NodeType
	Key       string
	Props     Props
	Children  []*Node
	Parent    *Node
	Component Component
}

// NewNode は新しいノードを作成します
func NewNode(nodeType NodeType, key string, props Props) *Node {
	return &Node{
		Type:     nodeType,
		Key:      key,
		Props:    props,
		Children: []*Node{},
	}
}

// AddChild は子ノードを追加します
func (n *Node) AddChild(child *Node) {
	child.Parent = n
	n.Children = append(n.Children, child)
}

// RemoveChild は子ノードを削除します
func (n *Node) RemoveChild(key string) bool {
	for i, child := range n.Children {
		if child.Key == key {
			// 子ノードの親参照をクリア
			n.Children[i].Parent = nil
			// 子ノードを削除
			n.Children = append(n.Children[:i], n.Children[i+1:]...)
			return true
		}
	}
	return false
}

// FindChild は指定されたキーを持つ子ノードを検索します
func (n *Node) FindChild(key string) *Node {
	for _, child := range n.Children {
		if child.Key == key {
			return child
		}
	}
	return nil
}

// PrintTree はノードツリーを再帰的に出力します（デバッグ用）
func (n *Node) PrintTree(indent int) {
	indentation := ""
	for i := 0; i < indent; i++ {
		indentation += "  "
	}

	fmt.Printf("%s- %s (Type: %s, Props: %v)\n", indentation, n.Key, n.Type, n.Props)
	
	for _, child := range n.Children {
		child.PrintTree(indent + 1)
	}
}

// Clone はノードの深いコピーを作成します
func (n *Node) Clone() *Node {
	clone := &Node{
		Type:      n.Type,
		Key:       n.Key,
		Props:     n.Props.Clone(),
		Children:  make([]*Node, 0, len(n.Children)),
		Component: n.Component,
	}

	for _, child := range n.Children {
		childClone := child.Clone()
		childClone.Parent = clone
		clone.Children = append(clone.Children, childClone)
	}

	return clone
}
