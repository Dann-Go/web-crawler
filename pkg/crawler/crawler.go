package crawler

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
	"time"
)

func Crawl(url string, depth int) interface{} {
	type data interface{}

	storage := make([]data, 0)

	c := colly.NewCollector(
		colly.AllowedDomains(url))
	colly.MaxDepth(depth)
	c.OnHTML("article", func(element *colly.HTMLElement) {
		metaTags := element.DOM.ParentsUntil("~").Find("meta")
		metaTags.Each(func(_ int, selection *goquery.Selection) {
		})
		storage = append(storage, element.Response.Body)
	})

	c.OnHTML("a[href]", func(element *colly.HTMLElement) {
		element.Request.Visit(element.Attr("href"))
	})

	c.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		RandomDelay: 1 * time.Second,
	})

	c.OnRequest(func(request *colly.Request) {
		fmt.Println("Crawl on page", request.URL.String())

	})

	c.Visit("http://" + url)

	return storage
}
