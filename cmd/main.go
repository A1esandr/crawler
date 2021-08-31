package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/A1esandr/crawler"
)

func main() {
	flag.Parse()
	links, err := crawler.New().Run()
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
