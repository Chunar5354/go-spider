package spider

import "fmt"

func CreateProcess(host string) (Processor, error) {
	client, err := NewClient(host)
	if err != nil {
		fmt.Printf("%v", err)
		return func(req Request) (ParseResult, error) {
			return ParseResult{}, err
		}, err
	}

	return func(req Request) (ParseResult, error) {
		workRequest := SerializeRequest(req)
		// workResult := WorkParseResult{}
		var workResult WorkParseResult
		// fmt.Println(workRequest)
		err := client.Call("WorkService.Process", workRequest, &workResult)
		if err != nil {
			return ParseResult{}, err
		}
		return DeserializeResult(workResult), nil
	}, nil
}
