package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	cm "github.com/rpc_rabbitmq/common"
	"github.com/rpc_rabbitmq/messages"

	"github.com/streadway/amqp"
)

func fib(n int) int {
	if n == 0 {
		return 0
	} else if n == 1 {
		return 1
	} else {
		return fib(n-1) + fib(n-2)
	}
}

func main() {
	hostServ := flag.String("host", "localhost", " RabbitMQ server name")
	flag.Parse()
	connStr := fmt.Sprintf("amqp://guest:guest@%v:5672/", *hostServ)
	conn, err := amqp.Dial(connStr)
	cm.FailOnError(err, "Failed connect to RabbitMQ : "+connStr)
	defer conn.Close()

	ch, err := conn.Channel()
	cm.FailOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		cm.RPC_QUEUE_NAME, // name
		false,             // durable
		false,             // delete when unused
		false,             // exclusive
		false,             // no-wait
		nil,               // arguments
	)
	cm.FailOnError(err, "Failed to declare a queue")

	err = ch.Qos(
		2,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	cm.FailOnError(err, "Failed to set QoS")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	cm.FailOnError(err, "Failed to register a consumer")

	forever := make(chan bool)

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		for d := range msgs {
			go func(ctx context.Context, msg amqp.Delivery) {
				respMsg := messages.MsgResponse{0, 0, ""}
				var reqMsg messages.MsgRequest
				log.Printf("Receive request:%s", string(msg.Body))
				if err := json.Unmarshal(msg.Body, &reqMsg); err != nil {
					respMsg.ErrorText = fmt.Sprintf("Failed to convert msg.Body:%v to MsqRequest err:%s",
						msg.Body, err)
					cm.FailOnError(err, "Failed to convert body to MsqRequest")
				} else {
					log.Printf("Parsing request for calculate fib(%d)", reqMsg.Value)
					respMsg.Value = reqMsg.Value
					chFib := make(chan int)
					go func() {
						chFib <- fib(reqMsg.Value)
						close(chFib)
					}()

					select {
					case respMsg.Result = <-chFib:
						log.Printf("Calc fib(%d) = %d\n", reqMsg.Value, respMsg.Result)
					case <-time.After(5 * time.Second):
						respMsg.ErrorText = fmt.Sprintf("Timed out! Can`t calc fib(%v)", reqMsg.Value)
						log.Printf("%s\n", respMsg.ErrorText)
					case <-ctx.Done():
						respMsg.ErrorText = fmt.Sprintf("Context is done! Can`t calc fib(%v)", reqMsg.Value)
					}

					//fmt.Printf("Publish resp to:%v corId:%v", msg.ReplyTo, msg.CorrelationId)
				}

				respBody, err := json.Marshal(&respMsg)
				if err != nil {
					log.Println("Can`t marshal response to json")
					respBody = []byte("Can`t marshal response to json")
				}
				err = ch.Publish(
					"",          // exchange
					msg.ReplyTo, // routing key
					false,       // mandatory
					false,       // immediate
					amqp.Publishing{
						ContentType:   "text/plain",
						CorrelationId: msg.CorrelationId,
						Body:          respBody,
					})
				cm.FailOnError(err, "Failed to publish a message")

				msg.Ack(false)
			}(ctx, d)
		}
	}()

	log.Printf(" [*] Awaiting RPC requests")

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for sig := range c {
			cancel()
			_ = sig
			time.Sleep(1 * time.Second)
			close(forever)
			os.Exit(1)
		}
	}()

	<-forever
}
