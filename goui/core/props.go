package core

import (
	"fmt"
	"reflect"
)

// Props はコンポーネントのプロパティを表します
type Props map[string]interface{}

// Get はプロパティの値を取得します
func (p Props) Get(key string) (interface{}, bool) {
	value, exists := p[key]
	return value, exists
}

// GetString は文字列プロパティを取得します
func (p Props) GetString(key string, defaultValue string) string {
	if value, exists := p[key]; exists {
		if strValue, ok := value.(string); ok {
			return strValue
		}
	}
	return defaultValue
}

// GetInt は整数プロパティを取得します
func (p Props) GetInt(key string, defaultValue int) int {
	if value, exists := p[key]; exists {
		if intValue, ok := value.(int); ok {
			return intValue
		}
	}
	return defaultValue
}

// GetFloat は浮動小数点プロパティを取得します
func (p Props) GetFloat(key string, defaultValue float64) float64 {
	if value, exists := p[key]; exists {
		if floatValue, ok := value.(float64); ok {
			return floatValue
		}
		// int から float への変換を試みる
		if intValue, ok := value.(int); ok {
			return float64(intValue)
		}
	}
	return defaultValue
}

// GetBool は真偽値プロパティを取得します
func (p Props) GetBool(key string, defaultValue bool) bool {
	if value, exists := p[key]; exists {
		if boolValue, ok := value.(bool); ok {
			return boolValue
		}
	}
	return defaultValue
}

// Set はプロパティの値を設定します
func (p Props) Set(key string, value interface{}) {
	p[key] = value
}

// Merge は別のプロパティマップとマージします
func (p Props) Merge(other Props) Props {
	result := Props{}
	
	// 現在のプロパティをコピー
	for k, v := range p {
		result[k] = v
	}
	
	// 新しいプロパティを追加または上書き
	for k, v := range other {
		result[k] = v
	}
	
	return result
}

// Clone はプロパティの深いコピーを作成します
func (p Props) Clone() Props {
	clone := Props{}
	for k, v := range p {
		// 値が map[string]interface{} の場合は再帰的にクローン
		if nestedProps, ok := v.(Props); ok {
			clone[k] = nestedProps.Clone()
		} else if nestedMap, ok := v.(map[string]interface{}); ok {
			// 通常の map をコピー
			clonedMap := make(map[string]interface{})
			for mk, mv := range nestedMap {
				clonedMap[mk] = mv
			}
			clone[k] = clonedMap
		} else {
			// その他の値はそのままコピー
			clone[k] = v
		}
	}
	return clone
}

// String はプロパティの文字列表現を返します
func (p Props) String() string {
	return fmt.Sprintf("%v", map[string]interface{}(p))
}

// Equal は2つのプロパティが等しいかどうかを判定します
func (p Props) Equal(other Props) bool {
	if len(p) != len(other) {
		return false
	}
	
	for k, v1 := range p {
		v2, exists := other[k]
		if !exists {
			return false
		}
		
		// 値の比較
		if !reflect.DeepEqual(v1, v2) {
			return false
		}
	}
	
	return true
}
