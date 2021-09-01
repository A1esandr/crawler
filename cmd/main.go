package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/A1esandr/crawler"
)

var urlFlag = flag.String("url", "", "URL of the site, for example, https://golang.org")

func main() {
	flag.Parse()
	url := os.Getenv("URL")
	if len(url) == 0 {
		url = *urlFlag
	}
	links, err := crawler.New().Run(url)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("Got %d links\n", len(links))
	var result string
	for _, link := range links {
		fmt.Println(link)
		result += link + "\n"
	}
	err = os.WriteFile("links", []byte(result), 0644)
	fmt.Println(err)
}
