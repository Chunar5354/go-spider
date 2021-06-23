# go-spider

使用Go语言编写的分布式爬虫

[参考](https://www.bilibili.com/video/BV1XK411V7DT?p=1)

## 说明

主程序是engine.go中的Run()方法，它创建workers（goroutine）并不断地将新的请求添加到请求队列（channel）中，同时worker解析出的结果中包含下一层的页面url信息，以此来层层遍历，直到找到想要的信息

[![RmWTqs.md.png](https://z3.ax1x.com/2021/06/23/RmWTqs.md.png)](https://imgtu.com/i/RmWTqs)

爬取不同的网页只需要根据需求修改parse.go

增加了rpc功能，将页面的parse方法进行序列化，使得解析页面的worker与保存数据的server可以运行在不同的主机上

## 测试

### 豆瓣读书

在服务器端运行mySpider下的`rpcServer.go`，并在客户端运行`rpcClient.go`就可以将爬取结果保存在服务器端的elasticsearch中(豆瓣读书)

### B站

spiderBili目录下的程序是爬取Bilibili视频信息（播放量，弹幕数量，日期），在原来的mySpider基础上做了一些改动（主要是fetch和parse）

由于页面由大量的AJAX动态加载和JS渲染构成，所以通过`chromedp`实现动态页面的解析
