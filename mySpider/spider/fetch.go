package spider

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

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
