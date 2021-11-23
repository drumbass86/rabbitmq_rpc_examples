#!/bin/bash

docker run -d --hostname rabbitmq-server \
--rm --name rabbitmq -p 5672:5672 -p 15672:15672 \
-v /home/user/rabbitmq/rabbitmq.conf:/etc/rabbitmq/rabbitmq.conf \
rabbitmq:3.9-management

docker inspect --format='{{.NetworkSettings.Networks.bridge.IPAddress}}' rabbitmq
