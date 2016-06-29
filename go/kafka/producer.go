package main

import (
	"log"
	"os"
	"os/signal"
	"sync"

	"github.com/Shopify/sarama"
)

var (
	kafkaAddrs = []string{"localhost:9092"}
	topic      = "test_topic"
)

func main() {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	producer, err := sarama.NewAsyncProducer(kafkaAddrs, config)
	if err != nil {
		panic(err)
	}

	// trap SIGINT
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	var (
		wg                          sync.WaitGroup
		enqueued, successes, errors int
	)
	wg.Add(1)
	go func() {
		defer wg.Done()
		for err := range producer.Errors() {
			log.Println(err)
			errors++
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for _ = range producer.Successes() {
			successes++
		}
	}()

	for {
		message := &sarama.ProducerMessage{Topic: topic, Value: sarama.StringEncoder("testing 123")}
		select {
		case producer.Input() <- message:
			enqueued++
		case <-signals:
			producer.AsyncClose()
			break
		}
	}
	wg.Wait()
	log.Printf("Successfully produced: %d, errors: %d\n", successes, errors)
}
