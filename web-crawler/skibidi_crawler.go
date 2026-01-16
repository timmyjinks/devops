package main

import (
	"fmt"
	"golang.org/x/net/html"
	"strings"
	"sync"
	"time"
)

type SkibidiCrawler struct {
	originURL string
	duration  time.Duration

	wg sync.WaitGroup

	mu       sync.Mutex
	websites map[string]int
}

func NewSkibidiCrawler(url string) *SkibidiCrawler {
	return &SkibidiCrawler{
		originURL: url,
		websites:  make(map[string]int),
	}
}

func (c *SkibidiCrawler) Crawl() {
	h, err := parseHTML(c.originURL)
	if err != nil {
		fmt.Println(err)
		return
	}

	c.wg.Go(func() {
		c.proccess(h, c.originURL)
	})

	c.wg.Wait()

	for website, count := range c.websites {
		fmt.Printf("Website: %s Skibidi Count: %d\n", website, count)
	}
}

func (c *SkibidiCrawler) proccess(h *html.Node, url string) {
	if h.Data == "head" {
		return
	}

	fmt.Println(url)
	c.mu.Lock()
	_, ok := c.websites[url]
	c.mu.Unlock()
	if !ok {
		c.mu.Lock()
		c.websites[url] = 0
		c.mu.Unlock()
	}

	iter := h.ChildNodes()
	for tag := range iter {
		if strings.TrimSpace(tag.Data) == "" {
			continue
		}

		c.mu.Lock()
		c.websites[url] += strings.Count(strings.ToLower(tag.Data), "skibidi")
		c.mu.Unlock()

		c.wg.Go(func() {
			c.proccess(tag, url)
		})
	}

	if h.Data == "a" {
		attrs := h.Attr
		for _, attr := range attrs {
			if attr.Key != "href" {
				continue
			}

			domain := getDomain(c.originURL)

			u := attr.Val

			if u[0] != 'h' {
				u = fmt.Sprintf("http://%s/%s", domain, u)
			}

			if !strings.Contains(u, domain) {
				continue
			}

			c.mu.Lock()
			_, ok := c.websites[u]
			c.mu.Unlock()

			if ok {
				break
			}

			n, err := parseHTML(u)
			if err != nil {
				fmt.Println(err)
				return
			}

			c.wg.Go(func() {
				c.proccess(n, u)
			})
		}
}
