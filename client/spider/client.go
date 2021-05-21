package spider

import (
	"log"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
)

// 新建RPC客户端
func NewClient(host string) (*rpc.Client, error) {
	conn, err := net.Dial("tcp", host)
	if err != nil {
		return nil, err
	}

	client := jsonrpc.NewClient(conn)

	return client, nil
}

// 代替SaveItem函数，RPC服务中的数据保存函数，它会调用服务器提供的服务
func ClientSaveItem(host string) (chan Item, error) {
	outItem := make(chan Item)

	client, err := NewClient(host)
	if err != nil {
		return nil, err
	}

	go func() {
		count := 0
		for {
			item := <-outItem
			result := ""
			err := client.Call("ItemService.RemoteSaveItem", item, &result)
			if err != nil {
				log.Printf("ERROR: remote saving item: %v", err)
			}
			log.Printf("Item saver: Got %d, %v", count, item)
			count++
		}
	}()

	return outItem, nil
}
