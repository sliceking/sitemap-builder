package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/svwielga4/link-parser"
)

func main() {
	urlFlag := flag.String("url", "https://gophercises.com", "the url you want to build a sitemap for")
	flag.Parse()

	pages := get(*urlFlag)
	for _, href := range pages {
		fmt.Println(href)
	}
}

func get(urlStr string) []string {
	resp, err := http.Get(urlStr)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	reqURL := resp.Request.URL
	baseURL := &url.URL{
		Scheme: reqURL.Scheme,
		Host:   reqURL.Host,
	}
	base := baseURL.String()
	return filter(hrefs(resp.Body, base), withPrefix(base))
}

func hrefs(body io.Reader, base string) []string {
	links, _ := link.Parse(body)
	var ret []string
	for _, l := range links {
		switch {
		case strings.HasPrefix(l.Href, "/"):
			ret = append(ret, base+l.Href)
		case strings.HasPrefix(l.Href, "http"):
			ret = append(ret, l.Href)
		}
	}
	return ret
}

func filter(links []string, keepFn func(string) bool) []string {
	var ret []string
	for _, link := range links {
		if keepFn(link) {
			ret = append(ret, link)
		}
	}
	return ret
}

func withPrefix(pfx string) func(string) bool {
	return func(link string) bool {
		return strings.HasPrefix(link, pfx)
	}
}
