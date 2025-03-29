package core

import (
	"reflect"
	"sync"
)

// StateManager は状態管理を担当します
type StateManager struct {
	states     map[string]interface{}
	listeners  map[string][]StateChangeListener
	mutex      sync.RWMutex
}

// StateChangeListener は状態変更を監視するリスナーです
type StateChangeListener func(oldState, newState interface{})

// NewStateManager は新しい状態管理マネージャーを作成します
func NewStateManager() *StateManager {
	return &StateManager{
		states:    make(map[string]interface{}),
		listeners: make(map[string][]StateChangeListener),
	}
}

// GetState は指定されたキーの状態を取得します
func (sm *StateManager) GetState(key string) interface{} {
	sm.mutex.RLock()
	defer sm.mutex.RUnlock()
	
	return sm.states[key]
}

// SetState は状態を更新し、リスナーに通知します
func (sm *StateManager) SetState(key string, newState interface{}) {
	sm.mutex.Lock()
	oldState := sm.states[key]
	sm.states[key] = newState
	listeners := sm.listeners[key]
	sm.mutex.Unlock()
	
	// リスナーに通知
	for _, listener := range listeners {
		listener(oldState, newState)
	}
}

// AddListener は状態変更リスナーを追加します
func (sm *StateManager) AddListener(key string, listener StateChangeListener) {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()
	
	sm.listeners[key] = append(sm.listeners[key], listener)
}

// RemoveListener は状態変更リスナーを削除します
func (sm *StateManager) RemoveListener(key string, listener StateChangeListener) {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()
	
	listeners := sm.listeners[key]
	for i, l := range listeners {
		// 関数ポインタの比較は難しいので、この実装は単純化しています
		// 実際には、リスナーを識別するための別の方法が必要かもしれません
		if reflect.ValueOf(l).Pointer() == reflect.ValueOf(listener).Pointer() {
			sm.listeners[key] = append(listeners[:i], listeners[i+1:]...)
			break
		}
	}
}

// CreateState は新しい状態を作成し、更新関数を返します
func (sm *StateManager) CreateState(key string, initialValue interface{}) (
	getter func() interface{},
	setter func(interface{}),
) {
	sm.mutex.Lock()
	sm.states[key] = initialValue
	sm.mutex.Unlock()
	
	getter = func() interface{} {
		return sm.GetState(key)
	}
	
	setter = func(newValue interface{}) {
		sm.SetState(key, newValue)
	}
	
	return getter, setter
}
