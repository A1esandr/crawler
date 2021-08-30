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
		links    map[string]struct{}
		excluded map[string]struct{}
		selected map[string]struct{}
		from     map[string]map[string]struct{}
	}

	Crawler interface {
		Run() ([]string, error)
	}
)

var urlFlag = flag.String("url", "", "URL of the site, for example, https://golang.org")
var excludedFlag = flag.String("exclude", "", "URL of the site separated by commas to exclude from parsing")
var selectedFlag = flag.String("select", "", "URL of the site separated by commas to parse only")

func New() Crawler {
	return &crawler{
		links:    make(map[string]struct{}),
		excluded: make(map[string]struct{}),
		selected: make(map[string]struct{}),
		from:     make(map[string]map[string]struct{}),
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

	sel := os.Getenv("SELECTED")
	if len(sel) == 0 {
		sel = *selectedFlag
	}
	if len(sel) > 0 {
		sels := strings.Split(sel, ",")
		for _, s := range sels {
			c.selected[s] = struct{}{}
		}
	}

	excl := os.Getenv("EXCLUDED")
	if len(excl) == 0 {
		excl = *excludedFlag
	}
	if len(excl) > 0 {
		excls := strings.Split(excl, ",")
		for _, ex := range excls {
			c.excluded[ex] = struct{}{}
		}
	}

	fmt.Println("Started")

	page := c.get(url, 0)
	doc, err := html.Parse(bytes.NewReader(page))
	if err != nil {
		return nil, err
	}
	c.parse(doc, url)

	var selFlag bool
	if len(c.selected) > 0 {
		selFlag = true
	}
	for key := range c.links {
		k := key
		if selFlag {
			if _, ok := c.selected[k]; !ok {
				continue
			}
		}
		if _, ok := c.excluded[k]; ok {
			continue
		}
		if strings.HasPrefix(k, "/") {
			k = url + k
		}
		page := c.get(k, 0)
		doc, err := html.Parse(bytes.NewReader(page))
		if err != nil {
			return nil, err
		}
		c.parse(doc, k)
	}

	c.print()
	result := make([]string, len(c.links))
	i := 0
	for key := range c.links {
		result[i] = key
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
					if _, exist := c.from[url]; !exist {
						c.from[url] = make(map[string]struct{})
					}
					c.from[url][at.Val] = struct{}{}
				}
			}
		}
	}
	for nn := n.FirstChild; nn != nil; nn = nn.NextSibling {
		c.parse(nn, url)
	}
}

func (c *crawler) print() {
	for key := range c.links {
		fmt.Println(key)
	}
	for keyMap, val := range c.from {
		for key := range val {
			fmt.Println(key, "from", keyMap)
		}
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
