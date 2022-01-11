package main

import "github.com/Dann-Go/web-crawler/pkg/crawler"

func main() {
	_ = crawler.Crawl("go.dev", 6)
}
