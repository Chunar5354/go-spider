package spider

import (
	"regexp"
)

const baseUrl = "https://book.douban.com"

var tagRe = regexp.MustCompile(`<a href="(.*?)" class="tag">(.*?)</a>`)
var bookListRe = regexp.MustCompile(`<a href="(.*?)" title="(.*?)"`)
var titleRe = regexp.MustCompile(`<span property="v:itemreviewed">(.*?)</span>`)
var authorRe = regexp.MustCompile(`<span class="pl"> 作者</span>:[\d\D]*?<a.*?>(.*?)</a>`)
var publisherRe = regexp.MustCompile(`<span class="pl">出版社:</span> (.*?)<br/>`)
var pageRe = regexp.MustCompile(`<span class="pl">页数:</span> (.*?)<br/>`)
var priceRe = regexp.MustCompile(`<span class="pl">定价:</span> (.*?)<br/>`)
var scoreRe = regexp.MustCompile(`<strong class="ll rating_num " property="v:average"> (.*?) </strong>`)
var idRe = regexp.MustCompile(`https://book.douban.com/subject/(.*?)/`)

// 在主页找到各个标签
func ParseTag(contents []byte, _ string) ParseResult {
	// 正则表达式匹配，(.*?)
	match := tagRe.FindAllSubmatch(contents, -1)

	result := ParseResult{}
	// 将目标内容添加到新的Requests中
	for _, m := range match {
		curr := m // 防止先遍历后运行的情况出现，在循环内部新定义一个变量
		// result.Items = append(result.Items, string(m[2]))
		result.Requests = append(result.Requests, Request{
			Url:   baseUrl + string(curr[1]),
			Parse: NewFuncParser(ParseBookList, "ParseBookList"), // 爬取标签之后要爬取图书列表
		})
	}
	return result
}

// 在标签页面中找到每本书
func ParseBookList(contents []byte, _ string) ParseResult {
	match := bookListRe.FindAllSubmatch(contents, -1)
	result := ParseResult{}

	for _, m := range match {
		// result.Items = append(result.Items, string(m[2]))
		curr := m
		result.Requests = append(result.Requests, Request{
			Url:   string(curr[1]),
			Parse: NewFuncParser(ParseBookDetail, "ParseBookDetail"), // 将ParseBookDetail进行包装，这样就可以额外传递参数，并在Request中保持相同的ParseFunc类型
		})
	}
	return result
}

// 获取书籍的具体信息
func ParseBookDetail(contents []byte, url string) ParseResult {
	bookDetail := BookDetail{}
	bookDetail.Url = url
	bookDetail.Title = BookString(contents, titleRe)
	bookDetail.Author = BookString(contents, authorRe)
	bookDetail.Publisher = BookString(contents, publisherRe)
	bookDetail.Page = BookString(contents, pageRe)
	bookDetail.Price = BookString(contents, priceRe)
	bookDetail.Score = BookString(contents, scoreRe)

	result := ParseResult{
		Items: []Item{
			{
				Url:     url,
				Type:    "book",
				Id:      BookString([]byte(url), idRe),
				Payload: bookDetail,
			},
		},
	}
	return result
}

func BookString(contents []byte, re *regexp.Regexp) string {
	match := re.FindSubmatch(contents)
	if len(match) >= 2 {
		return string(match[1])
	} else {
		return ""
	}
}

// 用于占位的空函数
func NilParse(_ []byte, _ string) ParseResult {
	return ParseResult{}
}
