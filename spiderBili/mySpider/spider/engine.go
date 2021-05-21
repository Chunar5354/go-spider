package spider

import (
	"time"
)

type Processor func(Request) (ParseResult, error)

// 并发引擎
type ConcurrentEngine struct {
	WorkerCount      int
	Scheduler        Scheduler
	ItemChan         chan Item
	RequestProcessor Processor
}

// 主要运行逻辑
func (c *ConcurrentEngine) Run(seeds ...Request) {
	out := make(chan ParseResult) // 输出通道

	c.Scheduler.Run() // 调度器工作

	// 创建worker
	for i := 0; i < c.WorkerCount; i++ {
		c.CreateWorker(c.Scheduler.CreateWorkChan(), out, c.Scheduler)
	}

	// 将请求添加到输入队列
	for _, r := range seeds {
		c.Scheduler.Submit(r)
	}

	for {
		// 从输出队列取出结果
		result := <-out
		for _, item := range result.Items {
			// fmt.Printf("********Got item: %s", item)
			time.Sleep(time.Millisecond * 10)
			go func() { c.ItemChan <- item }()
		}
		// 将结果中包含的新的url放入输入队列
		for _, r := range result.Requests {
			c.Scheduler.Submit(r)
		}
		time.Sleep(time.Millisecond * 100)
	}
}

// 创建worker
func (c *ConcurrentEngine) CreateWorker(in chan Request, out chan ParseResult, s Scheduler) {
	// 每个worker是一个协程，不断等待新的请求添加到队列中
	go func() {
		for {
			s.WorkerReady(in) // 有空闲工人，将工人的channel添加到调度器的WorkChan
			request := <-in
			//result, err := DoWork(request) // 执行任务
			result, err := c.RequestProcessor(request)
			if err != nil {
				continue
			}
			out <- result
		}
	}()
}

// 爬取页面的执行
func DoWork(r Request) (ParseResult, error) {
	/*
		log.Printf("Fetching url :%s\n", r.Url)
		body, err := Fetch(r.Url)
		if err != nil {
			log.Printf("Fetch error: %s\n", r.Url)
			return ParseResult{}, err
		}
	*/

	parseResult := r.Parse.Parse(r.Url)
	return parseResult, nil
}
