package common

import "log"

const RPC_QUEUE_NAME = "rpc_queue"

func FailOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
