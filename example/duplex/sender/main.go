package main

import (
	"fmt"
	"log"
	"time"

	"github.com/joe-at-startupmedia/posix_mq"
	pmq_sender "github.com/joe-at-startupmedia/posix_mq/duplex/sender"
)

const maxSendTickNum = 10

func main() {
	err := pmq_sender.New("posix_mq_example_duplex", posix_mq.POSIX_MQ_DIR, posix_mq.Ownership{})
	if err != nil {
		log.Fatal("Sender: could not initialize: ", err)
	}
	defer pmq_sender.Close()

	count := 0
	for {
		count++
		request := fmt.Sprintf("Hello, World : %d\n", count)
		err := pmq_sender.Send([]byte(request), 0)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Sender: sent a new request: %s", request)

		msg, _, err := pmq_sender.WaitForResponse(time.Second)

		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Sender: got a response: %s\n", msg)

		if count >= maxSendTickNum {
			break
		}

		time.Sleep(1 * time.Second)
	}
}
