package spider

import (
	"errors"
	"log"
)

// RPC调用中，将Parse解析函数注册成服务
type SerializeParser struct {
	Name string
	Args interface{}
}

type WorkRequest struct {
	Url   string
	Parse SerializeParser
}

type WorkParseResult struct {
	Items    []Item
	Requests []WorkRequest
}

// 序列化
func SerializeResult(r ParseResult) WorkParseResult {
	result := WorkParseResult{Items: r.Items}

	for _, req := range r.Requests {
		result.Requests = append(result.Requests, SerializeRequest(req))
	}
	return result
}

func SerializeRequest(r Request) WorkRequest {
	name, args := r.Parse.Serialize()

	return WorkRequest{
		Url: r.Url,
		Parse: SerializeParser{
			Name: name,
			Args: args,
		},
	}
}

// 反序列化
func DeserializeResult(r WorkParseResult) ParseResult {
	result := ParseResult{
		Items: r.Items,
	}

	for _, req := range r.Requests {
		deReq, err := DeserializeRequest(req)
		if err != nil {
			log.Printf("Error deserializing: %v", err)
		}
		result.Requests = append(result.Requests, deReq)
	}
	return result
}

func DeserializeRequest(r WorkRequest) (Request, error) {
	parse, err := DeserializeParse(r.Parse)
	if err != nil {
		return Request{}, err
	}
	return Request{
		Url:   r.Url,
		Parse: parse,
	}, nil
}

func DeserializeParse(p SerializeParser) (Parser, error) {
	switch p.Name {
	case "ParseTag":
		return NewFuncParser(ParseTag, "ParseTag"), nil
	case "ParseSubTag":
		return NewFuncParser(ParseSubTag, "ParseSubTag"), nil
	case "ParsePages":
		return NewFuncParser(ParsePages, "ParsePages"), nil
	case "ParseSinglePage":
		return NewFuncParser(ParseSinglePage, "ParseSinglePage"), nil
	case "ParseDetail":
		return NewFuncParser(ParseDetail, "ParseDetail"), nil
	case "NilParse":
		return NewFuncParser(NilParse, "NilParse"), nil
	default:
		return nil, errors.New("Unknown parse name")
	}

}
