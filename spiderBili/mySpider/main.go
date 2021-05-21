package main

import (
	"spider"
)

func main() {
	/*
		request := spider.Request{
			Url:       "https://book.douban.com/",
			ParseFunc: spider.ParseTag,
		}

		saveItem, err := spider.SaveItem()
		if err != nil {
			panic(err)
		}

		e := spider.ConcurrentEngine{
			WorkerCount: 10,
			Scheduler:   &(spider.QueueScheduler{}),
			ItemChan:    saveItem,
		}
		e.Run(request)
	*/

	url := "https://www.bilibili.com/v/douga"
	result, err := spider.Fetch(url)
	if err != nil {
		panic(err)
	}
	// fmt.Println(string(result))
	spider.ParseTag(result, "")
}
