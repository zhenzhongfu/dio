package main

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/streadway/amqp"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

var (
	conn *amqp.Connection
	ch   *amqp.Channel
	err  error
)

func reconnect() (<-chan amqp.Delivery, error) {
	var err error
	conn, err = amqp.Dial("amqp://nana:nana@192.168.1.108:5672/")
	if err != nil {
		return nil, errors.New("Failed to connect to RabbitMQ")
	}

	ch, err = conn.Channel()
	if err != nil {
		return nil, errors.New("Failed to connect to RabbitMQ")
	}

	m := make(amqp.Table)
	m["x-max-length"] = int64(2)
	m["x-overflow"] = "reject-publish"
	q, err := ch.QueueDeclare(
		"hello33", // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		m,         // arguments
	)
	if err != nil {
		return nil, err
	}

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		m,      // args
	)
	if err != nil {
		return nil, errors.New("Failed to register a consumer")
	}

	return msgs, nil
}

func main() {
	/*
		//conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
		conn, err := amqp.Dial("amqp://nana:nana@192.168.1.108:5672/")
		failOnError(err, "Failed to connect to RabbitMQ")
		defer conn.Close()

		ch, err := conn.Channel()
		failOnError(err, "Failed to open a channel")
		defer ch.Close()
	*/
	msgs, err := reconnect()
	if err != nil {
		fmt.Println("errrrrrr:", err)
		return
	}
	defer conn.Close()
	defer ch.Close()

	forever := make(chan bool)

	go func() {
		for {
		RECONNECT:
			select {
			case msg, ok := <-msgs:
				if !ok {
					fmt.Printf("what???-------------------\n")
					for {
						time.Sleep(time.Second * 1)
						msgs, err = reconnect()
						if err != nil {
							fmt.Println(err)
						} else if err == nil {
							goto RECONNECT
						}
					}
				}
				time.Sleep(time.Second * 3)
				msg.Reject(false)
				log.Printf("Received a message: %s", msg.Body)
			}
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
