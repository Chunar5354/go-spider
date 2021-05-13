package main

import (
	"spider"
)

func main() {
	request := spider.Request{
		Url:   "https://book.douban.com/",
		Parse: spider.NewFuncParser(spider.ParseTag, "ParseTag"),
	}

	// 连接存储服务
	saveItem, err := spider.ClientSaveItem("1.15.140.88:12028")
	// saveItem, err := spider.ClientSaveItem(":12345")
	if err != nil {
		panic(err)
	}

	// 连接工作服务
	process, err := spider.CreateProcess("1.15.140.88:12029")

	e := spider.ConcurrentEngine{
		WorkerCount:      10,
		Scheduler:        &(spider.QueueScheduler{}),
		ItemChan:         saveItem,
		RequestProcessor: process,
	}
	e.Run(request)
}
