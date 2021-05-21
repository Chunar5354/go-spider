package spider

import (
	"log"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"

	"github.com/olivere/elastic/v7"
)

// 提供的RPC调用接口
type ItemService struct {
	Client *elastic.Client
}

type WorkService struct {
}

func (i *ItemService) RemoteSaveItem(item Item, result *string) error {
	err := Save(i.Client, item)
	if err != nil {
		return err
	}
	*result = "ok"
	return nil
}

func (w *WorkService) Process(req WorkRequest, workResult *WorkParseResult) error {
	request, err := DeserializeRequest(req)
	if err != nil {
		return err
	}

	// 执行任务
	result, err := DoWork(request)
	// 将结果序列化之后返回
	*workResult = SerializeResult(result)
	return nil
}

// RPC服务注册
func ServeRpc(host string, service interface{}) error {
	rpc.Register(service)

	listen, err := net.Listen("tcp", host)
	if err != nil {
		return err
	}

	for {
		conn, err := listen.Accept()
		if err != nil {
			log.Printf("accept error: %v", err)
			continue
		}
		go jsonrpc.ServeConn(conn)
	}
	return nil
}
