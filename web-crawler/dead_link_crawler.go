package main

import (
	"fmt"
	"golang.org/x/net/html"
	"strings"
	"sync"
	"time"
)

type DeadLinkCrawler struct {
	originURL string
	domain    string
	duration  time.Duration

	wg sync.WaitGroup

	vimu    sync.Mutex
	visited map[string]bool

	mu        sync.Mutex
	deadLinks []string
}

func getDomain(url string) string {
	return strings.Split(url, "/")[2]
}

func NewDeadLinkCrawler(url string) *DeadLinkCrawler {
	if url == "" {
		fmt.Println("URL is empty")
		return nil
	}

	return &DeadLinkCrawler{
		originURL: url,
		domain:    getDomain(url),
		visited:   make(map[string]bool),
	}
}

func (c *DeadLinkCrawler) Crawl() {
	h, err := parseHTML(c.originURL)
	if err != nil {
		fmt.Println(err)
		return
	}

	t := time.Now()

	c.wg.Go(func() {
		c.proccess(h)
	})

	c.wg.Wait()

	duration := time.Since(t)
	fmt.Println(duration.Seconds(), "seconds with go routines")

	printResult(c.deadLinks)
}

func (c *DeadLinkCrawler) proccess(n *html.Node) {
	if n.Data != "a" {
		iter := n.ChildNodes()
		for child := range iter {
			c.wg.Go(func() {
				c.proccess(child)
			})
		}
	}

	for _, attrb := range n.Attr {
		if attrb.Key != "href" {
			continue
		}

		url := attrb.Val
		if url[0] != 'h' {
			url = fmt.Sprintf("https://%s%s", c.domain, url)
		}

		if !c.validateLink(url) {
			continue
		}

		c.vimu.Lock()
		c.visited[url] = true
		c.vimu.Unlock()

		h, err := parseHTML(url)
		if err != nil {
			continue
		}

		c.wg.Go(func() {
			c.proccess(h)
		})
	}
}

func (c *DeadLinkCrawler) validateLink(url string) bool {
	c.vimu.Lock()
	_, ok := c.visited[url]
	c.vimu.Unlock()

	if ok {
		fmt.Printf("Already visited skipping: %s\n", url)
		return false
	}

	if c.isDeadLink(url) {
		fmt.Printf("Found dead link: %s\n", url)
		c.mu.Lock()
		c.deadLinks = append(c.deadLinks, url)
		c.mu.Unlock()
		return false
	}

	if !strings.Contains(url, c.domain) {
		return false
	}

	return true
}
