package spider

import (
	"fmt"
	"regexp"
)

const baseUrl = "https://www.bilibili.com"
const SpmIdFrom = "/?spm_id_from=333.5.b_7375626e6176.3#"

var tagRe = regexp.MustCompile(`tabindex="0"><a href="//www.bilibili.com/v(.*?)" class="name"><span>(.*?)<em>`)
var subTagRe = regexp.MustCompile(`<li class=""><a href="/v(.*?)">(.*?)<!----></a></li>`)
var pageCountRe = regexp.MustCompile(`<button class="pagination-btn">(.*?)</button>`)
var pageRe = regexp.MustCompile(`<div class="spread-module"><a href="(.*?)" target="_blank">`)
var detailRe = regexp.MustCompile(`<div class="video-data"><span.*?>(.*?)播放&nbsp;·&nbsp;</span><span.*?>(.*?)弹幕</span><span>(.*?)</span><!----></div>`)

// 在主页找到各个标签
func ParseTag(url string) ParseResult {
	waitExpression := `#app > div > div.storey-box.b-wrap > div.proxy-box`
	jsSelector := `document.querySelector("#internationalHeader > div.b-wrap > div")`
	contents, _ := FetchChrome(url, waitExpression, jsSelector)
	// contents, _ := Fetch(url) // 标签可以直接获取
	// 正则表达式匹配，(.*?)
	match := tagRe.FindAllSubmatch(contents, -1)
	// match := subTagRe.FindAllSubmatch(contents, -1)

	result := ParseResult{}
	// 将目标内容添加到新的Requests中
	for _, m := range match {
		curr := m // 防止先遍历后运行的情况出现，在循环内部新定义一个变量
		// result.Items = append(result.Items, string(m[2]))
		fmt.Printf("%s, %s\n", m[1], m[2])
		result.Requests = append(result.Requests, Request{
			Url:   baseUrl + "/v"  + string(curr[1]),
			Parse: NewFuncParser(ParseSubTag, "ParseSubTag"), // 爬取标签之后要爬取图书列表
		})
	}
	return result
}

// 找到各个子标签
func ParseSubTag(url string) ParseResult {
	waitExpression := `#internationalHeader > div.b-wrap > div`
	jsSelector := `document.querySelector("#subnav > ul")`
	contents, _ := FetchChrome(url, waitExpression, jsSelector)
	// contents, _ := Fetch(url) // 子标签可以直接获取
	// 正则表达式匹配，(.*?)
	match := subTagRe.FindAllSubmatch(contents, -1)

	result := ParseResult{}
	// 将目标内容添加到新的Requests中
	for _, m := range match {
		curr := m // 防止先遍历后运行的情况出现，在循环内部新定义一个变量
		// result.Items = append(result.Items, string(m[2]))
		fmt.Printf("%s, %s\n", m[1], m[2])
		result.Requests = append(result.Requests, Request{
			Url:   baseUrl + "/v" + string(curr[1]),
			Parse: NewFuncParser(ParsePages, "ParsePages"), // 爬取标签之后要爬取图书列表
		})
	}
	return result
}

// 解析页面
func ParsePages(url string) ParseResult {
	/*
	waitExpression := `#videolist_box > div.vd-list-cnt > div.pager.pagination > ul > li.page-item.last > button`
	jsPath := `document.querySelector("#videolist_box > div.vd-list-cnt > div.pager.pagination > ul > li.page-item.last > button")`
	contents, _ := FetchChrome(url, waitExpression, jsPath) // 总页数的信息需要通过动态方式获取
	// 正则表达式匹配，(.*?)
	match := pageCountRe.FindSubmatch(contents)
	fmt.Println(string(match[1]))
	*/
	fmt.Println("Temp: ", url)
	result := ParseResult{}
	result.Requests = append(result.Requests, Request{
		// Url:   "https:" + string(match[1]),
		Url:   url + "#/all/default/0/4/",
		Parse: NewFuncParser(ParseSinglePage, "ParseSinglePage"), // 爬取标签之后要爬取图书列表
	})
	return result
}

func ParseSinglePage(url string) ParseResult {
	waitExpression := `#videolist_box`
	jsSelector := `document.querySelector("#videolist_box")`
	contents, _ := FetchChrome(url, waitExpression, jsSelector)

	// 正则表达式匹配，(.*?)
	match := pageRe.FindAllSubmatch(contents, -1)

	result := ParseResult{}
	// 将目标内容添加到新的Requests中
	for _, m := range match {
		curr := m // 防止先遍历后运行的情况出现，在循环内部新定义一个变量
		// result.Items = append(result.Items, string(m[2]))
		fmt.Printf("%s\n", m[1])
		result.Requests = append(result.Requests, Request{
			Url:   "https:" + string(curr[1]),
			Parse: NewFuncParser(ParseDetail, "ParseDetail"), // 爬取标签之后要爬取图书列表
		})
	}
	return result
}

func ParseDetail(url string) ParseResult {
	waitExpression := `#viewbox_report`
	jsSelector := `document.querySelector("#viewbox_report")`
	contents, _ := FetchChrome(url, waitExpression, jsSelector)

	// 正则表达式匹配，(.*?)
	match := detailRe.FindSubmatch(contents)

	videoDetail := VideoDetail{}
	videoDetail.Url = url
	videoDetail.Play = string(match[1])
	videoDetail.Comment = string(match[2])
	videoDetail.Date = string(match[3])

	result := ParseResult{
		Items: []Item{
			{
				Url:     url,
				Type:    "video",
				Id:      url[31:],
				Payload: videoDetail,
			},
		},
	}

	return result
}

// 用于占位的空函数
func NilParse(string) ParseResult {
	return ParseResult{}
}
