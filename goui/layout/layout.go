package layout

import (
	"github.com/tak/goui/core"
)

// Size はサイズを表す構造体です
type Size struct {
	Width  float64
	Height float64
}

// Position は位置を表す構造体です
type Position struct {
	X float64
	Y float64
}

// Rect は矩形領域を表す構造体です
type Rect struct {
	Position Position
	Size     Size
}

// Constraints はレイアウト制約を表す構造体です
type Constraints struct {
	MinSize Size
	MaxSize Size
}

// NewConstraints は新しい制約を作成します
func NewConstraints(minWidth, minHeight, maxWidth, maxHeight float64) Constraints {
	return Constraints{
		MinSize: Size{Width: minWidth, Height: minHeight},
		MaxSize: Size{Width: maxWidth, Height: maxHeight},
	}
}

// FixedSize は固定サイズの制約を作成します
func FixedSize(width, height float64) Constraints {
	return Constraints{
		MinSize: Size{Width: width, Height: height},
		MaxSize: Size{Width: width, Height: height},
	}
}

// Tight は最小サイズと最大サイズが同じ制約を作成します
func Tight(size Size) Constraints {
	return Constraints{
		MinSize: size,
		MaxSize: size,
	}
}

// Loose は最小サイズが0で最大サイズが指定されたサイズの制約を作成します
func Loose(maxSize Size) Constraints {
	return Constraints{
		MinSize: Size{Width: 0, Height: 0},
		MaxSize: maxSize,
	}
}

// Constrain はサイズを制約内に収めます
func (c Constraints) Constrain(size Size) Size {
	width := size.Width
	if width < c.MinSize.Width {
		width = c.MinSize.Width
	}
	if width > c.MaxSize.Width {
		width = c.MaxSize.Width
	}

	height := size.Height
	if height < c.MinSize.Height {
		height = c.MinSize.Height
	}
	if height > c.MaxSize.Height {
		height = c.MaxSize.Height
	}

	return Size{Width: width, Height: height}
}

// LayoutManager はレイアウト計算を担当します
type LayoutManager struct {
	// レイアウト計算に必要な状態やキャッシュを保持
}

// NewLayoutManager は新しいレイアウトマネージャーを作成します
func NewLayoutManager() *LayoutManager {
	return &LayoutManager{}
}

// CalculateLayout はノードツリーのレイアウトを計算します
func (lm *LayoutManager) CalculateLayout(root *core.Node, constraints Constraints) map[string]Rect {
	// レイアウト計算の結果を格納するマップ
	layout := make(map[string]Rect)
	
	// 再帰的にレイアウトを計算
	lm.calculateNodeLayout(root, constraints, Position{X: 0, Y: 0}, layout)
	
	return layout
}

// calculateNodeLayout は単一ノードとその子のレイアウトを計算します
func (lm *LayoutManager) calculateNodeLayout(
	node *core.Node,
	constraints Constraints,
	position Position,
	layout map[string]Rect,
) Size {
	// ノードタイプに基づいてレイアウト計算を行う
	var size Size
	
	switch node.Type {
	case core.TextNodeType:
		// テキストノードのサイズ計算（実際にはフォントメトリクスが必要）
		text := node.Props.GetString("text", "")
		fontSize := node.Props.GetFloat("fontSize", 16.0)
		// 単純化のため、文字数とフォントサイズに基づく簡易計算
		size = Size{
			Width:  float64(len(text)) * fontSize * 0.6,
			Height: fontSize * 1.2,
		}
		
	case core.BoxNodeType:
		// ボックスノードのサイズ計算
		width := node.Props.GetFloat("width", 0)
		height := node.Props.GetFloat("height", 0)
		
		// 明示的なサイズが指定されている場合はそれを使用
		if width > 0 && height > 0 {
			size = Size{Width: width, Height: height}
		} else {
			// 子ノードに基づいてサイズを計算
			size = lm.calculateChildrenLayout(node, constraints, position, layout)
		}
		
		// パディングを追加
		padding := node.Props.GetFloat("padding", 0)
		size.Width += padding * 2
		size.Height += padding * 2
		
	case core.RowNodeType:
		// 行レイアウト（水平方向に子ノードを配置）
		size = lm.calculateRowLayout(node, constraints, position, layout)
		
	case core.ColumnNodeType:
		// 列レイアウト（垂直方向に子ノードを配置）
		size = lm.calculateColumnLayout(node, constraints, position, layout)
		
	case core.ContainerNodeType:
		// コンテナは単一の子ノードを持ち、そのサイズに合わせる
		if len(node.Children) > 0 {
			size = lm.calculateNodeLayout(node.Children[0], constraints, position, layout)
		} else {
			size = Size{Width: 0, Height: 0}
		}
		
	default:
		// カスタムノードやその他のノードタイプ
		// コンポーネントがある場合はそのレンダリング結果を使用
		if node.Component != nil {
			// コンポーネントのレンダリング結果に対してレイアウト計算
			renderedNode := node.Component.Render(node.Props)
			size = lm.calculateNodeLayout(renderedNode, constraints, position, layout)
		} else {
			// デフォルトでは子ノードに基づいてサイズを計算
			size = lm.calculateChildrenLayout(node, constraints, position, layout)
		}
	}
	
	// 制約に従ってサイズを調整
	size = constraints.Constrain(size)
	
	// レイアウト結果を記録
	layout[node.Key] = Rect{
		Position: position,
		Size:     size,
	}
	
	return size
}

