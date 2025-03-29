package core

// Component はUIコンポーネントのインターフェースです
type Component interface {
	// Render はコンポーネントのUIツリーを生成します
	Render(props Props) *Node
	
	// ShouldUpdate はコンポーネントが再レンダリングすべきかを判断します
	ShouldUpdate(oldProps, newProps Props) bool
}

// ComponentFunc は関数型コンポーネントを定義するための型です
type ComponentFunc func(props Props) *Node

// FunctionComponent は関数型コンポーネントをラップする構造体です
type FunctionComponent struct {
	renderFunc ComponentFunc
}

// NewFunctionComponent は新しい関数型コンポーネントを作成します
func NewFunctionComponent(renderFunc ComponentFunc) *FunctionComponent {
	return &FunctionComponent{
		renderFunc: renderFunc,
	}
}

// Render は関数型コンポーネントのレンダリングを実行します
func (fc *FunctionComponent) Render(props Props) *Node {
	return fc.renderFunc(props)
}

// ShouldUpdate は常に true を返します（最適化の余地あり）
func (fc *FunctionComponent) ShouldUpdate(oldProps, newProps Props) bool {
	// デフォルトでは常に更新
	// 最適化: プロパティの比較に基づいて判断することも可能
	return !oldProps.Equal(newProps)
}

// StatefulComponent は状態を持つコンポーネントのインターフェースです
type StatefulComponent interface {
	Component
	
	// GetState はコンポーネントの現在の状態を返します
	GetState() interface{}
	
	// SetState はコンポーネントの状態を更新します
	SetState(newState interface{})
	
	// Initialize はコンポーネントを初期化します
	Initialize(props Props)
	
	// Cleanup はコンポーネントのクリーンアップを行います
	Cleanup()
}
