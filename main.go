package main

import (
	"flag"
	"io"
	"net/http"
	"os"
)

func main() {
	urlFlag := flag.String("url", "https://stanwielga.com", "the url you want to build a sitemap for")
	flag.Parse()

	resp, err := http.Get(*urlFlag)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	io.Copy(os.Stdout, resp.Body)
	/*
		get the webpage
		parse all the links on the page
		build proper urls with our links
		filter out any links with a different domain
		find all pages BFS
		print out XML
	*/
}
