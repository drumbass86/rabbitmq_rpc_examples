# RPC examples for RabbitMQ

## Description
RPC example using RabbitMQ.
In this example use the [6](https://www.rabbitmq.com/tutorials/tutorial-six-go.html) tutorial RabbitMQ with some changes:
 - exchange between client and server on json messages - `messages/messages.go`
 - each request server processing in separate goroutine
 - time spending on calculating Fibonacci not more than 5 seconds.
 - Server handle Ctrl+C
 - 

## Usage
Starting server:
```
go run rpc_server.go -host=[localhost|ip_rabbitmq_daemon]
```

Starting client and sending reuest calculate Fibonacci(N):
```
go run rpc_client.go -host=[localhost|ip_rabbitmq_daemon] N
```
where N - 0 - ......