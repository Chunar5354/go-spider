package spider

// 为了注册成RPC服务，解析函数需要是实现了Parser接口的类
type Parser interface {
	Parse(contents []byte, url string) ParseResult
	Serialize() (name string, args interface{})
}

type ParseFunc func(contents []byte, url string) ParseResult

// 通用的页面解析
type FuncParser struct {
	Parser ParseFunc
	Name   string
}

func (f FuncParser) Parse(contents []byte, url string) ParseResult {
	return f.Parser(contents, url)
}

func (f FuncParser) Serialize() (name string, args interface{}) {
	return f.Name, nil
}

// 新建FuncParser
func NewFuncParser(p ParseFunc, name string) *FuncParser {
	return &FuncParser{
		Parser: p,
		Name:   name,
	}
}

// 空解析函数
type NilParser struct {
}

func (n *NilParser) Parse(contents []byte, url string) ParseResult {
	return ParseResult{}
}

func (n *NilParser) Serialize() (name string, args interface{}) {
	return "NilParse", nil
}
