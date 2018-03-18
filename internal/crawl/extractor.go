package crawl

import (
	"fmt"
	"io"
	"net/url"

	"golang.org/x/net/html"
)

func extractUrls(body io.Reader, base *url.URL) []*url.URL {
	var urls []*url.URL
	links := extractHref(body)
	for _, link := range links {
		linkedUrl, err := base.Parse(link)
		if err != nil {
			fmt.Println("Invalid link found:", link)
		} else {
			// Only crawl the original host
			if base.Hostname() == linkedUrl.Hostname() {
				// Remove fragments (#...) before appending
				linkedUrl.Fragment = ""
				urls = append(urls, linkedUrl)
			}
		}
	}
	return urls
}

func extractHref(body io.Reader) []string {
	var links []string
	z := html.NewTokenizer(body)
	for {
		tt := z.Next()

		switch tt {
		case html.ErrorToken:
			return links
		case html.StartTagToken, html.EndTagToken:
			token := z.Token()
			if "a" == token.Data {
				for _, attr := range token.Attr {
					if attr.Key == "href" {
						links = append(links, attr.Val)
					}
				}
			}
		}
	}
}
