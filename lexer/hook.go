package lexer

import (
	"regexp"
	"strings"

	"github.com/elevenzqx/flowchart-generator/tool"
)

// parenHook 括号类代码解析
func parenHook(t *TokenParser) (*Token, bool) {
	var char = t.Input[t.Index]
	// 检查一下是不是一个左圆括号
	if char == '(' {
		return &Token{Type: IgnoreType, Value: string([]byte{char})}, true
	}
	if char == ')' {
		return &Token{Type: IgnoreType, Value: string([]byte{char})}, true
	}
	// 检查一下是不是一个左圆括号
	if char == '{' {
		return &Token{Type: LeftBraceType}, true
	}
	if char == '}' {
		return &Token{Type: RightBraceType}, true
	}
	return nil, false
}

// spaceHook 空格类解析
func spaceHook(t *TokenParser) (*Token, bool) {
	var char = t.Input[t.Index]
	// 所以我们只是简单地检查是不是空格，如果是，那么我们直接进入下一个循环。
	if matched, err := regexp.MatchString("\\s", string([]byte{char})); matched && err == nil {
		return &Token{Type: IgnoreType, Value: string([]byte{char})}, true
	}
	return nil, false
}

// notesHook 注释解析
func notesHook(t *TokenParser) (*Token, bool) {
	var char = t.Input[t.Index]
	if byte('/') == char {
		if t.Index+1 >= len(t.Input) {
			return nil, false
		}
		if t.Input[t.Index+1] != byte('/') {
			return nil, false
		}
		value := make([]byte, 0)
		for t.Index += 2; t.Index < len(t.Input) && t.Input[t.Index] != byte('\n'); t.Index++ {
			value = append(value, t.Input[t.Index])
		}
		// 空的日志返回空
		if len(value) == 0 {
			return &Token{Type: IgnoreType}, true
		}
		return &Token{Type: NotesType, Value: strings.TrimSpace(string(value))}, true
	}
	return nil, false
}

// alphabetHook 字母解析
func alphabetHook(t *TokenParser) (*Token, bool) {
	var char = t.Input[t.Index]
	input := t.Input
	if tool.IsLowerAlpha(char) {
		value := make([]byte, 0)
		for ; t.Index < len(input) && tool.IsLowerAlpha(input[t.Index]); t.Index++ {
			value = append(value, input[t.Index])
		}
		if len(value) == 0 {
			t.Index--
			return nil, true
		}
		t.Index--
		switch string(value) {
		case ForType.String():
			return &Token{Type: ForType}, true
		case ElseType.String():
			return &Token{Type: ElseType}, true
		case IfType.String():
			return &Token{Type: IfType}, true
		case ReturnType.String():
			return &Token{Type: ReturnType}, true
		case InterfaceType.String():
			return &Token{Type: IgnoreType, Value: "interface{}"}, true
		default:
			return nil, true
		}
	}
	return nil, false
}
