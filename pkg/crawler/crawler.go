package crawler

import (
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"time"
)

func Crawl(url string, depth int, ch *amqp.Channel) {

	type task struct {
		Url string `json:"url"`
	}

	q, err := ch.QueueDeclare(
		"tasks",
		true,
		false,
		false,
		false,
		nil)

	if err != nil {
		log.Error(err)
	}

	c := colly.NewCollector(
		colly.AllowedDomains(url))
	colly.MaxDepth(depth)
	c.OnHTML("article", func(element *colly.HTMLElement) {
		metaTags := element.DOM.ParentsUntil("~").Find("meta")
		metaTags.Each(func(_ int, selection *goquery.Selection) {
		})
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

		send, err := json.Marshal(task{
			Url: request.URL.String(),
		})
		if err != nil {
			log.Error(err)
		}
		err = ch.Publish(
			"",
			q.Name,
			false,
			false,
			amqp.Publishing{
				DeliveryMode: amqp.Persistent,
				ContentType:  "text/plain",
				Body:         send,
			})
		if err != nil {
			log.Error(err)
		}
	})

	c.Visit("http://" + url)

}
