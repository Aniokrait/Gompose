package widgets

import (
	"fmt"
	"github.com/tak/goui/core"
)

// Text はテキストウィジェットを作成します
func Text(key string, text string, props core.Props) *core.Node {
	// プロパティの設定
	mergedProps := core.Props{
		"text": text,
	}
	
	// 追加のプロパティをマージ
	for k, v := range props {
		mergedProps[k] = v
	}
	
	return core.NewNode(core.TextNodeType, key, mergedProps)
}

// Box はボックスウィジェットを作成します
func Box(key string, props core.Props, children ...*core.Node) *core.Node {
	node := core.NewNode(core.BoxNodeType, key, props)
	
	// 子ノードを追加
	for _, child := range children {
		node.AddChild(child)
	}
	
	return node
}

// Row は水平方向に子ノードを配置するウィジェットを作成します
func Row(key string, props core.Props, children ...*core.Node) *core.Node {
	node := core.NewNode(core.RowNodeType, key, props)
	
	// 子ノードを追加
	for _, child := range children {
		node.AddChild(child)
	}
	
	return node
}

// Column は垂直方向に子ノードを配置するウィジェットを作成します
func Column(key string, props core.Props, children ...*core.Node) *core.Node {
	node := core.NewNode(core.ColumnNodeType, key, props)
	
	// 子ノードを追加
	for _, child := range children {
		node.AddChild(child)
	}
	
	return node
}

// Container は単一の子ノードを持つコンテナウィジェットを作成します
func Container(key string, props core.Props, child *core.Node) *core.Node {
	node := core.NewNode(core.ContainerNodeType, key, props)
	
	// 子ノードを追加
	if child != nil {
		node.AddChild(child)
	}
	
	return node
}

// Button はボタンウィジェットを作成します
func Button(key string, label string, onClick func(), props core.Props) *core.Node {
	// プロパティの設定
	mergedProps := core.Props{
		"label":   label,
		"onClick": onClick,
	}
	
	// 追加のプロパティをマージ
	for k, v := range props {
		mergedProps[k] = v
	}
	
	// ボタンはカスタムコンポーネントとして実装
	buttonComponent := core.NewFunctionComponent(func(props core.Props) *core.Node {
		label := props.GetString("label", "Button")
		
		// ボタンの基本構造を作成
		buttonNode := Box(fmt.Sprintf("%s-box", key), core.Props{
			"padding":       8.0,
			"borderRadius":  4.0,
			"backgroundColor": "#2196F3",
		}, 
			Text(fmt.Sprintf("%s-text", key), label, core.Props{
				"color":     "#FFFFFF",
				"fontSize":  16.0,
				"textAlign": "center",
			}),
		)
		
		return buttonNode
	})
	
	// ボタンノードを作成
	node := core.NewNode(core.CustomNodeType, key, mergedProps)
	node.Component = buttonComponent
	
	return node
}

// Input はテキスト入力ウィジェットを作成します
func Input(key string, value string, onChange func(string), props core.Props) *core.Node {
	// プロパティの設定
	mergedProps := core.Props{
		"value":    value,
		"onChange": onChange,
	}
	
	// 追加のプロパティをマージ
	for k, v := range props {
		mergedProps[k] = v
	}
	
	// 入力フィールドはカスタムコンポーネントとして実装
	inputComponent := core.NewFunctionComponent(func(props core.Props) *core.Node {
		value := props.GetString("value", "")
		
		// 入力フィールドの基本構造を作成
		inputNode := Box(fmt.Sprintf("%s-box", key), core.Props{
			"padding":       8.0,
			"borderRadius":  4.0,
			"borderWidth":   1.0,
			"borderColor":   "#CCCCCC",
			"backgroundColor": "#FFFFFF",
		}, 
			Text(fmt.Sprintf("%s-text", key), value, core.Props{
				"color":     "#000000",
				"fontSize":  16.0,
			}),
		)
		
		return inputNode
	})
	
	// 入力ノードを作成
	node := core.NewNode(core.CustomNodeType, key, mergedProps)
	node.Component = inputComponent
	
	return node
}

// Image は画像ウィジェットを作成します
func Image(key string, source string, props core.Props) *core.Node {
	// プロパティの設定
	mergedProps := core.Props{
		"source": source,
	}
	
	// 追加のプロパティをマージ
	for k, v := range props {
		mergedProps[k] = v
	}
	
	// 画像はカスタムコンポーネントとして実装
	imageComponent := core.NewFunctionComponent(func(props core.Props) *core.Node {
		// 実際のレンダリングではここで画像を読み込む処理が必要
		return Box(fmt.Sprintf("%s-box", key), props)
	})
	
	// 画像ノードを作成
	node := core.NewNode(core.CustomNodeType, key, mergedProps)
	node.Component = imageComponent
	
	return node
}

// Spacer は指定されたサイズのスペースを作成します
func Spacer(key string, width, height float64) *core.Node {
	return Box(key, core.Props{
		"width":  width,
		"height": height,
	})
}

// Divider は区切り線を作成します
func Divider(key string, isHorizontal bool, props core.Props) *core.Node {
	// プロパティの設定
	mergedProps := core.Props{
		"backgroundColor": "#CCCCCC",
	}
	
	// 追加のプロパティをマージ
	for k, v := range props {
		mergedProps[k] = v
	}
	
	if isHorizontal {
		mergedProps["height"] = 1.0
	} else {
		mergedProps["width"] = 1.0
	}
	
	return Box(key, mergedProps)
}
