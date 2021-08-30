package main

import (
	"flag"
	"fmt"

	"github.com/A1esandr/crawler"
)

func main() {
	flag.Parse()
	links, err := crawler.New().Run()
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, link := range links {
		fmt.Println(link)
	}
}
