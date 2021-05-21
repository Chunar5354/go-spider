package spider

// 抽象一个调度器接口
type Scheduler interface {
	Submit(Request) // 将请求发送到队列
	WorkerReady(chan Request)
	Run()
	CreateWorkChan() chan Request
}

type SimpleScheduler struct {
	WorkChan chan Request
}

// 队列式调度器
type QueueScheduler struct {
	RequestChan chan Request
	WorkChan    chan chan Request // WorkChan中传递的是每一个worker的通道
}

func (q *QueueScheduler) Submit(r Request) {
	q.RequestChan <- r
}

func (q *QueueScheduler) WorkerReady(w chan Request) {
	q.WorkChan <- w
}

func (q *QueueScheduler) CreateWorkChan() chan Request {
	return make(chan Request)
}

func (q *QueueScheduler) Run() {
	// 调度器的主通道
	q.RequestChan = make(chan Request)
	q.WorkChan = make(chan chan Request)

	go func() {
		var requestQ []Request   // 请求等待区
		var workQ []chan Request // worker等待区
		for {
			var activeRequest Request
			var activeWork chan Request

			if len(requestQ) > 0 && len(workQ) > 0 {
				activeRequest = requestQ[0]
				activeWork = workQ[0]
			}

			// 在这里阻塞直到：
			// 1.有新的请求到来：添加到请求等待区
			// 2.有新的worker准备好：添加到worker等待区
			// 3.有已经就绪的请求和工作者：分配一个工作
			select {
			case r := <-q.RequestChan:
				requestQ = append(requestQ, r)
			case w := <-q.WorkChan:
				workQ = append(workQ, w)
			case activeWork <- activeRequest:
				requestQ = requestQ[1:]
				workQ = workQ[1:]
			}
		}
	}()
}

func (s *SimpleScheduler) Submit(r Request) {
	go func() { s.WorkChan <- r }()
}

func (s *SimpleScheduler) WorkerReady(c chan Request) {
	return
}

func (s *SimpleScheduler) CreateWorkChan() chan Request {
	return s.WorkChan
}

func (s *SimpleScheduler) Run() {
	s.WorkChan = make(chan Request)
}
