package main

import (
	"log"
	"spider"

	"github.com/olivere/elastic/v7"
)

func main() {
	serveRpc(":12008")
	log.Fatal(spider.ServeRpc(":12009", &spider.WorkService{}))
}

func serveRpc(host string) error {
	client, err := elastic.NewClient(
		elastic.SetSniff(false),
		elastic.SetURL("http://127.0.0.1:9200"),
	)

	if err != nil {
		return err
	}

	// 为什么是指针？
	return spider.ServeRpc(host, &spider.ItemService{
		Client: client,
	})
}
