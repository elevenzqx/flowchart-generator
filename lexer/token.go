package lexer

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

// TokenType token 的类型信息
type TokenType int

// 类型定义
const (
	KeywordType    TokenType = iota // 关键词
	ForType                         // 循环语句
	IfType                          // 条件语句
	ElseType                        // 条件分支语句
	InterfaceType                   // interface 类型解析
	LeftBraceType                   // 左括号类型
	RightBraceType                  // 右括号类型
	IgnoreType                      // 可忽略的类型
	NotesType                       // 注释类型
	ReturnType                      // 返回类型
)

func (tt TokenType) String() string {
	switch tt {
	case KeywordType:
		return "keyword"
	case LeftBraceType:
		return "{"
	case RightBraceType:
		return "}"
	case NotesType:
		return "notes"
	case ForType:
		return "for"
	case InterfaceType:
		return "interface{}"
	case IfType:
		return "if"
	case ElseType:
		return "else"
	case ReturnType:
		return "return"
	case IgnoreType:
		return ""
	}
	panic("unexpected token type")
}

// Token 源码解析最小单元
type Token struct {
	Type  TokenType
	Value string
}

// TokenHook 处理 token 的 hook
type TokenHook func(parse *TokenParser) (*Token, bool)

// TokenParser 解析器
type TokenParser struct {
	Input  []byte
	Tokens []*Token
	Index  int
	hooks  []TokenHook
}

// New 新建 token 解析器
func New(input []byte) *TokenParser {
	return &TokenParser{
		Input:  input,
		Tokens: []*Token{},
		hooks: []TokenHook{
			parenHook, spaceHook, notesHook, alphabetHook,
		},
	}
}

// Tokenizer token 解析
func (t *TokenParser) Tokenizer() {
	if t.Index >= len(t.Input) {
		return
	}
	for ; t.Index < len(t.Input); t.Index++ {
		for _, hook := range t.hooks {
			if token, ok := hook(t); ok {
				if token == nil {
					break
				}
				if token.Type == IgnoreType {
					break
				}
				t.Tokens = append(t.Tokens, token)
				break
			}
		}
	}
}

// Register 注册 token 的处理钩子
func (t *TokenParser) Register(hook ...TokenHook) {
	t.hooks = append(t.hooks, hook...)
}

// String 显示 token 解析内容
func (t *TokenParser) String() {
	parenPosition := 0
	for _, token := range t.Tokens {
		switch token.Type {
		case NotesType:
			for i := 0; i < parenPosition; i++ {
				print("  ")
			}
			fmt.Println("// " + token.Value)
		case LeftBraceType:
			for i := 0; i < parenPosition; i++ {
				print("  ")
			}
			fmt.Println(token.Type.String())
			parenPosition++
		case RightBraceType:
			parenPosition--
			for i := 0; i < parenPosition; i++ {
				print("  ")
			}
			fmt.Println(token.Type.String())
		case IfType:
			for i := 0; i < parenPosition; i++ {
				print("  ")
			}
			fmt.Println(token.Type.String())
		default:
			for i := 0; i < parenPosition; i++ {
				print("  ")
			}
			fmt.Println(token.Type.String())
		}
	}
}

// FromFile 从文件中加载并解析 token
func FromFile(path string) *TokenParser {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Fatalln(err)
		}
	}(file)
	content, err := ioutil.ReadAll(file)
	if err != nil {
		return nil
	}
	parser := New(content)
	parser.Tokenizer()
	return parser
}

// Analyze 分析内容
func Analyze(content string) *TokenParser {
	parser := New([]byte(content))
	parser.Tokenizer()
	parser.String()
	return parser
}
