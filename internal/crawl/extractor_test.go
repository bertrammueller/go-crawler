package crawl

import (
	"net/url"
	"strings"
	"testing"
)

func TestExtractsHrefsFromHtml(t *testing.T) {
	s := "<h1><a href=\"http://www.monzo.com/foo\">link1</a><a href=\"/bar\">link2</a></h1>"
	body := strings.NewReader(s)
	links := extractHref(body)
	if len(links) != 2 {
		t.Error("Expected 2 links to be extracted, got", len(links))
	}
	if links[0] != "http://www.monzo.com/foo" || links[1] != "/bar" {
		t.Error("Failed to extract links")
	}
}

func TestReturnsUniqueUrls(t *testing.T) {
	base, _ := url.Parse("http://www.monzo.com")
	urls := uniqUrls(base, []string{"/foo", "/bar", "/foo"})
	if len(urls) != 2 {
		t.Error("Wrong number of Urls returned")
	}
	if urls[0].Path != "/foo" || urls[1].Path != "/bar" {
		t.Error("Wrong Urls returned")
	}
}

func TestDiscardsFragments(t *testing.T) {
	base, _ := url.Parse("http://www.monzo.com")
	urls := uniqUrls(base, []string{"/foo#bar"})
	if len(urls) != 1 || urls[0].Path != "/foo" {
		t.Error("Fragments not discarded")
	}
}
