package main

import (
	"fmt"
	"golang.org/x/net/html"
	"net/http"
)

type Crawler interface {
	Crawl()
}

type CrawlerService struct {
	crawler Crawler
}

func NewCrawlerService(crawler *Crawler, url string) *CrawlerService {
	return &CrawlerService{
		crawler: *crawler,
	}

}

func parseHTML(url string) (*html.Node, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	h, err := html.Parse(resp.Body)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	return h, nil
}

func printResult(result []string) {
	fmt.Println("Dead Links:")
	for _, link := range result {
		fmt.Println(link)
	}
}

func (c *DeadLinkCrawler) isDeadLink(url string) bool {
	resp, err := http.Get(url)
	if err != nil {
		return true
	}

	if resp.StatusCode >= 400 {
		return true
	}

	fmt.Printf("Found link: %s\n", url)
	return false
}