// calculateChildrenLayout は子ノードのレイアウトを計算します（デフォルト実装）
func (lm *LayoutManager) calculateChildrenLayout(
	node *core.Node,
	constraints Constraints,
	position Position,
	layout map[string]Rect,
) Size {
	if len(node.Children) == 0 {
		return Size{Width: 0, Height: 0}
	}
	
	// 単純に最初の子ノードのサイズを使用（実際にはより複雑なロジックが必要）
	return lm.calculateNodeLayout(node.Children[0], constraints, position, layout)
}

// calculateRowLayout は行レイアウトを計算します
func (lm *LayoutManager) calculateRowLayout(
	node *core.Node,
	constraints Constraints,
	position Position,
	layout map[string]Rect,
) Size {
	if len(node.Children) == 0 {
		return Size{Width: 0, Height: 0}
	}
	
	spacing := node.Props.GetFloat("spacing", 0)
	x := position.X
	maxHeight := 0.0
	
	// 各子ノードを水平方向に配置
	for _, child := range node.Children {
		// 子ノードの制約を計算
		childConstraints := Constraints{
			MinSize: Size{Width: 0, Height: 0},
			MaxSize: Size{
				Width:  constraints.MaxSize.Width - (x - position.X),
				Height: constraints.MaxSize.Height,
			},
		}
		
		// 子ノードのレイアウトを計算
		childPos := Position{X: x, Y: position.Y}
		childSize := lm.calculateNodeLayout(child, childConstraints, childPos, layout)
		
		// 次の子ノードのX座標を更新
		x += childSize.Width + spacing
		
		// 最大の高さを記録
		if childSize.Height > maxHeight {
			maxHeight = childSize.Height
		}
	}
	
	// 行の合計サイズを返す
	return Size{
		Width:  x - position.X - spacing, // 最後のスペースを除く
		Height: maxHeight,
	}
}

// calculateColumnLayout は列レイアウトを計算します
func (lm *LayoutManager) calculateColumnLayout(
	node *core.Node,
	constraints Constraints,
	position Position,
	layout map[string]Rect,
) Size {
	if len(node.Children) == 0 {
		return Size{Width: 0, Height: 0}
	}
	
	spacing := node.Props.GetFloat("spacing", 0)
	y := position.Y
	maxWidth := 0.0
	
	// 各子ノードを垂直方向に配置
	for _, child := range node.Children {
		// 子ノードの制約を計算
		childConstraints := Constraints{
			MinSize: Size{Width: 0, Height: 0},
			MaxSize: Size{
				Width:  constraints.MaxSize.Width,
				Height: constraints.MaxSize.Height - (y - position.Y),
			},
		}
		
		// 子ノードのレイアウトを計算
		childPos := Position{X: position.X, Y: y}
		childSize := lm.calculateNodeLayout(child, childConstraints, childPos, layout)
		
		// 次の子ノードのY座標を更新
		y += childSize.Height + spacing
		
		// 最大の幅を記録
		if childSize.Width > maxWidth {
			maxWidth = childSize.Width
		}
	}
	
	// 列の合計サイズを返す
	return Size{
		Width:  maxWidth,
		Height: y - position.Y - spacing, // 最後のスペースを除く
	}
}
