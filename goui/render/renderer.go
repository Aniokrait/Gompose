package render

import (
	"github.com/tak/goui/core"
	"github.com/tak/goui/layout"
)

// RenderTarget はレンダリング先のインターフェースです
type RenderTarget interface {
	// Clear はレンダリング領域をクリアします
	Clear()
	
	// DrawRect は矩形を描画します
	DrawRect(rect layout.Rect, props core.Props)
	
	// DrawText はテキストを描画します
	DrawText(text string, rect layout.Rect, props core.Props)
	
	// Flush はレンダリング結果を出力します
	Flush()
}

// Renderer はUIツリーのレンダリングを担当します
type Renderer struct {
	target RenderTarget
	layoutManager *layout.LayoutManager
}

// NewRenderer は新しいレンダラーを作成します
func NewRenderer(target RenderTarget) *Renderer {
	return &Renderer{
		target: target,
		layoutManager: layout.NewLayoutManager(),
	}
}

// Render はUIツリーをレンダリングします
func (r *Renderer) Render(root *core.Node, constraints layout.Constraints) {
	// レイアウト計算
	layoutResult := r.layoutManager.CalculateLayout(root, constraints)
	
	// レンダリング領域をクリア
	r.target.Clear()
	
	// ノードツリーを再帰的にレンダリング
	r.renderNode(root, layoutResult)
	
	// レンダリング結果を出力
	r.target.Flush()
}

// renderNode は単一ノードとその子をレンダリングします
func (r *Renderer) renderNode(node *core.Node, layout map[string]layout.Rect) {
	// ノードのレイアウト情報を取得
	rect, exists := layout[node.Key]
	if !exists {
		return
	}
	
	// ノードタイプに基づいてレンダリング
	switch node.Type {
	case core.TextNodeType:
		text := node.Props.GetString("text", "")
		r.target.DrawText(text, rect, node.Props)
		
	case core.BoxNodeType:
		// ボックスの背景を描画
		r.target.DrawRect(rect, node.Props)
		
		// 子ノードをレンダリング
		for _, child := range node.Children {
			r.renderNode(child, layout)
		}
		
	case core.RowNodeType, core.ColumnNodeType, core.ContainerNodeType:
		// コンテナタイプのノードは自身は描画せず、子ノードのみレンダリング
		for _, child := range node.Children {
			r.renderNode(child, layout)
		}
		
	case core.CustomNodeType:
		// カスタムノードの場合、コンポーネントがあればそのレンダリング結果を使用
		if node.Component != nil {
			renderedNode := node.Component.Render(node.Props)
			r.renderNode(renderedNode, layout)
		} else {
			// コンポーネントがない場合は通常のコンテナとして扱う
			for _, child := range node.Children {
				r.renderNode(child, layout)
			}
		}
	}
}

// ConsoleRenderTarget はコンソールへのレンダリングを行うターゲットです
// （デバッグ用の簡易実装）
type ConsoleRenderTarget struct {
	width  int
	height int
	buffer [][]rune
}

// NewConsoleRenderTarget は新しいコンソールレンダリングターゲットを作成します
func NewConsoleRenderTarget(width, height int) *ConsoleRenderTarget {
	buffer := make([][]rune, height)
	for i := range buffer {
		buffer[i] = make([]rune, width)
		for j := range buffer[i] {
			buffer[i][j] = ' '
		}
	}
	
	return &ConsoleRenderTarget{
		width:  width,
		height: height,
		buffer: buffer,
	}
}

// Clear はレンダリング領域をクリアします
func (c *ConsoleRenderTarget) Clear() {
	for i := range c.buffer {
		for j := range c.buffer[i] {
			c.buffer[i][j] = ' '
		}
	}
}

// DrawRect は矩形を描画します
func (c *ConsoleRenderTarget) DrawRect(rect layout.Rect, props core.Props) {
	// 矩形の境界を描画
	x1 := int(rect.Position.X)
	y1 := int(rect.Position.Y)
	x2 := int(rect.Position.X + rect.Size.Width - 1)
	y2 := int(rect.Position.Y + rect.Size.Height - 1)
	
	// 範囲チェック
	if x1 < 0 {
		x1 = 0
	}
	if y1 < 0 {
		y1 = 0
	}
	if x2 >= c.width {
		x2 = c.width - 1
	}
	if y2 >= c.height {
		y2 = c.height - 1
	}
	
	// 上下の境界線
	for x := x1; x <= x2; x++ {
		if y1 >= 0 && y1 < c.height {
			c.buffer[y1][x] = '-'
		}
		if y2 >= 0 && y2 < c.height {
			c.buffer[y2][x] = '-'
		}
	}
	
	// 左右の境界線
	for y := y1; y <= y2; y++ {
		if x1 >= 0 && x1 < c.width {
			c.buffer[y][x1] = '|'
		}
		if x2 >= 0 && x2 < c.width {
			c.buffer[y][x2] = '|'
		}
	}
	
	// 角
	if x1 >= 0 && y1 >= 0 && x1 < c.width && y1 < c.height {
		c.buffer[y1][x1] = '+'
	}
	if x2 >= 0 && y1 >= 0 && x2 < c.width && y1 < c.height {
		c.buffer[y1][x2] = '+'
	}
	if x1 >= 0 && y2 >= 0 && x1 < c.width && y2 < c.height {
		c.buffer[y2][x1] = '+'
	}
	if x2 >= 0 && y2 >= 0 && x2 < c.width && y2 < c.height {
		c.buffer[y2][x2] = '+'
	}
}

// DrawText はテキストを描画します
func (c *ConsoleRenderTarget) DrawText(text string, rect layout.Rect, props core.Props) {
	x := int(rect.Position.X)
	y := int(rect.Position.Y)
	
	// 範囲チェック
	if y < 0 || y >= c.height {
		return
	}
	
	// テキストを描画
	for i, char := range text {
		posX := x + i
		if posX < 0 || posX >= c.width {
			continue
		}
		c.buffer[y][posX] = char
	}
}

// Flush はレンダリング結果を出力します
func (c *ConsoleRenderTarget) Flush() {
	for _, line := range c.buffer {
		println(string(line))
	}
}
