package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"

	"github.com/streadway/amqp"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

type RMQMessage struct {
	message            string
	queue              string
	queueDurable       bool
	queueAutoDelete    bool
	queueInternal      bool
	queueNoWait        bool
	exchangeName       string
	exchangeKind       string
	exchangeDurable    bool
	exchangeAutoDelete bool
	exchangeInternal   bool
	exchangeNoWait     bool
	exchangeArgs       amqp.Table
}

func sendMessage(conn *amqp.Connection, msg RMQMessage) {
	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	if msg.exchangeName != "" {
		err := ch.ExchangeDeclare(msg.exchangeName, msg.exchangeKind, msg.exchangeDurable, msg.exchangeAutoDelete, msg.exchangeInternal, msg.exchangeNoWait, msg.exchangeArgs)
		failOnError(err, "Failed to ExchangeDeclare")
	}

	q, err := ch.QueueDeclare(
		msg.queue,           // name
		msg.queueDurable,    // durable
		msg.queueAutoDelete, // delete when unused
		msg.queueInternal,   // exclusive
		msg.queueNoWait,     // no-wait
		nil,                 // arguments
	)
	failOnError(err, "Failed to declare a queue")

	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(msg.message),
		})
	failOnError(err, "Failed to publish a message")
	log.Printf(" [x] Sent %s", msg.message)

}

func main() {
	rmqURL := flag.String("url", "amqps://admin:Admin123@localhost1.localdomain:5672", "rabbitmq url in the format <proto>://username:password@hostname:port/ - proto can be amqp or amqps")
	queueName := flag.String("queue", "test-rmq", "The queue name you want to create. Default: test-rmq")
	message := flag.String("msg", "Hello this is test-rmq prog", "Message to send. Default: 'Hello this is test-rmq prog'")
	exchange := flag.String("ex", "", "exchange name. Default: empty that is dont use it")
	exchangeKind := flag.String("exkind", "direct", `exchange kind. Default: direct. Possible values:  "fanout", "topic" and "headers".`)
	exchangeDurabe := flag.Bool("exdurable", false, "exchange durable. Default: false")
	exchangeAutoDelete := flag.Bool("exautodelete", false, "exchange auto delete. Default: false")
	exchangeInternal := flag.Bool("exinternal", false, "exchange internal. Default: false")
	exchangeNoWait := flag.Bool("exnowait", false, "exchange nowait. Default: false")

	queueDurabe := flag.Bool("qdurable", false, "queue durable. Default: false")
	queueAutoDelete := flag.Bool("qautodelete", false, "queue auto delete. Default: false")
	queueInternal := flag.Bool("qinternal", false, "queue internal. Default: false")
	queueNoWait := flag.Bool("qnowait", false, "queue nowait. Default: false")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "test-rmq: Test the rabbitmq protocol by sending a message\n")
		flag.PrintDefaults()
	}

	flag.Parse()

	_u, err := url.Parse(*rmqURL)
	proto := _u.Scheme
	var conn *amqp.Connection

	if proto == "amqps" {
		conn, err = amqp.DialTLS(*rmqURL, &tls.Config{InsecureSkipVerify: true})
	} else {
		conn, err = amqp.Dial(*rmqURL)
	}
	failOnError(err, fmt.Sprintf("Failed to connect to RabbitMQ using the URL '%s'", *rmqURL))
	defer conn.Close()
	var myMsg RMQMessage

	myMsg = RMQMessage{
		message:            *message,
		queue:              *queueName,
		queueAutoDelete:    *queueAutoDelete,
		queueDurable:       *queueDurabe,
		queueInternal:      *queueInternal,
		queueNoWait:        *queueNoWait,
		exchangeName:       *exchange,
		exchangeKind:       *exchangeKind,
		exchangeAutoDelete: *exchangeAutoDelete,
		exchangeDurable:    *exchangeDurabe,
		exchangeInternal:   *exchangeInternal,
		exchangeNoWait:     *exchangeNoWait,
	}

	sendMessage(conn, myMsg)

	// sendMessage("admin:Admin123@localhost1.localdomain:5672", "Hello")
	// sendMessage("admin:Admin123@localhost2.localdomain:5672", "Hello")
}
