package spider

import (
	"context"
	"log"

	"github.com/olivere/elastic/v7"
)

// 本地保存
func SaveItem() (chan Item, error) {
	outItem := make(chan Item)

	// 创建es客户端
	client, err := elastic.NewClient(
		elastic.SetSniff(false),
		elastic.SetURL("http://127.0.0.1:9200"),
	)
	if err != nil {
		return nil, err
	}

	go func() {
		count := 0
		for {
			item := <-outItem
			Save(client, item)
			// log.Printf("Item saver: Got %d, %v", count, item)
			count++
		}
	}()

	return outItem, nil
}

// 保存操作的主要执行逻辑
func Save(client *elastic.Client, item Item) error {
	ctx := context.Background()
	indexName := "video"

	indexService := client.Index().Index(indexName).Type(item.Type).BodyJson(item)

	if item.Id != "" {
		indexService.Id(item.Id)
	}

	_, err := indexService.Do(ctx)

	if err != nil {
		return err
	}
	log.Printf("Item saver: Got %s", item)
	return nil
}
