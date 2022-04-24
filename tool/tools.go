// Package tool 工具类
package tool

// IsLowerAlpha 是否小写字母
func IsLowerAlpha(char byte) bool {
	return char >= byte('a') && char <= byte('z')
}
