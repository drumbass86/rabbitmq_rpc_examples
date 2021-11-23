# RPC examples for RabbitMQ

## Description
RPC example using RabbitMQ.
In this example use the [6](https://www.rabbitmq.com/tutorials/tutorial-six-go.html) tutorial RabbitMQ with some changes:
 - exchange between client and server on json messages - `messages/messages.go`
 - each request server processing in separate goroutine
 - time spending on calculating Fibonacci not more than 5 seconds
 - server handle Ctrl+C

## Usage
0. In `docker_rabbitmq.sh` set mount path to our `cfg/rabbitmq.conf`
1. Starting RabbitMQ server:
```
./docker_rabbitmq.sh
```

2. Starting rpc_server:
```
go run rpc_server.go -host=[localhost|ip_rabbitmq_daemon]
```

3. Starting client and sending request to calculate Fibonacci(N):
```
go run rpc_client.go -host=[localhost|ip_rabbitmq_daemon] 34
```
where 34 - param value N  in range [0 - .....]

4. In terminal client see:
   ```
   2021/11/23 10:14:54 Receive response:{"val":34,"res":5702887,"error":""}
   2021/11/23 10:14:54 Calculate Fibonacci(34) = 5702887
   ```

5. Well done

