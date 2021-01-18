package main

import (
	"crypto/tls"
	"fmt"
	"github.com/streadway/amqp"
	"log"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func sendMessage(amqpHost string, msg string) {
	//conn, err := amqp.DialTLS("amqps://admin:Admin123@localhost1:5672/", &tls.Config{ InsecureSkipVerify: true } )
	conn, err := amqp.DialTLS(fmt.Sprintf("amqps://%s/", amqpHost), &tls.Config{InsecureSkipVerify: true})
	failOnError(err, fmt.Sprintf("Failed to connect to RabbitMQ to %s", amqpHost))
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"hello", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	failOnError(err, "Failed to declare a queue")

	body := fmt.Sprintf("%s %s!", msg, amqpHost)
	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	failOnError(err, "Failed to publish a message")
	log.Printf(" [x] Sent %s", body)

}

func main() {
	sendMessage("admin:Admin123@localhost1.localdomain:5672", "Hello")
	sendMessage("admin:Admin123@localhost2.localdomain:5672", "Hello")
}
