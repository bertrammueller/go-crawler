package crawl

import (
	"bufio"
	"fmt"
	"net/url"
	"os"
)

type sitemap struct {
	edges map[*url.URL][]*url.URL
}

type WriteString interface {
	WriteString(s string) (int, error)
}

func newSitemap() sitemap {
	return sitemap{
		edges: make(map[*url.URL][]*url.URL),
	}
}

func (s *sitemap) addCrawledPage(base *url.URL, links []*url.URL) {
	s.edges[base] = links
}

func (s *sitemap) createDotFile() {
	fmt.Println("Writing DOT file..")
	f, err := os.Create("sitemap.dot")
	if err != nil {
		fmt.Println("Error writing sitemap", err.Error())
	}
	defer f.Close()
	w := bufio.NewWriter(f)
	s.writeDotFile(w)
	w.Flush()
}

func (s *sitemap) writeDotFile(w WriteString) {
	w.WriteString("digraph d {\n")
	w.WriteString("graph [	fontname = \"Helvetica-Oblique\", size = \"10,10\" ];\n")
	for v, e := range s.edges {
		if len(e) == 0 {
			w.WriteString(getNodeName(v) + "\n")
		} else {
			for _, edgeUrl := range e {
				w.WriteString(getNodeName(v) + " -> " + getNodeName(edgeUrl) + "\n")
			}
		}
	}
	w.WriteString("}")
}

func getNodeName(u *url.URL) string {
	if u.Path == "" || u.Path == "/" {
		return "\"" + u.String() + "\""
	}
	return "\"" + u.Path + "\""
}
