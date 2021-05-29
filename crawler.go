package main

import (
	"flag"
	"log"
	"os"
)

type (
	crawler struct {
	}

	Crawler interface {
		Run()
	}
)

var urlFlag = flag.String("url", "", "URL of the site, for example, https://golang.org")

func main() {
	New().Run()
}

func New() Crawler {
	return &crawler{}
}

func (c *crawler) Run() {
	url := os.Getenv("URL")
	if len(url) == 0 {
		url = *urlFlag
	}
	if len(url) == 0 {
		log.Fatal("no site url found")
	}
}
