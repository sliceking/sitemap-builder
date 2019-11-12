package main

import (
	"flag"
	"fmt"
)

func main() {
	urlFlag := flag.String("url", "https://stanwielga.com", "the url you want to build a sitemap for")
	flag.Parse()
	fmt.Println(*urlFlag)

	/*
		get the webpage
		parse all the links on the page
		build proper urls with our links
		filter out any links with a different domain
		find all pages BFS
		print out XML
	*/
}
