package main

import (
	"encoding/json"
	"github.com/gocolly/colly"
	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"sync"
)

func main() {
	type task struct {
		Url string `json:"url"`
	}

	type collection struct {
		Url   string `json:"url"`
		Title string `json:"title"`
	}

	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Error(err)
	}

	ch, err := conn.Channel()
	if err != nil {
		log.Error(err)
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

	ch.Qos(
		1,
		0,
		false)

	msgs, err := ch.Consume(
		q.Name,
		"",
		false,
		false,
		false,
		false,
		nil)

	if err != nil {
		log.Error(err)
	}

	q, err = ch.QueueDeclare(
		"collector",
		true,
		false,
		false,
		false,
		nil)

	collector := colly.NewCollector(
		colly.AllowedDomains("go.dev"),
	)
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func(c colly.Collector, queue amqp.Queue) {
		for d := range msgs {
			res := task{}
			coll := collection{}
			json.Unmarshal(d.Body, &res)

			c.OnHTML("title", func(element *colly.HTMLElement) {
				title := element.Text
				coll.Title = title
			})

			c.Visit(res.Url)

			coll.Url = res.Url

			log.Println(coll)

			send, err := json.Marshal(coll)
			if err != nil {
				log.Error(err)
			}
			err = ch.Publish(
				"",
				q.Name,
				false,
				false,
				amqp.Publishing{
					ContentType: "text/plain",
					Body:        send,
				})
			if err != nil {
				log.Error(err)
			}
			d.Ack(false)
		}
	}(*collector, q)
	wg.Wait()
}
