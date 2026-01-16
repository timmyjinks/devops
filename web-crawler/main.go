package main

func main() {
	url := "http://localhost:8080"
	crawler := NewSkibidiCrawler(url)
	crawler.Crawl()
}
