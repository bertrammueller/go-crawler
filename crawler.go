package main

import (
	"flag"

	crawl "github.com/bertgit/crawler/internal"
)

func main() {
	url := flag.String("url", "http://www.monzo.com", "url to crawl")
	flag.Parse()
	crawler := new(crawl.Crawler)
	crawler.Run(*url)
}
