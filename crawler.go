package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	"golang.org/x/net/html"
)

type (
	crawler struct {
		links    map[string]struct{}
		excluded map[string]struct{}
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
	return &crawler{links: make(map[string]struct{})}
}

func (c *crawler) Run() {
	url := os.Getenv("URL")
	if len(url) == 0 {
		url = *urlFlag
	}
	if len(url) == 0 {
		log.Fatal("no site url found")
	}

	excl := os.Getenv("EXCLUDED")
	if len(excl) > 0 {
		excls := strings.Split(excl, ",")
		for _, ex := range excls {
			c.excluded[ex] = struct{}{}
		}
	}

	page := c.get(url, 0)
	doc, err := html.Parse(bytes.NewReader(page))
	if err != nil {
		log.Fatal(err)
	}
	c.parse(doc)

	for key := range c.links {
		k := key
		_, ok := c.excluded[k]
		if strings.HasPrefix(k, "#") || ok {
			continue
		}
		if strings.HasPrefix(k, "/") {
			k = url + k
		}
		page := c.get(k, 0)
		doc, err := html.Parse(bytes.NewReader(page))
		if err != nil {
			log.Fatal(err)
		}
		c.parse(doc)
	}

	c.print()
}

func (c *crawler) get(url string, count int) []byte {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalf("error get %s: %s", url, err.Error())
	}
	if resp == nil {
		log.Fatalf("nil response from %s", url)
	}
	if resp.StatusCode != http.StatusOK && count < 3 {
		log.Println("Error loading", url)
		time.Sleep(time.Duration(300+rand.Intn(1000)) * time.Millisecond)
		return c.get(url, count+1)
	}
	defer func() {
		closeErr := resp.Body.Close()
		if closeErr != nil {
			log.Fatalf("error close response body %s", closeErr.Error())
		}
	}()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("error resp response body %s", err.Error())
	}
	fmt.Println("Loaded", url)
	return data
}

func (c *crawler) parse(n *html.Node) {
	if n.Type == html.ElementNode && n.Data == "a" {
		for _, at := range n.Attr {
			if at.Key == "href" {
				if _, ok := c.links[at.Val]; !ok {
					c.links[at.Val] = struct{}{}
				}
			}
		}
	}
	for nn := n.FirstChild; nn != nil; nn = nn.NextSibling {
		c.parse(nn)
	}
}

func (c *crawler) print() {
	for key := range c.links {
		fmt.Println(key)
	}
}
