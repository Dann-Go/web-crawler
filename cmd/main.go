package main

import (
	"github.com/Dann-Go/web-crawler/pkg/crawler"
	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Error(err)
	}

	ch, err := conn.Channel()
	if err != nil {
		log.Error(err)
	}
	crawler.Crawl("go.dev", 6, ch)
	defer ch.Close()
}
