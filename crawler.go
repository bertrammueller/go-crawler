package main

import (
	"flag"

	crawler "github.com/bertgit/crawler/internal"
)

func main() {
	domain := flag.String("domain", "monzo.com", "Domain to crawl")
	flag.Parse()
	crawler := new(crawler.Crawler)
	crawler.Run(*domain)
}
