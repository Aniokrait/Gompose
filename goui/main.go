package main

import (
	"fmt"
	"github.com/tak/goui/core"
	"github.com/tak/goui/layout"
	"github.com/tak/goui/render"
	"github.com/tak/goui/widgets"
)

func main() {
	// 状態管理マネージャーを作成
	stateManager := core.NewStateManager()
	
	// カウンター状態を作成
	counterGetter, counterSetter := stateManager.CreateState("counter", 0)
	
	// カウンターを増加させる関数
	incrementCounter := func() {
		currentValue := counterGetter().(int)
		counterSetter(currentValue + 1)
		fmt.Printf("カウンター: %d\n", currentValue + 1)
		
		// UIを再レンダリング
		renderUI(counterGetter().(int))
	}
	
	// カウンターを減少させる関数
	decrementCounter := func() {
		currentValue := counterGetter().(int)
		if currentValue > 0 {
			counterSetter(currentValue - 1)
			fmt.Printf("カウンター: %d\n", currentValue - 1)
			
			// UIを再レンダリング
			renderUI(counterGetter().(int))
		}
	}
	
	// 初期UIをレンダリング
	renderUI(counterGetter().(int))
	
	// インクリメントボタンをシミュレート
	fmt.Println("\n[インクリメントボタンがクリックされました]")
	incrementCounter()
	
	// インクリメントボタンをシミュレート
	fmt.Println("\n[インクリメントボタンがクリックされました]")
	incrementCounter()
	
	// デクリメントボタンをシミュレート
	fmt.Println("\n[デクリメントボタンがクリックされました]")
	decrementCounter()
}

// renderUI はUIツリーを構築してレンダリングします
func renderUI(counterValue int) {
	// ルートノードを作成
	root := buildUI(counterValue)
	
	// コンソールレンダリングターゲットを作成
	target := render.NewConsoleRenderTarget(80, 20)
	
	// レンダラーを作成
	renderer := render.NewRenderer(target)
	
	// UIをレンダリング
	constraints := layout.NewConstraints(0, 0, 80, 20)
	renderer.Render(root, constraints)
	
	// UIツリーを出力（デバッグ用）
	fmt.Println("\n--- UI Tree ---")
	root.PrintTree(0)
	fmt.Println("---------------")
}

// buildUI はUIツリーを構築します
func buildUI(counterValue int) *core.Node {
	// カウンターアプリのUI
	return widgets.Column("root", core.Props{
		"spacing": 1.0,
	},
		// タイトル
		widgets.Text("title", "カウンターアプリ", core.Props{
			"fontSize": 20.0,
			"textAlign": "center",
		}),
		
		// カウンター表示
		widgets.Box("counter-box", core.Props{
			"padding": 2.0,
			"borderRadius": 4.0,
			"backgroundColor": "#E0E0E0",
		},
			widgets.Text("counter-text", fmt.Sprintf("カウント: %d", counterValue), core.Props{
				"fontSize": 18.0,
				"textAlign": "center",
			}),
		),
		
		// ボタン行
		widgets.Row("buttons", core.Props{
			"spacing": 2.0,
		},
			// デクリメントボタン
			widgets.Button("decrement", "-", func() {}, core.Props{
				"backgroundColor": "#F44336",
			}),
			
			// インクリメントボタン
			widgets.Button("increment", "+", func() {}, core.Props{
				"backgroundColor": "#4CAF50",
			}),
		),
	)
}

// CounterApp はカウンターアプリのコンポーネントです
type CounterApp struct {
	counter int
}

// NewCounterApp は新しいカウンターアプリを作成します
func NewCounterApp() *CounterApp {
	return &CounterApp{
		counter: 0,
	}
}

// Render はカウンターアプリのUIを生成します
func (ca *CounterApp) Render(props core.Props) *core.Node {
	return buildUI(ca.counter)
}

// ShouldUpdate はコンポーネントが更新すべきかを判断します
func (ca *CounterApp) ShouldUpdate(oldProps, newProps core.Props) bool {
	return true
}

// GetState は現在の状態を返します
func (ca *CounterApp) GetState() interface{} {
	return ca.counter
}

// SetState は状態を更新します
func (ca *CounterApp) SetState(newState interface{}) {
	if newCounter, ok := newState.(int); ok {
		ca.counter = newCounter
	}
}

// Initialize はコンポーネントを初期化します
func (ca *CounterApp) Initialize(props core.Props) {
	// 初期化処理
}

// Cleanup はコンポーネントのクリーンアップを行います
func (ca *CounterApp) Cleanup() {
	// クリーンアップ処理
}
