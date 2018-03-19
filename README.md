# Web Crawler in Go
### Design
A pool of workers pops tasks from a url frontier channel. Each worker fetches the html body, extracts links and adds them back on a response channel. The crawler runs them through a url cache and adds the filtered links back into the url frontier.  
The output is a file in DOT format (sitemap.dot), which can be used to produce graphs similar to `sitemap-reduced.csv`

### Limitations
* Only respects html hrefs
* Inpolite, doesn't consider robots.txt
* Doesn't detect crawler traps (infinite generated loops, ...)
* Single machine
* Lots more compared to serious crawlers...

### Build & Run
Install in your go path and run  
`go build && ./go-crawler`

### Test
`go test`

### Usage
<pre>Usage of ./crawler:  
  -timeout int  
        Timeout in seconds per http request (default 10)  
  -url string  
        url to crawl (default "http://www.monzo.com")  
  -workers int  
        Number of workers (default 100)
</pre>

### Performance
i5 dual core, 5Ghz Wifi:  
Crawls 305 pages from www.monzo.com in ~3s

Number of workers doesn't improve much beyond 50 for above test case

### Sitemap
To create a sitemap visualization, install  
http://www.graphviz.org  
and run  
`dot -Tsvg sitemap.dot -O`

The resulting sitemap can be very big. The example provided in this repo (`sitemap-reduced.svg`) is reduced to a single link per website.