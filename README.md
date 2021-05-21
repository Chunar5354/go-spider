# go-spider

使用Go语言编写的分布式爬虫

根据[教程](https://www.bilibili.com/video/BV1XK411V7DT?p=1)

在服务器端运行mySpider下的`rpcServer.go`，并在客户端运行`rpcClient.go`就可以将爬取结果保存在服务器端的elasticsearch中(豆瓣图书)

spiderBili目录下的程序是爬取Bilibili视频信息（播放量，弹幕数量，日期），在原来的mySpider基础上做了一些改动（主要是fetch和parse）