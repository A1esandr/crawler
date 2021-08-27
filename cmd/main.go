package main

import (
	"flag"
	"github.com/A1esandr/crawler"
)

func main() {
	flag.Parse()
	crawler.New().Run()
}
