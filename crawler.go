package crawler

import (
	"bytes"
	"errors"
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
		links map[string]struct{}
	}

	Crawler interface {
		Run() ([]string, error)
	}
)

var urlFlag = flag.String("url", "", "URL of the site, for example, https://golang.org")

func New() Crawler {
	return &crawler{
		links: make(map[string]struct{}),
	}
}

func (c *crawler) Run() ([]string, error) {
	url := os.Getenv("URL")
	if len(url) == 0 {
		url = *urlFlag
	}
	if len(url) == 0 {
		return nil, errors.New("no site url found")
	}
	fmt.Println("Started")

	page := c.get(url, 0)
	doc, err := html.Parse(bytes.NewReader(page))
	if err != nil {
		return nil, err
	}
	c.parse(doc, url)
	result := make([]string, len(c.links))
	i := 0
	for key := range c.links {
		if strings.HasPrefix(key, "/") {
			result[i] = url + key
		} else {
			result[i] = key
		}
		i++
	}
	return result, nil
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
		if count == 2 {
			log.Println("Not downloaded", url)
			return []byte{}
		}
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

func (c *crawler) parse(n *html.Node, url string) {
	if n.Type == html.ElementNode && n.Data == "a" {
		for _, at := range n.Attr {
			if at.Key == "href" {
				if _, ok := c.links[at.Val]; !ok && c.allowed(at.Val) {
					c.links[at.Val] = struct{}{}
				}
			}
		}
	}
	for nn := n.FirstChild; nn != nil; nn = nn.NextSibling {
		c.parse(nn, url)
	}
}

func (c *crawler) allowed(url string) bool {
	if strings.HasPrefix(url, "#") {
		return false
	}
	if strings.HasPrefix(url, "mailto:") {
		return false
	}
	return true
}
