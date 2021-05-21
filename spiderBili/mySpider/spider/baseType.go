package spider

// 用于输入队列的请求信息
type Request struct {
	Url   string
	Parse Parser
}

// 用于输出队列的结果信息
type ParseResult struct {
	Requests []Request
	Items    []Item
}

// 结果的具体内容
type Item struct {
	Url     string
	Type    string      // 用于elasticsearch中的Type字段
	Id      string      // 用于elasticsearch中的Id字段
	Payload interface{} // 实际上是BookDetail
}

// 每本书的具体信息
type VideoDetail struct {
	Url     string `json:"url"`
	Play    string `json:"play"`
	Comment string `json:"comment"`
	Date    string `json:"date"`
}
