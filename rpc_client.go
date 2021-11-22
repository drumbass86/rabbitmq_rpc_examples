package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"

	cm "github.com/rpc_rabbitmq/common"
	"github.com/rpc_rabbitmq/messages"
	"github.com/streadway/amqp"
)

func intFromArgs(args []string) int {
	const defaultN = 30
	if (len(args) < 3) || args[2] == "" {
		return defaultN
	} else {
		n, err := strconv.Atoi(args[2])
		if err != nil {
			log.Fatalf("Can`t convert from arg:%v to int\n", args[3])
			n = defaultN
		}
		return n
	}
}

func randCorrId(l int) string {
	bytes := make([]byte, l)
	for i := 0; i < l; i++ {
		bytes[i] = byte(randInt(35, 90))
	}
	return string(bytes)
}

func randInt(min, max int) int {
	return min + rand.Intn(max-min)
}

func fibonacciRPC(n int) (res int, err error) {
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
		"",    // name
		false, // durable
		false, // delete when unused
		true,  // exclusive
		false, // noWait
		nil,   // arguments
	)
	cm.FailOnError(err, "Failed to declare q queue")

	msgResp, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	cm.FailOnError(err, "Failed to register a consumer")

	corrId := randCorrId(32)

	reqMsg := messages.MsgRequest{n}
	bodyReq, err := json.Marshal(reqMsg)
	if err != nil {
		log.Println("Can`t json marshal request")
		return
	}
	err = ch.Publish(
		"",
		cm.RPC_QUEUE_NAME,
		false,
		false,
		amqp.Publishing{
			ContentType:   "text/plain",
			CorrelationId: corrId,
			ReplyTo:       q.Name,
			Body:          bodyReq,
		},
	)
	cm.FailOnError(err, "Failed to publish a message")

	for resp := range msgResp {
		if corrId == resp.CorrelationId {
			log.Printf("Receive response:%v", string(resp.Body))
			var respMsg messages.MsgResponse
			if err = json.Unmarshal(resp.Body, &respMsg); err != nil {
				cm.FailOnError(err, "Failed to unmarshal json response")
			} else {
				if respMsg.ErrorText == "" && respMsg.Value == reqMsg.Value {
					log.Printf("Calculate Fibonacci(%v) = %v", respMsg.Value, respMsg.Result)
				} else {
					log.Printf("Error(server):%s! Calculate Fibonacci(%v)", respMsg.ErrorText, reqMsg.Value)
					res = respMsg.Result
				}
			}

			break
		}
	}
	return
}

func main() {

	rand.Seed(time.Now().UTC().UnixNano())
	n := intFromArgs(os.Args)
	_, err := fibonacciRPC(n)
	cm.FailOnError(err, "Failed to handle RPC request")
	//log.Printf(" Send to server calculate Fibonacci(%v) = %v", n, res)
}
