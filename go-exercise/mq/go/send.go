package main

import (
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

func main() {
	conn, err := amqp.Dial("amqp://nana:nana@192.168.1.108:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

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
	failOnError(err, "Failed to declare a queue")

	confirms := ch.NotifyPublish(make(chan amqp.Confirmation, 1))
	if err := ch.Confirm(false); err != nil {
		log.Fatalf("confirm.select destination: %s", err)
	}

	for {
		//time.Sleep(time.Second * 1)
		body := "Hello World!"
		err = ch.Publish(
			"",     // exchange
			q.Name, // routing key
			false,  // mandatory
			false,  // immediate
			amqp.Publishing{
				DeliveryMode: amqp.Persistent,
				Timestamp:    time.Now(),
				ContentType:  "text/plain",
				Body:         []byte(body),
			})
		log.Printf(" [x] Sent %s", body)
		fmt.Println(err)

		// 这里判断是否收到。
		/*
			type Confirmation struct {
				DeliveryTag uint64 // A 1 based counter of publishings from when the channel was put in Confirm mode
				Ack         bool   // True when the server successfully received the publishing
			}
		*/
		if confirmed := <-confirms; confirmed.Ack {
			fmt.Println("aaa----", confirmed)
		} else {
			fmt.Println("bbb---- ", confirmed)
		}

	}
	failOnError(err, "Failed to publish a message")
}
