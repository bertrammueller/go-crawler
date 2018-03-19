package crawl

import (
	"net/url"
	"strings"
	"testing"
)

func TestSitemapAddCrawledPages(t *testing.T) {
	s := newSitemap()
	u1, _ := url.Parse("www.google.com")
	u2, _ := url.Parse("/foo")
	u3, _ := url.Parse("/bar")
	s.addCrawledPage(u1, []*url.URL{u2, u3})
	e := s.edges[u1]
	if len(e) != 2 || e[0] != u2 || e[1] != u3 {
		t.Error("Edges are not added to sitemap")
	}
}

func TestSitemapWritesDotFile(t *testing.T) {
	s := newSitemap()
	u1, _ := url.Parse("google")
	u2, _ := url.Parse("/foo")
	u3, _ := url.Parse("/bar")
	s.addCrawledPage(u1, []*url.URL{u2, u3})
	builder := strings.Builder{}
	s.writeDotFile(&builder)
	if builder.String() != "digraph d {\n\"google\" -> \"/foo\"\n\"google\" -> \"/bar\"\n}" {
		t.Error("Unexspected dot file output", builder.String())
	}
}
