package crawler

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	neturl "net/url"
	"strings"
	"time"

	"golang.org/x/net/html"
)

type (
	crawler struct {
		links map[string]struct{}
	}

	Crawler interface {
		Run(url string) ([]string, error)
	}
)

func New() Crawler {
	return &crawler{
		links: make(map[string]struct{}),
	}
}

func (c *crawler) Run(url string) ([]string, error) {
	if len(url) == 0 {
		return nil, errors.New("no site url found")
	}
	fmt.Printf("Loading %s \n", url)

	page, err := c.get(url, 0)
	if err != nil {
		return nil, err
	}
	doc, err := html.Parse(bytes.NewReader(page))
	if err != nil {
		return nil, err
	}
	c.parse(doc, url)
	u, err := neturl.Parse(url)
	if err != nil {
		return nil, err
	}
	result := make([]string, len(c.links))
	i := 0
	for key := range c.links {
		if strings.HasPrefix(key, "/") {
			result[i] = u.Scheme + "://" + u.Host + key
		} else {
			result[i] = key
		}
		i++
	}
	return result, nil
}

func (c *crawler) get(url string, count int) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error get %s: %s", url, err.Error())
	}
	if resp == nil {
		return nil, fmt.Errorf("nil response from %s", url)
	}
	if resp.StatusCode != http.StatusOK && count < 3 {
		if count == 2 {
			return []byte{}, fmt.Errorf("not downloaded", url)
		}
		log.Println("Error loading", url)
		time.Sleep(time.Duration(300+rand.Intn(1000)) * time.Millisecond)
		return c.get(url, count+1)
	}
	defer func() {
		closeErr := resp.Body.Close()
		if closeErr != nil {
			log.Printf("error close response body %s \n", closeErr.Error())
		}
	}()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error resp response body %s", err.Error())
	}
	fmt.Println("Loaded", url)
	return data, nil
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
