package spider

import (
	"bufio"
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/chromedp/chromedp"
	"golang.org/x/net/html/charset"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
)

// 发起请求获得页面数据
func Fetch(url string) ([]byte, error) {
	time.Sleep(time.Millisecond * 100)
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("Error: get url: %s", err)
	}
	// 设置请求头
	req.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/90.0.4430.93 Safari/537.36")
	// 发送请求
	resp, err := client.Do(req)

	// 转码
	bodyReader := bufio.NewReader(resp.Body)
	e := DetermineEncoding(bodyReader)
	// e.NewDecoder()返回一个Encoder对象，它包含transform.Transformer接口，实现了Transform方法，将一种格式的字节数组转换为另一种格式的字节数组
	utf8Reader := transform.NewReader(bodyReader, e.NewDecoder()) // 转换为utf-8编码

	result, err := ioutil.ReadAll(utf8Reader)
	if err != nil {
		panic(err)
	}
	return result, err
}

// 获得页面的编码格式
func DetermineEncoding(r *bufio.Reader) encoding.Encoding {
	bytes, err := r.Peek(1024)
	if err != nil {
		log.Printf("fetch error: %v", err)
		return unicode.UTF8
	}

	// e是encoding.Encoding接口，它实现了NewDecoder和NewEncoder两个方法，分别返回Decoder和Encoder对象
	e, _, _ := charset.DetermineEncoding(bytes, "")
	return e
}

// 通过chromedp获取动态内容
func FetchChrome(url, waitExpression, jsSelector string) ([]byte, error) {
	time.Sleep(time.Millisecond * 500)
	options := []chromedp.ExecAllocatorOption{
		chromedp.Flag("headless", true),    // debug使用
		chromedp.Flag("disable-gpu", true), // debug使用
		chromedp.Flag("blink-settings", "imagesEnabled=false"),
		chromedp.UserAgent(`Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/90.0.4430.212 Safari/537.36`),
	}
	options = append(chromedp.DefaultExecAllocatorOptions[:], options...)

	c, _ := chromedp.NewExecAllocator(context.Background(), options...)

	// create context
	chromeCtx, cancel := chromedp.NewContext(c, chromedp.WithLogf(log.Printf))
	defer cancel()
	// 执行一个空task, 用提前创建Chrome实例
	chromedp.Run(chromeCtx, make([]chromedp.Action, 0, 1)...)
	timeoutCtx, cancel := context.WithTimeout(chromeCtx, 20*time.Second)
	defer cancel()

	log.Printf("Chrome visit page %s\n", url)

	var htmlContent string
	err := chromedp.Run(timeoutCtx,
		chromedp.Navigate(url),
		chromedp.WaitVisible(waitExpression),
		chromedp.OuterHTML(jsSelector, &htmlContent, chromedp.ByJSPath),
	)
	if err != nil {
		panic(err)
	}
	return []byte(htmlContent), nil
}
